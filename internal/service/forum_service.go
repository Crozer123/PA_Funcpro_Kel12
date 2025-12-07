package service

import (
	"context"
	"github.com/Dzox13524/PA_Funcpro_Kel12/internal/domain"
	"github.com/Dzox13524/PA_Funcpro_Kel12/internal/repository"
	"github.com/google/uuid"
)

type CreateQuestionFunc func(ctx context.Context, userID, title, content, category string) (domain.Question, error)
type GetFeedFunc func(ctx context.Context, currentUserID string) ([]domain.Question, error)
type GetQuestionDetailFunc func(ctx context.Context, id, currentUserID string) (domain.Question, error)
type AddAnswerFunc func(ctx context.Context, userID, questionID, content string) (domain.Answer, error)
type ToggleLikeFunc func(ctx context.Context, userID, questionID string) (bool, int64, error)
type ToggleFavFunc func(ctx context.Context, userID, questionID string) (bool, error)

func NewCreateQuestion(repo repository.CreateQuestionRepoFunc) CreateQuestionFunc {
	return func(ctx context.Context, userID, title, content, category string) (domain.Question, error) {
		newQ := domain.Question{
			ID:       uuid.New().String(),
			UserID:   userID,
			Title:    title,
			Content:  content,
			Category: category,
		}
		return repo(ctx, newQ)
	}
}

func NewGetFeed(repo repository.GetAllQuestionsRepoFunc) GetFeedFunc {
	return func(ctx context.Context, currentUserID string) ([]domain.Question, error) {
		return repo(ctx, currentUserID)
	}
}

func NewGetQuestionDetail(repo repository.GetQuestionByIDRepoFunc) GetQuestionDetailFunc {
	return func(ctx context.Context, id, currentUserID string) (domain.Question, error) {
		return repo(ctx, id, currentUserID)
	}
}

func NewAddAnswer(repo repository.CreateAnswerRepoFunc) AddAnswerFunc {
	return func(ctx context.Context, userID, questionID, content string) (domain.Answer, error) {
		newAns := domain.Answer{
			ID:         uuid.New().String(),
			UserID:     userID,
			QuestionID: questionID,
			Content:    content,
		}
		return repo(ctx, newAns)
	}
}

func NewToggleLike(repo repository.ToggleQuestionLikeFunc) ToggleLikeFunc {
	return func(ctx context.Context, userID, questionID string) (bool, int64, error) {
		return repo(ctx, userID, questionID)
	}
}

func NewToggleFav(repo repository.ToggleFavoriteFunc) ToggleFavFunc {
	return func(ctx context.Context, userID, questionID string) (bool, error) {
		return repo(ctx, userID, questionID)
	}
}