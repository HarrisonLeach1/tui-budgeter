package main

import (
	"github.com/rivo/tview"
)

func main() {
	app := tview.NewApplication()
	expensesBreakdown := setupExpensesBreakdownTable()
	incomeBreakdown := setupIncomeBreakdownTable()

	flex := tview.NewFlex().
		AddItem(tview.NewBox().SetBorder(true).SetTitle("Monthly Budgets"), 20, 1, false).
		AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
			AddItem(tview.NewBox().SetBorder(true), 0, 2, false).
			AddItem(tview.NewBox().SetBorder(true).SetTitle("Expenses"), 0, 1, false).
			AddItem(expensesBreakdown, 0, 2, false), 0, 2, false).
		AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
			AddItem(tview.NewBox().SetBorder(true), 0, 2, false).
			AddItem(tview.NewBox().SetBorder(true).SetTitle("Income"), 0, 1, false).
			AddItem(incomeBreakdown, 0, 2, false), 0, 2, false)

	if err := app.SetRoot(flex, true).EnableMouse(true).Run(); err != nil {
		panic(err)
	}
}

func setupExpensesBreakdownTable() *tview.Box {
	return tview.NewBox().SetBorder(true).SetTitle("Expenses Breakdown")
}

func setupIncomeBreakdownTable() *tview.Box {
	return tview.NewBox().SetBorder(true).SetTitle("Income Breakdown")
}
