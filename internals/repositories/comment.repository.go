package repositories

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/raihaninkam/finalPhase3/internals/models"
	"github.com/redis/go-redis/v9"
)

type CommentRepository struct {
	db  *pgxpool.Pool
	rdb *redis.Client
}

func NewCommentRepository(db *pgxpool.Pool, rdb *redis.Client) *CommentRepository {
	return &CommentRepository{
		db:  db,
		rdb: rdb,
	}
}

func (cr *CommentRepository) CreateComment(ctx context.Context, comment *models.Comment) (*models.Comment, error) {
	sql := `INSERT INTO comments (user_id, post_id, content, created_at, updated_at) 
	        VALUES ($1, $2, $3, now(), now()) 
	        RETURNING id, created_at, updated_at`

	err := cr.db.QueryRow(ctx, sql, comment.UserId, comment.PostId, comment.Content).
		Scan(&comment.Id, &comment.CreatedAt, &comment.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to create comment: %w", err)
	}

	// Invalidate cache
	cr.rdb.Del(ctx, fmt.Sprintf("comments:post:%s", comment.PostId))

	return comment, nil
}

func (cr *CommentRepository) GetPostComments(ctx context.Context, postID string) ([]models.CommentWithUser, error) {
	// Try cache first
	cacheKey := fmt.Sprintf("comments:post:%s", postID)
	cached, err := cr.rdb.Get(ctx, cacheKey).Result()
	if err == nil {
		var comments []models.CommentWithUser
		if err := json.Unmarshal([]byte(cached), &comments); err == nil {
			return comments, nil
		}
	}

	// Get from database
	sql := `
		SELECT 
			c.id, 
			c.user_id, 
			c.post_id, 
			c.content, 
			c.created_at,
			u.name as user_name,
			COALESCE(u.avatar_url, '') as user_avatar
		FROM comments c
		JOIN users u ON c.user_id = u.id
		WHERE c.post_id = $1
		ORDER BY c.created_at ASC
	`

	rows, err := cr.db.Query(ctx, sql, postID)
	if err != nil {
		return nil, fmt.Errorf("failed to get comments: %w", err)
	}
	defer rows.Close()

	var comments []models.CommentWithUser
	for rows.Next() {
		var comment models.CommentWithUser
		if err := rows.Scan(
			&comment.Id,
			&comment.UserId,
			&comment.PostId,
			&comment.Content,
			&comment.CreatedAt,
			&comment.UserName,
			&comment.UserAvatar,
		); err != nil {
			return nil, fmt.Errorf("failed to scan comment: %w", err)
		}
		comments = append(comments, comment)
	}

	// Cache the result
	if len(comments) > 0 {
		commentsJSON, _ := json.Marshal(comments)
		cr.rdb.Set(ctx, cacheKey, commentsJSON, 5*time.Minute)
	}

	return comments, nil
}

func (cr *CommentRepository) DeleteComment(ctx context.Context, commentID, userID string) error {
	sql := `DELETE FROM comments WHERE id = $1 AND user_id = $2`

	result, err := cr.db.Exec(ctx, sql, commentID, userID)
	if err != nil {
		return fmt.Errorf("failed to delete comment: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("comment not found or unauthorized")
	}

	return nil
}
