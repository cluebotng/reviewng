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
	"context"
	"encoding/base64"
	"encoding/json"
	"github.com/cluebotng/reviewng/db"
	"github.com/dghubble/oauth1"
	"io/ioutil"
	"net/http"
	"strings"
)

func (app *App) LoginHandler(w http.ResponseWriter, r *http.Request) {
	// Redirect to the login page
	requestToken, requestSecret, err := app.oauth.RequestToken()
	if err != nil {
		panic(err)
	}

	// Store the random secret for this handshake
	session := app.getSessionStore(r)
	session.Values["oauth.request-secret"] = requestSecret
	if err := session.Save(r, w); err != nil {
		panic(err)
	}

	authorizationURL, err := app.oauth.AuthorizationURL(requestToken)
	if err != nil {
		panic(err)
	}

	http.Redirect(w, r, authorizationURL.String(), http.StatusFound)
}

func (app *App) LoginCallbackHandler(w http.ResponseWriter, r *http.Request) {
	requestToken := r.URL.Query().Get("oauth_token")
	verifier := r.URL.Query().Get("oauth_verifier")
	session := app.getSessionStore(r)

	// Get the secret from storage
	requestSecret := ""
	if val, ok := session.Values["oauth.request-secret"]; ok {
		requestSecret = val.(string)
	}

	// Get an access token using the data passed + our initial secret
	accessToken, accessSecret, err := app.oauth.AccessToken(requestToken, requestSecret, verifier)
	if err != nil {
		panic(err)
	}

	// Using the access token fetch the identity
	token := oauth1.NewToken(accessToken, accessSecret)
	httpClient := app.oauth.Client(context.Background(), token)
	resp, _ := httpClient.Get("https://en.wikipedia.org/w/index.php?title=Special:OAuth/identify")
	body, _ := ioutil.ReadAll(resp.Body)
	if err := resp.Body.Close(); err != nil {
		panic(err)
	}

	identity := map[string]interface{}{}
	s, _ := base64.StdEncoding.DecodeString(strings.Split(string(body), ".")[1])

	// TODO: Verify the signature of the JWT using our private key
	if err := json.Unmarshal(s, &identity); err != nil {
		panic(err)
	}

	// We're done with the initial secret
	delete(session.Values, "oauth.request-secret")

	// Lookup the user by name from the identity data
	var user *db.User
	user, err = app.dbh.LookupUserByName(identity["username"].(string))
	if err != nil {
		panic(err)
	}

	// Need to create a user
	if user == nil {
		if err := app.dbh.CreateUser(db.User{
			Username: identity["username"].(string),
			Admin:    false,
			Approved: false,
		}); err != nil {
			panic(err)
		}

		// Lookup the user
		user, err = app.dbh.LookupUserByName(identity["username"].(string))
		if err != nil {
			panic(err)
		}
	}

	if err := app.setAuthenticatedUser(r, w, user); err != nil {
		panic(err)
	}
	http.Redirect(w, r, "/", http.StatusFound)
}

func (app *App) LogoutHandler(w http.ResponseWriter, r *http.Request) {
	if err := app.clearSessionData(r, w); err != nil {
		panic(err)
	}
	http.Redirect(w, r, "/", http.StatusFound)
}
