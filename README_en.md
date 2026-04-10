# golang-wails-panel

## ▍Features
This project is a desktop application with the following features:

1. Project Panel  
    You can place multiple projects under the [projects](release/projects/) directory. The application will automatically scan all projects in this directory and display them in the panel based on each project's `project.json` configuration file. Users can start, stop, or install each project directly from the panel for easy management.

    | Field   | Description |
    |---------|---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
    | title   | Project name, shown in the panel. |
    | type    | Project type, supports two types:<br>• pm2: For Node.js or similar projects managed by pm2. Commands are usually pnpm/yarn/npm scripts.<br>• exe: For Windows executables (.exe). start/stop are Windows commands. |
    | key     | Unique project identifier, usually for internal use.<br>• pm2：The PM2 process name (must match the "name" field in PM2).。<br>• exe：The executable process name, usually the filename. |
    | start   | Command to start the project. |
    | stop    | Command to stop the project. |
    | install | Command to install or initialize the project (optional). |


    > The panel adapts to different project types and provides one-click operations for each.

    - **`pm2` Project Example:**
        ```json
        {
            "title": "Nodejs Project Example",
            "type": "pm2",
            "key": "pm2 process name",
            "start": "pnpm start",        // or npm start (must be defined in package.json)
            "stop": "pnpm stop",          // or npm stop (must be defined in package.json)
            "install": "pnpm install"     // or npm install
        }
        ```

    - **`exe` Project Example:**
        ```json
        {
            "title": "Go Project Example",
            "type": "exe",
            "key": "process name, usually the executable filename",
            "start": "start /b app.exe",
            "stop": "taskkill /IM app.exe /F",
            "install": ""   // leave blank if installation is not required
        }
        ```

2. Log Viewer  
Automatically scans all log files under the [log](storage/log/) directory. You can also select additional files to view. The viewer supports reading and searching very large files efficiently, making log analysis and debugging easier.

3. System Tray (`Systray`)  
Supports minimizing to the system tray and running in the background, with a quick action menu (e.g., show/hide window, open logs folder, quit the app) so you can manage the application without occupying the taskbar.

---

## ▍Directory Structure
```
├── app.go                      # Wails register
├── go.mod
├── go.sum
├── main.go                     # Go main entry point
├── wails.json                  # Wails config
├── README.md                   # Project documentation
├── frontend/                   # Frontend code (see frontend/README.md for details)
├── src/                        # Backend Go code
│   ├── core/                   # Core features (cmd, config, def, helper, logger, pm2)
│   └── service/                # Services (log_viewer, panel, etc.)
├── release/                    # Release resources
│   ├── golang-wails-panel.exe  # Compiled Go executable
│   └── projects/               # Projects directory
└── storage/                    # Local storage (log/, etc.)
```

## ▍Dependencies

| Package | Description |
|---|---|
| [Wails](https://github.com/wailsapp/wails) | A Go-based desktop application framework (version 2). |
| [fyne.io/systray](https://github.com/fyne-io/systray) | Enables system tray functionality. |
| [caarlos0/env](https://github.com/caarlos0/env) | Maps environment variables to Go structs. |
| [golobby/dotenv](https://github.com/golobby/dotenv) | Loads environment variables from `.env` files. |


## ▍Environment Variables
All configuration is located in [src/core/config](/src/core/config) and initialized automatically by [src/core/config/main.go](/src/core/config/main.go), which also loads `.env`.

- Example:
    ```go
    import (
        "fmt"
        "gbase/src/core/config"
    )

    fmt.Println(config.App.Name)
    ```

### Common Environment Variables (.env)

| Name | Purpose | Default Value (Production) | Notes |
|---|---|---|---|
| `APP_MODE` | Runtime mode (`release` / `debug` / `test`) | `release` | Use `release` in production |
| `PROJECT_BASE_PATH` | Project directory to scan | `./projects` | Set the directory to scan directly. Both relative and absolute paths are supported, such as `./release/projects` or `C:\myDir\projects` |

## ▍Logging
If you only want to write log messages to files, use `logger.Log`. If you want logs to be displayed in the frontend as well, use `runtime`.

| Logger | Behavior |
| ----------- | ------------------------------------------------------------ |
| `logger.Log` | In `release` mode, logs are written to files only. In other modes, logs are also output to the console. |
| `runtime` | Uses the Wails Logger Adapter to wrap the existing logger. In addition to executing the custom `logger.Log`, it also sends log events to the frontend in real time through `runtime.EventsEmit`. |

- Generated log path: the application automatically creates the `log` directory under `storage` in the project root (default: `./storage/log`).
- Filename format: `log-YYYY-MM-DD.log`. New log file generated daily (with `UTF-8 BOM`).
- Characteristics:
    - Automatically rotates the log file when the date changes.
    - Removes `ANSI` color escape codes to keep the log content as plain text.
    - `logger.Log` is actually a wrapper around `*log.Logger` from the Go standard library, so you can use standard `log` methods and behavior directly:
        - `Print/Printf/Println`, `Fatal/Fatalf`, `Panic/Panicf`, and similar methods are available.
        - The default flags are `log.LstdFlags | log.Lshortfile`, which include timestamps and the caller file:line.
    - Expired logs are automatically compressed to `.gz` to save disk space.
    - If you need to change the format or add a log rotate strategy, you can modify it directly in [src/core/logger](/src/core/logger).
- Wails and frontend log events:  
    - [wails.go](src/core/logger/wails.go) defines `RegisterCtx` and `NewWailsLog()`. The backend registers this adapter with Wails during startup.
    - The `adapter` forwards each log entry to `logger.Log` and also emits an event through `runtime.EventsEmit` so the frontend can display logs in real time.
- Example:
    ```go
    import (
        "gbase/src/core/logger"
        "github.com/wailsapp/wails/v2/pkg/runtime"
    )

    // Write a regular log
    logger.Log.Println("Server started")

    // Wails log, can be displayed in the frontend (requires ctx)
    runtime.LogErrorf(s.ctx, "Server started")
    ```

## ▍Development Commands

- Installation  
    Install `Wails CLI` and tidy `Go` dependencies:
    ```bash
    go install github.com/wailsapp/wails/v2/cmd/wails@latest
    go mod tidy
    ```

- Run  
    Start development mode (auto-reloads on code changes):
    ```bash
    wails dev
    ```

## ▍ Build  
Compile the project into an executable:
```bash
wails build
```

The `Wails` build flow has two stages:

1. Default output directory `bin`  
    After running `wails build`, the compiled executable is first generated under `bin/` (`Wails` always places the build output inside `bin/`).
    - Windows: `bin/golang-wails-panel.exe`
    - Linux: `bin/golang-wails-panel`

2. Automatically sync to the `release` directory  
    This project uses `postBuildHooks` configured in `wails.json` to run a script after the build finishes, then copies the executable to `release/` for distribution.

    Different operating systems call different scripts:
    - Windows: `scripts/post-build.ps1`
    - Linux: `scripts/post-build.sh`

### Build Linux Executable with Docker on Windows
If you are on a `Windows` machine and want to quickly generate a `Linux` executable, this project includes a `Docker` build flow. It will:
- Run `wails build` inside a `Linux` container
- Export the build output directly into the project's `release/` directory
- Automatically run `docker builder prune -af` after the build to clean the build cache
- Avoid keeping the final `image`

#### Prerequisites
- `Docker Desktop` is installed
- `Docker buildx` works correctly

#### Basic Usage
Run this at the project root:
```powershell
.\scripts\build-linux-docker.ps1
```

The default output path is:
```text
release/golang-wails-panel
```

If you need to keep outputs for multiple architectures at the same time, specify `-OutputDir` yourself, for example `release/linux-arm64`.

#### Available Parameters
```powershell
.\scripts\build-linux-docker.ps1 -Arch amd64
.\scripts\build-linux-docker.ps1 -Arch arm64
.\scripts\build-linux-docker.ps1 -Arch amd64 -NoCache
.\scripts\build-linux-docker.ps1 -Arch amd64 -UseWebkit241
```

Parameter details:
| Parameter | Description |
|---|---|
| `-Arch` | Target architecture. Supported values: `amd64`, `arm64`. |
| `-OutputDir` | Custom output directory. Default is `release`. If you want to keep multiple architecture outputs, use different directories manually. |
| `-NoCache` | Do not use Docker build cache. |
| `-UseWebkit241` | Use `libwebkit2gtk-4.1-dev` for newer `Linux` environments and add `-tags webkit2_41`. |

#### Use the Dockerfile Directly
If you do not want to use the PowerShell script, you can run:
```powershell
docker buildx build --pull `
    --build-arg TARGETARCH=amd64 `
    -f .\scripts\Dockerfile.linux-build `
    --output type=local,dest=.\release `
    .
docker builder prune -af
```

If the target `Linux` distribution requires `webkit2gtk-4.1`, use:
```powershell
docker buildx build --pull `
    --build-arg TARGETARCH=amd64 `
    --build-arg WEBKIT_PKG=libwebkit2gtk-4.1-dev `
    --build-arg WAILS_TAGS=webkit2_41 `
    -f .\scripts\Dockerfile.linux-build `
    --output type=local,dest=.\release `
    .
docker builder prune -af
```


## ▍Testing  
Place test files directly under the `test` directory (do not create subdirectories), then run:
```
go test -v ./test/...
```


## ▍Additional Notes

- projects (For Development Only)
    If you are developing a Go project under the `projects` directory, you can create a shortcut (symlink) to the build or release directory for easier access in the panel. 

    > ⚠️ Please use the command line to create symlinks (do not use right-click shortcuts). Only use this during development; do not keep symlinks in production.**

    Example command to create a directory symlink (run in `projects` directory):
    ```cmd
    mklink /D your_link target_folder
    mklink /D aaa path1\path2\release
    ```
