package ui

import "github.com/VladimirMarkelov/clui"

func RenderSetupPage() {
	view := clui.AddWindow(0, 0, 60, 60, "Setup Guide")
	view.SetPack(clui.Vertical)

	viewWidth, _ := view.Size()

	textFrame := clui.CreateFrame(view, viewWidth*2, 60, clui.BorderThin, clui.AutoSize)
	textFrame.SetPack(clui.Vertical)
	textFrame.SetGaps(clui.KeepValue, 1)
	textFrame.SetPaddings(1, 1)

	ar := readFileIntoStringArray("docs/setup.md")

	text := clui.CreateTextView(textFrame, clui.AutoSize, 60, clui.AutoSize)
	text.AddText(ar)
	text.SetBackColor(clui.ColorBlack)
	text.SetActiveBackColor(clui.ColorBlack)
	text.SetActive(false)
	text.SetTextColor(clui.ColorWhite)
	text.SetActiveTextColor(clui.ColorWhite)
	text.SetWordWrap(true)
	text.SetAlign(clui.AlignCenter)
	text.SetGaps(clui.KeepValue, 1)

}
