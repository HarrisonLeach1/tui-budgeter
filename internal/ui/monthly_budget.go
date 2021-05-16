package ui

import (
	"fmt"
	"math"
	"strconv"

	"github.com/HarrisonLeach1/tui-budgeter/internal/api"
	"github.com/HarrisonLeach1/tui-budgeter/internal/api/models"
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

	openingBalance, err := strconv.ParseFloat(bankReport.Rows[1].Rows[0].Cells[1].Value, 64)
	if err != nil {
		panic("what")
	}
	closingBalance, err := strconv.ParseFloat(bankReport.Rows[1].Rows[0].Cells[4].Value, 64)
	if err != nil {
		panic("the")
	}

	summaryList := clui.CreateListBox(summary, 10, 10, clui.AutoSize)
	summaryList.AddItem(fmt.Sprintf("Increase in savings: %%%.2f", ((closingBalance-openingBalance)/openingBalance)*100))
	summaryList.AddItem(fmt.Sprintf("Actual Savings: $%.2f", actualIncome-actualExpenses))
	summaryList.AddItem(fmt.Sprintf("Planned Savings: $%.2f", budgetedIncome-budgetedExpenses))
	summaryList.SetTextColor(clui.ColorWhite)
	summaryList.SetActiveTextColor(clui.ColorWhiteBold)
	summaryList.SetBackColor(clui.ColorBlack)
	summaryList.SetActiveBackColor(clui.ColorBlack)

	updatePlannedSavings := func(amount float64) {
		summaryList.RemoveItem(2)
		summaryList.AddItem(fmt.Sprintf("Planned Savings: $%.2f", (budgetedIncome - budgetedExpenses)))
		summaryList.Draw()
	}

	maxExp := math.Max(budgetedExpenses, actualExpenses)
	expBar := createFramedProgressBar(left, int(budgetedExpenses), int(maxExp), "Planned Expenses", fmt.Sprintf("$%.2f", budgetedExpenses))
	actExpBar := createFramedProgressBar(left, int(actualExpenses), int(maxExp), "Actual Expenses", fmt.Sprintf("$%.2f", actualExpenses))

	updatePlannedExpenses := func(amount float64) {
		budgetedExpenses = budgetedExpenses + amount
		maxExp := math.Max(budgetedExpenses, actualExpenses)
		expBar.SetLimits(0, int(maxExp))
		actExpBar.SetLimits(0, int(maxExp))
		expBar.SetValue(int(budgetedExpenses))
		expBar.SetTitle(fmt.Sprintf("$%.2f", budgetedExpenses))
		expBar.Draw()
		actExpBar.Draw()
	}

	createTable(left, "Expenses Breakdown", pnlReport, budgetReport, "Less Operating Expenses", updatePlannedSavings, updatePlannedExpenses)

	right := clui.CreateFrame(view, viewWidth, clui.AutoSize, clui.BorderNone, clui.AutoSize)
	right.SetPack(clui.Vertical)
	right.SetGaps(clui.KeepValue, 1)
	right.SetPaddings(1, 1)

	createBalanceChart(right, openingBalance, closingBalance)

	maxInc := math.Max(budgetedIncome, actualIncome)
	incBar := createFramedProgressBar(right, int(budgetedIncome), int(maxInc), "Planned Income", fmt.Sprintf("$%.2f", budgetedIncome))
	actIncBar := createFramedProgressBar(right, int(actualIncome), int(maxInc), "Actual Income", fmt.Sprintf("$%.2f", actualIncome))

	updatePlannedIncome := func(amount float64) {
		budgetedIncome = budgetedIncome + amount
		max := math.Max(budgetedIncome, actualIncome)
		incBar.SetValue(int(budgetedIncome))
		incBar.SetLimits(0, int(max))
		actIncBar.SetLimits(0, int(max))
		incBar.SetTitle(fmt.Sprintf("$%.2f", budgetedIncome))
		incBar.Draw()
		actIncBar.Draw()
	}

	createTable(right, "Income Breakdown", pnlReport, budgetReport, "Income", updatePlannedSavings, updatePlannedIncome)
	return nil
}

func createBalanceChart(parent *clui.Frame, openingBalance float64, closingBalance float64) {

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

	data := []clui.BarData{
		{Value: openingBalance, Title: "Opening Balance", Fg: clui.ColorBlue},
		{Value: closingBalance, Title: "Closing Balance", Fg: clui.ColorGreen},
	}
	barChart.SetData(data)
}

func createFramedProgressBar(parent *clui.Frame, value int, maxValue int, title string, valueLabel string) *clui.ProgressBar {

	barFrame := clui.CreateFrame(parent, 1, 1, clui.BorderThin, clui.AutoSize)
	barFrame.SetPack(clui.Vertical)
	barFrame.SetGaps(clui.KeepValue, 1)
	barFrame.SetPaddings(1, 1)
	barFrame.SetTitle(title)

	bar := clui.CreateProgressBar(barFrame, 1, 1, clui.AutoSize)
	bar.SetLimits(0, maxValue)
	bar.SetTitle(valueLabel)
	bar.SetTitleColor(clui.ColorBlack)
	bar.SetValue(value)
	bar.SetBackColor(clui.ColorWhite)
	bar.SetActiveBackColor(clui.ColorCyan)
	bar.SetTextColor(clui.ColorBlack)
	return bar
}

func createTable(parent *clui.Frame, title string, pnlReport models.Report, budgetReport models.Report, sectionTitle string, updatePlannedSavings func(float64), updateBar func(float64)) {
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

	budgetSection := findReportSection(budgetReport, sectionTitle, true)
	pnlMap := *createCategoryToValueMap(pnlReport, sectionTitle)
	table.SetRowCount(len(budgetSection.Rows))

	isIncome := sectionTitle == "Income"

	values := make([][4]float64, len(budgetSection.Rows))
	categories := make([]string, len(budgetSection.Rows))

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
		categories[i] = category
		values[i][1] = budgetValue
		values[i][2] = actualValue
		values[i][3] = (budgetValue) - (actualValue)
	}

	table.OnDrawCell(func(info *clui.ColumnDrawInfo) {
		if info.Col == 0 {
			info.Text = categories[info.Row]
		} else if info.Col == 3 {
			diff := values[info.Row][3]
			// TODO: simplify, what da heck is going on here im too tired
			if diff < 0 && !isIncome || diff >= 0 && isIncome {
				if !isIncome {
					diff *= -1
				}
				info.Text = fmt.Sprintf("-$%.2f", diff)
			} else {
				if isIncome {
					diff *= -1
				}
				info.Text = fmt.Sprintf("$%.2f", diff)
			}

			if info.Text[0:1] == "-" {
				info.Fg = clui.ColorRed
			} else {
				info.Fg = clui.ColorGreen
			}

		} else {
			info.Text = fmt.Sprintf("$%.2f", values[info.Row][info.Col])
		}
	})

	table.OnAction(func(ev clui.TableEvent) {
		btns := []string{"Close", "Dismiss"}
		var action string
		switch ev.Action {
		case clui.TableActionEdit:
			c := ev.Col
			r := ev.Row

			// only allow editing of budget column
			if c != 1 {
				return
			}

			oldVal := values[r][c] // ignore the dollar sign
			dlg := clui.CreateEditDialog(
				fmt.Sprintf("Editing budget for %s", categories[r]), "New value", fmt.Sprintf("%.2f", oldVal),
			)
			dlg.OnClose(func() {
				switch dlg.Result() {
				case clui.DialogButton1:
					newText := dlg.EditResult()

					newFloat, err := strconv.ParseFloat(newText, 64)
					if err != nil {
						// TODO
						fmt.Errorf("input validation error")
					}
					values[r][c] = newFloat
					values[r][3] = newFloat - values[r][2]

					if isIncome {
						updateBar(newFloat - oldVal)
						updatePlannedSavings(newFloat - oldVal)
					} else {
						updateBar(newFloat - oldVal)
						updatePlannedSavings(oldVal - newFloat)
					}

					clui.PutEvent(clui.Event{Type: clui.EventRedraw})
				}
			})
			return
		default:
			action = "Unknown action"
		}

		dlg := clui.CreateConfirmationDialog(
			"<c:blue>"+action,
			"Click any button or press <c:yellow>SPACE<c:> to close the dialog",
			btns, clui.DialogButton1)
		dlg.OnClose(func() {})
	})

}

func createCategoryToValueMap(report models.Report, sectionTitle string) *map[string]string {
	m := make(map[string]string)
	section := findReportSection(report, sectionTitle, true)
	for _, row := range section.Rows {
		m[row.Cells[0].Value] = row.Cells[1].Value
	}
	return &m
}

func findReportSection(report models.Report, title string, excludeSummary bool) models.ReportRow {
	for _, section := range report.Rows {
		if section.Title == title {
			if excludeSummary {
				return filterSectionByRowType(section, "Row")
			}
			return section
		}
	}
	return models.ReportRow{}
}

func filterSectionByRowType(section models.ReportRow, rowTypeName string) models.ReportRow {
	newRows := models.ReportRow{}
	for _, row := range section.Rows {
		if row.RowType == rowTypeName {
			newRows.Rows = append(newRows.Rows, row)
		}
	}
	return newRows
}

func findSectionTotal(report models.Report, title string) float64 {
	section := findReportSection(report, title, false)
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
