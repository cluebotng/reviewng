package db

import (
	"encoding/json"
)

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

type TrainingData struct {
	Current struct {
		Id      int    `json:"id"`
		Comment string `json:"comment"`
		User    struct {
			Name               string `json:"name"`
			EditCount          int    `json:"edit_count"`
			DistinctPagesCount int    `json:"distinct_pages_count"`
			WarningCount       int    `json:"warning_count"`
			RegistrationTime   int    `json:"registration_time"`
		} `json:"user"`
		Minor     bool   `json:"minor"`
		Timestamp int    `json:"timestamp"`
		Text      string `json:"text"`
	} `json:"current"`
	Previous struct {
		Id      int    `json:"id"`
		Comment string `json:"comment"`
		User    struct {
			Name string `json:"name"`
		} `json:"user"`
		Minor     bool   `json:"minor"`
		Timestamp int    `json:"timestamp"`
		Text      string `json:"text"`
	} `json:"previous"`
	Page struct {
		Title                string `json:"title"`
		Namespace            string `json:"namespace"`
		Creator              string `json:"creator"`
		CreationTime         int    `json:"creation_time"`
		RecentEditCount      int    `json:"recent_edit_count"`
		RecentReversionCount int    `json:"recent_reversion_count"`
	} `json:"page"`
	IsVandalism bool `json:"is_vandalism"`
}

func (db *Db) StoreTrainingDataForEdit(id int, td *TrainingData) error {
	jsonTrainingData, err := json.Marshal(td)
	if err != nil {
		panic(err)
	}

	if _, err := db.db.Query("REPLACE INTO edit_training_data (edit_id, training_data) VALUES (?, ?)", id, string(jsonTrainingData)); err != nil {
		return err
	}

	return nil
}

func (db *Db) GetTrainingDataEdits() (map[int]*TrainingData, error) {
	results, err := db.db.Query("SELECT edit_id, training_data FROM edit_training_data")
	if err != nil {
		return nil, err
	}

	trainingDataEdits := map[int]*TrainingData{}
	for results.Next() {
		row := struct {
			EditId       int
			TrainingData []byte
		}{}
		if err := results.Scan(&row.EditId, &row.TrainingData); err != nil {
			return nil, err
		}

		trainingDataEdit := TrainingData{}
		if err := json.Unmarshal(row.TrainingData, &trainingDataEdit); err != nil {
			panic(err)
		}
		trainingDataEdits[row.EditId] = &trainingDataEdit
	}

	if err := results.Close(); err != nil {
		return nil, err
	}

	return trainingDataEdits, nil
}

func (db *Db) GetTrainingDataByEditId(id int) (*TrainingData, error) {
	results, err := db.db.Query("SELECT training_data FROM edit_training_data WHERE edit_id = ?", id)
	if err != nil {
		return nil, err
	}

	if !results.Next() {
		return nil, nil
	}

	rawData := []byte{}
	if err := results.Scan(&rawData); err != nil {
		return nil, err
	}

	trainingData := TrainingData{}
	if err := json.Unmarshal(rawData, &trainingData); err != nil {
		panic(err)
	}
	return &trainingData, nil
}
