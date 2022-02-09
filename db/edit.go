package db

import (
	"context"
	"log"
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

type Edit struct {
	Id                              int
	Required                        int
	Classification                  int
	UserClassificationsVandalism    int
	UserClassificationsConstructive int
	UserClassificationsSkipped      int
}

func (edit *Edit) ReviewedClassification() int {
	sum := edit.UserClassificationsConstructive + edit.UserClassificationsVandalism + edit.UserClassificationsSkipped
	max := MaxInt(edit.UserClassificationsConstructive, MaxInt(edit.UserClassificationsVandalism, edit.UserClassificationsSkipped))

	if max < edit.Required {
		return EDIT_CLASSIFICATION_UNKNOWN
	}
	if 2*edit.UserClassificationsSkipped > sum {
		return EDIT_CLASSIFICATION_SKIPPED
	}
	if edit.UserClassificationsConstructive >= 3*edit.UserClassificationsVandalism {
		return EDIT_CLASSIFICATION_CONSTRUCTIVE
	}
	if edit.UserClassificationsVandalism >= 3*edit.UserClassificationsConstructive {
		return EDIT_CLASSIFICATION_VANDALISM
	}
	return EDIT_CLASSIFICATION_UNKNOWN
}

func MaxInt(x, y int) int {
	if x < y {
		return y
	}
	return x
}

func (db *Db) CreateEdit(id int, eg *EditGroup, required, classification int) error {
	ctx := context.Background()
	tx, err := db.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	if _, err := tx.ExecContext(ctx, "INSERT INTO edit (id, required, classification) VALUES (?, ?, ?)", id, required, classification); err != nil {
		if err := tx.Rollback(); err != nil {
			log.Fatal(err)
		}
		return err
	}

	if _, err := tx.ExecContext(ctx, "INSERT INTO edit_edit_group (edit_id, edit_group_id) VALUES (?, ?)", id, eg.Id); err != nil {
		if err := tx.Rollback(); err != nil {
			log.Fatal(err)
		}
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}

func (db *Db) LookupEditById(id int) (*Edit, error) {
	results, err := db.db.Query("SELECT id, required, classification "+
		"COUNT(DISTINCT user_classification_vandalism.id) AS user_classifications_vandalism, "+
		"COUNT(DISTINCT user_classification_constructive.id) AS user_classifications_constructive, "+
		"COUNT(DISTINCT user_classification_skipped.id) AS user_classifications_skipped "+
		"FROM edit "+
		"LEFT JOIN user_classification AS user_classification_vandalism ON (user_classification_vandalism.edit_id = edit.id AND user_classification_vandalism.classification = 0) "+
		"LEFT JOIN user_classification AS user_classification_constructive ON (user_classification_constructive.edit_id = edit.id AND user_classification_constructive.classification = 1) "+
		"LEFT JOIN user_classification AS user_classification_skipped ON (user_classification_skipped.edit_id = edit.id AND user_classification_skipped.classification = 2) "+
		"WHERE edit.id = ? GROUP BY edit.id, edit.required, edit.classification", id)
	if err != nil {
		return nil, err
	}

	if !results.Next() {
		return nil, nil
	}

	edit := &Edit{}
	if err := results.Scan(&edit.Id, &edit.Required, &edit.Classification, &edit.UserClassificationsVandalism, &edit.UserClassificationsConstructive, &edit.UserClassificationsSkipped); err != nil {
		return nil, err
	}

	if err := results.Close(); err != nil {
		return nil, err
	}

	return edit, nil
}

func (db *Db) LookupEditsByGroupId(id int) ([]*Edit, error) {
	results, err := db.db.Query("SELECT edit.id, edit.required, edit.classification, "+
		"COUNT(DISTINCT user_classification_vandalism.id) AS user_classifications_vandalism, "+
		"COUNT(DISTINCT user_classification_constructive.id) AS user_classifications_constructive, "+
		"COUNT(DISTINCT user_classification_skipped.id) AS user_classifications_skipped "+
		"FROM edit "+
		"INNER JOIN edit_edit_group ON (edit_edit_group.edit_id = edit.id) "+
		"INNER JOIN edit_group ON (edit_group.id = edit_edit_group.edit_group_id) "+
		"LEFT JOIN user_classification AS user_classification_vandalism ON (user_classification_vandalism.edit_id = edit.id AND user_classification_vandalism.classification = 0) "+
		"LEFT JOIN user_classification AS user_classification_constructive ON (user_classification_constructive.edit_id = edit.id AND user_classification_constructive.classification = 1) "+
		"LEFT JOIN user_classification AS user_classification_skipped ON (user_classification_skipped.edit_id = edit.id AND user_classification_skipped.classification = 2) "+
		"WHERE edit_group.id = ? "+
		"GROUP BY edit.id, edit.required, edit.classification", id)
	if err != nil {
		return nil, err
	}

	edits := []*Edit{}
	for results.Next() {
		edit := Edit{}
		if err := results.Scan(&edit.Id, &edit.Required, &edit.Classification, &edit.UserClassificationsVandalism, &edit.UserClassificationsConstructive, &edit.UserClassificationsSkipped); err != nil {
			return nil, err
		}
		edits = append(edits, &edit)
	}

	if err := results.Close(); err != nil {
		return nil, err
	}

	return edits, nil
}

func (db *Db) FetchAllEdits() ([]*Edit, error) {
	results, err := db.db.Query("SELECT edit.id, edit.required, edit.classification, " +
		"COUNT(DISTINCT user_classification_vandalism.id) AS user_classifications_vandalism, " +
		"COUNT(DISTINCT user_classification_constructive.id) AS user_classifications_constructive, " +
		"COUNT(DISTINCT user_classification_skipped.id) AS user_classifications_skipped " +
		"FROM edit " +
		"INNER JOIN edit_edit_group ON (edit_edit_group.edit_id = edit.id) " +
		"INNER JOIN edit_group ON (edit_group.id = edit_edit_group.edit_group_id) " +
		"LEFT JOIN user_classification AS user_classification_vandalism ON (user_classification_vandalism.edit_id = edit.id AND user_classification_vandalism.classification = 0) " +
		"LEFT JOIN user_classification AS user_classification_constructive ON (user_classification_constructive.edit_id = edit.id AND user_classification_constructive.classification = 1) " +
		"LEFT JOIN user_classification AS user_classification_skipped ON (user_classification_skipped.edit_id = edit.id AND user_classification_skipped.classification = 2) " +
		"GROUP BY edit.id, edit.required, edit.classification")
	if err != nil {
		return nil, err
	}

	edits := []*Edit{}
	for results.Next() {
		edit := &Edit{}
		if err := results.Scan(&edit.Id, &edit.Required, &edit.Classification, &edit.UserClassificationsVandalism, &edit.UserClassificationsConstructive, &edit.UserClassificationsSkipped); err != nil {
			return nil, err
		}
		edits = append(edits, edit)
	}

	if err := results.Close(); err != nil {
		return nil, err
	}

	return edits, nil
}

func (db *Db) CalculateEditStatus(edit *Edit) (int, error) {
	sum := edit.UserClassificationsConstructive + edit.UserClassificationsVandalism + edit.UserClassificationsSkipped
	max := MaxInt(edit.UserClassificationsConstructive, MaxInt(edit.UserClassificationsVandalism, edit.UserClassificationsSkipped))

	if sum == 0 {
		return EDIT_STATUS_NOT_DONE, nil
	}

	if max >= edit.Required {
		return EDIT_STATUS_DONE, nil
	}

	return EDIT_STATUS_PARTIAL, nil
}
