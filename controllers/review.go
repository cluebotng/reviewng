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
	"html/template"
	"net/http"
)

func (app *App) ReviewHandler(w http.ResponseWriter, r *http.Request) {
	// Not logged in, send to the login page
	user := app.getAuthenticatedUser(r)
	if user == nil {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	// Not approved, return an error
	if !user.Approved {
		http.Error(w, "Forbidden", 403)
		return
	}

	// Render the template - most the logic is in javascript
	t, err := template.ParseFS(app.fsTemplates, "templates/review.tmpl")
	if err != nil {
		panic(err)
	}

	if err := t.Execute(w, struct{ User *db.User }{User: user}); err != nil {
		panic(err)
	}
}
