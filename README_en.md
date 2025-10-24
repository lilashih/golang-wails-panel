# golang-wails-panel

## ▍Features
This project is a desktop application with two main features:

### 1. Project Panel
You can place multiple projects under the [projects](release/projects/) directory. The application will automatically scan all projects in this directory and display them in the panel based on each project's `project.json` configuration file. Users can start, stop, or install each project directly from the panel for easy management.

#### `project.json` Example

**`pm2` Project Example:**
```json
{
  "title": "Nodejs Project Example",
  "type": "pm2",
  "key": "pm2 process name",
  "start": "pnpm start",
  "stop": "pnpm stop",
  "install": "pnpm install"
}
```

**`exe` Project Example:**
```json
{
  "title": "Go Project Example",
  "type": "exe",
  "key": "process name, usually the executable filename",
  "start": "start /b app.exe",
  "stop": "taskkill /IM app.exe /F",
  "install": ""
}
```

| Field   | Description |
|---------|---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| title   | Project name, shown in the panel. |
| type    | Project type, supports two types:<br>• pm2: For Node.js or similar projects managed by pm2. Commands are usually pnpm/yarn/npm scripts.<br>• exe: For Windows executables (.exe). start/stop are Windows commands. |
| key     | Unique project identifier, usually for internal use. |
| start   | Command to start the project. |
| stop    | Command to stop the project. |
| install | Command to install or initialize the project (optional). |

> The panel adapts to different project types and provides one-click operations for each.

### 2. Log Viewer
Automatically scans all log files under the [log](storage/log/) directory. You can also select additional files to view. The viewer supports reading and searching very large files efficiently, making log analysis and debugging easier.

### 3. System Tray (Systray)
Supports minimizing to the system tray and running in the background, with a quick action menu (e.g., show/hide window, open logs folder, quit the app) so you can manage the application without occupying the taskbar.

---

## ▍Directory Structure
```
├── app.go                  # Go main entry point
├── go.mod
├── go.sum
├── main.go                 # Go main entry point
├── wails.json              # Wails config
├── README.md               # Project documentation
├── frontend/               # Frontend code (see frontend/README.md for details)
├── src/                    # Backend Go code
│   ├── core/               # Core features (cmd, config, helper, logger, pm2)
│   ├── def/                # Definitions
│   └── service/            # Services (log_viewer, panel, etc.)
├── release/                # Release resources
│   ├── bin/                # Compiled Go executables
│   └── projects/           # Projects directory
├── storage/                # Data storage
│   └── log/                # System log files
```

## ▍Development Commands

### Installation
Install Wails CLI and tidy Go dependencies:
```bash
go install github.com/wailsapp/wails/v2/cmd/wails@latest
go mod tidy
```

### Run
Start development mode (auto-reloads on code changes):
```bash
wails dev
```

### Build
Compile the project into an executable:
```bash
wails build
```

## ▍Additional Notes

### projects (For Development Only)
If you are developing a Go project under the `projects` directory, you can create a shortcut (symlink) to the build or release directory for easier access in the panel. 

> ⚠️ Please use the command line to create symlinks (do not use right-click shortcuts). Only use this during development; do not keep symlinks in production.**

Example command to create a directory symlink (run in `projects` directory):
```cmd
mklink /D your_link target_folder
mklink /D aaa path1\path2\release
```
