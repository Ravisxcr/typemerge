import { useEffect, useMemo, useState, type FormEvent } from "react"
import {
  CheckCircle2Icon,
  ChevronRightIcon,
  CircleIcon,
  FileCode2Icon,
  FileSpreadsheetIcon,
  FileTextIcon,
  FolderOpenIcon,
  Layers3Icon,
  LoaderCircleIcon,
  PlayIcon,
  ScrollTextIcon,
  Settings2Icon,
  SparklesIcon,
  TerminalSquareIcon,
} from "lucide-react"

import { Button } from "@/components/ui/button"
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card"
import { Checkbox } from "@/components/ui/checkbox"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"
import { cn } from "@/lib/utils"
import {
  DefaultState,
  Generate,
  SelectFile,
  SelectOutputDirectory,
} from "../wailsjs/go/gui/App"
import type { app, gui } from "../wailsjs/go/models"

const emptyState: gui.State = {
  templatePath: "",
  csvPath: "",
  metadataPath: "",
  outputDir: "",
  typstBinary: "typst",
  pdf: true,
  png: false,
  svg: false,
}

type PathField = "templatePath" | "csvPath" | "metadataPath" | "outputDir"
type FileKind = "template" | "csv" | "metadata"
type Format = "pdf" | "png" | "svg"

const fileFields: Array<{
  field: Exclude<PathField, "outputDir">
  kind: FileKind
  label: string
  description: string
  placeholder: string
  icon: typeof FileTextIcon
  optional?: boolean
}> = [
  {
    field: "templatePath",
    kind: "template",
    label: "Typst template",
    description: "The document layout used for every row",
    placeholder: "Choose a .typ template",
    icon: FileCode2Icon,
  },
  {
    field: "csvPath",
    kind: "csv",
    label: "CSV data",
    description: "Each row becomes a generated document",
    placeholder: "Choose a .csv data file",
    icon: FileSpreadsheetIcon,
  },
  {
    field: "metadataPath",
    kind: "metadata",
    label: "Shared metadata",
    description: "Key-value data available to every document",
    placeholder: "Choose a metadata file",
    icon: ScrollTextIcon,
    optional: true,
  },
]

const formats: Array<{
  id: Format
  label: string
  description: string
}> = [
  { id: "pdf", label: "PDF", description: "Ready to share" },
  { id: "png", label: "PNG", description: "Raster image" },
  { id: "svg", label: "SVG", description: "Vector graphic" },
]

function errorMessage(error: unknown) {
  return error instanceof Error ? error.message : String(error)
}

function fileName(path: string) {
  return path.split(/[\\/]/).pop() || path
}

function App() {
  const [form, setForm] = useState<gui.State>(emptyState)
  const [files, setFiles] = useState<app.GeneratedFile[]>([])
  const [log, setLog] = useState("Ready when you are.")
  const [loading, setLoading] = useState(true)
  const [generating, setGenerating] = useState(false)

  useEffect(() => {
    DefaultState()
      .then(setForm)
      .catch((error: unknown) => setLog(`Error: ${errorMessage(error)}`))
      .finally(() => setLoading(false))
  }, [])

  const selectedFormats = useMemo(
    () => formats.filter(({ id }) => form[id]),
    [form],
  )

  const requiredInputsReady = Boolean(
    form.templatePath && form.csvPath && form.outputDir,
  )

  const updatePath = (field: PathField, value: string) => {
    setForm((current) => ({ ...current, [field]: value }))
  }

  const selectFile = async (
    field: Exclude<PathField, "outputDir">,
    kind: FileKind,
  ) => {
    try {
      const path = await SelectFile(kind)
      if (path) updatePath(field, path)
    } catch (error) {
      setLog(`Error: ${errorMessage(error)}`)
    }
  }

  const selectOutput = async () => {
    try {
      const path = await SelectOutputDirectory()
      if (path) updatePath("outputDir", path)
    } catch (error) {
      setLog(`Error: ${errorMessage(error)}`)
    }
  }

  const generate = async (event: FormEvent<HTMLFormElement>) => {
    event.preventDefault()
    setGenerating(true)
    setLog("Generating documents...")

    try {
      const result = await Generate(form)
      const generated = result.files ?? []
      setFiles(generated)
      setLog(
        generated.length
          ? generated.map((file) => `Created  ${file.path}`).join("\n")
          : "No files were generated.",
      )
    } catch (error) {
      setFiles([])
      setLog(`Error: ${errorMessage(error)}`)
    } finally {
      setGenerating(false)
    }
  }

  return (
    <main className="app-shell">
      <aside className="sidebar">
        <div className="brand">
          <div className="brand-mark">
            <Layers3Icon className="size-5" />
          </div>
          <div>
            <p className="brand-name">typemerge</p>
            <p className="brand-version">Document studio</p>
          </div>
        </div>

        <nav className="sidebar-nav" aria-label="Workspace">
          <button className="nav-item nav-item-active" type="button">
            <SparklesIcon className="size-4" />
            Generate
          </button>
          <div className="nav-item nav-item-muted">
            <Settings2Icon className="size-4" />
            Configuration
          </div>
        </nav>

        <div className="sidebar-status">
          <div className="flex items-center gap-2">
            <span className="status-dot" />
            <span>Local workspace</span>
          </div>
          <p>Your files stay on this device.</p>
        </div>
      </aside>

      <div className="workspace">
        <header className="workspace-header">
          <div>
            <p className="eyebrow">Document generator</p>
            <h1>Create a new batch</h1>
            <p className="header-copy">
              Merge your data into polished, consistent Typst documents.
            </p>
          </div>
          <div className="header-badge">
            <CircleIcon className="size-2 fill-current" />
            Typst ready
          </div>
        </header>

        <form className="workspace-grid" onSubmit={generate}>
          <div className="main-column">
            <Card className="panel-card">
              <CardHeader className="section-heading">
                <div className="step-number">1</div>
                <div>
                  <CardTitle>Source files</CardTitle>
                  <CardDescription>
                    Select the template and data for this batch.
                  </CardDescription>
                </div>
              </CardHeader>
              <CardContent className="grid gap-3">
                {fileFields.map(
                  ({
                    field,
                    kind,
                    label,
                    description,
                    placeholder,
                    icon: Icon,
                    optional,
                  }) => (
                    <div className="file-field" key={field}>
                      <div className="file-icon">
                        <Icon className="size-5" />
                      </div>
                      <div className="min-w-0 flex-1">
                        <div className="mb-1 flex items-center gap-2">
                          <Label htmlFor={field}>{label}</Label>
                          {optional && (
                            <span className="optional-tag">Optional</span>
                          )}
                        </div>
                        <p className="field-description">{description}</p>
                        <Input
                          className="mt-2"
                          id={field}
                          value={form[field]}
                          placeholder={placeholder}
                          required={!optional}
                          onChange={(event) =>
                            updatePath(field, event.target.value)
                          }
                        />
                      </div>
                      <Button
                        className="browse-button"
                        type="button"
                        variant="outline"
                        size="sm"
                        onClick={() => selectFile(field, kind)}
                      >
                        Browse
                      </Button>
                    </div>
                  ),
                )}
              </CardContent>
            </Card>

            <Card className="panel-card">
              <CardHeader className="section-heading">
                <div className="step-number">2</div>
                <div>
                  <CardTitle>Output settings</CardTitle>
                  <CardDescription>
                    Choose where and how your documents are created.
                  </CardDescription>
                </div>
              </CardHeader>
              <CardContent className="grid gap-5">
                <div className="grid gap-2">
                  <Label htmlFor="outputDir">Destination folder</Label>
                  <div className="flex gap-2">
                    <Input
                      id="outputDir"
                      value={form.outputDir}
                      placeholder="Choose an output directory"
                      required
                      onChange={(event) =>
                        updatePath("outputDir", event.target.value)
                      }
                    />
                    <Button
                      type="button"
                      variant="outline"
                      onClick={selectOutput}
                    >
                      <FolderOpenIcon className="size-4" />
                      Browse
                    </Button>
                  </div>
                </div>

                <fieldset className="grid gap-2.5">
                  <legend className="mb-2 text-sm font-medium">
                    Export formats
                  </legend>
                  <div className="format-grid">
                    {formats.map(({ id, label, description }) => (
                      <Label
                        className={cn(
                          "format-option",
                          form[id] && "format-option-selected",
                        )}
                        htmlFor={id}
                        key={id}
                      >
                        <Checkbox
                          id={id}
                          checked={form[id]}
                          onCheckedChange={(checked) =>
                            setForm((current) => ({
                              ...current,
                              [id]: checked === true,
                            }))
                          }
                        />
                        <span>
                          <strong>{label}</strong>
                          <small>{description}</small>
                        </span>
                      </Label>
                    ))}
                  </div>
                </fieldset>

                <div className="grid gap-2">
                  <Label htmlFor="typstBinary">Typst executable</Label>
                  <div className="relative">
                    <TerminalSquareIcon className="input-leading-icon" />
                    <Input
                      className="pl-10 font-mono"
                      id="typstBinary"
                      value={form.typstBinary}
                      required
                      onChange={(event) =>
                        setForm((current) => ({
                          ...current,
                          typstBinary: event.target.value,
                        }))
                      }
                    />
                  </div>
                </div>
              </CardContent>
            </Card>
          </div>

          <aside className="run-column">
            <Card className="run-card">
              <CardHeader>
                <CardTitle>Batch summary</CardTitle>
                <CardDescription>
                  Review your setup before generating.
                </CardDescription>
              </CardHeader>
              <CardContent className="grid gap-5">
                <div className="summary-list">
                  <SummaryRow
                    label="Template"
                    value={
                      form.templatePath
                        ? fileName(form.templatePath)
                        : "Not selected"
                    }
                    ready={Boolean(form.templatePath)}
                  />
                  <SummaryRow
                    label="Data source"
                    value={
                      form.csvPath ? fileName(form.csvPath) : "Not selected"
                    }
                    ready={Boolean(form.csvPath)}
                  />
                  <SummaryRow
                    label="Formats"
                    value={
                      selectedFormats.map(({ label }) => label).join(", ") ||
                      "None"
                    }
                    ready={selectedFormats.length > 0}
                  />
                  <SummaryRow
                    label="Destination"
                    value={
                      form.outputDir ? fileName(form.outputDir) : "Not selected"
                    }
                    ready={Boolean(form.outputDir)}
                  />
                </div>

                <Button
                  className="h-12 w-full"
                  type="submit"
                  disabled={
                    loading ||
                    generating ||
                    !requiredInputsReady ||
                    selectedFormats.length === 0
                  }
                >
                  {generating ? (
                    <LoaderCircleIcon className="size-4 animate-spin" />
                  ) : (
                    <PlayIcon className="size-4 fill-current" />
                  )}
                  {generating ? "Generating..." : "Generate documents"}
                </Button>
                <p className="run-hint">
                  One document will be created for each CSV row.
                </p>
              </CardContent>
            </Card>

            <Card className="activity-card">
              <CardHeader className="pb-0">
                <div className="flex items-center justify-between">
                  <CardTitle className="flex items-center gap-2">
                    <TerminalSquareIcon className="size-4 text-primary" />
                    Activity
                  </CardTitle>
                  {files.length > 0 && (
                    <span className="success-tag">
                      {files.length} created
                    </span>
                  )}
                </div>
              </CardHeader>
              <CardContent>
                <pre className="activity-log">{log}</pre>
                {files.length > 0 && (
                  <div className="output-list">
                    {files.slice(0, 3).map((file) => (
                      <div className="output-file" key={file.path}>
                        <FileTextIcon className="size-4" />
                        <span>{fileName(file.path)}</span>
                        <ChevronRightIcon className="ml-auto size-4" />
                      </div>
                    ))}
                  </div>
                )}
              </CardContent>
            </Card>
          </aside>
        </form>
      </div>
    </main>
  )
}

function SummaryRow({
  label,
  value,
  ready,
}: {
  label: string
  value: string
  ready: boolean
}) {
  return (
    <div className="summary-row">
      {ready ? (
        <CheckCircle2Icon className="size-4 shrink-0 text-primary" />
      ) : (
        <CircleIcon className="size-4 shrink-0 text-muted-foreground/50" />
      )}
      <div className="min-w-0">
        <span>{label}</span>
        <strong className={cn(!ready && "text-muted-foreground")}>
          {value}
        </strong>
      </div>
    </div>
  )
}

export default App
