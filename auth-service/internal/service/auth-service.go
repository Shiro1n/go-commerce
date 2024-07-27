package service

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/go-redis/redis/v8"
	"github.com/shiro1n/go-commerce/auth-service/internal/config"
	"github.com/shiro1n/go-commerce/auth-service/internal/model"
	"github.com/shiro1n/go-commerce/auth-service/internal/utils"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/net/context"
)

type AuthService interface {
	AuthenticateUser(username, password string) (*model.User, error)
	RegisterUser(username, email, password string) (*model.User, error)
	UpdateUser(user *model.User) error
	CreateTokens(userId int) (*model.TokenDetails, error)
	VerifyToken(tokenStr string) (*jwt.Token, error)
	RefreshTokens(refreshToken string) (*model.TokenDetails, error)
	InvalidateRefreshToken(refreshToken string) error
	ValidateToken(tokenStr string) (string, error)
	StoreUserInRedis(user *model.User) error
	GetUserFromRedis(userId int) (*model.User, error)
}

type authService struct {
	Config      config.Config
	RedisClient *redis.Client
}

func NewAuthService(cfg config.Config) AuthService {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "redis:6379", // Docker service name for Redis
		Password: "",           // no password set
		DB:       0,            // use default DB
	})

	return &authService{Config: cfg, RedisClient: rdb}
}

func (s *authService) AuthenticateUser(username, password string) (*model.User, error) {
	user, err := s.getUserByUsername(username)
	if err != nil {
		return nil, err
	}

	if !checkPasswordHash(password, user.Password) {
		return nil, errors.New("invalid username or password")
	}

	return user, nil
}

func (s *authService) RegisterUser(username, email, password string) (*model.User, error) {
	user, err := s.getUserByUsername(username)
	if user != nil {
		return nil, errors.New("username already taken")
	}

	hashedPassword, err := hashPassword(password)
	if err != nil {
		return nil, err
	}

	newUser := &model.User{
		Username: username,
		Email:    email,
		Password: hashedPassword,
	}

	err = s.createUser(newUser)
	if err != nil {
		return nil, err
	}

	return newUser, nil
}

func (s *authService) UpdateUser(user *model.User) error {
	// Update user in primary database
	err := s.updateUserInPrimaryDatabase(user)
	if err != nil {
		return err
	}

	// Update user in Redis
	err = s.StoreUserInRedis(user)
	if err != nil {
		return err
	}

	return nil
}

func (s *authService) updateUserInPrimaryDatabase(user *model.User) error {
	data, err := json.Marshal(user)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPut, fmt.Sprintf("%s/users/%d", s.Config.UserServiceURL, user.ID), bytes.NewBuffer(data))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return errors.New("failed to update user in primary database")
	}

	return nil
}

func (s *authService) CreateTokens(userId int) (*model.TokenDetails, error) {
	userIdStr := strconv.Itoa(userId)

	td, err := utils.CreateToken(userIdStr, s.Config.JWTSecret) // Updated to use utils package
	if err != nil {
		return nil, err
	}

	err = s.storeRefreshToken(td.RefreshToken, td.RtExpires)
	if err != nil {
		return nil, err
	}

	return td, nil
}

func (s *authService) VerifyToken(tokenStr string) (*jwt.Token, error) {
	return utils.VerifyToken(tokenStr, s.Config.JWTSecret) // Updated to use utils package
}

func (s *authService) RefreshTokens(refreshToken string) (*model.TokenDetails, error) {
	token, err := s.VerifyToken(refreshToken)
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	userIdStr, ok := claims["user_id"].(string)
	if !ok {
		return nil, errors.New("invalid token")
	}

	_, err = s.RedisClient.Get(context.Background(), refreshToken).Result()
	if err == redis.Nil {
		return nil, errors.New("invalid refresh token")
	} else if err != nil {
		return nil, err
	}

	userId, _ := strconv.Atoi(userIdStr)
	return s.CreateTokens(userId)
}

func (s *authService) InvalidateRefreshToken(refreshToken string) error {
	_, err := s.RedisClient.Del(context.Background(), refreshToken).Result()
	return err
}

func (s *authService) storeRefreshToken(refreshToken string, expires int64) error {
	return s.RedisClient.Set(context.Background(), refreshToken, true, time.Until(time.Unix(expires, 0))).Err()
}

func (s *authService) getUserByUsername(username string) (*model.User, error) {
	resp, err := http.Get(fmt.Sprintf("%s/users/%s", s.Config.UserServiceURL, username))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("user not found")
	}

	var user model.User
	err = json.NewDecoder(resp.Body).Decode(&user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (s *authService) createUser(user *model.User) error {
	data, err := json.Marshal(user)
	if err != nil {
		return err
	}

	resp, err := http.Post(fmt.Sprintf("%s/register", s.Config.UserServiceURL), "application/json", bytes.NewBuffer(data))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return errors.New("failed to create user")
	}

	return nil
}

func (s *authService) ValidateToken(tokenStr string) (string, error) {
	token, err := s.VerifyToken(tokenStr)
	if err != nil {
		return "", err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return "", errors.New("invalid token")
	}

	userId, ok := claims["user_id"].(string)
	if !ok {
		return "", errors.New("invalid token")
	}

	return userId, nil
}

func (s *authService) StoreUserInRedis(user *model.User) error {
	userData, err := json.Marshal(user)
	if err != nil {
		return err
	}

	ctx := context.Background()
	err = s.RedisClient.Set(ctx, s.getUserRedisKey(user.ID), userData, time.Hour*24).Err() // Store user data for 24 hours
	if err != nil {
		return err
	}
	return nil
}

func (s *authService) GetUserFromRedis(userId int) (*model.User, error) {
	ctx := context.Background()
	userData, err := s.RedisClient.Get(ctx, s.getUserRedisKey(userId)).Result()
	if err == redis.Nil {
		return nil, nil // User not found in Redis
	} else if err != nil {
		return nil, err
	}

	var user model.User
	err = json.Unmarshal([]byte(userData), &user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *authService) getUserRedisKey(userId int) string {
	return "user:" + strconv.Itoa(userId)
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
