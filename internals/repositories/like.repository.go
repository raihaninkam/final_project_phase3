package repositories

import (
	"context"
	"errors"
	"fmt"
	"log"

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

func (lr *LikeRepository) GetPostLikesCount(ctx context.Context, postID string) (int, error) {
	sql := `SELECT COUNT(*) FROM likes WHERE post_id = $1`

	var count int
	err := lr.db.QueryRow(ctx, sql, postID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to get likes count: %w", err)
	}

	return count, nil
}

// CONCURRENCY SAFE: Update like counter dengan database transaction

// LikePost dengan concurrency safe increment
func (lr *LikeRepository) LikePost(ctx context.Context, userID, postID string) error {
	// Gunakan transaction untuk memastikan atomicity
	tx, err := lr.db.Begin(ctx)
	if err != nil {
		log.Println("Failed to begin transaction:", err.Error())
		return err
	}
	defer tx.Rollback(ctx)

	// Insert like
	query := `INSERT INTO likes (user_id, post_id) VALUES ($1, $2)`
	_, err = tx.Exec(ctx, query, userID, postID)
	if err != nil {
		log.Println("Anda sudah like postingan ini:", err.Error())
		return err
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		log.Println("Failed to commit transaction:", err.Error())
		return err
	}

	return nil
}

// // UnlikePost dengan concurrency safe decrement
func (lr *LikeRepository) UnlikePost(ctx context.Context, userID, postID string) error {
	// Gunakan transaction untuk memastikan atomicity
	tx, err := lr.db.Begin(ctx)
	if err != nil {
		log.Println("Failed to begin transaction:", err.Error())
		return err
	}
	defer tx.Rollback(ctx)

	// Delete like
	query := `DELETE FROM likes WHERE user_id = $1 AND post_id = $2`
	result, err := tx.Exec(ctx, query, userID, postID)
	if err != nil {
		log.Println("Failed to delete like:", err.Error())
		return err
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return errors.New("like not found")
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		log.Println("Failed to commit transaction:", err.Error())
		return err
	}

	return nil
}
