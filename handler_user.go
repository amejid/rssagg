package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"

	"github.com/amejid/rssagg/internal/database"
)

func (apiCfg *apiConfig) handlerCreateUser(
	w http.ResponseWriter,
	r *http.Request,
) {
	type parameters struct {
		Name string `json:"name"`
	}

	decoder := json.NewDecoder(r.Body)

	var params parameters
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(
			w,
			http.StatusBadRequest,
			fmt.Sprintf("Error parsing JSON: %v", err),
		)
		return
	}

	user, err := apiCfg.DB.CreateUser(r.Context(), database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      params.Name,
	})
	if err != nil {
		respondWithError(
			w,
			http.StatusBadRequest,
			fmt.Sprintf("Could not create user: %v", err),
		)
		return
	}

	respondWithJSON(w, http.StatusCreated, databaseUserToUser(user))
}

func (apiCfg *apiConfig) handlerGetUser(
	w http.ResponseWriter,
	r *http.Request,
	user database.User,
) {
	respondWithJSON(w, http.StatusOK, databaseUserToUser(user))
}

func (apiCfg *apiConfig) handlerGetPostsForUser(
	w http.ResponseWriter,
	r *http.Request,
	user database.User,
) {
	posts, err := apiCfg.DB.GetPostsForUser(
		r.Context(),
		database.GetPostsForUserParams{
			UserID: user.ID,
			Limit:  10,
		},
	)
	if err != nil {
		respondWithError(
			w,
			http.StatusBadRequest,
			fmt.Sprintf("could not get posts: %v", err),
		)
		return
	}

	respondWithJSON(w, http.StatusOK, databasePostsToPosts(posts))
}
