package service

import (
	"auth-service/model"
	"auth-service/pkg/repository"
	"context"
	"crypto/sha1"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"math/rand"
	"time"
)

const (
	salt            = "vieurwnbv984hbiwu"
	signingKey      = "qrkjk#4#%35FSFJsdlja#4353qrcjk#4#%35"
	tokenTTL        = 10 * time.Hour
	refreshTokenTTL = 600 * time.Hour
)

type AuthService struct {
	repo repository.Authorization
}

type tokenClaims struct {
	jwt.StandardClaims
	UserId primitive.ObjectID `json:"id"`
}

type Tokens struct {
	AccessToken  string
	RefreshToken string
}

func NewAuthService(repo repository.Authorization) *AuthService {
	return &AuthService{repo: repo}
}

func (s *AuthService) CreateUser(ctx context.Context, user model.User) error {
	user.Password = generatePasswordHash(user.Password)
	return s.repo.CreateUser(ctx, user)
}

func (s *AuthService) SignIn(ctx context.Context, email, password string) (Tokens, error) {
	user, err := s.repo.GetUser(ctx, email, generatePasswordHash(password))
	if err != nil {
		return Tokens{}, err
	}

	return s.createSession(ctx, user.ID)
}

func (s *AuthService) RefreshTokens(ctx context.Context, refreshToken string) (Tokens, error) {
	user, err := s.repo.GetByRefreshToken(ctx, refreshToken)
	if err != nil {
		return Tokens{}, err
	}

	return s.createSession(ctx, user.ID)
}

func generatePasswordHash(password string) string {
	hash := sha1.New()
	hash.Write([]byte(password))

	return fmt.Sprintf("%x", hash.Sum([]byte(salt)))
}

func (s *AuthService) createSession(ctx context.Context, userId primitive.ObjectID) (Tokens, error) {
	var (
		res Tokens
		err error
	)

	res.AccessToken, err = newJWT(userId)
	if err != nil {
		return res, err
	}
	res.RefreshToken, err = newRefreshToken()

	res.RefreshToken += res.AccessToken[len(res.AccessToken)-8 : len(res.AccessToken)]

	if err != nil {
		return res, err
	}
	session := model.Session{
		RefreshToken: res.RefreshToken,
		ExpiresAt:    time.Now().Add(refreshTokenTTL),
	}

	err = s.repo.SetSession(ctx, userId, session)

	return res, err
}

func newJWT(userId primitive.ObjectID) (string, error) {
	token := jwt.NewWithClaims(
		jwt.SigningMethodHS512, &tokenClaims{
			jwt.StandardClaims{
				ExpiresAt: time.Now().Add(tokenTTL).Unix(),
				IssuedAt:  time.Now().Unix(),
			},
			userId,
		},
	)

	return token.SignedString([]byte(signingKey))
}

func newRefreshToken() (string, error) {
	b := make([]byte, 32)

	s := rand.NewSource(time.Now().Unix())
	r := rand.New(s)

	if _, err := r.Read(b); err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", b), nil
}
