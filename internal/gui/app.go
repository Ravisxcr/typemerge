package gui

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

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
	dialogOptions := runtime.OpenDialogOptions{
		Title:            "Select file",
		DefaultDirectory: defaultDialogDirectory(),
	}
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
		DefaultDirectory:     defaultDialogDirectory(),
		CanCreateDirectories: true,
	})
}

func (a *App) SelectProjectDirectory() (ProjectSetup, error) {
	directory, err := runtime.OpenDirectoryDialog(a.ctx, runtime.OpenDialogOptions{
		Title:                "Open project folder",
		DefaultDirectory:     defaultDialogDirectory(),
		CanCreateDirectories: true,
	})
	if err != nil || directory == "" {
		return ProjectSetup{}, err
	}
	return discoverProject(directory)
}

func discoverProject(directory string) (ProjectSetup, error) {
	setup := ProjectSetup{
		Directory: directory,
		Name:      filepath.Base(directory),
		OutputDir: filepath.Join(directory, "output"),
	}
	entries, err := os.ReadDir(directory)
	if err != nil {
		return ProjectSetup{}, fmt.Errorf("read project folder: %w", err)
	}
	sort.Slice(entries, func(i, j int) bool {
		return strings.ToLower(entries[i].Name()) < strings.ToLower(entries[j].Name())
	})
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		path := filepath.Join(directory, entry.Name())
		name := strings.ToLower(entry.Name())
		switch {
		case setup.CSVPath == "" && filepath.Ext(name) == ".csv":
			setup.CSVPath = path
		case setup.MetadataPath == "" && (name == "metadata.txt" || name == "meta.txt"):
			setup.MetadataPath = path
		case setup.TemplatePath == "" && filepath.Ext(name) == ".typ":
			setup.TemplatePath = path
		}
	}
	return setup, nil
}

func (a *App) UseTemplatePreset(directory, preset string) (string, error) {
	if directory == "" {
		return "", fmt.Errorf("choose a project folder before using a template preset")
	}
	var contents string
	switch preset {
	case "salary":
		contents = salarySlipTemplate
	case "invoice":
		contents = invoiceTemplate
	default:
		return "", fmt.Errorf("unknown template preset %q", preset)
	}

	templateDir := filepath.Join(directory, ".typemerge", "templates")
	if err := os.MkdirAll(templateDir, 0o755); err != nil {
		return "", fmt.Errorf("create template directory: %w", err)
	}
	path := filepath.Join(templateDir, preset+".typ")
	if err := os.WriteFile(path, []byte(contents), 0o644); err != nil {
		return "", fmt.Errorf("write template preset: %w", err)
	}
	return path, nil
}

const salarySlipTemplate = `#set page(paper: "a4", margin: 22mm)
#set text(font: "Libertinus Serif", size: 10pt)
#let employee = "{{ get .Employee "name" }}"
#align(center)[
  #text(size: 19pt, weight: "bold")[{{ get .Meta "company_name" }}]
  #v(4pt)
  #text(fill: gray)[Salary slip for {{ get .Meta "month" }} {{ get .Meta "year" }}]
]
#v(20pt)
#grid(
  columns: (1fr, 1fr),
  gutter: 12pt,
  [*Employee*], [#employee],
  [*Employee ID*], [{{ get .Employee "employee_id" }}],
  [*Payment date*], [{{ get .Meta "payment_date" }}],
  [*Net salary*], [{{ money (get .Employee "net_salary") }}],
)
`

const invoiceTemplate = `#set page(paper: "a4", margin: 22mm)
#set text(font: "Libertinus Serif", size: 10pt)
#text(size: 22pt, weight: "bold")[INVOICE]
#v(8pt)
{{ get .Meta "company_name" }} \
{{ get .Meta "address" }}
#v(22pt)
#grid(
  columns: (1fr, 1fr),
  gutter: 12pt,
  [*Invoice number*], [{{ get .Employee "invoice_number" }}],
  [*Customer*], [{{ get .Employee "name" }}],
  [*Date*], [{{ get .Employee "date" }}],
  [*Amount due*], [{{ money (get .Employee "amount") }}],
)
`

func Run(service *core.Service) error {
	application := NewApp(service)
	return wails.Run(&options.App{
		Title:                    "typemerge",
		Width:                    1280,
		Height:                   800,
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
