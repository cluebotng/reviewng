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
	"fmt"
	"github.com/cluebotng/reviewng/db"
	"io/ioutil"
	"net/http"
)

func downloadTrainingDataForEdit(editId int, editIsVandalism bool) (*db.TrainingData, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("https://cluebotng.toolforge.org/api/?action=training.data&include_text=1&rev_id=%d", editId), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "ClueBot NG Review NG/1.0")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("Error returned from API: %d", resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	trainingData := &db.TrainingData{}
	if err := json.Unmarshal(body, &trainingData); err != nil {
		return nil, err
	}
	trainingData.IsVandalism = editIsVandalism

	if err := resp.Body.Close(); err != nil {
		return nil, err
	}

	return trainingData, nil
}

func (app *App) ApiTrainingImportHandler(w http.ResponseWriter, r *http.Request) {
	app.trainingSync.Lock()
	defer app.trainingSync.Unlock()

	// Find all edits that we should have data for
	allEditGroups, err := app.dbh.FetchAllEditGroups()
	if err != nil {
		panic(err)
	}

	editsRequiringTrainingData := map[int]bool{}
	for _, editGroup := range allEditGroups {
		editGroupEdits, err := app.dbh.LookupEditsByGroupId(editGroup.Id)
		if err != nil {
			panic(err)
		}

		for _, e := range editGroupEdits {
			if e.ReviewedClassification() != db.EDIT_CLASSIFICATION_UNKNOWN {
				editsRequiringTrainingData[e.Id] = e.ReviewedClassification() == db.EDIT_CLASSIFICATION_VANDALISM
			}
		}
	}

	// Find all edits we have training data for
	allEditsWithTrainingData, err := app.dbh.GetTrainingDataEdits()
	if err != nil {
		panic(err)
	}

	editsMissingTrainingData := map[int]bool{}
	for editId, isVandalism := range editsRequiringTrainingData {
		if _, ok := allEditsWithTrainingData[editId]; !ok {
			editsMissingTrainingData[editId] = isVandalism
		}
	}

	// Cache training data for all missing edits
	for editId, isVandalism := range editsMissingTrainingData {
		if trainingData, err := downloadTrainingDataForEdit(editId, isVandalism); err == nil {
			if err := app.dbh.StoreTrainingDataForEdit(editId, trainingData); err != nil {
				panic(err)
			}
		}
	}
}
