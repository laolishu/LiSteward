package main

import (
	"embed"
	"fmt"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"

	"LiSteward/config"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	// Load application config (optional)
	cfg, err := config.Load()
	if err != nil {
		fmt.Printf("加载配置失败: %v, 使用内置默认值\n", err)
	}

	// Create an instance of the app structure
	app := NewApp()

	// Create application with options
	appTitle := "LiSteward"
	width := 1280
	height := 800
	minw := 800
	minh := 600
	if cfg != nil {
		if cfg.ProductName != "" {
			appTitle = cfg.ProductName
		}
		width = cfg.DefaultWindow.Width
		height = cfg.DefaultWindow.Height
		minw = cfg.DefaultWindow.MinWidth
		minh = cfg.DefaultWindow.MinHeight
	}

	err = wails.Run(&options.App{
		Title:     appTitle,
		Width:     width,
		Height:    height,
		MinWidth:  minw,
		MinHeight: minh,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1},
		OnStartup:        app.startup,
		Bind: []interface{}{
			app,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
