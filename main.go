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
		Name:        "Dash",
		Description: "A menubar dashboard",
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

	// Add app window
	addAppWindow := app.Window.NewWithOptions(application.WebviewWindowOptions{
		Name:   "Applications",
		Title:  "Application",
		Width:  700,
		Height: 500,
		Hidden: true,
		URL:    "/apps",
	})
	addAppWindow.RegisterHook(events.Common.WindowClosing, func(e *application.WindowEvent) {
		addAppWindow.Hide()
		e.Cancel()
	})

	// Right-click context menu
	menu := app.NewMenu()
	menu.Add("Apps...").OnClick(func(ctx *application.Context) {
		addAppWindow.Show()
		addAppWindow.Focus()
		addAppWindow.Center()
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
