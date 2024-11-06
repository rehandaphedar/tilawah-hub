package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"git.sr.ht/~rehandaphedar/tilawah-hub/internal/db"
	"git.sr.ht/~rehandaphedar/tilawah-hub/internal/models"
	"git.sr.ht/~rehandaphedar/tilawah-hub/internal/sqlc"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
)

// UpdateRecitationTiming godoc
//
//	@Tags		RecitationTiming
//	@Accept		json
//	@Produce	json
//
//	@Param		X-CSRF-TOKEN	header		string			true	"CSRF Token"
//
//	@Param		slug			path		string			true	"Slug"
//	@Param		verse_key		path		string			true	"Verse Key"
//	@Param		request			body		models.Timing	true	"Update Recitation Timing"
//
//	@Success	200				{object}	models.Timing
//	@Failure	400				{object}	models.Error
//	@Failure	401				{object}	models.Error
//	@Router		/recitation-timings/{slug}/{verse_key} [post]
func UpdateRecitationTiming(w http.ResponseWriter, r *http.Request) {
	reciter := r.Context().Value("username").(string)
	slug := chi.URLParam(r, "slug")
	verseKey := chi.URLParam(r, "verse_key")

	existingRecitationFile, err := db.Queries.RecitationFileSelectRecitationFile(context.Background(), sqlc.RecitationFileSelectRecitationFileParams{
		Reciter:  reciter,
		Slug:     slug,
		VerseKey: verseKey,
	})
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, render.M{
			"message": "Error checking status of recitation file",
			"error":   err.Error(),
		})
		return
	}
	if existingRecitationFile.LafzizeProcessing {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, render.M{
			"message": "The recitation is currently being lafzized",
			"error":   "",
		})
		return
	}

	var timing models.Timing
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&timing)
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, render.M{
			"message": "Error parsing request JSON",
			"error":   err.Error(),
		})
		log.Println(err)
		return
	}
	defer r.Body.Close()

	baseDir := filepath.Join("data", "uploads", reciter, slug)
	err = os.MkdirAll(baseDir, 0755)
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, render.M{
			"message": "Error creating recitation directory",
			"error":   err.Error(),
		})
		return
	}

	timingsFilepath := filepath.Join(baseDir, verseKey+".json")
	err = os.RemoveAll(timingsFilepath)
	if err != nil {
		log.Printf("Error removing existing timings file: %v\n", err)
	}

	fileHandler, err := os.Create(timingsFilepath)
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, render.M{
			"message": "Error creating recitation file",
			"error":   err.Error(),
		})
		return
	}
	defer fileHandler.Close()

	jsonData, err := json.Marshal(timing)
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, render.M{
			"message": "Error saving JSON file",
			"error":   err.Error(),
		})
		return
	}

	err = os.WriteFile(timingsFilepath, jsonData, 0644)
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, render.M{
			"message": "Error saving JSON file",
			"error":   err.Error(),
		})
		return
	}

	_, err = db.Queries.RecitationFileUpdateRecitationFile(context.Background(), sqlc.RecitationFileUpdateRecitationFileParams{
		Reciter:           reciter,
		Slug:              slug,
		VerseKey:          verseKey,
		HasTimings:        true,
		LafzizeProcessing: false,
	})
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, render.M{
			"message": "Error updating status of recitation file",
			"error":   err.Error(),
		})
		return
	}

	render.JSON(w, r, timing)
}

// DeleteRecitationTiming godoc
//
//	@Tags		RecitationTiming
//	@Produce	json
//
//	@Param		X-CSRF-TOKEN	header		string	true	"CSRF Token"
//
//	@Param		slug			path		string	true	"Slug"
//	@Param		verse_key		path		string	true	"Verse Key"
//
//	@Success	200				{object}	models.Timing
//	@Failure	400				{object}	models.Error
//	@Failure	401				{object}	models.Error
//	@Router		/recitation-timings/{slug}/{verse_key} [delete]
func DeleteRecitationTiming(w http.ResponseWriter, r *http.Request) {
	reciter := r.Context().Value("username").(string)
	slug := chi.URLParam(r, "slug")
	verseKey := chi.URLParam(r, "verse_key")

	existingRecitationFile, err := db.Queries.RecitationFileSelectRecitationFile(context.Background(), sqlc.RecitationFileSelectRecitationFileParams{
		Reciter:  reciter,
		Slug:     slug,
		VerseKey: verseKey,
	})
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, render.M{
			"message": "Error checking status of recitation file",
			"error":   err.Error(),
		})
		return
	}
	if existingRecitationFile.LafzizeProcessing {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, render.M{
			"message": "The recitation is currently being lafzized",
			"error":   "",
		})
		return
	}

	baseDir := filepath.Join("data", "uploads", reciter, slug)
	timingsFilepath := filepath.Join(baseDir, fmt.Sprintf("%s.json", verseKey))

	file, err := os.Open(timingsFilepath)
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, render.M{
			"message": "Error opening recitation timing",
			"error":   err.Error(),
		})
		return
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	var timing models.Timing
	err = decoder.Decode(&timing)
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, render.M{
			"message": "Error reading recitation timing",
			"error":   err.Error(),
		})
		return
	}

	err = os.RemoveAll(timingsFilepath)
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, render.M{
			"message": "Error deleting recitation timing",
			"error":   err.Error(),
		})
		return
	}

	_, err = db.Queries.RecitationFileUpdateRecitationFile(context.Background(), sqlc.RecitationFileUpdateRecitationFileParams{
		Reciter:           reciter,
		Slug:              slug,
		VerseKey:          verseKey,
		HasTimings:        false,
		LafzizeProcessing: false,
	})
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, render.M{
			"message": "Error updating status of recitation file.",
			"error":   err.Error(),
		})
		return
	}

	render.JSON(w, r, timing)
}
