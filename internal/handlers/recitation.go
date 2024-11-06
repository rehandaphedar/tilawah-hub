package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"

	"git.sr.ht/~rehandaphedar/tilawah-hub/internal/db"
	"git.sr.ht/~rehandaphedar/tilawah-hub/internal/sqlc"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
)

type createRecitationDTO struct {
	Slug string `json:"slug"`
}

type updateRecitationDTO struct {
	Slug string `json:"slug"`
}

// CreateRecitation godoc
//
//	@Tags		Recitation
//	@Accept		json
//	@Produce	json
//
//	@Param		X-CSRF-TOKEN	header		string				true	"CSRF Token"
//
//	@Param		request			body		createRecitationDTO	true	"Create Recitation"
//
//	@Success	200				{object}	sqlc.Recitation
//	@Failure	400				{object}	models.Error
//	@Failure	401				{object}	models.Error
//	@Failure	500				{object}	models.Error
//	@Router		/recitations [post]
func CreateRecitation(w http.ResponseWriter, r *http.Request) {
	reciter := r.Context().Value("username").(string)

	var request sqlc.RecitationCreateRecitationParams
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, render.M{
			"message": "Could not parse request body",
			"error":   err.Error(),
		})
		return
	}

	request.Reciter = reciter
	request.Name = request.Slug

	recitation, err := db.Queries.RecitationCreateRecitation(context.Background(), request)
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, render.M{
			"message": "Error creating recitation",
			"error":   err.Error(),
		})
		return
	}

	recitationDir := filepath.Join("data", "static", recitation.Reciter, recitation.Slug)
	err = os.MkdirAll(recitationDir, 0755)
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, render.M{
			"message": "Error creating recitation directory",
			"error":   err.Error(),
		})
		return
	}

	render.JSON(w, r, recitation)
}

// GetRecitations godoc
//
//	@Tags		Recitation
//	@Produce	json
//
//	@Success	200	{object}	[]sqlc.Recitation
//	@Failure	500	{object}	models.Error
//	@Router		/recitations [get]
func GetRecitations(w http.ResponseWriter, r *http.Request) {
	recitations, err := db.Queries.RecitationSelectRecitations(context.Background())
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, render.M{
			"message": "Error querying recitations",
			"error":   err.Error(),
		})
		return
	}

	render.JSON(w, r, recitations)
}

// GetRecitation godoc
//
//	@Tags		Recitation
//	@Produce	json
//
//	@Param		reciter	path		string	true	"Reciter"
//	@Param		slug	path		string	true	"Slug"
//
//	@Success	200		{object}	sqlc.Recitation
//	@Failure	500		{object}	models.Error
//	@Router		/recitations/{reciter}/{slug} [get]
func GetRecitation(w http.ResponseWriter, r *http.Request) {
	request := sqlc.RecitationSelectRecitationParams{
		Reciter: chi.URLParam(r, "reciter"),
		Slug:    chi.URLParam(r, "slug"),
	}
	recitation, err := db.Queries.RecitationSelectRecitation(context.Background(), request)
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, render.M{
			"message": "Error querying recitation",
			"error":   err.Error(),
		})
		return
	}

	render.JSON(w, r, recitation)
}

// UpdateRecitation godoc
//
//	@Tags		Recitation
//	@Accept		json
//	@Produce	json
//
//	@Param		X-CSRF-TOKEN	header		string				true	"CSRF Token"
//
//	@Param		slug			path		string				true	"Slug"
//	@Param		request			body		updateRecitationDTO	true	"Update Recitation"
//
//	@Success	200				{object}	sqlc.Recitation
//	@Failure	400				{object}	models.Error
//	@Failure	401				{object}	models.Error
//	@Failure	500				{object}	models.Error
//	@Router		/recitations/{slug} [put]
func UpdateRecitation(w http.ResponseWriter, r *http.Request) {
	reciter := r.Context().Value("username").(string)
	slug := chi.URLParam(r, "slug")

	selectRequest := sqlc.RecitationSelectRecitationParams{
		Reciter: reciter,
		Slug:    slug,
	}

	existingRecitation, err := db.Queries.RecitationSelectRecitation(context.Background(), selectRequest)
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, render.M{
			"message": "Error querying recitation",
			"error":   err.Error(),
		})
		return
	}

	updatedRecitationData := &sqlc.RecitationUpdateRecitationParams{
		Reciter: reciter,
		Slug:    slug,
		Name:    existingRecitation.Name,
	}

	var updateRequest sqlc.RecitationUpdateRecitationParams
	if err := json.NewDecoder(r.Body).Decode(&updateRequest); err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, render.M{
			"message": "Could not parse request body",
			"error":   err.Error(),
		})
		return
	}

	if updateRequest.Name != "" {
		updatedRecitationData.Name = updateRequest.Name
	}

	updatedRecitation, err := db.Queries.RecitationUpdateRecitation(context.Background(), *updatedRecitationData)
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, render.M{
			"message": "Unexpected error",
			"error":   err.Error(),
		})
		return
	}

	render.JSON(w, r, updatedRecitation)
}

// DeleteRecitation godoc
//
//	@Tags		Recitation
//	@Produce	json
//
//	@Param		X-CSRF-TOKEN	header		string	true	"CSRF Token"
//
//	@Param		slug			path		string	true	"Slug"
//
//	@Success	200				{object}	sqlc.Recitation
//	@Failure	400				{object}	models.Error
//	@Failure	401				{object}	models.Error
//	@Failure	500				{object}	models.Error
//	@Router		/recitations/{slug} [delete]
func DeleteRecitation(w http.ResponseWriter, r *http.Request) {
	reciter := r.Context().Value("username").(string)
	slug := chi.URLParam(r, "slug")

	request := sqlc.RecitationDeleteRecitationParams{
		Reciter: reciter,
		Slug:    slug,
	}

	deletedRecitation, err := db.Queries.RecitationDeleteRecitation(context.Background(), request)
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, render.M{
			"message": "Unexpected error",
			"error":   err.Error(),
		})
		return
	}

	recitationDir := filepath.Join("data", "uploads", reciter, slug)
	err = os.RemoveAll(recitationDir)
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, render.M{
			"message": "Error deleting recitation directory",
			"error":   err.Error(),
		})
		return
	}

	render.JSON(w, r, deletedRecitation)
}
