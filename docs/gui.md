# GUI

Run the Fyne frontend with:

```bash
go run ./cmd/typemerge-gui
```

On Linux, native builds require the desktop development packages used by Fyne,
including OpenGL, X11, and Xxf86vm headers/libraries. Headless CI uses Fyne's
`ci` build tag and does not require a display server.

The GUI provides file and output-directory pickers, compiled-format controls,
generation status, and an output summary. It is a thin frontend over
`internal/app.Service`; parsing, rendering, naming, and Typst compilation remain
in the shared core.
