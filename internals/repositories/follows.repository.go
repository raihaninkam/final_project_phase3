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

type FollowRepository struct {
	db  *pgxpool.Pool
	rdb *redis.Client
}

func NewFollowsRepository(db *pgxpool.Pool, rdb *redis.Client) *FollowRepository {
	return &FollowRepository{
		db:  db,
		rdb: rdb,
	}
}

func (fr *FollowRepository) GetFollowing(ctx context.Context, userID string) ([]models.UserProfile, error) {
	// Try cache first
	cacheKey := fmt.Sprintf("following:%s", userID)
	cached, err := fr.rdb.Get(ctx, cacheKey).Result()
	if err == nil {
		var users []models.UserProfile
		if err := json.Unmarshal([]byte(cached), &users); err == nil {
			return users, nil
		}
	}

	// Get from database
	sql := `
		SELECT u.id, u.name, u.avatar_url, u.bio
		FROM follows f
		JOIN users u ON f.following_id = u.id
		WHERE f.follower_id = $1
		ORDER BY f.created_at DESC
	`

	rows, err := fr.db.Query(ctx, sql, userID)
	if err != nil {
		log.Printf("Query error: %v", err) // Log error
		return nil, fmt.Errorf("failed to get following: %w", err)
	}
	defer rows.Close()

	var users []models.UserProfile
	for rows.Next() {
		var user models.UserProfile
		if err := rows.Scan(&user.Id, &user.Name, &user.Avatar, &user.Bio); err != nil {
			log.Printf("Scan error: %v", err) // Log error
			return nil, fmt.Errorf("failed to scan user: %w", err)
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		log.Printf("Rows error: %v", err) // Log error
		return nil, err
	}

	// Cache the result
	if len(users) > 0 {
		usersJSON, _ := json.Marshal(users)
		fr.rdb.Set(ctx, cacheKey, usersJSON, 5*time.Minute)
	}

	return users, nil
}

func (fr *FollowRepository) Follow(ctx context.Context, followerID, followingID string) error {
	sql := `INSERT INTO follows (follower_id, following_id, created_at) 
	        VALUES ($1, $2, now())`

	_, err := fr.db.Exec(ctx, sql, followerID, followingID)
	if err != nil {
		return fmt.Errorf("failed to follow user: %w", err)
	}

	// Invalidate cache
	fr.rdb.Del(ctx, fmt.Sprintf("following:%s", followerID))
	fr.rdb.Del(ctx, fmt.Sprintf("followers:%s", followingID))

	return nil
}
