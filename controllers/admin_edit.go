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
	"github.com/cluebotng/reviewng/db"
	"github.com/gorilla/mux"
	"html/template"
	"net/http"
	"strconv"
)

func (app *App) AdminEditDetailsHandler(w http.ResponseWriter, r *http.Request) {
	// Not logged in, send to the login page
	user := app.getAuthenticatedUser(r)
	if user == nil {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	// Not an admin, return an error
	if !user.Admin {
		http.Error(w, "Forbidden", 403)
		return
	}

	// Decode the request
	editId, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		panic(err)
	}

	// Lookup the edit
	edit, err := app.dbh.LookupEditById(editId)
	if err != nil {
		panic(err)
	}
	if edit == nil {
		http.Error(w, "Edit Not Found", 404)
		return
	}

	// Lookup all users
	allUsers, err := app.dbh.FetchAllUsers()
	if err != nil {
		panic(err)
	}

	userNamesById := map[int]string{}
	for _, user := range allUsers {
		userNamesById[user.Id] = user.Username
	}

	// Map all user classifications
	userClassifications, err := app.dbh.LookupUserClassificationsByEditId(edit.Id)
	if err != nil {
		panic(err)
	}

	type userEditClassification struct {
		Username       string
		Classification string
		Comment        string
	}
	editUserClassifications := []userEditClassification{}
	for _, userClassification := range userClassifications {
		editUserClassifications = append(editUserClassifications, userEditClassification{
			Username:       userNamesById[userClassification.UserId],
			Classification: ConvertClassificationToHumanString(userClassification.Classification),
			Comment:        userClassification.Comment,
		})
	}

	t, err := template.New("details.tmpl").Funcs(template.FuncMap{
		"classificationToHuman": ConvertClassificationToHumanString,
	}).ParseFS(app.fsTemplates, "templates/admin/details.tmpl")
	if err != nil {
		panic(err)
	}

	if err := t.Execute(w, struct {
		Edit                  *db.Edit
		CurrentClassification int
		UserClassifications   []userEditClassification
	}{
		Edit:                  edit,
		CurrentClassification: edit.ReviewedClassification(),
		UserClassifications:   editUserClassifications,
	}); err != nil {
		panic(err)
	}
}
