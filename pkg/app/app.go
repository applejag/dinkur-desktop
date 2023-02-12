package app

import (
	"context"
	"embed"
	"fmt"

	"fyne.io/systray"
	"github.com/dinkur/dinkur-desktop/internal/wailsutil"
	"github.com/dinkur/dinkur-desktop/pkg/config"
	"github.com/dinkur/dinkur/pkg/dinkur"
	"github.com/dinkur/dinkur/pkg/dinkurclient"
	"github.com/iver-wharf/wharf-core/v2/pkg/logger"
	"github.com/wailsapp/wails/v2"
	wailslogger "github.com/wailsapp/wails/v2/pkg/logger"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/options/linux"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

var Assets embed.FS
var IconBytes []byte

var log = logger.NewScoped("app")

func Run(cfg *config.Config) error {
	app := New(cfg)

	// Create application with options
	return wails.Run(&options.App{
		Title:  "Dinkur desktop",
		Width:  480,
		Height: 640,
		AssetServer: &assetserver.Options{
			Assets: Assets,
		},
		BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1},
		OnStartup:        app.onStartup,
		OnShutdown:       app.onShutdown,
		Bind: []interface{}{
			app,
		},
		Linux: &linux.Options{
			Icon: IconBytes,
		},
		HideWindowOnClose:  !cfg.ExitOnWindowClose,
		Logger:             wailsutil.Logger{WharfLogger: logger.NewScoped("Wails")},
		LogLevel:           wailslogger.DEBUG,
		LogLevelProduction: wailslogger.INFO,
	})
}

// App struct
type App struct {
	cfg    *config.Config
	ctx    context.Context
	dinkur dinkur.Client

	trayCheckOut *systray.MenuItem
}

// New creates a new App application struct
func New(cfg *config.Config) *App {
	return &App{cfg: cfg}
}

// onStartup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) onStartup(ctx context.Context) {
	a.ctx = ctx
	systray.Run(a.onSystrayReady, a.onSystrayExit)
}

func (a *App) onShutdown(ctx context.Context) {
	if err := a.cfg.Save(); err != nil {
		log.Error().WithError(err).Message("Failed to save config before exiting.")
	}
	systray.Quit()
	if a.dinkur != nil {
		if err := a.dinkur.Close(); err != nil {
			log.Error().WithError(err).Message("Failed to close connection to Dinkur database.")
		}
		a.dinkur = nil
	}
}

func (a *App) onSystrayReady() {
	systray.SetTemplateIcon(IconBytes, IconBytes)
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
