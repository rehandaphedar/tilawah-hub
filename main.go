package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"git.sr.ht/~rehandaphedar/tilawah-hub/internal/db"
	"git.sr.ht/~rehandaphedar/tilawah-hub/internal/handlers"
	"git.sr.ht/~rehandaphedar/tilawah-hub/internal/middlewares"
	"git.sr.ht/~rehandaphedar/tilawah-hub/internal/validators"
	"git.sr.ht/~rehandaphedar/tilawah-hub/pkg/config"
	"github.com/go-chi/chi"
	"github.com/spf13/viper"
)

//	@title		tilawah-hub
//	@version	0.0.1

func main() {
	router := chi.NewRouter()

	err := os.MkdirAll("data", 0755)
	if err != nil {
		log.Fatalf("Error creating data directory: %v", err)
	}

	config.Load()
	db.Connect()
	validators.Initialise()

	router.Group(func(r chi.Router) {
		r.Post("/register", handlers.Register)
		r.Post("/login", handlers.Login)

		r.Get("/users", handlers.GetUsers)
		r.Get("/users/{username}", handlers.GetUser)
	})

	router.Group(func(r chi.Router) {
		r.Use(middlewares.Auth)

		r.Post("/logout", handlers.Logout)
		r.Put("/user", handlers.UpdateUser)
		r.Delete("/user", handlers.DeleteUser)
	})

	router.Group(func(r chi.Router) {
		r.Get("/recitations", handlers.GetRecitations)
		r.Get("/recitations/{reciter}/{slug}", handlers.GetRecitation)
	})

	router.Group(func(r chi.Router) {
		r.Use(middlewares.Auth)

		r.Post("/recitations", handlers.CreateRecitation)
		r.Put("/recitations/{slug}", handlers.UpdateRecitation)
		r.Delete("/recitations/{slug}", handlers.DeleteRecitation)
	})

	router.Group(func(r chi.Router) {
		r.Get("/recitation-files/{reciter}/{slug}", handlers.GetRecitationFiles)
		r.Get("/recitation-files/{reciter}/{slug}/{verse_key}", handlers.GetRecitationFile)
	})

	router.Group(func(r chi.Router) {
		r.Use(middlewares.Auth)

		r.Post("/recitation-files/{slug}", handlers.CreateRecitationFile)
		r.Delete("/recitation-files/{slug}/{verse_key}", handlers.DeleteRecitationFile)

		r.Post("/recitation-timings/{slug}/{verse_key}", handlers.UpdateRecitationTiming)
		r.Delete("/recitation-timings/{slug}/{verse_key}", handlers.DeleteRecitationTiming)
	})

	router.Handle("/uploads/*", http.StripPrefix("/uploads/", http.FileServer(http.Dir(filepath.Join("data", "uploads")))))

	router.Group(func(r chi.Router) {
		r.Use(middlewares.Auth)

		r.Post("/lafzize/{slug}/{verse_key}", handlers.Lafzize)
	})

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", viper.GetInt("port")), router))
}
