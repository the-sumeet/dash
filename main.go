package main

import (
	"embed"
	"log"

	"github.com/wailsapp/wails/v3/pkg/application"
	"github.com/wailsapp/wails/v3/pkg/events"
)

//go:embed all:frontend/dist
var assets embed.FS

//go:embed build/appicon.png
var icon []byte

func main() {
	appService := NewAppService()

	app := application.New(application.Options{
		Name:        "tray",
		Description: "A system tray app",
		Services: []application.Service{
			application.NewService(appService),
		},
		Assets: application.AssetOptions{
			Handler: application.AssetFileServerFS(assets),
		},
		Mac: application.MacOptions{
			ActivationPolicy: application.ActivationPolicyAccessory,
			ApplicationShouldTerminateAfterLastWindowClosed: false,
		},
	})

	window := app.Window.NewWithOptions(application.WebviewWindowOptions{
		Width:         512,
		Height:        512,
		Frameless:     true,
		AlwaysOnTop:   true,
		Hidden:        true,
		DisableResize: true,
		URL:           "/",
		Mac: application.MacWindow{
			WindowLevel: application.MacWindowLevelPopUpMenu,
			CollectionBehavior: application.MacWindowCollectionBehaviorCanJoinAllSpaces |
				application.MacWindowCollectionBehaviorFullScreenAuxiliary,
		},
	})

	window.RegisterHook(events.Common.WindowClosing, func(e *application.WindowEvent) {
		window.Hide()
		e.Cancel()
	})

	// Settings window (full-size, opened from right-click menu)
	settingsWindow := app.Window.NewWithOptions(application.WebviewWindowOptions{
		Name:   "Settings",
		Title:  "Tray Settings",
		Width:  700,
		Height: 500,
		Hidden: true,
		URL:    "/settings",
	})
	settingsWindow.RegisterHook(events.Common.WindowClosing, func(e *application.WindowEvent) {
		settingsWindow.Hide()
		e.Cancel()
	})

	// Add app window
	addAppWIndow := app.Window.NewWithOptions(application.WebviewWindowOptions{
		Name:   "Applications",
		Title:  "Application",
		Width:  700,
		Height: 500,
		Hidden: true,
		URL:    "/apps",
	})
	addAppWIndow.RegisterHook(events.Common.WindowClosing, func(e *application.WindowEvent) {
		addAppWIndow.Hide()
		e.Cancel()
	})

	// Right-click context menu
	menu := app.NewMenu()
	menu.Add("Settings...").OnClick(func(ctx *application.Context) {
		settingsWindow.Show()
		settingsWindow.Focus()
		settingsWindow.Center()
	})
	menu.Add("Apps...").OnClick(func(ctx *application.Context) {
		addAppWIndow.Show()
		addAppWIndow.Focus()
		addAppWIndow.Center()
	})
	menu.AddSeparator()
	menu.Add("Quit").OnClick(func(ctx *application.Context) {
		app.Quit()
	})

	systemTray := app.SystemTray.New()
	systemTray.SetIcon(icon)
	systemTray.SetMenu(menu)
	systemTray.AttachWindow(window).WindowOffset(5)

	err := app.Run()
	if err != nil {
		log.Fatal(err)
	}
}
