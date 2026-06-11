package gui

import (
	"context"
	"fmt"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/runtime"

	"typemerge/frontend"
	core "typemerge/internal/app"
)

type App struct {
	ctx     context.Context
	service *core.Service
}

func NewApp(service *core.Service) *App {
	return &App{service: service}
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

func (a *App) DefaultState() State {
	return NewState()
}

func (a *App) Generate(state State) (core.Result, error) {
	if a.ctx == nil {
		return core.Result{}, fmt.Errorf("application is not ready")
	}
	return a.service.Generate(a.ctx, state.Options())
}

func (a *App) SelectFile(kind string) (string, error) {
	dialogOptions := runtime.OpenDialogOptions{Title: "Select file"}
	switch kind {
	case "template":
		dialogOptions.Title = "Select Typst template"
		dialogOptions.Filters = []runtime.FileFilter{
			{DisplayName: "Typst templates (*.typ)", Pattern: "*.typ"},
		}
	case "csv":
		dialogOptions.Title = "Select CSV data"
		dialogOptions.Filters = []runtime.FileFilter{
			{DisplayName: "CSV files (*.csv)", Pattern: "*.csv"},
		}
	case "metadata":
		dialogOptions.Title = "Select metadata"
		dialogOptions.Filters = []runtime.FileFilter{
			{DisplayName: "Text files (*.txt)", Pattern: "*.txt"},
			{DisplayName: "All files", Pattern: "*"},
		}
	default:
		return "", fmt.Errorf("unknown file kind %q", kind)
	}
	return runtime.OpenFileDialog(a.ctx, dialogOptions)
}

func (a *App) SelectOutputDirectory() (string, error) {
	return runtime.OpenDirectoryDialog(a.ctx, runtime.OpenDialogOptions{
		Title:                "Select output directory",
		CanCreateDirectories: true,
	})
}

func Run(service *core.Service) error {
	application := NewApp(service)
	return wails.Run(&options.App{
		Title:                    "typemerge",
		Width:                    960,
		Height:                   680,
		MinWidth:                 720,
		MinHeight:                560,
		BackgroundColour:         &options.RGBA{R: 244, G: 247, B: 251, A: 255},
		OnStartup:                application.startup,
		EnableDefaultContextMenu: true,
		AssetServer: &assetserver.Options{
			Assets: frontend.Assets,
		},
		Bind: []interface{}{
			application,
		},
	})
}
