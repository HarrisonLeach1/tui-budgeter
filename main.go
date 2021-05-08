package main

import (
	"fmt"
	"log"
	"math"

	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

func main() {
	if err := ui.Init(); err != nil {
		log.Fatalf("failed to initialize termui: %v", err)
	}
	defer ui.Close()

	// build layout
	grid := ui.NewGrid()
	termWidth, termHeight := ui.TerminalDimensions()
	grid.SetRect(0, 0, termWidth, termHeight)

	list := setupList()
	balancesBarGraph := setupBarChart()
	overview := setupOverview()
	xeroTui := widgets.NewParagraph()
	xeroTui.Text = `
	                          _         _ 
	__  _____ _ __ ___       | |_ _   _(_)
	\ \/ / _ \ '__/ _ \ _____| __| | | | |
	 >  <  __/ | | (_) |_____| |_| |_| | |
	/_/\_\___|_|  \___/       \__|\__,_|_|
   
	`
	xeroTui.PaddingLeft = 5
	xeroTui.TextStyle.Fg = ui.ColorCyan

	plannedExpensesGauge, actualExpensesGauge := setupExpenseGauges()
	plannedIncomeGauge, actualIncomeGauge := setupIncomeGauges()
	incomeBreakdown := setupIncomeTable()
	expensesBreakdown := setupExpensesTable()

	grid.Set(
		ui.NewCol(1.0/7, list),
		ui.NewCol(3.0/7,
			ui.NewRow(1.0/5, xeroTui),
			ui.NewRow(1.0/5, overview),
			ui.NewRow(1.0/10, plannedExpensesGauge),
			ui.NewRow(1.0/10, actualExpensesGauge),
			ui.NewRow(2.0/5, expensesBreakdown),
		),
		ui.NewCol(3.0/7,
			ui.NewRow(2.0/5, balancesBarGraph),
			ui.NewRow(1.0/10, plannedIncomeGauge),
			ui.NewRow(1.0/10, actualIncomeGauge),
			ui.NewRow(2.0/5, incomeBreakdown),
		),
	)
	ui.Render(grid)

	uiEvents := ui.PollEvents()

	previousKey := ""

	for {
		e := <-uiEvents
		switch e.ID {
		case "q", "<C-c>":
			return
		case "<Resize>":
			payload := e.Payload.(ui.Resize)
			grid.SetRect(0, 0, payload.Width, payload.Height)
			ui.Clear()
			ui.Render(grid)
		case "j", "<Down>":
			list.ScrollDown()
		case "k", "<Up>":
			list.ScrollUp()
		case "<C-d>":
			list.ScrollHalfPageDown()
		case "<C-u>":
			list.ScrollHalfPageUp()
		case "<C-f>":
			list.ScrollPageDown()
		case "<C-b>":
			list.ScrollPageUp()
		case "g":
			if previousKey == "g" {
				list.ScrollTop()
			}
		case "<Home>":
			list.ScrollTop()
		case "G", "<End>":
			list.ScrollBottom()
		}

		if previousKey == "g" {
			previousKey = ""
		} else {
			previousKey = e.ID
		}
		ui.Render(list)
	}
}

func setupList() *widgets.List {
	list := widgets.NewList()
	list.Title = "Monthly Budgets"
	list.Rows = []string{
		"Apr 2021",
		"Mar 2021",
		"Feb 2021",
		"Jan 2021",
	}

	list.TextStyle = ui.NewStyle(ui.ColorYellow)
	list.WrapText = false

	return list
}

func setupBarChart() *widgets.BarChart {
	bc := widgets.NewBarChart()
	bc.Data = []float64{1584, 3837}
	bc.Labels = []string{"Starting Balance", "Closing Balance"}
	bc.Title = "Balance Diff."
	bc.BarWidth = 10
	bc.BarColors = []ui.Color{ui.ColorWhite, ui.ColorWhite}
	bc.LabelStyles = []ui.Style{ui.NewStyle(ui.ColorWhite)}
	bc.NumStyles = []ui.Style{ui.NewStyle(ui.ColorBlack)}
	bc.NumFormatter = func(f float64) string {

		return fmt.Sprintf("$%.0f", f)
	}
	bc.BarGap = 10
	bc.PaddingBottom = 1
	bc.PaddingLeft = 10
	bc.PaddingRight = 10
	bc.PaddingTop = 1
	return bc
}

func setupExpenseGauges() (*widgets.Gauge, *widgets.Gauge) {
	plannedAmount := 1415.25
	actualAmount := 2979.45

	max := math.Max(plannedAmount, actualAmount)

	planned := widgets.NewGauge()
	planned.Title = "Planned Expenses"
	planned.Percent = int((plannedAmount / max) * 100)
	planned.Label = fmt.Sprintf("$%.2f", plannedAmount)

	actual := widgets.NewGauge()
	actual.Title = "Actual Expenses"
	actual.Percent = int((actualAmount / max) * 100)
	actual.Label = fmt.Sprintf("$%.2f", actualAmount)

	return planned, actual
}

func setupIncomeGauges() (*widgets.Gauge, *widgets.Gauge) {
	plannedAmount := 3113.56
	actualAmount := 5232.73

	max := math.Max(plannedAmount, actualAmount)

	planned := widgets.NewGauge()
	planned.Title = "Planned Income"
	planned.Percent = int((plannedAmount / max) * 100)
	planned.Label = fmt.Sprintf("$%.2f", plannedAmount)

	actual := widgets.NewGauge()
	actual.Title = "Actual Income"
	actual.Percent = int((actualAmount / max) * 100)
	actual.Label = fmt.Sprintf("$%.2f", actualAmount)

	return planned, actual
}

func setupOverview() *widgets.Table {
	table := widgets.NewTable()
	table.Title = "Overview for this Month"
	table.Rows = [][]string{
		{"Increase in Savings", "+142%"},
		{"Planned Savings", "$1,698.21"},
		{"Saved this Month", "$2,253.28"},
	}
	table.TextStyle = ui.NewStyle(ui.ColorWhite)
	table.PaddingTop = 1
	table.PaddingBottom = 1
	table.PaddingLeft = 1
	table.PaddingRight = 0

	return table
}

func setupExpensesTable() *widgets.Table {
	table := widgets.NewTable()
	table.Title = "Expenses Breakdown"
	table.Rows = [][]string{
		{"Category", "Planned", "Actual", "Diff"},
		{"Groceries", "$240.00", "$313.35", "-$73.35"},
		{"Rent", "$760.00", "$760.00", "$0.00"},
		{"Transportation", "$31.00", "54.20", "-$23.20"},
		{"Utilities", "$71.25", "$71.25", "$0.00"},
		{"Social", "$50.00", "$29.00", "$21.00"},
		{"Fitness", "$28.00", "$27.96", "$0.04"},
		{"Dining out", "$50.00", "$22.55", "$27.45"},
	}
	table.TextStyle = ui.NewStyle(ui.ColorWhite)
	// table.RowSeparator = true
	table.FillRow = true
	table.RowStyles[0] = ui.NewStyle(ui.ColorWhite, ui.ColorBlack, ui.ModifierBold)
	table.RowStyles[1] = ui.NewStyle(ui.ColorRed, ui.ColorBlack)
	table.RowStyles[2] = ui.NewStyle(ui.ColorWhite, ui.ColorBlack)
	table.RowStyles[3] = ui.NewStyle(ui.ColorRed, ui.ColorBlack)
	table.RowStyles[4] = ui.NewStyle(ui.ColorWhite, ui.ColorBlack)
	table.RowStyles[5] = ui.NewStyle(ui.ColorGreen, ui.ColorBlack)
	table.RowStyles[6] = ui.NewStyle(ui.ColorGreen, ui.ColorBlack)
	table.RowStyles[7] = ui.NewStyle(ui.ColorGreen, ui.ColorBlack)
	return table
}

func setupIncomeTable() *widgets.Table {
	table := widgets.NewTable()
	table.Title = "Income Breakdown"
	table.Rows = [][]string{
		{"Category", "Planned", "Actual", "Diff"},
		{"Salary", "$4000.00", "$4001.35", "$1.35"},
		{"Bonus", "$0.00", "$400.00", "$400.00"},
	}
	table.TextStyle = ui.NewStyle(ui.ColorWhite)
	table.RowSeparator = true
	table.FillRow = true
	table.RowStyles[0] = ui.NewStyle(ui.ColorWhite, ui.ColorBlack, ui.ModifierBold)
	table.RowStyles[1] = ui.NewStyle(ui.ColorGreen, ui.ColorBlack)
	table.RowStyles[2] = ui.NewStyle(ui.ColorGreen, ui.ColorBlack)
	return table
}
