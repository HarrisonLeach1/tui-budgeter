package ui

import (
	"embed"
	"fmt"
	"strings"

	"github.com/VladimirMarkelov/clui"
	"github.com/nsf/termbox-go"
)

var FS embed.FS

func RenderHomePage() {
	view := clui.AddWindow(0, 0, 60, 60, "Home Page")
	// view.OnScreenResize(func(ev clui.Event) {
	// 	view.SetSize(ev.Width, ev.Height)
	// 	view.ResizeChildren()
	// })
	view.SetPack(clui.Vertical)

	viewWidth, _ := view.Size()

	bannerFrame := clui.CreateFrame(view, viewWidth*2, 10, clui.BorderThin, clui.AutoSize)
	bannerFrame.SetPack(clui.Vertical)
	bannerFrame.SetGaps(clui.KeepValue, 1)
	bannerFrame.SetPaddings(1, 1)

	bt := []string{`     __        _       __              __           __           `,
		`    / /___  __(_)     / /_  __  ______/ /___ ____  / /____  _____`,
		`   / __/ / / / /_____/ __ \/ / / / __  / __ \/ _ \/ __/ _ \/ ___/`,
		`  / /_/ /_/ / /_____/ /_/ / /_/ / /_/ / /_/ /  __/ /_/  __/ /    `,
		`  \__/\__,_/_/     /_.___/\__,_/\__,_/\__, /\___/\__/\___/_/     `,
		`                                     /____/							`}

	banner := clui.CreateTextView(bannerFrame, clui.AutoSize, 10, clui.AutoSize)
	banner.AddText(bt)
	banner.SetBackColor(clui.ColorBlack)
	banner.SetTextColor(clui.ColorCyanBold)
	banner.SetGaps(clui.KeepValue, 1)
	banner.SetAlign(clui.AlignCenter)

	descFrame := clui.CreateFrame(view, viewWidth*2, 40, clui.BorderThin, clui.AutoSize)
	descFrame.SetPack(clui.Vertical)
	descFrame.SetGaps(clui.KeepValue, 1)
	descFrame.SetPaddings(1, 1)
	descFrame.SetTitle("Welcome to xui-budgeter!")

	ar := readFileIntoStringArray("README.md")

	desc := clui.CreateTextView(descFrame, clui.AutoSize, 40, clui.AutoSize)
	desc.AddText(ar)
	desc.SetBackColor(clui.ColorBlack)
	desc.SetActiveBackColor(clui.ColorBlack)
	desc.SetActive(false)
	desc.SetTextColor(clui.ColorWhite)
	desc.SetActiveTextColor(clui.ColorWhite)
	desc.SetWordWrap(true)
	desc.SetAlign(clui.AlignCenter)
	desc.SetGaps(clui.KeepValue, 1)

	menuFrame := clui.CreateFrame(view, viewWidth*2, 10, clui.BorderThin, clui.AutoSize)
	menuFrame.SetPack(clui.Vertical)
	menuFrame.SetGaps(clui.KeepValue, 1)
	menuFrame.SetPaddings(1, 1)
	menuFrame.SetTitle("Menu")

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
				RenderSetupPage()
			}
			return true
		}
		return false
	})

}

func readFileIntoStringArray(filePath string) []string {

	b, err := FS.ReadFile(filePath)
	if err != nil {
		fmt.Print(err)
	}

	str := string(b)
	ar := strings.Split(str, "\n")
	length := len(ar) * 2
	out := make([]string, length)

	for i := 0; i < len(out); i++ {
		if i%2 == 0 {
			out[i] = ar[i/2]
		} else {

			out[i] = " "
		}
	}
	return append([]string{" "}, out...)
}
