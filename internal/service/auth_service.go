package service

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/penonpaper/event-api/internal/domain"
	"golang.org/x/crypto/bcrypt"
)

type AuthUseCase struct {
	userRepo domain.UserRepository

	redisClient *redis.Client

	accessTokenSecret  string
	refreshTokenSecret string

	accessTokenTTL  time.Duration
	refreshTokenTTL time.Duration
}

func NewAuthUseCase(repo domain.UserRepository, accessTokenSecret, refreshTokenSecret string, accessTokenTTL, refreshTokenTTL time.Duration, rdb *redis.Client) domain.AuthService {
	return &AuthUseCase{
		userRepo:           repo,
		accessTokenSecret:  accessTokenSecret,
		refreshTokenSecret: refreshTokenSecret,
		accessTokenTTL:     accessTokenTTL,
		refreshTokenTTL:    refreshTokenTTL,
		redisClient:        rdb,
	}
}
func (a *AuthUseCase) generateAccessToken(userID string, userRole string) (string, error) {
	now := time.Now()

	claims := &domain.TokenClaims{
		Role: userRole,
		Type: "access",

		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   userID,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(a.accessTokenTTL)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(a.accessTokenSecret))

	if err != nil {
		return "", fmt.Errorf("failed to sign access token: %w", err)
	}
	return tokenString, nil
}

func (a *AuthUseCase) generateRefreshToken(ctx context.Context, userID string) (string, error) {
	now := time.Now()
	jti := uuid.NewString()
	claims := &domain.TokenClaims{
		Type: "refresh",

		RegisteredClaims: jwt.RegisteredClaims{
			ID:        jti,
			Subject:   userID,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(a.refreshTokenTTL)),
		},
	}
	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256, claims,
	)
	tokenString, err := token.SignedString([]byte(a.refreshTokenSecret))
	if err != nil {
		return "", fmt.Errorf("failed to sign refresh token: %w", err)
	}

	key := "refresh:" + jti
	err = a.redisClient.Set(
		ctx,
		key,
		userID,
		a.refreshTokenTTL,
	).Err()

	if err != nil {
		return "", fmt.Errorf("Failed to save refresh token in redis: %w", err)
	}
	return tokenString, nil
}

func (a *AuthUseCase) SignIn(ctx context.Context, email string, password string) (*domain.TokenPair, error) {
	user, err := a.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("%w: Failed to get data from db: %v", domain.ErrDataBase, err)

		//return nil, fmt.Errorf("Failed to get data from db: %w", err)
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))

	if err != nil {
		return nil, fmt.Errorf("%w: Failed to compare data: %v", domain.ErrConditions, err)

	}

	accessToken, err := a.generateAccessToken(user.ID, user.Role)
	if err != nil {
		return nil, fmt.Errorf("%w: Failed to generate accessToken: %v", domain.ErrAccessToken, err)

	}

	refreshToken, err := a.generateRefreshToken(ctx, user.ID)
	if err != nil {
		return nil, fmt.Errorf("%w: Failed to generate refreshToken: %v", domain.ErrRefreshToken, err)
	}

	return &domain.TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil

}
func (a *AuthUseCase) Refresh(ctx context.Context, refreshToken string) (*domain.TokenPair, error) {
	claims := &domain.TokenClaims{}

	token, err := jwt.ParseWithClaims(refreshToken,
		claims,
		func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				slog.Error("Unexpected signing method")
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(a.refreshTokenSecret), nil
		},
	)
	if err != nil {
		slog.Error("Invalid refresh token", "error", err)
		return nil, fmt.Errorf("%w: Invalid refresh token: %v", domain.ErrRefreshToken, err)
	}
	if !token.Valid {
		slog.Error("Invalid refresh token")
		return nil, fmt.Errorf("%w: Invalid refresh token: %v", domain.ErrRefreshToken, err)
	}

	if claims.Type != "refresh" {
		slog.Error("Invalid refresh token type", "Type", claims.Type)
		return nil, fmt.Errorf("%w: Invalid refresh token type: %v", domain.ErrRefreshToken, err)
	}
	userID := claims.Subject

	if userID == "" {
		slog.Error("UserID is empty")
		return nil, fmt.Errorf("%w: UserID is empty: %v", domain.ErrConditions, err)
	}

	// KEY - VALUE
	// "refresh:jti" - USERID
	jti := claims.ID

	if jti == "" {
		return nil, fmt.Errorf("%w: jti is empty: %v", domain.ErrConditions, err)
	}
	key := "refresh:" + jti

	userIdFromRedis, err := a.redisClient.Get(ctx, key).Result()
	if err == redis.Nil {
		slog.Error("UserID not found in redis")
		return nil, fmt.Errorf("%w: userID not found: %v", domain.ErrRedis, err)
	}

	if err != nil {
		slog.Error("UserID not found in redis", "error", err)

		return nil, fmt.Errorf("%w: Failed to get userID: %v", domain.ErrRedis, err)
	}

	if userIdFromRedis != userID {
		return nil, fmt.Errorf("%w: UserID dont match: %v", domain.ErrConditions, err)
	}

	err = a.redisClient.Del(ctx, key).Err()
	if err != nil {
		return nil, fmt.Errorf("%w: Failed to delete refresh token: %v", domain.ErrRedis, err)
	}

	parsedUserID, err := uuid.Parse(userID)
	if err != nil {
		slog.Error("Failed to convert to uuid", "error", err)
		return nil, fmt.Errorf("%w: Failed to contert to uuid: %v", domain.ErrInternal, err)
	}

	user, err := a.userRepo.GetByID(ctx, parsedUserID)
	if err != nil {
		slog.Error("Failed to get user role by ID", "error", err,
			"ID", parsedUserID)
		return nil, fmt.Errorf("%w: Failed to get role: %v", domain.ErrDataBase, err)
	}

	accessToken, err := a.generateAccessToken(userID, user.Role)
	if err != nil {
		slog.Error("Failed to generate access token", "error", err)
		return nil, fmt.Errorf("%w: Failed to generate access token: %v", domain.ErrAccessToken, err)
	}
	newRefreshToken, err := a.generateRefreshToken(ctx, userID)
	if err != nil {
		slog.Error("Failed to generate new refresh token", "error", err)
		return nil, fmt.Errorf("%w: Failed to generate new refresh token: %v", domain.ErrRefreshToken, err)
	}
	return &domain.TokenPair{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
	}, nil

}
func (a *AuthUseCase) Logout(ctx context.Context, refreshToken string) error {
	claims := &domain.TokenClaims{}

	token, err := jwt.ParseWithClaims(
		refreshToken,
		claims,
		func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf(
					"unexpected signing method: %v",
					token.Header["alg"],
				)
			}

			return []byte(a.refreshTokenSecret), nil
		},
	)

	if err != nil {
		return fmt.Errorf("%w: Invalid refresh token: %v", domain.ErrRefreshToken, err)
	}

	if !token.Valid {
		return fmt.Errorf("%w: Invalid refresh token: %v", domain.ErrRefreshToken, err)
	}

	if claims.Type != "refresh" {
		return fmt.Errorf("%w: Invalid refresh token type: %v", domain.ErrRefreshToken, err)

	}

	jti := claims.ID

	if jti == "" {
		return fmt.Errorf("%w: Refresh token jti is empty: %v", domain.ErrConditions, err)
	}

	key := "refresh:" + jti

	deleted, err := a.redisClient.Del(ctx, key).Result()

	if err != nil {
		return fmt.Errorf("%w: Failed to delete refresh token: %v", domain.ErrRedis, err)
	}

	if deleted == 0 {
		return fmt.Errorf("%w: Refresh token not found: %v", domain.ErrRedis, err)
	}

	slog.Info("Refresh token successfully deleted")

	return nil
}
