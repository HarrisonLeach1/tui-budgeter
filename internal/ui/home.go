package ui

import (
	"github.com/VladimirMarkelov/clui"
	"github.com/nsf/termbox-go"
)

func RenderHomePage() {
	view := clui.AddWindow(0, 0, 50, 50, "Home Page")
	// view.OnScreenResize(func(ev clui.Event) {
	// 	view.SetSize(ev.Width, ev.Height)
	// 	view.ResizeChildren()
	// })
	view.SetPack(clui.Vertical)

	viewWidth, _ := view.Size()

	banner := clui.CreateFrame(view, viewWidth*2, 40, clui.BorderThin, clui.AutoSize)
	banner.SetPack(clui.Vertical)
	banner.SetGaps(clui.KeepValue, 1)
	banner.SetPaddings(1, 1)

	text := []string{"Welcome to the budgeting tool!", " ", "This tool uses the Xero API. For instructions on how to get it up and running view the setup guide."}
	desc := clui.CreateTextView(banner, clui.AutoSize, clui.AutoSize, clui.AutoSize)
	desc.AddText(text)
	desc.SetBackColor(clui.ColorBlack)
	desc.SetActiveBackColor(clui.ColorBlack)
	desc.SetActive(false)
	desc.SetTextColor(clui.ColorWhite)
	desc.SetActiveTextColor(clui.ColorWhite)
	desc.SetWordWrap(true)
	desc.SetAlign(clui.AlignCenter)
	desc.SetGaps(clui.KeepValue, 1)
	desc.SetClipped(true)

	menuFrame := clui.CreateFrame(view, viewWidth*2, 10, clui.BorderThin, clui.AutoSize)
	menuFrame.SetPack(clui.Vertical)
	menuFrame.SetGaps(clui.KeepValue, 1)
	menuFrame.SetPaddings(1, 1)

	menu := clui.CreateListBox(menuFrame, clui.AutoSize, clui.AutoSize, clui.AutoSize)
	menu.AddItem("View Monthly Budgets")
	menu.AddItem("Setup Guide")
	menu.SetTextColor(clui.ColorWhite)
	menu.SetActiveTextColor(clui.ColorWhiteBold)
	menu.SetBackColor(clui.ColorBlack)
	menu.SetActiveBackColor(clui.ColorBlack)

	menu.OnKeyPress(func(key termbox.Key) bool {
		if key == termbox.KeyEnter {
			if menu.SelectedItem() == 0 {
				SelectReportPeriod()
			} else if menu.SelectedItem() == 1 {
				// Render help page
			}
			return true
		}
		return false
	})

}
