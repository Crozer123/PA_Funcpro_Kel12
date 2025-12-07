package repository

import (
	"context"
	"errors"
	"github.com/Dzox13524/PA_Funcpro_Kel12/internal/domain"
	"gorm.io/gorm"
)

type CreateQuestionRepoFunc func(ctx context.Context, question domain.Question) (domain.Question, error)
type GetAllQuestionsRepoFunc func(ctx context.Context, currentUserID string) ([]domain.Question, error)
type GetQuestionByIDRepoFunc func(ctx context.Context, id, currentUserID string) (domain.Question, error)
type CreateAnswerRepoFunc func(ctx context.Context, answer domain.Answer) (domain.Answer, error)

type ToggleQuestionLikeFunc func(ctx context.Context, userID, questionID string) (bool, int64, error)
type ToggleFavoriteFunc func(ctx context.Context, userID, questionID string) (bool, error)

func NewCreateQuestionRepository(db *gorm.DB) CreateQuestionRepoFunc {
	return func(ctx context.Context, question domain.Question) (domain.Question, error) {
		result := db.WithContext(ctx).Create(&question)
		return question, result.Error
	}
}

func NewGetAllQuestionsRepository(db *gorm.DB) GetAllQuestionsRepoFunc {
	return func(ctx context.Context, currentUserID string) ([]domain.Question, error) {
		var questions []domain.Question
		if err := db.WithContext(ctx).Preload("User").Order("created_at desc").Find(&questions).Error; err != nil {
			return nil, err
		}

		for i := range questions {
			var count int64
			db.Model(&domain.QuestionLike{}).Where("question_id = ?", questions[i].ID).Count(&count)
			questions[i].LikesCount = count

			if currentUserID != "" {
				var like domain.QuestionLike
				if err := db.Where("user_id = ? AND question_id = ?", currentUserID, questions[i].ID).First(&like).Error; err == nil {
					questions[i].IsLiked = true
				}
				var fav domain.Favorite
				if err := db.Where("user_id = ? AND question_id = ?", currentUserID, questions[i].ID).First(&fav).Error; err == nil {
					questions[i].IsFavorited = true
				}
			}
		}
		return questions, nil
	}
}

func NewGetQuestionByIDRepository(db *gorm.DB) GetQuestionByIDRepoFunc {
	return func(ctx context.Context, id, currentUserID string) (domain.Question, error) {
		var q domain.Question
		err := db.WithContext(ctx).
			Preload("User").
			Preload("Answers", func(db *gorm.DB) *gorm.DB {
				return db.Order("created_at asc").Preload("User")
			}).
			First(&q, "id = ?", id).Error
		
		if err != nil {
			return domain.Question{}, err
		}

		var qCount int64
		db.Model(&domain.QuestionLike{}).Where("question_id = ?", q.ID).Count(&qCount)
		q.LikesCount = qCount

		if currentUserID != "" {
			var like domain.QuestionLike
			if err := db.Where("user_id = ? AND question_id = ?", currentUserID, q.ID).First(&like).Error; err == nil {
				q.IsLiked = true
			}
			var fav domain.Favorite
			if err := db.Where("user_id = ? AND question_id = ?", currentUserID, q.ID).First(&fav).Error; err == nil {
				q.IsFavorited = true
			}
		}

		return q, nil
	}
}

func NewCreateAnswerRepository(db *gorm.DB) CreateAnswerRepoFunc {
	return func(ctx context.Context, answer domain.Answer) (domain.Answer, error) {
		result := db.WithContext(ctx).Create(&answer)
		return answer, result.Error
	}
}

func NewToggleQuestionLikeRepository(db *gorm.DB) ToggleQuestionLikeFunc {
	return func(ctx context.Context, userID, questionID string) (bool, int64, error) {
		var like domain.QuestionLike
		var isLiked bool

		err := db.WithContext(ctx).Where("user_id = ? AND question_id = ?", userID, questionID).First(&like).Error
		
		if err == nil {
			db.WithContext(ctx).Delete(&like)
			isLiked = false
		} else if errors.Is(err, gorm.ErrRecordNotFound) {
			newLike := domain.QuestionLike{UserID: userID, QuestionID: questionID}
			db.WithContext(ctx).Create(&newLike)
			isLiked = true
		} else {
			return false, 0, err
		}

		var count int64
		db.WithContext(ctx).Model(&domain.QuestionLike{}).Where("question_id = ?", questionID).Count(&count)
		
		return isLiked, count, nil
	}
}

func NewToggleFavoriteRepository(db *gorm.DB) ToggleFavoriteFunc {
	return func(ctx context.Context, userID, questionID string) (bool, error) {
		var fav domain.Favorite
		var isFavorited bool

		err := db.WithContext(ctx).Where("user_id = ? AND question_id = ?", userID, questionID).First(&fav).Error
		
		if err == nil {
			db.WithContext(ctx).Delete(&fav)
			isFavorited = false
		} else {
			newFav := domain.Favorite{UserID: userID, QuestionID: questionID}
			db.WithContext(ctx).Create(&newFav)
			isFavorited = true
		}

		return isFavorited, nil
	}
}