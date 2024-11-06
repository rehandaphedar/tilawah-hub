package handlers

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"

	"git.sr.ht/~rehandaphedar/tilawah-hub/internal/db"
	"git.sr.ht/~rehandaphedar/tilawah-hub/internal/sqlc"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/spf13/viper"
)

// Lafzize godoc
//
//	@Tags		lafzize
//	@Produce	json
//
//	@Param		X-CSRF-TOKEN	header		string	true	"CSRF Token"
//
//	@Param		slug			path		string	true	"Recitation slug"
//	@Param		verse_key		path		string	true	"Verse key of recitation"
//	@Success	200				{object}	sqlc.RecitationFileUpdateRecitationFileParams
//	@Failure	400				{object}	models.Error
//	@Failure	401				{object}	models.Error
//	@Router		/lafzize/{slug}/{verse_key} [post]
func Lafzize(w http.ResponseWriter, r *http.Request) {
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
			"message": "Error checking status of recitation",
			"error":   err.Error(),
		})
		return
	}

	if existingRecitationFile.LafzizeProcessing {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, render.M{
			"message": "The recitation is already being lafzized",
			"error":   "",
		})
		return
	}

	updatedRecitationFile, err := db.Queries.RecitationFileUpdateRecitationFile(context.Background(), sqlc.RecitationFileUpdateRecitationFileParams{
		Reciter:           reciter,
		Slug:              slug,
		VerseKey:          verseKey,
		HasTimings:        false,
		LafzizeProcessing: true,
	})
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, render.M{
			"message": "Error updating status of recitation file",
			"error":   err.Error(),
		})
		return
	}

	audioPath := filepath.Join("data", "uploads", reciter, slug, fmt.Sprintf("%s.mp3", verseKey))
	timingsPath := filepath.Join("data", "uploads", reciter, slug, fmt.Sprintf("%s.json", verseKey))

	file, err := os.Open(audioPath)
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, render.M{
			"message": "Recitation file does not exist.",
			"error":   err.Error(),
		})
		return
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile("file", filepath.Base(file.Name()))
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, render.M{
			"message": "Error making request to lafzize server",
			"error":   err.Error(),
		})
		return
	}

	_, err = io.Copy(part, file)
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, render.M{
			"message": "Error making request to lafzize server",
			"error":   err.Error(),
		})
		return
	}

	err = writer.WriteField("verse_key", verseKey)
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, render.M{
			"message": "Error making request to lafzize server",
			"error":   err.Error(),
		})
		return
	}

	err = writer.Close()
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, render.M{
			"message": "Error making request to lafzize server",
			"error":   err.Error(),
		})
		return
	}

	lafzizeRequest, err := http.NewRequest("POST", viper.GetString("lafzize_endpoint"), body)
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, render.M{
			"message": "Error making request to lafzize server",
			"error":   err.Error(),
		})
		return
	}

	lafzizeRequest.Header.Add("Content-Type", writer.FormDataContentType())
	client := &http.Client{}

	err = os.RemoveAll(timingsPath)
	if err != nil {
		log.Printf("Error removing possible existing timing file: %v\n", err)
	}

	go doAsyncRequest(client, lafzizeRequest, timingsPath, updatedRecitationFile)

	render.JSON(w, r, updatedRecitationFile)
}

func doAsyncRequest(client *http.Client, r *http.Request, timingsPath string, recitationFile sqlc.RecitationFile) {
	resp, err := client.Do(r)
	if err != nil {
		log.Printf("Error making request while lafzizing recitation %v: %v\n", timingsPath, err)
		return
	}

	file, err := os.Create(timingsPath)
	if err != nil {
		log.Printf("Error creating timings file while lafzizing recitation %v: %v\n", timingsPath, err)
		return
	}
	defer file.Close()

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		log.Printf("Error writing to timings file while lafzizing recitation %v: %v\n", timingsPath, err)
		return
	}
	err = resp.Body.Close()
	if err != nil {
		log.Printf("Error closing reader while lafzizing recitation %v: %v\n", timingsPath, err)
		return
	}

	_, err = db.Queries.RecitationFileUpdateRecitationFile(context.Background(), sqlc.RecitationFileUpdateRecitationFileParams{
		Reciter:           recitationFile.Reciter,
		Slug:              recitationFile.Slug,
		VerseKey:          recitationFile.VerseKey,
		HasTimings:        true,
		LafzizeProcessing: false,
	})
	if err != nil {
		log.Printf("Error updating status while lafzizing recitation %v: %v\n", timingsPath, err)
		return
	}
}
