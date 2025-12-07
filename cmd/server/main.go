package main

import (
	"log"
	"net/http"

	handle "github.com/Dzox13524/PA_Funcpro_Kel12/internal/handler"
	"github.com/Dzox13524/PA_Funcpro_Kel12/internal/middleware"
	"github.com/Dzox13524/PA_Funcpro_Kel12/internal/platform/database"
	"github.com/Dzox13524/PA_Funcpro_Kel12/internal/repository"
	"github.com/Dzox13524/PA_Funcpro_Kel12/internal/service"
	"github.com/Dzox13524/PA_Funcpro_Kel12/internal/domain"
)

func main() {
	db := database.NewConnection()

	db.AutoMigrate(
		&domain.User{}, 
		&domain.Question{}, 
		&domain.Answer{},
		&domain.QuestionLike{}, 
		&domain.Favorite{},
	)

	userRepoGetByID := repository.NewGetUserByIDRepository(db)
	userRepoGetByEmail := repository.NewGetUserByEmailRepository(db)
	userRepoCreate := repository.NewCreateUserRepository(db)
	userRepoUpdate := repository.NewUpdateUserRepository(db)

	createUserService := service.NewCreateUser(userRepoCreate, userRepoGetByEmail)
	getUserByIDService := service.NewGetUserByID(userRepoGetByID)
	loginService := service.NewLoginService(userRepoGetByEmail)
	updateUserService := service.NewUpdateUser(userRepoUpdate)

	qRepoCreate := repository.NewCreateQuestionRepository(db)
	qRepoGet := repository.NewGetAllQuestionsRepository(db)
	qRepoDetail := repository.NewGetQuestionByIDRepository(db)
	aRepoCreate := repository.NewCreateAnswerRepository(db)
	likeRepo := repository.NewToggleQuestionLikeRepository(db)
	favRepo := repository.NewToggleFavoriteRepository(db)

	svcCreateQ := service.NewCreateQuestion(qRepoCreate)
	svcGetFeed := service.NewGetFeed(qRepoGet)
	svcGetDetail := service.NewGetQuestionDetail(qRepoDetail)
	svcAddAns := service.NewAddAnswer(aRepoCreate)
	svcLike := service.NewToggleLike(likeRepo)
	svcFav := service.NewToggleFav(favRepo)

	log.SetFlags(0)
	mux := http.NewServeMux()

	mux.HandleFunc("POST /api/v1/auth/register", handle.HandleCreateUser(createUserService))
	mux.HandleFunc("POST /api/v1/auth/login", handle.HandleLogin(loginService))
	
	mux.HandleFunc("GET /api/v1/users/{id}", handle.HandleGetUserByID(getUserByIDService))
	mux.HandleFunc("GET /api/v1/users/me", middleware.AuthMiddleware(handle.HandleGetMe(getUserByIDService)))
	mux.HandleFunc("PATCH /api/v1/users/me", middleware.AuthMiddleware(handle.HandleUpdateMe(updateUserService)))

	mux.HandleFunc("GET /api/v1/questions", middleware.AuthMiddlewareOptional(handle.HandleGetFeed(svcGetFeed)))
	mux.HandleFunc("POST /api/v1/questions", middleware.AuthMiddleware(handle.HandleCreateQuestion(svcCreateQ)))
	mux.HandleFunc("GET /api/v1/questions/{id}", middleware.AuthMiddlewareOptional(handle.HandleGetQuestionDetail(svcGetDetail)))
	mux.HandleFunc("POST /api/v1/questions/{id}/answers", middleware.AuthMiddleware(handle.HandleAddAnswer(svcAddAns)))
	mux.HandleFunc("POST /api/v1/questions/{id}/like", middleware.AuthMiddleware(handle.HandleToggleLike(svcLike)))
	mux.HandleFunc("POST /api/v1/questions/{id}/favorite", middleware.AuthMiddleware(handle.HandleToggleFav(svcFav)))

	var finalHandler http.Handler = mux
	finalHandler = middleware.Logging(finalHandler)
	finalHandler = middleware.CORSMiddleware(finalHandler)

	log.Println("Server running on port :8080")
	http.ListenAndServe(":8080", finalHandler)
}