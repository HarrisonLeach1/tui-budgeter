package ui

import (
	"fmt"
	"strconv"

	"github.com/HarrisonLeach1/xero-tui/internal/api"
	"github.com/HarrisonLeach1/xero-tui/internal/api/models"
	"github.com/VladimirMarkelov/clui"
)

// Renders a window containing a budget report for the given period
//
// Params fromDate and toDate should be in the format "yyyy-mm-dd" e.g. "2021-05-31"
func RenderBudgetReport(fromDate string, toDate string) error {
	pnlReport, err := api.GetProfitAndLossStatement(fromDate, toDate)
	if err != nil {
		return err
	}
	budgetReport, err := api.GetBudgetSummary(fromDate, toDate)
	if err != nil {
		return err
	}
	bankReport, err := api.GetBankSummary(fromDate, toDate)
	if err != nil {
		return err
	}

	budgetedExpenses := findSectionTotal(budgetReport, "Less Operating Expenses")
	actualExpenses := findSectionTotal(pnlReport, "Less Operating Expenses")
	budgetedIncome := findSectionTotal(budgetReport, "Income")
	actualIncome := findSectionTotal(pnlReport, "Income")

	view := clui.AddWindow(0, 0, 60, 60, fmt.Sprintf("Monthly Budget for %s", fromDate[0:7]))

	viewWidth, _ := view.Size()

	left := clui.CreateFrame(view, viewWidth, clui.AutoSize, clui.BorderNone, clui.AutoSize)
	left.SetPack(clui.Vertical)
	left.SetGaps(clui.KeepValue, 1)
	left.SetPaddings(1, 1)
	left.SetActive(true)

	summary := clui.CreateFrame(left, clui.AutoSize, clui.AutoSize, clui.BorderThin, clui.AutoSize)
	summary.SetPack(clui.Vertical)
	summary.SetGaps(clui.KeepValue, 1)
	summary.SetPaddings(1, 1)

	summaryList := clui.CreateListBox(summary, 10, 10, clui.AutoSize)
	summaryList.AddItem("Increase in savings: %142")
	summaryList.AddItem(fmt.Sprintf("Planned Savings: $%.2f", budgetedIncome-budgetedExpenses))
	summaryList.AddItem(fmt.Sprintf("Actual Savings: $%.2f", actualIncome-actualExpenses))
	summaryList.SetTextColor(clui.ColorWhite)
	summaryList.SetActiveTextColor(clui.ColorWhiteBold)
	summaryList.SetBackColor(clui.ColorBlack)
	summaryList.SetActiveBackColor(clui.ColorBlack)

	createFramedProgressBar(left, 1415, 2980, "Planned Expenses", fmt.Sprintf("$%.2f", budgetedExpenses))
	createFramedProgressBar(left, 2980, 2980, "Actual Expenses", fmt.Sprintf("$%.2f", actualExpenses))

	createTable(left, "Expenses Breakdown", pnlReport, budgetReport, "Less Operating Expenses")

	right := clui.CreateFrame(view, viewWidth, clui.AutoSize, clui.BorderNone, clui.AutoSize)
	right.SetPack(clui.Vertical)
	right.SetGaps(clui.KeepValue, 1)
	right.SetPaddings(1, 1)

	createBalanceChart(right, bankReport)

	createFramedProgressBar(right, 1415, 2980, "Planned Income", fmt.Sprintf("$%.2f", budgetedIncome))
	createFramedProgressBar(right, 2980, 2980, "Actual Income", fmt.Sprintf("$%.2f", actualIncome))

	createTable(right, "Income Breakdown", pnlReport, budgetReport, "Income")
	return nil
}

func createBalanceChart(parent *clui.Frame, bankReport models.Report) {

	barChartFrame := clui.CreateFrame(parent, 10, 10, clui.BorderThin, clui.AutoSize)
	barChartFrame.SetPack(clui.Vertical)
	barChartFrame.SetGaps(clui.KeepValue, 1)
	barChartFrame.SetPaddings(1, 1)

	barChart := clui.CreateBarChart(barChartFrame, 10, 10, clui.AutoSize)
	barChart.SetShowTitles(true)
	barChart.SetShowMarks(true)
	barChart.SetValueWidth(10)
	barChart.SetMinBarWidth(20)
	barChart.SetBarGap(5)

	openingBalance, err := strconv.ParseFloat(bankReport.Rows[1].Rows[0].Cells[1].Value, 64)
	if err != nil {
		panic("what")
	}
	closingBalance, err := strconv.ParseFloat(bankReport.Rows[1].Rows[0].Cells[4].Value, 64)
	if err != nil {
		panic("the")
	}

	data := []clui.BarData{
		{Value: openingBalance, Title: "Opening Balance", Fg: clui.ColorBlue},
		{Value: closingBalance, Title: "Closing Balance", Fg: clui.ColorGreen},
	}
	barChart.SetData(data)
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

func createTable(parent *clui.Frame, title string, pnlReport models.Report, budgetReport models.Report, sectionTitle string) {
	tableFrame := clui.CreateFrame(parent, clui.AutoSize, clui.AutoSize, clui.BorderThin, clui.AutoSize)
	tableFrame.SetPack(clui.Vertical)
	tableFrame.SetGaps(clui.KeepValue, 1)
	tableFrame.SetPaddings(1, 1)
	tableFrame.SetTitle(title)

	table := clui.CreateTableView(tableFrame, clui.AutoSize, 30, clui.AutoSize)
	table.SetShowLines(true)
	table.SetShowRowNumber(true)
	cols := []clui.Column{
		clui.Column{Title: "Category", Width: 20, Alignment: clui.AlignLeft},
		clui.Column{Title: "Planned", Width: 10, Alignment: clui.AlignRight},
		clui.Column{Title: "Actual", Width: 10, Alignment: clui.AlignRight},
		clui.Column{Title: "Diff", Width: 10, Alignment: clui.AlignRight},
	}
	table.SetColumns(cols)

	budgetSection := findReportSection(budgetReport, sectionTitle)
	pnlMap := *createCategoryToValueMap(pnlReport, sectionTitle)
	table.SetRowCount(len(budgetSection.Rows))

	isIncome := sectionTitle == "Income"

	values := make([][4]string, len(budgetSection.Rows))
	for i, budgetRow := range budgetSection.Rows {
		category := budgetRow.Cells[0].Value
		actualValueStr := pnlMap[category]

		// If value for the given category is not in the pnl map
		// then no money has been assigned to this category
		actualValue := 0.00
		if actualValueStr != "" {
			val, err := strconv.ParseFloat(actualValueStr, 32)
			if err != nil {
				panic("I don't know what to do: " + actualValueStr + category)
			}
			actualValue = val
		}

		budgetValue, err := strconv.ParseFloat(budgetRow.Cells[1].Value, 32)
		if err != nil {
			panic("How do people even code in go?")
		}
		values[i][0] = category
		values[i][1] = fmt.Sprintf("$%.2f", budgetValue)
		values[i][2] = fmt.Sprintf("$%.2f", actualValue)
		diff := (budgetValue) - (actualValue)

		// TODO: simplify, what da heck is going on here im too tired
		if diff < 0 && !isIncome || diff >= 0 && isIncome {
			if !isIncome {
				diff *= -1
			}
			values[i][3] = fmt.Sprintf("-$%.2f", diff)
		} else {
			if isIncome {
				diff *= -1
			}
			values[i][3] = fmt.Sprintf("$%.2f", diff)
		}
	}
	table.OnDrawCell(func(info *clui.ColumnDrawInfo) {
		info.Text = values[info.Row][info.Col]
		if info.Col == 3 {
			if info.Text[0:1] == "-" {
				info.Fg = clui.ColorRed
			} else {
				info.Fg = clui.ColorGreen
			}
		}
	})

}

func createCategoryToValueMap(report models.Report, sectionTitle string) *map[string]string {
	m := make(map[string]string)
	section := findReportSection(report, sectionTitle)
	for _, row := range section.Rows {
		m[row.Cells[0].Value] = row.Cells[1].Value
	}
	return &m
}

func findReportSection(report models.Report, title string) models.ReportRow {
	for _, section := range report.Rows {
		if section.Title == title {
			return section
		}
	}
	return models.ReportRow{}
}

func findSectionTotal(report models.Report, title string) float64 {
	section := findReportSection(report, title)
	for _, row := range section.Rows {
		if row.RowType == "SummaryRow" {
			val, err := strconv.ParseFloat(row.Cells[1].Value, 32)
			if err != nil {
				panic("This is so sad")
			}
			return val
		}
	}
	return 0
}
