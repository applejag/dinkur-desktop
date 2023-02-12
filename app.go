package main

import (
	"context"
	"fmt"

	"fyne.io/systray"
	"github.com/dinkur/dinkur/pkg/dinkur"
	"github.com/dinkur/dinkur/pkg/dinkurclient"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// App struct
type App struct {
	ctx    context.Context
	dinkur dinkur.Client

	trayCheckOut *systray.MenuItem
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// onStartup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) onStartup(ctx context.Context) {
	a.ctx = ctx
	systray.Run(a.onSystrayReady, a.onSystrayExit)
}

func (a *App) onShutdown(ctx context.Context) {
	systray.Quit()
	if a.dinkur != nil {
		a.dinkur.Close()
		a.dinkur = nil
	}
}

func (a *App) onSystrayReady() {
	systray.SetTemplateIcon(iconBytes, iconBytes)
	systray.SetTitle("Dinkur desktop")
	systray.SetTooltip("Placeholder tooltip")

	a.trayCheckOut = systray.AddMenuItem("No active entry", "You have no active entry tracking time right now.")
	a.trayCheckOut.Disable()
	menuShow := systray.AddMenuItem("Show Dinkur", "Opens Dinkur when it has been hidden/closed.")
	go func() {
		for range menuShow.ClickedCh {
			runtime.Show(a.ctx)
		}
	}()
	systray.AddSeparator()
	menuQuit := systray.AddMenuItem("Quit Dinkur", "Exits Dinkur desktop")
	go func() {
		if _, ok := <-menuQuit.ClickedCh; ok {
			systray.Quit()
		}
	}()
}

func (a *App) onSystrayExit() {
	runtime.Quit(a.ctx)
}

// // ConnectDinkur tries to connect to a Dinkur daemon over gRPC
func (a *App) ConnectDinkur(serverAddr string) error {
	if a.dinkur != nil {
		a.dinkur.Close()
	}
	a.dinkur = dinkurclient.NewClient(serverAddr, dinkurclient.Options{})
	fmt.Println("Connecting to:", serverAddr)
	if err := a.dinkur.Connect(a.ctx); err != nil {
		a.dinkur = nil
		return err
	}
	if err := a.dinkur.Ping(a.ctx); err != nil {
		a.dinkur.Close()
		a.dinkur = nil
		return fmt.Errorf("ping: %w", err)
	}
	fmt.Println("Successfully connected:", serverAddr)
	return nil
}

func (a *App) GetActiveEntry() (*dinkur.Entry, error) {
	return a.dinkur.GetActiveEntry(a.ctx)
}
