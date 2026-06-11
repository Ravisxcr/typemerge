# GUI

The GUI uses React, TypeScript, Bun, Tailwind CSS, and shadcn/ui conventions.
Install [Bun](https://bun.sh/) before running the frontend tasks.

```bash
task frontend:install
```

Run the desktop frontend in development mode:

```bash
task dev-gui
```

Build a production application with:

```bash
task build-gui
```

The Taskfile invokes the pinned Wails v2.12.0 CLI through Go, so a global
Wails installation is not required. To add another shadcn component:

```bash
cd frontend
bunx --bun shadcn@latest add dialog
```

Wails uses WebView2 on Windows. Linux builds require GTK3 and WebKitGTK
development packages.

The GUI provides file and output-directory pickers, compiled-format controls,
generation status, and an output summary. It is a thin frontend over
`internal/app.Service`; parsing, rendering, naming, and Typst compilation remain
in the shared core.
