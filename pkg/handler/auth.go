package handler

import (
	"auth-service/model"
	"encoding/json"
	"net/http"
)

func (h *Handler) signUp(w http.ResponseWriter, r *http.Request) {
	var input model.User
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		newErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}
	err := h.services.Authorization.CreateUser(r.Context(), input)

	if err != nil {
		newErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(
		map[string]interface{}{
			"id": "all is ok",
		},
	)
}

type signInInput struct {
	Email    string `json:"email" binding:"required" bson:"email"`
	Password string `json:"password" binding:"required" bson:"password"`
}

func (h *Handler) signIn(w http.ResponseWriter, r *http.Request) {
	var input signInInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		newErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	result, err := h.services.Authorization.SignIn(r.Context(), input.Email, input.Password)

	if err != nil {
		newErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(
		map[string]interface{}{
			"accessToken":  result.AccessToken,
			"refreshToken": result.RefreshToken,
		},
	)
}

type refreshInput struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

func (h *Handler) refresh(w http.ResponseWriter, r *http.Request) {
	var input refreshInput

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		newErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	if !h.refreshTokenIdentity(input) {
		newErrorResponse(w, http.StatusBadRequest, "refresh token does not match access token")
		return
	}

	result, err := h.services.RefreshTokens(r.Context(), input.RefreshToken)

	if err != nil {
		newErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	json.NewEncoder(w).Encode(
		map[string]interface{}{
			"accessToken":  result.AccessToken,
			"refreshToken": result.RefreshToken,
		},
	)
}
