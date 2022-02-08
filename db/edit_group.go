package db

import (
	"sort"
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

type EditGroup struct {
	Id     int
	Name   string
	Weight int
}

func (db *Db) LookupEditGroupById(id int) (*EditGroup, error) {
	results, err := db.db.Query("SELECT id, name, weight FROM edit_group WHERE id = ?", id)
	if err != nil {
		return nil, err
	}

	if !results.Next() {
		return nil, nil
	}

	group := &EditGroup{}
	if err := results.Scan(&group.Id, &group.Name, &group.Weight); err != nil {
		return nil, err
	}

	if err := results.Close(); err != nil {
		return nil, err
	}

	return group, nil
}

func (db *Db) LookupEditGroupByName(name string) (*EditGroup, error) {
	results, err := db.db.Query("SELECT id, name, weight FROM edit_group WHERE name = ?", name)
	if err != nil {
		return nil, err
	}

	if !results.Next() {
		return nil, nil
	}

	group := &EditGroup{}
	if err := results.Scan(&group.Id, &group.Name, &group.Weight); err != nil {
		return nil, err
	}

	if err := results.Close(); err != nil {
		return nil, err
	}

	return group, nil
}

func (db *Db) FetchAllEditGroups() ([]*EditGroup, error) {
	results, err := db.db.Query("SELECT id, name, weight FROM edit_group")
	if err != nil {
		return nil, err
	}

	editGroups := []*EditGroup{}
	for results.Next() {
		editGroup := &EditGroup{}
		if err := results.Scan(&editGroup.Id, &editGroup.Name, &editGroup.Weight); err != nil {
			return nil, err
		}
		editGroups = append(editGroups, editGroup)
	}

	if err := results.Close(); err != nil {
		return nil, err
	}

	return editGroups, nil
}

func (db *Db) CalculateRandomPendingEditForUser(user *User) (*Edit, error) {
	allGroups, err := db.FetchAllEditGroups()
	if err != nil {
		return nil, err
	}

	// Sort groups by weight
	sort.Slice(allGroups[:], func(i, j int) bool {
		return allGroups[i].Weight < allGroups[j].Weight
	})

	// Get all user edit ids
	userClassifications, err := db.LookupUserClassificationsByUserId(user.Id)
	if err != nil {
		return nil, err
	}
	knownUserEdits := map[int]bool{}
	for _, userClassification := range userClassifications {
		knownUserEdits[userClassification.EditId] = true
	}

	for _, group := range allGroups {
		// Calculate if we have some edits in this selected group
		groupEdits, err := db.LookupEditsByGroupId(group.Id)
		if err != nil {
			return nil, err
		}
		for _, edit := range groupEdits {
			// First check if we don't already have a classification for this edit
			if _, ok := knownUserEdits[edit.Id]; ok {
				continue
			}

			// Return the edit if it still needs classifying
			if edit.ReviewedClassification() == EDIT_CLASSIFICATION_UNKNOWN {
				return edit, nil
			}
		}
	}
	return nil, nil
}
