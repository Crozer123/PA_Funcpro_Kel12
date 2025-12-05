package service

import (
	"context"
	"errors"
	"time"

	"github.com/Dzox13524/PA_Funcpro_Kel12/internal/domain"
	"github.com/Dzox13524/PA_Funcpro_Kel12/internal/repository"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

var jwtSecret = []byte("rahasia-negara-api")

type GetUserByIDFunc func(ctx context.Context, id string) (domain.User, error)
type CreateUserFunc func(ctx context.Context, name, email, password string) (domain.User, error)
type LoginFunc func(ctx context.Context, email, password string) (domain.AuthResponse, error)
type UpdateUserFunc func(ctx context.Context, id string, name string) (domain.User, error)

func NewGetUserByID(getUserRepo repository.GetUserByIDRepoFunc) GetUserByIDFunc {
	return func(ctx context.Context, id string) (domain.User, error) {
		return getUserRepo(ctx, id)
	}
}

func NewCreateUser(createRepo repository.CreateUserRepoFunc, getByEmailRepo repository.GetUserByEmailRepoFunc) CreateUserFunc {
	return func(ctx context.Context, name, email, password string) (domain.User, error) {
		_, err := getByEmailRepo(ctx, email)
		if err == nil {
			return domain.User{}, errors.New("email sudah terdaftar")
		}

		if len(password) < 8 {
			return domain.User{}, errors.New("password minimal 8 karakter")
		}

		hashPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			return domain.User{}, err
		}

		newUser := domain.User{
			ID:       uuid.New().String(),
			Name:     name,
			Email:    email,
			Password: string(hashPassword),
			Role:     "User",
		}

		return createRepo(ctx, newUser)
	}
}

func NewLoginService(getByEmailRepo repository.GetUserByEmailRepoFunc) LoginFunc {
	return func(ctx context.Context, email, password string) (domain.AuthResponse, error) {
		user, err := getByEmailRepo(ctx, email)
		if err != nil {
			return domain.AuthResponse{}, errors.New("email atau password salah")
		}

		err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
		if err != nil {
			return domain.AuthResponse{}, errors.New("email atau password salah")
		}

		claims := jwt.MapClaims{
			"user_id": user.ID,
			"role":    user.Role,
			"exp":     time.Now().Add(time.Hour * 72).Unix(), 
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		t, err := token.SignedString(jwtSecret)
		if err != nil {
			return domain.AuthResponse{}, err
		}

		return domain.AuthResponse{
			Token: t,
			User:  user,
		}, nil
	}
}

func NewUpdateUser(updateRepo repository.UpdateUserRepoFunc) UpdateUserFunc {
	return func(ctx context.Context, id string, name string) (domain.User, error) {
		updates := map[string]interface{}{}
		if name != "" {
			updates["name"] = name
		}
		updates["updated_at"] = time.Now()

		return updateRepo(ctx, id, updates)
	}
}