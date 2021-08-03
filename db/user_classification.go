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

type UserClassification struct {
	Id             int
	UserId         int
	Comment        string
	Classification int
	EditId         int
}

func (db *Db) CreateUserClassification(newUserClassification UserClassification) error {
	insert, err := db.db.Query("INSERT INTO user_classification (user_id, edit_id, comment, classification) VALUES (?, ?, ?, ?)", newUserClassification.UserId, newUserClassification.EditId, newUserClassification.Comment, newUserClassification.Classification)
	if err != nil {
		return err
	}

	if err := insert.Close(); err != nil {
		return err
	}
	return nil
}

func (db *Db) LookupUserClassificationsById(id int) (*UserClassification, error) {
	results, err := db.db.Query("SELECT id, user_id, comment, classification, edit_id FROM user_classification WHERE id = ?", id)
	if err != nil {
		return nil, err
	}

	if !results.Next() {
		return nil, nil
	}

	c := &UserClassification{}
	if err := results.Scan(&c.Id, &c.UserId, &c.Comment, &c.Classification, &c.EditId); err != nil {
		return nil, err
	}

	if err := results.Close(); err != nil {
		return nil, err
	}

	return c, nil
}

func (db *Db) LookupUserClassificationsByEditId(id int) ([]*UserClassification, error) {
	results, err := db.db.Query("SELECT id, user_id, comment, classification, edit_id FROM user_classification WHERE edit_id = ?", id)
	if err != nil {
		return nil, err
	}

	classifications := []*UserClassification{}
	for results.Next() {
		c := &UserClassification{}
		if err := results.Scan(&c.Id, &c.UserId, &c.Comment, &c.Classification, &c.EditId); err != nil {
			return nil, err
		}
		classifications = append(classifications, c)
	}

	if err := results.Close(); err != nil {
		return nil, err
	}

	return classifications, nil
}

func (db *Db) LookupUserClassificationsByUserId(id int) ([]*UserClassification, error) {
	results, err := db.db.Query("SELECT id, user_id, comment, classification, edit_id FROM user_classification WHERE user_id = ?", id)
	if err != nil {
		return nil, err
	}

	classifications := []*UserClassification{}
	for results.Next() {
		c := &UserClassification{}
		if err := results.Scan(&c.Id, &c.UserId, &c.Comment, &c.Classification, &c.EditId); err != nil {
			return nil, err
		}
		classifications = append(classifications, c)
	}

	if err := results.Close(); err != nil {
		return nil, err
	}

	return classifications, nil
}

func (db *Db) FetchAllUserClassifications() ([]*UserClassification, error) {
	results, err := db.db.Query("SELECT id, user_id, comment, classification, edit_id FROM user_classification")
	if err != nil {
		return nil, err
	}

	classifications := []*UserClassification{}
	for results.Next() {
		c := &UserClassification{}
		if err := results.Scan(&c.Id, &c.UserId, &c.Comment, &c.Classification, &c.EditId); err != nil {
			return nil, err
		}
		classifications = append(classifications, c)
	}

	if err := results.Close(); err != nil {
		return nil, err
	}

	return classifications, nil
}
