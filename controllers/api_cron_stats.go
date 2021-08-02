package controllers

import (
	"bytes"
	"github.com/cluebotng/reviewng/db"
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

	if _, err := w.Write(tpl.Bytes()); err != nil {
		panic(err)
	}
}
