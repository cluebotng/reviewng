package controllers

// MIT License
//
// Copyright (c) 2021 Damian Zaremba
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

import (
	"encoding/json"
	"github.com/cluebotng/reviewng/db"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

func (app *App) ApiUserClassificationListHandler(w http.ResponseWriter, r *http.Request) {
	// Not logged in, return an error
	user := app.getAuthenticatedUser(r)
	if user == nil {
		http.Error(w, "Unauthorized", 401)
		return
	}

	// Not an admin, return an error
	if !user.Admin {
		http.Error(w, "Forbidden", 403)
		return
	}

	// Get all classifications keyed by id
	allUserClassifications := map[int]*db.UserClassification{}
	userClassifications, err := app.dbh.FetchAllUserClassifications()
	if err != nil {
		panic(err)
	}

	for _, userClassification := range userClassifications {
		allUserClassifications[userClassification.Id] = userClassification
	}

	response, err := json.Marshal(allUserClassifications)
	if err != nil {
		panic(err)
	}

	w.Header().Set("Content-Type", "application/json")
	if _, err := w.Write(response); err != nil {
		panic(err)
	}
}

func (app *App) ApiUserClassificationCreateHandler(w http.ResponseWriter, r *http.Request) {
	// Not logged in, return an error
	user := app.getAuthenticatedUser(r)
	if user == nil {
		http.Error(w, "Unauthorized", 401)
		return
	}

	userClassification := struct {
		EditId         int    `json:"edit_id"`
		Classification int    `json:"classification"`
		Comment        string `json:"comment"`
		Confirmation   bool   `json:"confirmation"`
	}{}

	// Decode the request
	if err := json.NewDecoder(r.Body).Decode(&userClassification); err != nil {
		panic(err)
	}

	edit, err := app.dbh.LookupEditById(userClassification.EditId)
	if err != nil {
		panic(err)
	}

	if edit == nil {
		http.Error(w, "Not Found", 404)
		return
	}

	editClassification, err := app.dbh.CalculateEditClassification(edit)
	if err != nil {
		panic(err)
	}

	// Ask the user to confirm if the classification is statistically different
	requiresConfirmation := false
	if editClassification != db.EDIT_CLASSIFICATION_UNKNOWN && editClassification != userClassification.Classification {
		if !userClassification.Confirmation {
			requiresConfirmation = true
		}
	}

	if !requiresConfirmation {
		// We are all good - either confirmed or inline
		if err := app.dbh.CreateUserClassification(db.UserClassification{
			UserId:         user.Id,
			Comment:        userClassification.Comment,
			Classification: userClassification.Classification,
			EditId:         userClassification.EditId,
		}); err != nil {
			panic(err)
		}
	}

	response, err := json.Marshal(map[string]bool{"require_confirmation": requiresConfirmation})
	if err != nil {
		panic(err)
	}
	w.Header().Set("Content-Type", "application/json")
	if _, err := w.Write(response); err != nil {
		panic(err)
	}
}

func (app *App) ApiUserClassificationGetHandler(w http.ResponseWriter, r *http.Request) {
	// Not logged in, return an error
	user := app.getAuthenticatedUser(r)
	if user == nil {
		http.Error(w, "Unauthorized", 401)
		return
	}

	// Not an admin, return an error
	if !user.Admin {
		http.Error(w, "Forbidden", 403)
		return
	}

	// Decode the request
	getUserClassificationId, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		panic(err)
	}
	userClassification, err := app.dbh.LookupUserClassificationsByEditId(getUserClassificationId)
	if err != nil {
		panic(err)
	}

	response, err := json.Marshal(userClassification)
	if err != nil {
		panic(err)
	}

	w.Header().Set("Content-Type", "application/json")
	if _, err := w.Write(response); err != nil {
		panic(err)
	}
}
