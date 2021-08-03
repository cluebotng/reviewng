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
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

func (app *App) ApiReportImportHandler(w http.ResponseWriter, r *http.Request) {
	resp, _ := http.Get("https://cluebotng.toolforge.org/api.php")
	body, _ := ioutil.ReadAll(resp.Body)
	if err := resp.Body.Close(); err != nil {
		panic(err)
	}

	// Create a map of known ids
	knownEditIds := map[int]bool{}
	allEdits, err := app.dbh.FetchAllEdits()
	if err != nil {
		panic(err)
	}
	for _, edit := range allEdits {
		knownEditIds[edit.Id] = false
	}

	// Fetch the edit group we log these into
	eg, err := app.dbh.LookupEditGroupByName("Report Interface Import")
	if err != nil {
		panic(err)
	}

	// Create entries for everything we don't know about
	for _, line := range strings.Split(string(body), "\n") {
		newEditId, err := strconv.Atoi(line)
		if err != nil {
			panic(err)
		}

		// Already have this edit, ignore it
		if _, ok := knownEditIds[newEditId]; ok {
			continue
		}

		// Create a new entry
		if err := app.dbh.CreateEdit(newEditId, eg, 2, db.EDIT_CLASSIFICATION_CONSTRUCTIVE); err != nil {
			panic(err)
		}
	}
}
