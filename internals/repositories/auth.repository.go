package repositories

import (
	"context"
	"errors"
	"log"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/raihaninkam/finalPhase3/internals/models"
	"github.com/redis/go-redis/v9"
)

type AuthRepositories struct {
	db  *pgxpool.Pool
	rdb *redis.Client
}

func NewAuthRepository(db *pgxpool.Pool, rdb *redis.Client) *AuthRepositories {
	return &AuthRepositories{
		db:  db,
		rdb: rdb,
	}
}

func (ar *AuthRepositories) GetEmail(ctx context.Context, email string) (*models.User, error) {
	sql := `SELECT id, email, password, name, avatar_url, bio FROM users WHERE email =$1`

	var user models.User
	if err := ar.db.QueryRow(ctx, sql, email).Scan(&user.ID, &user.Email, &user.Password, &user.Name, &user.AvatarUrl, &user.Bio); err != nil {
		if err == pgx.ErrNoRows {
			return nil, errors.New("user not found")
		}
		log.Println("Internal Server Error.\nCause: ", err.Error())
		return nil, err
	}
	return &user, nil
}

func (a *AuthRepositories) CheckEmailExists(rctx context.Context, email string) (bool, error) {
	sql := "SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)"

	var exists bool
	if err := a.db.QueryRow(rctx, sql, email).Scan(&exists); err != nil {
		log.Println("Error checking email existence:", err.Error())
		return false, err
	}

	return exists, nil
}

func (ar *AuthRepositories) CreateAccount(ctx context.Context, user *models.User) error {
	sql := `INSERT INTO users (email, password) VALUES ($1, $2) RETURNING id`

	var userId string
	if err := ar.db.QueryRow(ctx, sql, user.Email, user.Password).Scan(&userId); err != nil {
		log.Println("Failed to create account.\nCause: ", err.Error())
		return err
	}

	user.ID = userId
	return nil
}

func (ar *AuthRepositories) UpdateAccount(ctx context.Context, updateData *models.UserUpdate) (*models.User, error) {

	query := `
		UPDATE users 
		SET name = COALESCE(NULLIF($1, ''), name),
			bio = COALESCE(NULLIF($2, ''), bio),
			avatar_url = COALESCE(NULLIF($3, ''), avatar_url),
			updated_at = NOW()
		WHERE id = $4
		RETURNING id, email, password, name, avatar_url, bio
	`

	var user models.User
	err := ar.db.QueryRow(ctx, query,
		updateData.Name,
		updateData.Bio,
		updateData.AvatarUrl,
		updateData.ID,
	).Scan(
		&user.ID,
		&user.Email,
		&user.Password,
		&user.Name,
		&user.AvatarUrl,
		&user.Bio,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, errors.New("user not found")
		}
		log.Println("Failed to update account.\nCause: ", err.Error())
		return nil, err
	}

	return &user, nil
}
