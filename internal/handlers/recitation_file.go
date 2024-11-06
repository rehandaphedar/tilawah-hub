package handlers

import (
	"context"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"

	"git.sr.ht/~rehandaphedar/tilawah-hub/internal/db"
	"git.sr.ht/~rehandaphedar/tilawah-hub/internal/sqlc"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
)

// CreateRecitationFile godoc
//
//	@Tags		RecitationFile
//	@Accept		multipart/form-data
//	@Produce	json
//
//	@Param		X-CSRF-TOKEN	header		string	true	"CSRF Token"
//
//	@Param		slug			path		string	true	"Slug"
//	@Param		file			formData	file	true	"File"
//	@Param		verse_key		formData	string	true	"Verse Key"
//
//	@Success	200				{object}	sqlc.RecitationFile
//	@Failure	400				{object}	models.Error
//	@Failure	401				{object}	models.Error
//	@Failure	500				{object}	models.Error
//	@Router		/recitation-files/{slug}/ [post]
func CreateRecitationFile(w http.ResponseWriter, r *http.Request) {
	reciter := r.Context().Value("username").(string)
	slug := chi.URLParam(r, "slug")

	err := r.ParseMultipartForm(4 << 20) // 4MB max file size
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, render.M{
			"message": "Error parsing multipart form",
			"error":   err.Error(),
		})
		return
	}

	verseKey := r.FormValue("verse_key")

	var request sqlc.RecitationFileCreateRecitationFileParams

	request.Reciter = reciter
	request.Slug = slug
	request.VerseKey = verseKey

	recitationFile, err := db.Queries.RecitationFileCreateRecitationFile(context.Background(), request)
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, render.M{
			"message": "Error creating recitation file",
			"error":   err.Error(),
		})
		return
	}

	file, _, err := r.FormFile("file")
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, render.M{
			"message": "Error retrieving uploaded file",
			"error":   err.Error(),
		})
		return
	}
	defer file.Close()

	baseDir := filepath.Join("data", "uploads", request.Reciter, request.Slug)
	err = os.MkdirAll(baseDir, 0755)
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, render.M{
			"message": "Error creating recitation directory",
			"error":   err.Error(),
		})
		return
	}

	rawFilepath := filepath.Join(baseDir, request.VerseKey+".raw")
	transcodedFilepath := filepath.Join(baseDir, request.VerseKey+".mp3")

	fileHandler, err := os.Create(rawFilepath)
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, render.M{
			"message": "Error creating recitation file",
			"error":   err.Error(),
		})
		return
	}
	defer fileHandler.Close()

	if _, err := io.Copy(fileHandler, file); err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, render.M{
			"message": "Error copying recitation file contents",
			"error":   err.Error(),
		})
		return
	}
	cmd := exec.Command("ffmpeg", "-i", rawFilepath, transcodedFilepath)

	_, err = cmd.CombinedOutput()
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, render.M{
			"message": "Error transcoding recitation file",
			"error":   err.Error(),
		})
		return
	}

	err = os.RemoveAll(rawFilepath)
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, render.M{
			"message": "Error deleting raw recitation file",
			"error":   err.Error(),
		})
		return
	}

	render.JSON(w, r, recitationFile)
}

// GetRecitationFiles godoc
//
//	@Tags		RecitationFile
//	@Produce	json
//
//	@Param		reciter	path		string	true	"Reciter"
//	@Param		slug	path		string	true	"Slug"
//
//	@Success	200		{object}	[]sqlc.RecitationFile
//	@Failure	500		{object}	models.Error
//	@Router		/recitation-files/{reciter}/{slug} [get]
func GetRecitationFiles(w http.ResponseWriter, r *http.Request) {
	request := sqlc.RecitationFileSelectRecitationFilesParams{
		Reciter: chi.URLParam(r, "reciter"),
		Slug:    chi.URLParam(r, "slug"),
	}
	recitationFiles, err := db.Queries.RecitationFileSelectRecitationFiles(context.Background(), request)
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, render.M{
			"message": "Error querying recitation files",
			"error":   err.Error(),
		})
		return
	}

	render.JSON(w, r, recitationFiles)
}

// GetRecitationFile godoc
//
//	@Tags		RecitationFile
//	@Produce	json
//
//	@Param		reciter		path		string	true	"Reciter"
//	@Param		slug		path		string	true	"Slug"
//	@Param		verse_key	path		string	true	"Verse key"
//
//	@Success	200			{object}	sqlc.RecitationFile
//	@Failure	500			{object}	models.Error
//	@Router		/recitation-files/{reciter}/{slug}/{verse_key} [get]
func GetRecitationFile(w http.ResponseWriter, r *http.Request) {
	request := sqlc.RecitationFileSelectRecitationFileParams{
		Reciter:  chi.URLParam(r, "reciter"),
		Slug:     chi.URLParam(r, "slug"),
		VerseKey: chi.URLParam(r, "verse_key"),
	}
	recitationFile, err := db.Queries.RecitationFileSelectRecitationFile(context.Background(), request)
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, render.M{
			"message": "Error querying recitation file",
			"error":   err.Error(),
		})
		return
	}

	render.JSON(w, r, recitationFile)
}

// DeleteRecitationFile godoc
//
//	@Tags		RecitationFile
//	@Produce	json
//
//	@Param		X-CSRF-TOKEN	header		string	true	"CSRF Token"
//
//	@Param		slug			path		string	true	"Slug"
//	@Param		verse_key		path		string	true	"Verse Key"
//
//	@Success	200				{object}	sqlc.RecitationFile
//	@Failure	400				{object}	models.Error
//	@Failure	401				{object}	models.Error
//	@Router		/recitation-files/{slug}/{verse_key} [delete]
func DeleteRecitationFile(w http.ResponseWriter, r *http.Request) {
	reciter := r.Context().Value("username").(string)
	slug := chi.URLParam(r, "slug")
	verseKey := chi.URLParam(r, "verse_key")

	request := sqlc.RecitationFileDeleteRecitationFileParams{
		Reciter:  reciter,
		Slug:     slug,
		VerseKey: verseKey,
	}

	deletedRecitationFile, err := db.Queries.RecitationFileDeleteRecitationFile(context.Background(), request)
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, render.M{
			"message": "Unexpected error",
			"error":   err.Error(),
		})
		return
	}

	render.JSON(w, r, deletedRecitationFile)

	audioFilepath := filepath.Join("data", "uploads", reciter, slug, verseKey+".mp3")
	timingsFilepath := filepath.Join("data", "uploads", reciter, slug, verseKey+".json")

	err = os.RemoveAll(audioFilepath)
	if err != nil {
		log.Printf("Error deleting audio file: %v\n", err)
	}
	os.RemoveAll(timingsFilepath)
	if err != nil {
		log.Printf("Error deleting timings file: %v\n", err)
	}
}
