package ui

import (
	"github.com/HarrisonLeach1/xero-tui/internal/api"
	"github.com/VladimirMarkelov/clui"
)

// Renders a window containing a budget report for the given period
//
// Params fromDate and toDate should be in the format "yyyy-mm-dd" e.g. "2021-05-31"
func RenderBudgetReport(fromDate string, toDate string) error {
	report, err := api.GetProfitAndLossStatement(fromDate, toDate)
	if err != nil {
		return err
	}
	view := clui.AddWindow(0, 0, 50, 50, "Monthly Budget")

	viewWidth, _ := view.Size()

	left := clui.CreateFrame(view, viewWidth, clui.AutoSize, clui.BorderNone, clui.AutoSize)
	left.SetPack(clui.Vertical)
	left.SetGaps(clui.KeepValue, 1)
	left.SetPaddings(1, 1)
	left.SetActive(true)

	topLeft := clui.CreateFrame(left, 10, 10, clui.BorderNone, clui.AutoSize)
	topLeft.SetPack(clui.Vertical)
	topLeft.SetGaps(clui.KeepValue, 1)
	topLeft.SetPaddings(1, 1)

	summary := clui.CreateFrame(topLeft, clui.AutoSize, clui.AutoSize, clui.BorderThin, clui.AutoSize)
	summary.SetPack(clui.Vertical)
	summary.SetGaps(clui.KeepValue, 1)
	summary.SetPaddings(1, 1)

	summaryList := clui.CreateListBox(summary, 10, 10, clui.AutoSize)
	summaryList.AddItem("Increase in savings: %142")
	summaryList.AddItem("Planned Savings: $1428")

	createFramedProgressBar(left, 1415, 2980, "Planned Expenses", "$1415.25")
	createFramedProgressBar(left, 2980, 2980, "Actual Expenses", "$2979.45")

	createTable(left, "Expenses Breakdown")

	right := clui.CreateFrame(view, viewWidth, clui.AutoSize, clui.BorderNone, clui.AutoSize)
	right.SetPack(clui.Vertical)
	right.SetGaps(clui.KeepValue, 1)
	right.SetPaddings(1, 1)

	barChartFrame := clui.CreateFrame(right, 10, 10, clui.BorderThin, clui.AutoSize)
	barChartFrame.SetPack(clui.Vertical)
	barChartFrame.SetGaps(clui.KeepValue, 1)
	barChartFrame.SetPaddings(1, 1)

	barChart := clui.CreateBarChart(barChartFrame, 10, 10, clui.AutoSize)
	barChart.SetShowTitles(true)
	barChart.SetShowMarks(true)
	barChart.SetValueWidth(10)
	barChart.SetMinBarWidth(20)
	barChart.SetBarGap(5)

	data := []clui.BarData{
		{Value: 1584, Title: "Opening Balance", Fg: clui.ColorBlue},
		{Value: 3837, Title: "Closing Balance", Fg: clui.ColorGreen},
	}
	barChart.SetData(data)

	createFramedProgressBar(right, 1415, 2980, "Planned Income", "$1415.25")
	createFramedProgressBar(right, 2980, 2980, "Actual Income", "$2979.45")

	createTable(right, "Income Breakdown")
	return nil
}

func createFramedProgressBar(parent *clui.Frame, value int, maxValue int, title string, valueLabel string) {

	barFrame := clui.CreateFrame(parent, 1, 1, clui.BorderThin, clui.AutoSize)
	barFrame.SetPack(clui.Vertical)
	barFrame.SetGaps(clui.KeepValue, 1)
	barFrame.SetPaddings(1, 1)
	barFrame.SetTitle(title)

	bar := clui.CreateProgressBar(barFrame, 1, 1, clui.AutoSize)
	bar.SetLimits(0, maxValue)
	bar.SetTitle(valueLabel)
	bar.SetValue(value)
}

func createTable(parent *clui.Frame, title string) {
	tableFrame := clui.CreateFrame(parent, clui.AutoSize, clui.AutoSize, clui.BorderThin, clui.AutoSize)
	tableFrame.SetPack(clui.Vertical)
	tableFrame.SetGaps(clui.KeepValue, 1)
	tableFrame.SetPaddings(1, 1)
	tableFrame.SetTitle(title)

	table := clui.CreateTableView(tableFrame, 10, 10, clui.AutoSize)
	table.SetShowLines(true)
	table.SetShowRowNumber(true)
	cols := []clui.Column{
		clui.Column{Title: "Category", Width: 10, Alignment: clui.AlignLeft},
		clui.Column{Title: "Planned", Width: 10, Alignment: clui.AlignRight},
		clui.Column{Title: "Actual", Width: 10, Alignment: clui.AlignRight},
		clui.Column{Title: "Diff", Width: 10, Alignment: clui.AlignRight},
	}

	table.SetColumns(cols)
}
