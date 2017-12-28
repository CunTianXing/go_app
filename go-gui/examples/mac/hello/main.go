package main

import (
	"github.com/murlokswarm/app"
	_ "github.com/murlokswarm/mac"
)

var (
	win app.Contexter
)

func main() {
	app.OnLaunch = func() {
		appMenu := &AppMainMenu{}    // Creates the AppMainMenu component.
		if menuBar, ok := app.MenuBar(); ok { // Mounts the AppMainMenu component into the application menu bar.
			menuBar.Mount(appMenu)
		}

		appMenuDock := &AppMainMenu{} // Creates another AppMainMenu.
		if dock, ok := app.Dock(); ok {
			dock.Mount(appMenuDock)
		}

		win = newMainWindow() // Create the main window.
	}

	app.OnReopen = func() {
		if win != nil {
			return
		}

		win = newMainWindow()
	}

	app.Run()
}

func newMainWindow() app.Contexter {
	// Creates a window context.
	win := app.NewWindow(app.Window{
		Title:          "Hello World",
		Width:          1280,
		Height:         720,
		TitlebarHidden: true,
		OnClose: func() bool {
			win = nil
			return true
		},
	})

	hello := &Hello{} // Creates a Hello component.
	win.Mount(hello)  // Mounts the Hello component into the window context.
	return win
}
