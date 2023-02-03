package main

import (
	"context"
	"fmt"

	"github.com/dinkur/dinkur/pkg/dinkur"
	"github.com/dinkur/dinkur/pkg/dinkurclient"
)

// App struct
type App struct {
	ctx    context.Context
	dinkur dinkur.Client
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

func (a *App) shutdown(ctx context.Context) {
	if a.dinkur != nil {
		a.dinkur.Close()
		a.dinkur = nil
	}
}

// Greet returns a greeting for the given name
func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s, It's crazy time!", name)
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
