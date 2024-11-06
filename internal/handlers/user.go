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

type updateUserDTO struct {
	Displayname string `json:"displayname"`
}

// GetUsers godoc
//
//	@Tags		User
//	@Produce	json
//
//	@Success	200	{object}	[]sqlc.UserSelectUsersRow
//	@Failure	500	{object}	models.Error
//	@Router		/users [get]
func GetUsers(w http.ResponseWriter, r *http.Request) {
	users, err := db.Queries.UserSelectUsers(context.Background())
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, render.M{
			"message": "Error querying users",
			"error":   err.Error(),
		})
		return
	}

	render.JSON(w, r, users)
}

// GetUser godoc
//
//	@Tags		User
//	@Produce	json
//
//	@Param		username	path		string	true	"Username"
//
//	@Success	200			{object}	[]sqlc.UserSelectUsersRow
//	@Failure	400			{object}	models.Error
//	@Router		/users/{username} [get]
func GetUser(w http.ResponseWriter, r *http.Request) {
	user, err := db.Queries.UserSelectUser(context.Background(), chi.URLParam(r, "username"))
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, render.M{
			"message": "User does not exist",
			"error":   err.Error(),
		})
		return
	}

	render.JSON(w, r, user)
}

// UpdateUser godoc
//
//	@Tags		User
//	@Accept		json
//	@Produce	json
//
//	@Param		X-CSRF-TOKEN	header		string			true	"CSRF Token"
//
//	@Param		request			body		updateUserDTO	true	"Update Recitation"
//
//	@Success	200				{object}	sqlc.UserUpdateUserRow
//	@Failure	400				{object}	models.Error
//	@Failure	401				{object}	models.Error
//	@Router		/user [put]
func UpdateUser(w http.ResponseWriter, r *http.Request) {
	username := r.Context().Value("username").(string)

	existingUser, err := db.Queries.UserSelectUser(context.Background(), username)
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, render.M{
			"message": "User does not exist",
			"error":   err.Error(),
		})
		return
	}

	updatedUserData := &sqlc.UserUpdateUserParams{
		Username:    existingUser.Username,
		Displayname: existingUser.Displayname,
	}

	var request sqlc.UserUpdateUserParams
	err = json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, render.M{
			"message": "Could not parse request body",
			"error":   err.Error(),
		})
		return
	}

	if request.Displayname != "" {
		updatedUserData.Displayname = request.Displayname
	}

	updatedUser, err := db.Queries.UserUpdateUser(context.Background(), *updatedUserData)
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, render.M{
			"message": "Error while updating user",
			"error":   err.Error(),
		})
		return
	}

	render.JSON(w, r, updatedUser)
}

// DeleteUser godoc
//
//	@Tags		User
//	@Produce	json
//
//	@Param		X-CSRF-TOKEN	header		string	true	"CSRF Token"
//
//	@Success	200				{object}	sqlc.UserDeleteUserRow
//	@Failure	400				{object}	models.Error
//	@Failure	401				{object}	models.Error
//	@Router		/user [delete]
func DeleteUser(w http.ResponseWriter, r *http.Request) {
	username := r.Context().Value("username").(string)

	_, err := db.Queries.UserSelectUser(context.Background(), username)
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, render.M{
			"message": "User does not exist",
			"error":   err.Error(),
		})
		return
	}

	deletedUser, err := db.Queries.UserDeleteUser(context.Background(), username)
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, render.M{
			"message": "Error deleting user",
			"error":   err.Error(),
		})
		return
	}

	userDir := filepath.Join("data", "uploads", username)
	err = os.RemoveAll(userDir)
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, render.M{
			"message": "Error deleting user directory",
			"error":   err.Error(),
		})
		return
	}

	render.JSON(w, r, deletedUser)
}
