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
	"embed"
	"github.com/cluebotng/reviewng/cfg"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"net/http"
	"os"
	"time"
)

type NoIndexFileSystem struct{ fs http.FileSystem }

func (nfs NoIndexFileSystem) Open(path string) (http.File, error) {
	f, err := nfs.fs.Open(path)
	if err != nil {
		return nil, err
	}
	s, err := f.Stat()
	if err != nil {
		return nil, err
	}
	if s.IsDir() {
		return nil, os.ErrNotExist // toHTTPError converts this to a 404
	}
	return f, nil
}

type App struct {
	config       *cfg.Config
	router       *mux.Router
	sessionStore *sessions.CookieStore
	fsTemplates  *embed.FS
	fsStatic     *embed.FS
}

func NewApp(cfg *cfg.Config, fsTemplates, fsStatic *embed.FS) *App {
	session := sessions.NewCookieStore([]byte(cfg.Session.SecretKey))
	app := App{
		config:       cfg,
		sessionStore: session,
		fsTemplates:  fsTemplates,
		fsStatic:     fsStatic,
	}
	return &app
}

func (app *App) initializeRoutes() {
	app.router = mux.NewRouter()
	app.router.PathPrefix("/static/").Handler(http.FileServer(NoIndexFileSystem{http.FS(app.fsStatic)})).Methods("GET")
	app.router.HandleFunc("/", app.WelcomeHandler).Methods("GET")
}

func (app *App) RunForever(addr string) {
	app.initializeRoutes()
	server := &http.Server{
		Addr:         addr,
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      app.router,
	}
	if err := server.ListenAndServe(); err != nil {
		panic(err)
	}
}
