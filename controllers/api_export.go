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
	"encoding/xml"
	"github.com/cluebotng/reviewng/db"
	"net/http"
)

type User struct {
	Key             int
	Nick            string
	Classifications int
}

type EditGroup struct {
	Key      int
	Name     string
	Weight   int
	Edits    []Edit `xml:"Edits>Edit,omitempty"`
	Reviewed []Edit `xml:"Reviewed>Edit,omitempty"`
	Done     []Edit `xml:"Done>Edit,omitempty"`
}

type Edit struct {
	Key                    int
	Id                     int
	Weight                 int
	Required               int
	Constructive           int
	Skipped                int
	Vandalism              int
	OriginalClassification string
	RealClassification     string
	Comments               []string `xml:"Comments>Comment,omitempty"`
	Users                  []string `xml:"Users>User,omitempty"`
}

type Data struct {
	EditGroups []EditGroup `xml:"EditGroups>EditGroup,omitempty"`
	Users      []User      `xml:"Users>User,omitempty"`
}

func calculateDataDump(app *App, done bool) Data {
	// Fetch user data
	userData := []User{}
	allUsers, err := app.dbh.FetchAllUsers()
	if err != nil {
		panic(err)
	}

	userNameById := map[int]string{}
	for _, user := range allUsers {
		userNameById[user.Id] = user.Username
		if !done {
			userClassifications, err := app.dbh.CalculateTotalUserClassifications(user)
			if err != nil {
				panic(err)
			}
			userData = append(userData, User{
				Key:             user.Id,
				Nick:            user.Username,
				Classifications: userClassifications,
			})
		}
	}

	// Fetch edit group data
	editGroups := []EditGroup{}
	allEditGroups, err := app.dbh.FetchAllEditGroups()
	if err != nil {
		panic(err)
	}
	for _, editGroup := range allEditGroups {
		editGroupEdits, err := app.dbh.LookupEditsByGroupId(editGroup.Id)
		if err != nil {
			panic(err)
		}

		allEdits, reviewedEdits, doneEdits := []Edit{}, []Edit{}, []Edit{}
		for _, e := range editGroupEdits {
			editClassification, err := app.dbh.CalculateEditClassification(e)
			if err != nil {
				panic(err)
			}

			constructive, skipped, vandalism := 0, 0, 0
			allComments, allUsers := []string{}, []string{}
			userEditClassification, err := app.dbh.LookupUserClassificationsByEditId(e.Id)
			if err != nil {
				panic(err)
			}
			for _, uec := range userEditClassification {
				if uec.Comment != "" {
					allComments = append(allComments, uec.Comment)
				}
				if val, ok := userNameById[uec.UserId]; ok {
					allUsers = append(allUsers, val)
				}
				if uec.Classification == db.EDIT_CLASSIFICATION_CONSTRUCTIVE {
					constructive += 1
					continue
				}
				if uec.Classification == db.EDIT_CLASSIFICATION_SKIPPED {
					skipped += 1
					continue
				}
				if uec.Classification == db.EDIT_CLASSIFICATION_VANDALISM {
					vandalism += 1
					continue
				}
			}

			edit := Edit{
				Key:                    e.Id,
				Id:                     e.Id,
				Weight:                 editGroup.Weight,
				Required:               e.Required,
				Constructive:           constructive,
				Skipped:                skipped,
				Vandalism:              vandalism,
				OriginalClassification: ConvertClassificationToString(e.Classification),
				RealClassification:     ConvertClassificationToString(editClassification),
				Comments:               allComments,
				Users:                  allUsers,
			}
			allEdits = append(allEdits, edit)
			if constructive+skipped+vandalism > 0 {
				reviewedEdits = append(reviewedEdits, edit)
			}
			if editClassification != db.EDIT_CLASSIFICATION_UNKNOWN {
				doneEdits = append(doneEdits, edit)
			}
		}

		eg := EditGroup{
			Key:    editGroup.Id,
			Name:   editGroup.Name,
			Weight: editGroup.Weight,
			Done:   doneEdits,
		}
		if !done {
			eg.Edits = allEdits
			eg.Reviewed = reviewedEdits
		}
		editGroups = append(editGroups, eg)
	}

	data := Data{EditGroups: editGroups}
	if !done {
		data.Users = userData
	}
	return data
}

func (app *App) ApiExportDumpHandler(w http.ResponseWriter, r *http.Request) {
	if _, err := w.Write([]byte(xml.Header)); err != nil {
		panic(err)
	}
	encoder := xml.NewEncoder(w)
	encoder.Indent("", "\t")

	data := calculateDataDump(app, false)
	if err := encoder.Encode(&data); err != nil {
		panic(err)
	}
}

func (app *App) ApiExportDoneHandler(w http.ResponseWriter, r *http.Request) {
	if _, err := w.Write([]byte(xml.Header)); err != nil {
		panic(err)
	}
	encoder := xml.NewEncoder(w)
	encoder.Indent("", "\t")

	data := calculateDataDump(app, true)
	if err := encoder.Encode(&data); err != nil {
		panic(err)
	}
}
