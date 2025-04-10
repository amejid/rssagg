package main

import (
	"fmt"
	"net/http"

	"github.com/amejid/rssagg/internal/auth"
	"github.com/amejid/rssagg/internal/database"
)

type authedHandler func(w http.ResponseWriter, r *http.Request, user database.User)

func (apiCfg *apiConfig) middlewareAuth(
	handler authedHandler,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		apiKey, err := auth.GetAPIKey(r.Header)
		if err != nil {
			respondWithError(
				w,
				http.StatusForbidden,
				fmt.Sprintf("Auth error: %v", err),
			)
			return
		}

		user, err := apiCfg.DB.GetUserByAPIKey(r.Context(), apiKey)
		if err != nil {
			respondWithError(
				w,
				http.StatusBadRequest,
				fmt.Sprintf("Could not get user: %v", err),
			)
			return
		}

		handler(w, r, user)
	}
}
