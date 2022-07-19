package db

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
	"database/sql"
	"fmt"
	"github.com/cluebotng/reviewng/cfg"
	_ "github.com/go-sql-driver/mysql"
	"time"
)

type Db struct {
	db *sql.DB
}

func NewDb(cfg *cfg.Config) (*Db, error) {
	url := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", cfg.Db.User, cfg.Db.Pass, cfg.Db.Host, cfg.Db.Port, cfg.Db.Name)

	database, err := sql.Open("mysql", url)
	if err != nil {
		return nil, err
	}
	database.SetConnMaxLifetime(time.Second * 10)
	database.SetMaxOpenConns(50)
	database.SetMaxIdleConns(1)

	db := Db{db: database}
	return &db, nil
}
