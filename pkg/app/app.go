package app

import (
	"context"
	"embed"
	"fmt"
	"time"

	"fyne.io/systray"
	"github.com/dinkur/dinkur-desktop/internal/wailsutil"
	"github.com/dinkur/dinkur-desktop/pkg/config"
	"github.com/dinkur/dinkur/pkg/dinkur"
	"github.com/dinkur/dinkur/pkg/dinkurdb"
	"github.com/dinkur/dinkur/pkg/timeutil"
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

var log = logger.NewScoped("Dinkur desktop")

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
	opt := dinkurdb.Options{
		MkdirAll: cfg.Sqlite.Mkdir,
	}
	return &App{
		cfg:    cfg,
		dinkur: dinkurdb.NewClient(cfg.Sqlite.Path, opt),
	}
}

// onStartup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) onStartup(ctx context.Context) {
	a.ctx = ctx
	go systray.Run(a.onSystrayReady, a.onSystrayExit)
	a.ConnectDinkur()
}

func (a *App) onShutdown(ctx context.Context) {
	if err := a.cfg.Save(); err != nil {
		log.Error().WithError(err).Message("Failed to save config before exiting.")
	}
	systray.Quit()
	a.DisconnectDinkur()
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

func (a *App) ConnectDinkur() error {
	if err := a.dinkur.Connect(a.ctx); err != nil {
		log.Error().WithError(err).Message("Failed to connect to Dinkur database.")
		return err
	}
	if err := a.dinkur.Ping(a.ctx); err != nil {
		log.Error().WithError(err).Message("Failed to ping Dinkur database.")
		a.dinkur.Close()
		return fmt.Errorf("ping: %w", err)
	}
	log.Info().Message("Successfully connected to Dinkur!")
	return nil
}

func (a *App) DisconnectDinkur() error {
	err := a.dinkur.Close()
	if err != nil {
		log.Error().WithError(err).Message("Failed to close connection to Dinkur database.")
	}
	return err
}

func (a *App) GetActiveEntry() (*dinkur.Entry, error) {
	return a.dinkur.GetActiveEntry(a.ctx)
}

func (a *App) GetEntriesForDay(day time.Time) ([]dinkur.Entry, error) {
	span := timeutil.Day(day)
	log.Debug().WithString("day", day.Format(time.DateOnly)).
		WithTime("start", *span.Start).
		WithTime("end", *span.End).
		Message("Getting entries for day.")
	entries, err := a.dinkur.GetEntryList(context.Background(), dinkur.SearchEntry{
		Limit: 100,
		Start: span.Start,
		End:   span.End,
	})
	log.Debug().WithError(err).
		WithInt("count", len(entries)).
		Message("Got entries response.")
	return entries, err
}
