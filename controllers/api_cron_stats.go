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
	"bytes"
	"github.com/cluebotng/reviewng/db"
	"github.com/cluebotng/reviewng/wikipedia"
	"html/template"
	"net/http"
)

type userContributionStat struct {
	Username           string
	Admin              bool
	EditCount          int
	AccuracyCount      int
	AccuracyPercentage float32
}

type editGroupStat struct {
	Name       string
	Weight     int
	Partial    int
	NotStarted int
	Done       int
}

func calculateUserContributionStats(app *App) []userContributionStat {
	stats := []userContributionStat{}

	allUsers, err := app.dbh.FetchAllUsers()
	if err != nil {
		panic(err)
	}

	for _, user := range allUsers {
		userTotal, err := app.dbh.CalculateTotalUserClassifications(user)
		if err != nil {
			panic(err)
		}
		userAccuracy, err := app.dbh.CalculateUserClassificationAccuracy(user)
		if err != nil {
			panic(err)
		}
		stats = append(stats, userContributionStat{
			Username:           user.Username,
			EditCount:          userTotal,
			Admin:              user.Admin,
			AccuracyCount:      userAccuracy.EditCount,
			AccuracyPercentage: userAccuracy.Percentage,
		})
	}

	return stats
}

func calculateEditGroupStats(app *App) []editGroupStat {
	stats := []editGroupStat{}

	allEditGroups, err := app.dbh.FetchAllEditGroups()
	if err != nil {
		panic(err)
	}

	for _, editGroup := range allEditGroups {
		totalDone, totalPartial, totalNotStarted := 0, 0, 0

		edits, err := app.dbh.LookupEditsByGroupId(editGroup.Id)
		if err != nil {
			panic(err)
		}
		for _, edit := range edits {
			editStatus, err := app.dbh.CalculateEditStatus(edit)
			if err != nil {
				panic(err)
			}
			if editStatus == db.EDIT_STATUS_DONE {
				totalDone += 1
			} else if editStatus == db.EDIT_STATUS_PARTIAL {
				totalPartial += 1
			} else if editStatus == db.EDIT_STATUS_NOT_DONE {
				totalNotStarted += 1
			}
		}

		stats = append(stats, editGroupStat{
			Name:       editGroup.Name,
			Weight:     editGroup.Weight,
			Partial:    totalPartial,
			NotStarted: totalNotStarted,
			Done:       totalDone,
		})
	}

	return stats
}

func (app *App) ApiCronStatsHandler(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFS(app.fsTemplates, "templates/stats.tmpl")
	if err != nil {
		panic(err)
	}

	var tpl bytes.Buffer
	if err := t.Execute(&tpl, struct {
		EditGroups []editGroupStat
		AllUsers   []userContributionStat
	}{
		EditGroups: calculateEditGroupStats(app),
		AllUsers:   calculateUserContributionStats(app),
	}); err != nil {
		panic(err)
	}

	if app.config.App.UpdateStats {
		if app.config.Wikipedia.Username != "" {
			if err := wikipedia.UpdatePageWithCredentials(tpl.String(), app.config.Wikipedia.Username, app.config.Wikipedia.Password); err != nil {
				panic(err)
			}
		} else {
			if err := wikipedia.UpdatePage(""); err != nil {
				panic(err)
			}
		}
	}

	if _, err := w.Write(tpl.Bytes()); err != nil {
		panic(err)
	}
}
