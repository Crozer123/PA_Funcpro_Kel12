package handle

import (
	"encoding/json"
	"net/http"

	"github.com/Dzox13524/PA_Funcpro_Kel12/internal/domain"
	"github.com/Dzox13524/PA_Funcpro_Kel12/internal/middleware"
	"github.com/Dzox13524/PA_Funcpro_Kel12/internal/response"
	"github.com/Dzox13524/PA_Funcpro_Kel12/internal/service"
)

type CreateUserRequest struct {
	Name     string
	Email    string
	Password string
}

func HandleCreateUser(createUser service.CreateUserFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req CreateUserRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			response.WriteJSON(w, http.StatusBadRequest, "Format JSON tidak valid!", nil)
			return
		}

		newUser, err := createUser(r.Context(), req.Name, req.Email, req.Password)
		if err != nil {
			response.WriteJSON(w, http.StatusConflict, err.Error(), nil)
			return
		}

		response.WriteJSON(w, http.StatusCreated, "success_created", newUser)
	}
}

func HandleLogin(loginService service.LoginFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req domain.LoginRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			response.WriteJSON(w, http.StatusBadRequest, "Invalid Request", nil)
			return
		}

		authResp, err := loginService(r.Context(), req.Email, req.Password)
		if err != nil {
			response.WriteJSON(w, http.StatusUnauthorized, err.Error(), nil)
			return
		}

		response.WriteJSON(w, http.StatusOK, "login_success", authResp)
	}
}

func HandleGetMe(getUserByID service.GetUserByIDFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := middleware.GetUserIDFromContext(r.Context())
		if userID == "" {
			response.WriteJSON(w, http.StatusUnauthorized, "Unauthorized", nil)
			return
		}

		user, err := getUserByID(r.Context(), userID)
		if err != nil {
			response.WriteJSON(w, http.StatusNotFound, "User not found", nil)
			return
		}

		response.WriteJSON(w, http.StatusOK, "success", user)
	}
}

func HandleUpdateMe(updateUser service.UpdateUserFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := middleware.GetUserIDFromContext(r.Context())
		if userID == "" {
			response.WriteJSON(w, http.StatusUnauthorized, "Unauthorized", nil)
			return
		}

		var req domain.UpdateUserRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			response.WriteJSON(w, http.StatusBadRequest, "Invalid JSON", nil)
			return
		}

		updatedUser, err := updateUser(r.Context(), userID, req.Name)
		if err != nil {
			response.WriteJSON(w, http.StatusInternalServerError, err.Error(), nil)
			return
		}

		response.WriteJSON(w, http.StatusOK, "profile_updated", updatedUser)
	}
}

func HandleGetUserByID(getUserByID service.GetUserByIDFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		user, err := getUserByID(r.Context(), id)
		if err != nil {
			response.WriteJSON(w, http.StatusNotFound, "User tidak ditemukan", nil)
			return
		}
		response.WriteJSON(w, http.StatusOK, "success", user)
	}
}