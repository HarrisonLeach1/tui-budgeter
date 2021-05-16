package ui

import (
	"fmt"
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
			RenderMonthlyBudgetReport()
		}
		// ask the composer to repaint all windows
		clui.PutEvent(clui.Event{Type: clui.EventRedraw})
	})
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
			periods = append([]string{fmt.Sprintf("%s %d", time.Month(month), i)}, periods...)
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
