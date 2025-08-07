package utils

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/redis/go-redis/v9"
)

const BL = "blacklist:%s"

type CustomClaims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

type JWTService struct {
	secretKey string
	cache     redis.UniversalClient
}

func NewJWTService(secretKey string, rdb redis.UniversalClient) *JWTService {
	return &JWTService{
		secretKey: secretKey,
		cache:     rdb,
	}
}

func (j *JWTService) GenerateJWT(username string, expireDuration time.Duration) (string, error) {
	// Create custom claims
	claims := CustomClaims{
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expireDuration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "elotus-system",
			Subject:   username,
		},
	}

	// Create token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign token with secret key
	tokenString, err := token.SignedString([]byte(j.secretKey))
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return tokenString, nil
}

// ValidateJWT validates token signature, expiration, and blacklist status
func (j *JWTService) ValidateJWT(tokenString string) (*CustomClaims, error) {
	// First check if token is blacklisted (fast Redis check)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	key := fmt.Sprintf(BL, tokenString)
	exists := j.cache.Exists(ctx, key)
	if exists.Err() != nil {
		return nil, fmt.Errorf("failed to check blacklist: %w", exists.Err())
	}

	if exists.Val() > 0 {
		return nil, errors.New("TOKEN_BLACKLISTED")
	}

	// Parse and validate token signature and expiration
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Validate signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(j.secretKey), nil
	})

	if err != nil {
		// Check if error is due to expiration
		if errors.Is(err, jwt.ErrTokenExpired) {
			// Automatically blacklist expired token to prevent reuse
			go j.blacklistTokenAsync(tokenString) // Async to not block response
			return nil, errors.New("TOKEN_EXPIRED")
		}
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	// Extract claims
	claims, ok := token.Claims.(*CustomClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token claims")
	}

	return claims, nil
}

// GetJWTPayload extracts payload without validation (useful for expired tokens)
func (j *JWTService) GetJWTPayload(tokenString string) (*CustomClaims, error) {
	// Parse token without validation
	token, _, err := new(jwt.Parser).ParseUnverified(tokenString, &CustomClaims{})
	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	claims, ok := token.Claims.(*CustomClaims)
	if !ok {
		return nil, fmt.Errorf("invalid token claims")
	}

	return claims, nil
}

// BlacklistToken manually adds a token to blacklist (for logout)
func (j *JWTService) BlacklistToken(tokenString string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	// Calculate appropriate TTL based on token expiration
	claims, err := j.GetJWTPayload(tokenString)
	if err != nil {
		return fmt.Errorf("failed to get token payload: %w", err)
	}

	var ttl time.Duration
	if claims.ExpiresAt != nil {
		remaining := time.Until(claims.ExpiresAt.Time)
		if remaining > 0 {
			ttl = remaining + time.Hour // Add buffer
		} else {
			ttl = time.Hour // Minimum TTL for already expired tokens
		}
	} else {
		ttl = 24 * time.Hour // Default TTL for tokens without expiration
	}

	key := fmt.Sprintf(BL, tokenString)
	// Store minimal value to save memory
	result := j.cache.SetNX(ctx, key, "1", ttl)
	if result.Err() != nil {
		return fmt.Errorf("failed to blacklist token: %w", result.Err())
	}

	return nil
}

// blacklistTokenAsync blacklists token asynchronously (internal use)
func (j *JWTService) blacklistTokenAsync(tokenString string) {
	// Don't block the main response, just try to blacklist
	if err := j.BlacklistToken(tokenString); err != nil {
		// Could log this error for monitoring
		// log.Printf("Failed to auto-blacklist expired token: %v", err)
	}
}

// LogoutUser blacklists a user's token (for logout functionality)
func (j *JWTService) LogoutUser(tokenString string) error {
	return j.BlacklistToken(tokenString)
}

// IsTokenExpired checks if token is naturally expired (separate from blacklist)
func (j *JWTService) IsTokenExpired(tokenString string) (bool, error) {
	claims, err := j.GetJWTPayload(tokenString)
	if err != nil {
		return false, err
	}

	if claims.ExpiresAt == nil {
		return false, nil // Token has no expiration
	}

	return time.Now().After(claims.ExpiresAt.Time), nil
}

// IsTokenBlacklisted checks only blacklist status
func (j *JWTService) IsTokenBlacklisted(tokenString string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	key := fmt.Sprintf(BL, tokenString)
	exists := j.cache.Exists(ctx, key)
	if exists.Err() != nil {
		return false, fmt.Errorf("redis error: %w", exists.Err())
	}

	return exists.Val() > 0, nil
}
