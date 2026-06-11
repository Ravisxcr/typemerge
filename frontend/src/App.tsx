import { useEffect, useMemo, useState, type ReactNode } from "react"
import {
  CheckIcon,
  ChevronDownIcon,
  CircleIcon,
  FileCode2Icon,
  FileSpreadsheetIcon,
  FileTextIcon,
  FolderIcon,
  FolderOpenIcon,
  LoaderCircleIcon,
  MoreHorizontalIcon,
  PanelLeftIcon,
  PlayIcon,
  PlusIcon,
  Settings2Icon,
  XIcon,
} from "lucide-react"

import { Button } from "@/components/ui/button"
import { Checkbox } from "@/components/ui/checkbox"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"
import { cn } from "@/lib/utils"
import {
  DefaultState,
  Generate,
  SelectFile,
  SelectOutputDirectory,
  SelectProjectDirectory,
  UseTemplatePreset,
} from "../wailsjs/go/gui/App"
import type { app, gui } from "../wailsjs/go/models"

type Format = "pdf" | "png" | "svg"
type PathField = "templatePath" | "csvPath" | "metadataPath" | "outputDir"
type TemplatePreset = "project" | "salary" | "invoice" | "custom"

type Project = {
  id: string
  name: string
  directory: string
  preset: TemplatePreset
  state: gui.State
}

const storageKey = "typemerge.projects.v1"
const formats: Array<{ id: Format; label: string }> = [
  { id: "pdf", label: "PDF" },
  { id: "png", label: "PNG" },
  { id: "svg", label: "SVG" },
]

const blankState: gui.State = {
  templatePath: "",
  csvPath: "",
  metadataPath: "",
  outputDir: "",
  typstBinary: "typst",
  pdf: true,
  png: false,
  svg: false,
}

function newID() {
  return `${Date.now()}-${Math.random().toString(16).slice(2)}`
}

function errorMessage(error: unknown) {
  return error instanceof Error ? error.message : String(error)
}

function fileName(path: string) {
  return path.split(/[\\/]/).pop() || path
}

function pathIn(directory: string, name: string) {
  if (!directory) return name
  const separator = directory.includes("\\") ? "\\" : "/"
  return `${directory.replace(/[\\/]$/, "")}${separator}${name}`
}

function App() {
  const [projects, setProjects] = useState<Project[]>([])
  const [activeID, setActiveID] = useState("")
  const [files, setFiles] = useState<app.GeneratedFile[]>([])
  const [log, setLog] = useState("Choose a project folder to begin.")
  const [generating, setGenerating] = useState(false)
  const [inspectorOpen, setInspectorOpen] = useState(true)

  useEffect(() => {
    const stored = localStorage.getItem(storageKey)
    if (stored) {
      try {
        const restored = JSON.parse(stored) as Project[]
        if (restored.length) {
          setProjects(restored)
          setActiveID(restored[0].id)
          return
        }
      } catch {
        localStorage.removeItem(storageKey)
      }
    }

    DefaultState()
      .then((state) => {
        const project = {
          id: newID(),
          name: "Untitled project",
          directory: "",
          preset: "project" as const,
          state,
        }
        setProjects([project])
        setActiveID(project.id)
      })
      .catch((error: unknown) => setLog(`Error: ${errorMessage(error)}`))
  }, [])

  useEffect(() => {
    if (projects.length) {
      localStorage.setItem(storageKey, JSON.stringify(projects))
    }
  }, [projects])

  const active = projects.find(({ id }) => id === activeID) ?? projects[0]
  const selectedFormats = useMemo(
    () => formats.filter(({ id }) => active?.state[id]),
    [active],
  )
  const ready = Boolean(
    active?.state.templatePath &&
      active.state.csvPath &&
      active.state.metadataPath &&
      active.state.outputDir &&
      selectedFormats.length,
  )

  const updateProject = (change: (project: Project) => Project) => {
    setProjects((current) =>
      current.map((project) => (project.id === active?.id ? change(project) : project)),
    )
  }

  const updateState = (field: keyof gui.State, value: string | boolean) => {
    updateProject((project) => ({
      ...project,
      state: { ...project.state, [field]: value },
    }))
  }

  const openProject = async () => {
    try {
      const setup = await SelectProjectDirectory()
      if (!setup.directory) return
      const state = await DefaultState()
      const project: Project = {
        id: newID(),
        name: setup.name || "Untitled project",
        directory: setup.directory,
        preset: "project",
        state: {
          ...state,
          templatePath: setup.templatePath,
          csvPath: setup.csvPath,
          metadataPath: setup.metadataPath,
          outputDir: setup.outputDir,
        },
      }
      setProjects((current) => [
        ...current.filter((item) => item.directory || item.name !== "Untitled project"),
        project,
      ])
      setActiveID(project.id)
      setFiles([])
      setLog(`Opened ${setup.directory}`)
    } catch (error) {
      setLog(`Error: ${errorMessage(error)}`)
    }
  }

  const addBlankProject = async () => {
    const state = await DefaultState().catch(() => ({ ...blankState }))
    const project: Project = {
      id: newID(),
      name: `Untitled ${projects.length + 1}`,
      directory: "",
      preset: "project",
      state,
    }
    setProjects((current) => [...current, project])
    setActiveID(project.id)
  }

  const removeProject = (id: string) => {
    setProjects((current) => {
      const next = current.filter((project) => project.id !== id)
      setActiveID((selected) => (selected === id ? (next[0]?.id ?? "") : selected))
      return next
    })
  }

  const selectPath = async (
    field: Exclude<PathField, "outputDir">,
    kind: "template" | "csv" | "metadata",
  ) => {
    try {
      const path = await SelectFile(kind)
      if (path) updateState(field, path)
    } catch (error) {
      setLog(`Error: ${errorMessage(error)}`)
    }
  }

  const selectOutput = async () => {
    try {
      const path = await SelectOutputDirectory()
      if (path) updateState("outputDir", path)
    } catch (error) {
      setLog(`Error: ${errorMessage(error)}`)
    }
  }

  const selectPreset = async (preset: TemplatePreset) => {
    if (!active) return
    if (preset === "custom") {
      const path = await SelectFile("template")
      if (!path) return
      updateProject((project) => ({
        ...project,
        preset,
        state: { ...project.state, templatePath: path },
      }))
      return
    }
    if (preset !== "project") {
      try {
        const path = await UseTemplatePreset(active.directory, preset)
        updateProject((project) => ({
          ...project,
          preset,
          state: { ...project.state, templatePath: path },
        }))
      } catch (error) {
        setLog(`Error: ${errorMessage(error)}`)
      }
      return
    }
    updateProject((project) => ({
      ...project,
      preset,
      state: {
        ...project.state,
        templatePath: pathIn(project.directory, "template.typ"),
      },
    }))
  }

  const generate = async () => {
    if (!active) return
    setGenerating(true)
    setLog(`Rendering ${active.name}...`)
    try {
      const result = await Generate(active.state)
      const generated = result.files ?? []
      setFiles(generated)
      setLog(
        generated.length
          ? `Complete. Created ${generated.length} files in ${active.state.outputDir}`
          : "No files were generated.",
      )
    } catch (error) {
      setFiles([])
      setLog(`Error: ${errorMessage(error)}`)
    } finally {
      setGenerating(false)
    }
  }

  if (!active) return <main className="loading-screen">Loading workspace...</main>

  return (
    <main className="studio-shell">
      <header className="titlebar">
        <div className="app-brand">
          <div className="brand-tile">Tm</div>
          <strong>TypeMerge</strong>
        </div>
        <nav className="menu-bar" aria-label="Application menu">
          <span>File</span><span>Edit</span><span>Project</span><span>View</span><span>Help</span>
        </nav>
        <div className="titlebar-actions">
          <button type="button" onClick={() => setInspectorOpen((value) => !value)}>
            <PanelLeftIcon className="size-4" />
          </button>
          <Button size="sm" onClick={generate} disabled={!ready || generating}>
            {generating ? <LoaderCircleIcon className="size-4 animate-spin" /> : <PlayIcon className="size-4 fill-current" />}
            Generate
          </Button>
        </div>
      </header>

      <section className={cn("studio-grid", !inspectorOpen && "inspector-hidden")}>
        <aside className="project-panel">
          <div className="panel-heading">
            <span>Projects</span>
            <button type="button" title="New project" onClick={addBlankProject}>
              <PlusIcon className="size-4" />
            </button>
          </div>
          <button className="open-project" type="button" onClick={openProject}>
            <FolderOpenIcon className="size-4" />
            Open folder
          </button>
          <div className="project-list">
            {projects.map((project) => (
              <button
                className={cn("project-item", project.id === active.id && "active")}
                key={project.id}
                type="button"
                onClick={() => setActiveID(project.id)}
              >
                <span className="project-icon"><FileCode2Icon className="size-4" /></span>
                <span className="project-copy">
                  <strong>{project.name}</strong>
                  <small>{project.directory || "No folder selected"}</small>
                </span>
                {projects.length > 1 && (
                  <span
                    className="remove-project"
                    role="button"
                    tabIndex={0}
                    onClick={(event) => {
                      event.stopPropagation()
                      removeProject(project.id)
                    }}
                  >
                    <XIcon className="size-3" />
                  </span>
                )}
              </button>
            ))}
          </div>
          <div className="panel-section">
            <div className="section-label">Project files</div>
            <FileTreeRow icon={FileCode2Icon} label={fileName(active.state.templatePath) || "Template missing"} ready={Boolean(active.state.templatePath)} />
            <FileTreeRow icon={FileSpreadsheetIcon} label={fileName(active.state.csvPath) || "CSV missing"} ready={Boolean(active.state.csvPath)} />
            <FileTreeRow icon={FileTextIcon} label={fileName(active.state.metadataPath) || "Metadata missing"} ready={Boolean(active.state.metadataPath)} />
            <FileTreeRow icon={FolderIcon} label="output" ready={Boolean(active.state.outputDir)} />
          </div>
        </aside>

        <section className="editor-area">
          <div className="document-tabs">
            <div className="document-tab active">
              <FileCode2Icon className="size-3.5" />
              {fileName(active.state.templatePath) || "template.typ"}
              <XIcon className="size-3" />
            </div>
            <button type="button"><PlusIcon className="size-4" /></button>
          </div>

          <div className="canvas">
            <div className="canvas-toolbar">
              <span>100%</span><ChevronDownIcon className="size-3" />
              <i />
              <span>A4</span>
              <MoreHorizontalIcon className="size-4" />
            </div>
            <div className="page">
              <div className="page-accent" />
              <div className="page-header">
                <div>
                  <span className="page-kicker">TYPEMERGE DOCUMENT</span>
                  <h2>{active.name}</h2>
                  <p>Generated from structured project data</p>
                </div>
                <div className="page-logo">Tm</div>
              </div>
              <div className="page-rule" />
              <div className="preview-grid">
                <PreviewField label="Template" value={fileName(active.state.templatePath) || "Select a template"} />
                <PreviewField label="Data source" value={fileName(active.state.csvPath) || "Select CSV data"} />
                <PreviewField label="Metadata" value={fileName(active.state.metadataPath) || "Select metadata"} />
                <PreviewField label="Export" value={selectedFormats.map((format) => format.label).join(", ") || "No format"} />
              </div>
              <div className="preview-table">
                <div className="preview-table-head"><span>FIELD</span><span>MERGED VALUE</span></div>
                {["name", "document_id", "date", "amount"].map((field, index) => (
                  <div className="preview-table-row" key={field}>
                    <span>{field}</span><span>{`{{ row.${field} }}`}</span><b>{index + 1}</b>
                  </div>
                ))}
              </div>
              <div className="page-footer">Preview layout · Final content is created from each CSV row</div>
            </div>
          </div>

          <div className="activity-strip">
            <div className="activity-title">
              <CircleIcon className={cn("size-2 fill-current", files.length ? "success" : "")} />
              Activity
            </div>
            <code>{log}</code>
            {files.length > 0 && <span>{files.length} files</span>}
          </div>
        </section>

        {inspectorOpen && (
          <aside className="inspector">
            <div className="inspector-title">
              <div><Settings2Icon className="size-4" /><span>Properties</span></div>
              <MoreHorizontalIcon className="size-4" />
            </div>

            <InspectorSection title="Project">
              <Label htmlFor="project-name">Name</Label>
              <Input
                id="project-name"
                value={active.name}
                onChange={(event) => updateProject((project) => ({ ...project, name: event.target.value }))}
              />
              <Label>Workspace folder</Label>
              <PathControl value={active.directory} placeholder="Choose a project folder" onBrowse={openProject} />
            </InspectorSection>

            <InspectorSection title="Template">
              <div className="preset-grid">
                {([
                  ["project", "Project", "template.typ"],
                  ["salary", "Salary slip", "salary-slip.typ"],
                  ["invoice", "Invoice", "invoice.typ"],
                  ["custom", "Custom", "Choose a file"],
                ] as const).map(([id, title, copy]) => (
                  <button
                    className={cn("preset-card", active.preset === id && "selected")}
                    key={id}
                    type="button"
                    onClick={() => selectPreset(id)}
                  >
                    <FileCode2Icon className="size-4" />
                    <span><strong>{title}</strong><small>{copy}</small></span>
                    {active.preset === id && <CheckIcon className="check size-3" />}
                  </button>
                ))}
              </div>
            </InspectorSection>

            <InspectorSection title="Data files">
              <Label>CSV data</Label>
              <PathControl value={active.state.csvPath} placeholder="Required .csv file" onBrowse={() => selectPath("csvPath", "csv")} />
              <Label>Metadata</Label>
              <PathControl value={active.state.metadataPath} placeholder="Required metadata file" onBrowse={() => selectPath("metadataPath", "metadata")} />
            </InspectorSection>

            <InspectorSection title="Export">
              <div className="format-row">
                {formats.map(({ id, label }) => (
                  <Label className={cn("format-chip", active.state[id] && "selected")} htmlFor={id} key={id}>
                    <Checkbox id={id} checked={active.state[id]} onCheckedChange={(checked) => updateState(id, checked === true)} />
                    {label}
                  </Label>
                ))}
              </div>
              <Label>Output folder</Label>
              <PathControl value={active.state.outputDir} placeholder="output" onBrowse={selectOutput} />
            </InspectorSection>

            <div className="inspector-run">
              <Button onClick={generate} disabled={!ready || generating}>
                {generating ? <LoaderCircleIcon className="size-4 animate-spin" /> : <PlayIcon className="size-4 fill-current" />}
                {generating ? "Generating..." : "Generate documents"}
              </Button>
              <small>{ready ? "Project is ready to export" : "Complete the missing project files"}</small>
            </div>
          </aside>
        )}
      </section>
    </main>
  )
}

function FileTreeRow({ icon: Icon, label, ready }: { icon: typeof FileTextIcon; label: string; ready: boolean }) {
  return (
    <div className={cn("file-tree-row", !ready && "missing")}>
      <Icon className="size-3.5" /><span>{label}</span>
      <CircleIcon className={cn("ml-auto size-1.5 fill-current", ready && "ready")} />
    </div>
  )
}

function PreviewField({ label, value }: { label: string; value: string }) {
  return <div><span>{label}</span><strong>{value}</strong></div>
}

function InspectorSection({ title, children }: { title: string; children: ReactNode }) {
  return (
    <section className="inspector-section">
      <div className="section-label">{title}</div>
      {children}
    </section>
  )
}

function PathControl({ value, placeholder, onBrowse }: { value: string; placeholder: string; onBrowse: () => void }) {
  return (
    <div className="path-control">
      <span title={value}>{value ? fileName(value) : placeholder}</span>
      <button type="button" onClick={onBrowse}><FolderOpenIcon className="size-3.5" /></button>
    </div>
  )
}

export default App
