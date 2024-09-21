package main

import (
	"context"
	"fmt"

	"DiyHanconDemo/demo"
	"DiyHanconDemo/device"
)

// App struct
type App struct {
	ctx  context.Context
	demo *demo.Demo
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called at application startup
func (a *App) startup(ctx context.Context) {
	// Perform your setup here
	a.ctx = ctx
	// Initialize the device
	dev := device.New(ctx)
	// Initialize the demo
	app := demo.New(dev)
	// Return the app
	a.demo = app
	go app.Run(ctx)
}

// domReady is called after the front-end dom has been loaded
func (a App) domReady(ctx context.Context) {
	// Add your action here
}

// shutdown is called at application termination
func (a *App) shutdown(ctx context.Context) {
	// Perform your teardown here
}

// Greet returns a greeting for the given name
func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s!", name)
}
