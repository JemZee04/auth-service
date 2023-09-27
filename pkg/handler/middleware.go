package handler

func (h *Handler) refreshTokenIdentity(inputTokens refreshInput) bool {
	return inputTokens.RefreshToken[len(inputTokens.RefreshToken)-8:len(inputTokens.RefreshToken)] == inputTokens.AccessToken[len(inputTokens.AccessToken)-8:len(inputTokens.AccessToken)]
}

//import (
//	"net/http"
//	"strings"
//)
//
//const (
//	//authorizationHeader = "Authorization"
//	userCtx             = "userId"
//)
//
//func (h *Handler) userIdentity(w http.ResponseWriter, r *http.Request) {
//	header := r.Header
//	if header == "" {
//		newErrorResponse(w, http.StatusUnauthorized, "empty auth header")
//		return
//	}
//
//	headerParts := strings.Split(header, " ")
//	if len(headerParts) != 2 {
//		newErrorResponse(w, http.StatusUnauthorized, "invalid auth header")
//	}
//
//	userId, err := h.services.Authorization.ParseToken(headerParts[1])
//	if err != nil {
//		newErrorResponse(w, http.StatusUnauthorized, err.Error())
//		return
//	}
//
//	w.Set(userCtx, userId)
//	r.
//}
