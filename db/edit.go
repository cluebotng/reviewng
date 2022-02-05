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

type Edit struct {
	Id             int
	Required       int
	Classification int
}

func MaxInt(x, y int) int {
	if x < y {
		return y
	}
	return x
}

func (db *Db) CreateEdit(id int, eg *EditGroup, required, classification int) error {
	insert, err := db.db.Query("INSERT INTO edit (id, required, classification) VALUES (?, ?, ?); INSERT INTO edit_edit_group (edit_id, edit_group_id) VALUES (?, ?)", id, required, classification, id, eg.Id)
	if err != nil {
		return err
	}

	if err := insert.Close(); err != nil {
		return err
	}
	return nil
}

func (db *Db) LookupEditById(id int) (*Edit, error) {
	results, err := db.db.Query("SELECT id, required, classification FROM edit WHERE id = ?", id)
	if err != nil {
		return nil, err
	}

	if !results.Next() {
		return nil, nil
	}

	edit := &Edit{}
	if err := results.Scan(&edit.Id, &edit.Required, &edit.Classification); err != nil {
		return nil, err
	}

	if err := results.Close(); err != nil {
		return nil, err
	}

	return edit, nil
}

func (db *Db) LookupEditsByGroupId(id int) ([]*Edit, error) {
	results, err := db.db.Query("SELECT edit.id, edit.required, edit.classification FROM edit "+
		"INNER JOIN edit_edit_group ON (edit_edit_group.edit_id = edit.id) "+
		"WHERE edit_group.id = ?", id)
	if err != nil {
		return nil, err
	}

	edits := []*Edit{}
	for results.Next() {
		edit := Edit{}
		if err := results.Scan(&edit.Id, &edit.Required, &edit.Classification); err != nil {
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
	results, err := db.db.Query("SELECT id, required, classification FROM edit")
	if err != nil {
		return nil, err
	}

	edits := []*Edit{}
	for results.Next() {
		edit := &Edit{}
		if err := results.Scan(&edit.Id, &edit.Required, &edit.Classification); err != nil {
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
	classifications, err := db.LookupUserClassificationsByEditId(edit.Id)
	if err != nil {
		return -1, err
	}

	constructive, vandalism, skipped := 0, 0, 0
	for _, classification := range classifications {
		if classification.Classification == EDIT_CLASSIFICATION_CONSTRUCTIVE {
			constructive += 1
		} else if classification.Classification == EDIT_CLASSIFICATION_VANDALISM {
			vandalism += 1
		} else if classification.Classification == EDIT_CLASSIFICATION_SKIPPED {
			skipped += 1
		}
	}
	sum := constructive + vandalism + skipped
	max := MaxInt(constructive, MaxInt(vandalism, skipped))

	if sum == 0 {
		return EDIT_STATUS_NOT_DONE, nil
	}

	if max >= edit.Required {
		return EDIT_STATUS_DONE, nil
	}

	return EDIT_STATUS_PARTIAL, nil
}

func (db *Db) CalculateEditClassification(edit *Edit) (int, error) {
	classifications, err := db.LookupUserClassificationsByEditId(edit.Id)
	if err != nil {
		return -1, err
	}

	constructive, vandalism, skipped := 0, 0, 0
	for _, classification := range classifications {
		if classification.Classification == EDIT_CLASSIFICATION_CONSTRUCTIVE {
			constructive += 1
		} else if classification.Classification == EDIT_CLASSIFICATION_VANDALISM {
			vandalism += 1
		} else if classification.Classification == EDIT_CLASSIFICATION_SKIPPED {
			skipped += 1
		}
	}

	sum := constructive + vandalism + skipped
	max := MaxInt(constructive, MaxInt(vandalism, skipped))
	if max < edit.Required {
		return EDIT_CLASSIFICATION_UNKNOWN, nil
	}
	if 2*skipped > sum {
		return EDIT_CLASSIFICATION_SKIPPED, nil
	}
	if constructive >= 3*vandalism {
		return EDIT_CLASSIFICATION_CONSTRUCTIVE, nil
	}
	if vandalism >= 3*constructive {
		return EDIT_CLASSIFICATION_VANDALISM, nil
	}
	return EDIT_CLASSIFICATION_UNKNOWN, nil
}
