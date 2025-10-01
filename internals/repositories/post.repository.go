package repositories

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/raihaninkam/finalPhase3/internals/models"
	"github.com/redis/go-redis/v9"
)

type PostRepository struct {
	db  *pgxpool.Pool
	rdb *redis.Client
}

func NewPostRepository(db *pgxpool.Pool, rdb *redis.Client) *PostRepository {
	return &PostRepository{
		db:  db,
		rdb: rdb,
	}
}

func (pr *PostRepository) CreatePost(ctx context.Context, post *models.Posts) (*models.Posts, error) {
	sql := `INSERT INTO posts (user_id, content_text, image_url, created_at) 
	        VALUES ($1, $2, $3, now()) 
	        RETURNING id, created_at`

	err := pr.db.QueryRow(ctx, sql, post.UserId, post.Content, post.ImageUrl).Scan(&post.Id, &post.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to create post: %w", err)
	}

	// Invalidate cache after creating new post
	pr.rdb.Del(ctx, "posts:all")

	return post, nil
}

func (pr *PostRepository) GetPostByID(ctx context.Context, id int) (*models.Posts, error) {
	// Try to get from cache first
	cacheKey := fmt.Sprintf("post:%d", id)
	cached, err := pr.rdb.Get(ctx, cacheKey).Result()
	if err == nil {
		var post models.Posts
		if err := json.Unmarshal([]byte(cached), &post); err == nil {
			return &post, nil
		}
	}

	// If not in cache, get from database
	sql := `SELECT id, content_text, image_url, created_at FROM posts WHERE id = $1`

	var post models.Posts
	err = pr.db.QueryRow(ctx, sql, id).Scan(&post.Id, &post.Content, &post.ImageUrl, &post.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to get post: %w", err)
	}

	// Store in cache
	postJSON, _ := json.Marshal(post)
	pr.rdb.Set(ctx, cacheKey, postJSON, 1*time.Hour)

	return &post, nil
}

func (pr *PostRepository) GetAllPosts(ctx context.Context) ([]models.Posts, error) {
	// Try to get from cache first
	cacheKey := "posts:all"
	cached, err := pr.rdb.Get(ctx, cacheKey).Result()
	if err == nil {
		var posts []models.Posts
		if err := json.Unmarshal([]byte(cached), &posts); err == nil {
			return posts, nil
		}
	}

	// If not in cache, get from database
	sql := `SELECT id, content_text, image_url, created_at FROM posts ORDER BY created_at DESC`

	rows, err := pr.db.Query(ctx, sql)
	if err != nil {
		return nil, fmt.Errorf("failed to get posts: %w", err)
	}
	defer rows.Close()

	var posts []models.Posts
	for rows.Next() {
		var post models.Posts
		if err := rows.Scan(&post.Id, &post.Content, &post.ImageUrl, &post.CreatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan post: %w", err)
		}
		posts = append(posts, post)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating posts: %w", err)
	}

	// Store in cache
	postsJSON, _ := json.Marshal(posts)
	pr.rdb.Set(ctx, cacheKey, postsJSON, 5*time.Minute)

	return posts, nil
}

func (pr *PostRepository) UpdatePost(ctx context.Context, id int, post *models.Posts) (*models.Posts, error) {
	sql := `UPDATE posts SET content_text = $1, image_url = $2 WHERE id = $3 RETURNING id, content_text, image_url, created_at`

	err := pr.db.QueryRow(ctx, sql, post.Content, post.ImageUrl, id).Scan(&post.Id, &post.Content, &post.ImageUrl, &post.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to update post: %w", err)
	}

	// Invalidate cache
	pr.rdb.Del(ctx, fmt.Sprintf("post:%d", id))
	pr.rdb.Del(ctx, "posts:all")

	return post, nil
}

func (pr *PostRepository) DeletePost(ctx context.Context, id int) error {
	sql := `DELETE FROM posts WHERE id = $1`

	result, err := pr.db.Exec(ctx, sql, id)
	if err != nil {
		return fmt.Errorf("failed to delete post: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("post not found")
	}

	// Invalidate cache
	pr.rdb.Del(ctx, fmt.Sprintf("post:%d", id))
	pr.rdb.Del(ctx, "posts:all")

	return nil
}

func (pr *PostRepository) GetFollowingPosts(ctx context.Context, userID string, limit, offset int) ([]models.PostWithUser, error) {
	// Try cache first
	cacheKey := fmt.Sprintf("feed:%s:%d:%d", userID, limit, offset)
	cached, err := pr.rdb.Get(ctx, cacheKey).Result()
	if err == nil {
		var posts []models.PostWithUser
		if err := json.Unmarshal([]byte(cached), &posts); err == nil {
			return posts, nil
		}
	}

	// Get from database
	sql := `
		SELECT 
			p.id, 
			p.user_id, 
			p.content_text, 
			p.image_url, 
			p.created_at,
			u.name as user_name,
			COALESCE(u.avatar_url, '') as user_avatar
		FROM posts p
		JOIN users u ON p.user_id = u.id
		WHERE p.user_id IN (
			SELECT following_id 
			FROM follows 
			WHERE follower_id = $1
		)
		ORDER BY p.created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := pr.db.Query(ctx, sql, userID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get following posts: %w", err)
	}
	defer rows.Close()

	var posts []models.PostWithUser
	for rows.Next() {
		var post models.PostWithUser
		if err := rows.Scan(
			&post.Id,
			&post.UserId,
			&post.Content,
			&post.ImageUrl,
			&post.CreatedAt,
			&post.UserName,
			&post.UserAvatar,
		); err != nil {
			return nil, fmt.Errorf("failed to scan post: %w", err)
		}
		posts = append(posts, post)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	// Cache the result
	if len(posts) > 0 {
		postsJSON, _ := json.Marshal(posts)
		pr.rdb.Set(ctx, cacheKey, postsJSON, 2*time.Minute)
	}

	return posts, nil
}

// GetPopularPosts mendapatkan postingan dengan interaksi tinggi
func (pr *PostRepository) GetPopularPosts(ctx context.Context, limit, offset int) ([]*models.Posting, error) {
	query := `
		SELECT 
			p.id, p.user_id, p.content_text, p.image_url, p.created_at, p.updated_at,
			u.name as user_name, u.avatar_url as user_avatar_url,
			COUNT(DISTINCT l.id) as like_count,
			COUNT(DISTINCT c.id) as comment_count,
			COUNT(DISTINCT f.follower_id) as follower_count,
			(COUNT(DISTINCT l.id) * 1.0 + COUNT(DISTINCT c.id) * 2.0 + COUNT(DISTINCT f.follower_id) * 0.5) as popularity_score
		FROM posts p
		INNER JOIN users u ON p.user_id = u.id
		LEFT JOIN likes l ON p.id = l.post_id
		LEFT JOIN comments c ON p.id = c.post_id
		LEFT JOIN follows f ON p.user_id = f.following_id
		WHERE p.created_at >= NOW() - INTERVAL '7 days'
		GROUP BY p.id, u.id
		HAVING COUNT(DISTINCT l.id) + COUNT(DISTINCT c.id) > 0
		ORDER BY popularity_score DESC, p.created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := pr.db.Query(ctx, query, limit, offset)
	if err != nil {
		log.Println("Failed to get popular posts:", err.Error())
		return nil, err
	}
	defer rows.Close()

	var posts []*models.Posting
	for rows.Next() {
		var post models.Posting
		var popularityScore float64

		err := rows.Scan(
			&post.ID,
			&post.UserID,
			&post.Content,
			&post.ImageUrl,
			&post.CreatedAt,
			&post.UpdatedAt,
			&post.UserName,
			&post.UserAvatarUrl,
			&post.LikeCount,
			&post.CommentCount,
			&post.FollowerCount,
			&popularityScore,
		)
		if err != nil {
			log.Println("Failed to scan post:", err.Error())
			continue
		}
		posts = append(posts, &post)
	}

	return posts, nil
}
