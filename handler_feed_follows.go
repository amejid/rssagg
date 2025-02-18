package main

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"net/http"
	"time"

	"github.com/google/uuid"

	"github.com/amejid/rssagg/internal/database"
)

func (apiCfg *apiConfig) handlerCreateFeedFollow(
	w http.ResponseWriter,
	r *http.Request,
	user database.User,
) {
	type parameters struct {
		FeedID uuid.UUID `json:"feed_id"`
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

	feedFollow, err := apiCfg.DB.CreateFeedFollow(
		r.Context(),
		database.CreateFeedFollowParams{
			ID:        uuid.New(),
			CreatedAt: time.Now().UTC(),
			UpdatedAt: time.Now().UTC(),
			UserID:    user.ID,
			FeedID:    params.FeedID,
		},
	)
	if err != nil {
		respondWithError(
			w,
			http.StatusBadRequest,
			fmt.Sprintf("Could not create feed follow: %v", err),
		)
		return
	}

	respondWithJSON(
		w,
		http.StatusCreated,
		databaseFeedFollowToFeedFollow(feedFollow),
	)
}

func (apiCfg *apiConfig) handlerGetFeedFollows(
	w http.ResponseWriter,
	r *http.Request,
	user database.User,
) {
	feedFollows, err := apiCfg.DB.GetFeedFollows(
		r.Context(), user.ID,
	)
	if err != nil {
		respondWithError(
			w,
			http.StatusBadRequest,
			fmt.Sprintf("Could not get feed follows: %v", err),
		)
		return
	}

	respondWithJSON(
		w,
		http.StatusOK,
		databaseFeedFollowsToFeedFollows(feedFollows),
	)
}

func (apiCfg *apiConfig) handlerDeleteFeedFollow(
	w http.ResponseWriter,
	r *http.Request,
	user database.User,
) {
	feedFollowIDStr := chi.URLParam(r, "feedFollowID")
	feedFollowID, err := uuid.Parse(feedFollowIDStr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Could not parse feed follow id: %v", err))
		return
	}

	err = apiCfg.DB.DeleteFeedFollow(r.Context(), database.DeleteFeedFollowParams{
		ID:     feedFollowID,
		UserID: user.ID,
	})
	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Could not delete feed follow: %v", err))
		return
	}

	respondWithJSON(w, http.StatusNoContent, struct {
	}{})
}
