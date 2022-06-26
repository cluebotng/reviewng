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
	"github.com/cluebotng/reviewng/cache"
	"github.com/cluebotng/reviewng/cfg"
	"github.com/cluebotng/reviewng/db"
	"github.com/dghubble/oauth1"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"net/http"
	"os"
	"sync"
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
	cacheStore   *cache.InMemoryStorage
	dbh          *db.Db
	oauth        *oauth1.Config
	fsTemplates  *embed.FS
	fsStatic     *embed.FS
	trainingSync sync.Mutex
}

func NewApp(cfg *cfg.Config, fsTemplates, fsStatic *embed.FS) *App {
	oauth := oauth1.Config{
		ConsumerKey:    cfg.OAuth.Token,
		ConsumerSecret: cfg.OAuth.Secret,
		CallbackURL:    "oob",
		Endpoint: oauth1.Endpoint{
			RequestTokenURL: "https://en.wikipedia.org/w/index.php?title=Special:OAuth/initiate",
			AuthorizeURL:    "https://en.wikipedia.org/w/index.php?title=Special:OAuth/authorize",
			AccessTokenURL:  "https://en.wikipedia.org/w/index.php?title=Special:OAuth/token",
		},
	}

	dbh, err := db.NewDb(cfg)
	if err != nil {
		panic(err)
	}

	session := sessions.NewCookieStore([]byte(cfg.Session.SecretKey))
	memoryCache := cache.NewInMemoryStorage()
	app := App{
		config:       cfg,
		sessionStore: session,
		cacheStore:   memoryCache,
		dbh:          dbh,
		oauth:        &oauth,
		fsTemplates:  fsTemplates,
		fsStatic:     fsStatic,
	}
	return &app
}

func (app *App) initializeRoutes() {
	app.router = mux.NewRouter()
	app.router.PathPrefix("/static/").Handler(http.FileServer(NoIndexFileSystem{http.FS(app.fsStatic)})).Methods("GET")

	app.router.HandleFunc("/login/callback", app.LoginCallbackHandler).Methods("GET")
	app.router.HandleFunc("/login", app.LoginHandler).Methods("GET")
	app.router.HandleFunc("/logout", app.LogoutHandler).Methods("GET")

	app.router.HandleFunc("/api/user", app.ApiUserListHandler).Methods("GET")
	app.router.HandleFunc("/api/user", app.ApiUserCreateHandler).Methods("POST")
	app.router.HandleFunc("/api/user/{id}", app.ApiUserGetHandler).Methods("GET")
	app.router.HandleFunc("/api/user/{id}", app.ApiUserUpdateHandler).Methods("UPDATE")

	app.router.HandleFunc("/api/edit-group", app.ApiEditGroupListHandler).Methods("GET")
	app.router.HandleFunc("/api/edit-group", app.ApiEditGroupCreateHandler).Methods("POST")
	app.router.HandleFunc("/api/edit-group/{id}", app.ApiEditGroupGetHandler).Methods("GET")
	app.router.HandleFunc("/api/edit-group/{id}", app.ApiEditGroupUpdateHandler).Methods("UPDATE")

	app.router.HandleFunc("/api/edit", app.ApiEditListHandler).Methods("GET")
	app.router.HandleFunc("/api/edit", app.ApiEditCreateHandler).Methods("POST")
	app.router.HandleFunc("/api/edit/next", app.ApiEditNextHandler).Methods("GET")
	app.router.HandleFunc("/api/edit/{id}", app.ApiEditGetHandler).Methods("GET")
	app.router.HandleFunc("/api/edit/{id}", app.ApiEditUpdateHandler).Methods("UPDATE")

	app.router.HandleFunc("/api/user-classification", app.ApiUserClassificationListHandler).Methods("GET")
	app.router.HandleFunc("/api/user-classification", app.ApiUserClassificationCreateHandler).Methods("POST")
	app.router.HandleFunc("/api/user-classification/{id}", app.ApiUserClassificationGetHandler).Methods("GET")

	app.router.HandleFunc("/api/cron/stats", app.ApiCronStatsHandler).Methods("GET")
	app.router.HandleFunc("/api/report/import", app.ApiReportImportHandler).Methods("GET")
	app.router.HandleFunc("/api/report/export", app.ApiReportExportHandler).Methods("GET")

	app.router.HandleFunc("/api/training/import", app.ApiTrainingImportHandler).Methods("GET")
	app.router.HandleFunc("/api/training/data/{id}", app.ApiTrainingDataHandler).Methods("GET")

	app.router.HandleFunc("/api/export/done", app.ApiExportDoneHandler).Methods("GET")
	app.router.HandleFunc("/api/export/done.json", app.ApiExportDoneJsonHandler).Methods("GET")
	app.router.HandleFunc("/api/export/dump", app.ApiExportDumpHandler).Methods("GET")
	app.router.HandleFunc("/api/export/dump.json", app.ApiExportDumpJsonHandler).Methods("GET")
	app.router.HandleFunc("/api/export/trainer.json", app.ApiExportTrainerJsonHandler).Methods("GET")
	app.router.HandleFunc("/api/config", app.ApiConfigHandler).Methods("GET")

	app.router.HandleFunc("/", app.WelcomeHandler).Methods("GET")
	app.router.HandleFunc("/review", app.ReviewHandler).Methods("GET")

	app.router.HandleFunc("/admin", app.AdminHandler).Methods("GET")
	app.router.HandleFunc("/admin/users", app.AdminUsersHandler).Methods("GET")
	app.router.HandleFunc("/admin/edit-groups", app.AdminEditGroupsHandler).Methods("GET")
	app.router.HandleFunc("/admin/edit-groups/{id}", app.AdminEditGroupDetailHandler).Methods("GET")
	app.router.HandleFunc("/admin/details/{id}", app.AdminEditDetailsHandler).Methods("GET")
}

func (app *App) RunForever(addr string) {
	app.initializeRoutes()
	server := &http.Server{
		Addr:         addr,
		WriteTimeout: time.Second * 120,
		ReadTimeout:  time.Second * 10,
		IdleTimeout:  time.Second * 60,
		Handler:      app.router,
	}
	if err := server.ListenAndServe(); err != nil {
		panic(err)
	}
}
