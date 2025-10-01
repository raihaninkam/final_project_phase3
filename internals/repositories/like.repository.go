package repositories

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

type LikeRepository struct {
	db  *pgxpool.Pool
	rdb *redis.Client
}

func NewLikeRepository(db *pgxpool.Pool, rdb *redis.Client) *LikeRepository {
	return &LikeRepository{
		db:  db,
		rdb: rdb,
	}
}

func (lr *LikeRepository) LikePost(ctx context.Context, userID, postID string) error {
	sql := `INSERT INTO likes (user_id, post_id, created_at) 
	        VALUES ($1, $2, now())
	        ON CONFLICT (user_id, post_id) DO NOTHING`

	result, err := lr.db.Exec(ctx, sql, userID, postID)
	if err != nil {
		return fmt.Errorf("failed to like post: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("post already liked")
	}

	// Invalidate cache
	lr.rdb.Del(ctx, fmt.Sprintf("likes:post:%s", postID))

	return nil
}

func (lr *LikeRepository) UnlikePost(ctx context.Context, userID, postID string) error {
	sql := `DELETE FROM likes WHERE user_id = $1 AND post_id = $2`

	result, err := lr.db.Exec(ctx, sql, userID, postID)
	if err != nil {
		return fmt.Errorf("failed to unlike post: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("like not found")
	}

	// Invalidate cache
	lr.rdb.Del(ctx, fmt.Sprintf("likes:post:%s", postID))

	return nil
}

func (lr *LikeRepository) GetPostLikesCount(ctx context.Context, postID string) (int, error) {
	sql := `SELECT COUNT(*) FROM likes WHERE post_id = $1`

	var count int
	err := lr.db.QueryRow(ctx, sql, postID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to get likes count: %w", err)
	}

	return count, nil
}
