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
	"github.com/gorilla/sessions"
	"net/http"
)

func (app *App) getSessionStore(r *http.Request) *sessions.Session {
	session, _ := app.sessionStore.Get(r, "cluebotng-review")
	return session
}

func (app *App) getAuthenticatedUser(r *http.Request) *db.User {
	session := app.getSessionStore(r)
	if userId, ok := session.Values["user.id"]; ok {
		if user, err := app.dbh.LookupUserById(userId.(int)); err == nil {
			return user
		}
	}
	return nil
}

func (app *App) setAuthenticatedUser(r *http.Request, w http.ResponseWriter, user *db.User) error {
	session := app.getSessionStore(r)
	session.Values["user.id"] = user.Id
	return session.Save(r, w)
}

func (app *App) clearSessionData(r *http.Request, w http.ResponseWriter) error {
	session := app.getSessionStore(r)
	session.Values = map[interface{}]interface{}{}
	return session.Save(r, w)
}

func ConvertClassificationToString(classification int) string {
	if classification == db.EDIT_CLASSIFICATION_VANDALISM {
		return "V"
	}
	if classification == db.EDIT_CLASSIFICATION_CONSTRUCTIVE {
		return "C"
	}
	if classification == db.EDIT_CLASSIFICATION_SKIPPED {
		return "S"
	}
	return "U"
}
