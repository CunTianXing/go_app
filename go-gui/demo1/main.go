package main

import (
	"fmt"

	"github.com/murlokswarm/app"
)

var win app.Windower

func main() {
	//icon := filepath.Join(app.Resources(), "like.png")
	dock, _ := app.Dock()
	//dock.SetBadge("Hello")
	//dock.SetIcon(icon)
	fmt.Println(dock)
	app.OnReopen = func() {
		if win != nil {
			return
		}
		win = newMainWindow()
		win.Mount(&Home{})
	}
	app.Run()

}

func newMainWindow() app.Windower {
	return app.NewWindow(app.Window{
		Title:           "test",
		TitlebarHidden:  true,
		Width:           1280,
		Height:          768,
		BackgroundColor: "#21252b",
		OnClose: func() bool {
			win = nil
			return true
		},
	})
}
