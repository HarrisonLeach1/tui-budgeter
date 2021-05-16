package ui

import (
	"fmt"
	"strconv"
	"time"

	"github.com/VladimirMarkelov/clui"
)

func SelectReportPeriod() {
	periods := generatePeriods()
	dlgType := clui.SelectDialogList

	selDlg := clui.CreateSelectDialog("Choose a period", periods, 0, dlgType)
	selDlg.OnClose(func() {
		switch selDlg.Result() {
		case clui.DialogButton1:
			idx := selDlg.Value()
			selected := periods[idx]
			RenderBudgetReport(getStartAndEndDates(selected))
		}
		// ask the composer to repaint all windows
		clui.PutEvent(clui.Event{Type: clui.EventRedraw})
	})
}

// Given a selected month period e.g. "2021-05"
// Generates the start and end dates e.g. "2021-05-01" and "2021-05-31"
func getStartAndEndDates(selected string) (string, string) {
	month := selected[len(selected)-2:]
	year := selected[0:4]
	y, err := strconv.Atoi(year)
	if err != nil {
		panic("Could not parse selected year value")
	}
	m, err := strconv.Atoi(month)
	if err != nil {
		panic("Could not parse selected month value")
	}
	fromDate := time.Date(y, time.Month(m), 1, 0, 0, 0, 0, time.Local)
	toDate := fromDate.AddDate(0, 1, 0).Add(time.Nanosecond * -1)
	return fromDate.Format("2006-01-02"), toDate.Format("2006-01-02")

}

func generatePeriods() []string {

	var periods []string
	endYear, endMonthString, _ := time.Now().Date()

	// HACK: Only generate periods onwards from Jan 2021
	startMonthString := time.January
	startYear := 2021

	for i := startYear; i <= endYear; i++ {
		var endMonth int
		if i != endYear {
			endMonth = 11
		} else {
			endMonth = int(endMonthString) - 1
		}

		var startMonth int
		if i == startYear {
			startMonth = int(startMonthString) - 1
		} else {
			startMonth = 0
		}

		for j := startMonth; j <= endMonth; j = getNextMonth(j) {
			month := j + 1
			monthStr := strconv.Itoa(month)
			if month < 10 {
				monthStr = "0" + monthStr
			}
			periods = append([]string{fmt.Sprintf("%d-%s", i, monthStr)}, periods...)
		}
	}
	return periods
}

// Uses modulo to keep month numbers under 12
func getNextMonth(num int) int {
	if num > 12 {
		val := num % 12
		if val == 0 {
			return 11
		}
		return val
	}
	return num + 1
}
