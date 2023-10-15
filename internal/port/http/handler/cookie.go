package handler

import (
	"net/http"
	"time"
)

func (h *Handler) createCookie(token string) (*http.Cookie, error) {
	return &http.Cookie{
		Name:     h.GetAppName() + "_token",
		Value:    token,
		Expires:  time.Now().Add(h.GetTokenLifeTime() * time.Hour),
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	}, nil
}
