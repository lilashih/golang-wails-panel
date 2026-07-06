# golang-wails-panel

`golang-wails-panel` 是一個以 Wails3、Go、Vue 3 建立的桌面管理工具，主要用途是集中管理專案、檢視應用程式日誌，並透過桌面視窗提供穩定的操作介面。

本專案目前採用 Wails3 架構，後端以 `application.Service` 註冊服務，前端使用 Wails3 generated bindings 呼叫 Go 方法。

## 主要功能

### 專案面板

專案面板會掃描 `PROJECT_BASE_PATH` 指定的專案目錄，讀取每個子專案底下的 `project.json`，並顯示可操作的專案清單。

目前支援的操作：

- 讀取專案清單與執行狀態。
- 開啟專案資料夾。
- 啟動專案。
- 停止專案。
- 執行安裝或重建依賴指令。
- 操作過程會寫入 log，並同步顯示於前端即時日誌面板。

後端會依 `project.json` 中符合目前作業系統的設定建立 runner。Windows、Linux、macOS 等平台差異由後端處理，前端不需要判斷作業系統。

### 日誌檢視器

日誌檢視器會讀取 `storage/log` 底下的 `.log` 與 `.txt` 檔案，提供大型檔案分批載入與查找能力。

目前支援：

- 列出日誌檔案。
- 設定目前檔案。
- 分批讀取檔案開頭與後續內容。
- 計算總行數。
- 搜尋、篩選、標記文字。
- tail follow，即時追蹤目前 log 檔新增內容。
- 開啟 log 檔所在資料夾。

### 開發指南

前端提供 README 閱讀器，可透過後端 API 取得可閱讀文件清單與內容。

### 即時日誌

後端可透過 `logger.Info`、`logger.Warn`、`logger.Error`、`logger.Debug` 寫入檔案並同步事件到前端。前端由全域 log store 統一接收 `log` event，不在各頁面重複註冊事件。

## 技術架構

| 類別 | 技術 |
| --- | --- |
| 桌面框架 | Wails3 |
| 後端 | Go |
| 後端服務註冊 | `github.com/wailsapp/wails/v3/pkg/application` |
| 前端 | Vue 3、TypeScript、Vite |
| 前端狀態 | Pinia |
| 前端路由 | Vue Router |
| 樣式 | Tailwind CSS 4、CSS theme token |
| 圖示 | lucide-vue-next |
| 大型清單 | vue3-virtual-scroller |
| 設定 | `.env`、`src/core/config` |

## 目錄架構

```text
.
├── main.go                         # Wails3 應用程式入口
├── Taskfile.yml                    # Wails3 task 與開發指令
├── build/
│   └── config.yml                  # Wails3 專案與 dev mode 設定
├── docs/                           # 開發規畫與設計文件
├── frontend/                       # Vue 3 前端
│   ├── bindings/                   # Wails3 產生的前後端 bindings
│   ├── src/
│   │   ├── components/             # 共用元件
│   │   ├── composables/            # 前端 composable
│   │   ├── pages/                  # 頁面（Panel、LogViewer、DevGuide）
│   │   ├── router/                 # 路由設定
│   │   ├── stores/                 # Pinia store
│   │   ├── styles/                 # Tailwind 與主題樣式
│   │   └── types/                  # TypeScript 型別
│   └── package.json
├── src/
│   ├── app/                        # Wails3 app 組裝層
│   ├── core/                       # 核心共用能力
│   │   ├── cmd/                    # shell、explore 等跨平台指令
│   │   ├── config/                 # 環境變數與設定
│   │   ├── def/                    # 共用常數
│   │   ├── helper/                 # 共用 helper
│   │   └── logger/                 # 日誌與前端事件
│   └── service/                    # 後端服務層
│       ├── app_info/               # 應用程式資訊服務
│       ├── log_viewer/             # 日誌檢視服務
│       ├── panel/                  # 專案面板服務
│       └── readme_reader/          # README 閱讀服務
├── storage/
│   └── log/                        # 應用程式日誌
└── release/                        # 發佈或本機執行資源
```

## 後端架構

### `main.go`

`main.go` 只負責 Wails3 應用程式啟動流程：

- 註冊 Wails3 event 型別。
- 建立後端 app 組裝物件。
- 建立 `application.App`。
- 註冊 services。
- 設定 Wails3 logger。
- 設定嵌入的前端 assets。
- 建立主視窗。
- 啟動 `app.Run()`。

前端 assets 由以下方式嵌入：

```go
//go:embed all:frontend/dist
var assets embed.FS
```

### `src/app`

`src/app` 是 Wails3 後端組裝層，避免讓 `main.go` 直接依賴所有業務細節。

| 檔案 | 職責 |
| --- | --- |
| `app.go` | 建立 service instance，連接 Wails app 與事件 emitter |
| `services.go` | 將 Go service 包裝成 `application.Service` |
| `events.go` | 註冊 Wails3 typed events |

目前註冊的 services：

- `PanelService`
- `AppInfoService`
- `LogViewerService`
- `ReadmeReaderService`

目前註冊的 events：

- `time`
- `log`
- `logViewer:tail`

### `src/core`

`src/core` 放跨功能共用能力，不放特定頁面的業務邏輯。

| 目錄 | 說明 |
| --- | --- |
| `cmd/explore` | 跨平台開啟資料夾 |
| `cmd/shell` | 跨平台執行 shell 指令、timeout、環境變數補強 |
| `config` | 載入 `.env` 與預設設定 |
| `def` | 共用格式與常數 |
| `helper` | 路徑、檔案、時間、結構轉換等 helper |
| `logger` | 每日日誌、壓縮、前端事件、Wails3 slog adapter |

### `src/service`

`src/service` 依功能切分 package。每個 service 對應一組前端可呼叫的後端能力。

新增服務時應遵守：

- 新增於 `src/service/<feature>`。
- service 型別命名為 `Service`。
- 若需要啟動或關閉生命週期，實作 Wails3 service lifecycle。
- 在 `src/app/app.go` 建立 instance。
- 在 `src/app/services.go` 註冊為 `application.NewService(...)`。
- 若需要推送事件，透過小型 interface 或 app 組裝層注入，不要讓核心模組直接綁死前端頁面。

Wails3 service lifecycle 範例：

```go
func (s *Service) ServiceStartup(ctx context.Context, _ application.ServiceOptions) error {
    s.ctx = ctx
    return nil
}

func (s *Service) ServiceShutdown() error {
    return nil
}
```

## 專案面板設定

### 專案目錄

專案清單由 `PROJECT_BASE_PATH` 指定。預設值是：

```env
PROJECT_BASE_PATH=./projects
```

路徑規則：

- 可使用絕對路徑，如 C:\projects。
- 可使用相對路徑，相對於 runtime base path。
- `APP_MODE=release` 時，runtime base path 優先使用執行檔所在目錄。
- 非 release 模式會使用目前工作目錄，方便 `wails3 dev` 與 IDE 開發。
- 支援一般資料夾與 symlink。

### `project.json`

每個專案資料夾底下需放置 `project.json`。此檔案必須是陣列格式，且每筆設定可對應不同作業系統。

必要欄位：

| 欄位 | 說明 |
| --- | --- |
| `os` | 作業系統，常用值為 `windows`、`linux`、`darwin` |
| `title` | 前端顯示名稱 |
| `key` | runner 判斷執行狀態用的唯一識別 |
| `type` | runner 類型，目前支援 `pm2`、`exe` |
| `start` | 啟動指令 |
| `stop` | 停止指令 |
| `install` | 安裝或重建依賴指令，可為空字串 |

Windows 與 Linux 共用設定範例：

```json
[
  {
    "os": "windows",
    "title": "Nodejs 專案設定範例",
    "type": "pm2",
    "key": "my-node-app",
    "start": "pnpm start",
    "stop": "pnpm stop",
    "install": "pnpm install"
  },
  {
    "os": "linux",
    "title": "Nodejs 專案設定範例",
    "type": "pm2",
    "key": "my-node-app",
    "start": "pnpm start",
    "stop": "pnpm stop",
    "install": "pnpm install"
  }
]
```

`exe` 類型範例：

```json
[
  {
    "os": "windows",
    "title": "執行檔專案設定範例",
    "type": "exe",
    "key": "app_windows.exe",
    "start": "start /B app_windows.exe",
    "stop": "taskkill /IM app_windows.exe /F",
    "install": "app_windows.exe migrate"
  },
  {
    "os": "linux",
    "title": "執行檔專案設定範例",
    "type": "exe",
    "key": "app_linux",
    "start": "nohup ./app_linux > ./app.log 2>&1 &",
    "stop": "pkill -f app_linux",
    "install": "./app_linux migrate"
  }
]
```

### runner 行為

| 類型 | 適用情境 | 執行狀態判斷 |
| --- | --- | --- |
| `pm2` | Node.js、前端服務、由 PM2 管理的常駐服務 | 執行 `pm2 jlist`，比對 `name == key` 且狀態為 `online` |
| `exe` | 一般執行檔或自行管理的常駐 process | Windows 透過 task list 判斷；Unix-like 透過 process list 判斷 |

啟動與停止指令有 30 秒 timeout。若指令本身是常駐服務，必須自行背景化並結束，否則 Panel 會等待到 timeout。

建議：

- Windows 常駐啟動可使用 `start /B ...`。
- Linux 常駐啟動可使用 `nohup ... > app.log 2>&1 &`。
- `pm2` 的 `key` 必須與 PM2 process name 一致。
- `exe` 的 `key` 應填可辨識的 process 名稱，不要只填顯示名稱。

## 日誌系統

### 日誌目錄

日誌預設寫入：

```text
storage/log
```

檔名格式：

```text
log-YYYY-MM-DD.log
```

日誌特性：

- 每日自動換檔。
- `LOGGER_COMPRESS=true` 時，過期 log 會自動壓縮為 `.gz`。
- 寫入前會移除 ANSI 顏色控制碼。
- `logger.Log` 為標準 library `*log.Logger`的包裝，因此可以直接使用 `log` 的方法與行為：
    - 可直接使用 `Print/Printf/Println`、`Fatal/Fatalf`、`Panic/Panicf` 等方法。
    - 預設使用 `flags` `log.LstdFlags | log.Lshortfile`（含時間戳與呼叫來源檔案:行號）。  
- `logger.Info`、`logger.Warn`、`logger.Error`、`logger.Debug` 會寫檔並推送 `log` event 給前端。

### 後端使用方式

只寫入檔案：

```go
logger.Log.Println("背景任務完成")
```

寫入檔案並同步前端：

```go
logger.Info("正在啟動服務：%s", name)
logger.Error("啟動失敗：%v", err)
```

Wails3 framework logger 由 `logger.NewSlogLogger(slog.LevelInfo)` 接到本專案 logger，讓框架層 log 也能進入既有日誌系統。

## 環境變數

設定由 `src/core/config` 初始化，會先套用預設值，再嘗試讀取 `.env`。

`.env` 不存在時，程式仍會用預設值啟動。

| 名稱 | 預設值 | 說明 |
| --- | --- | --- |
| `APP_ID` | `Panel` | 應用程式 ID，使用者設定目錄 fallback 會用到 |
| `APP_NAME` | `Panel` | 應用程式名稱 |
| `APP_MODE` | `release` | 執行模式，常用 `debug`、`release` |
| `APP_BASE_PATH` | `./` | 應用程式基準路徑設定 |
| `PROJECT_BASE_PATH` | `./projects` | 專案面板掃描的專案根目錄 |
| `PROJECT_JSON` | `project.json` | 專案設定檔名稱 |
| `LOGGER_COMPRESS` | `true` | 是否壓縮過期日誌 |

本機開發常見 `.env`：

```env
APP_MODE=debug
PROJECT_BASE_PATH=C:\projects
LOGGER_COMPRESS=false
```

## 前端架構

前端位於 `frontend`，使用 Vue 3、TypeScript、Vite、Pinia、Vue Router 與 Tailwind CSS 4。

### 主要結構

```text
frontend/src
├── App.vue
├── main.ts
├── components/
│   ├── layout/                    # AppShell、SidebarNav
│   ├── log/                       # LiveLogPanel、LogLevelBadge
│   └── ui/                        # Button、EmptyState、StatusBadge 等
├── composables/
│   └── useWailsEvents.ts          # 統一註冊 Wails3 events
├── pages/
│   ├── PanelPage.vue
│   └── LogViewerPage.vue
├── router/
│   ├── index.ts
│   └── routes.ts
├── stores/
│   ├── loading.ts
│   ├── logs.ts
│   └── projects.ts
├── styles/
│   ├── app.css
│   └── style.css
└── types/
```

### 路由

目前路由：

| 路徑 | 頁面 |
| --- | --- |
| `/` | 重新導向 `/panel` |
| `/panel` | 專案面板 |
| `/log-viewer` | 日誌檢視器 |
| `/dev-guide` | 開發指南 |

### Wails3 bindings

前端透過 `frontend/bindings` 呼叫後端服務。以 Panel 為例：

```ts
import {
  InstallProject,
  ListProjects,
  OpenProject,
  StartProject,
  StopProject,
} from "@bindings/gbase/src/service/panel/service.js";
```

呼叫後端服務建議集中在 store 或 API wrapper，不要散落在大量小元件內。

### Wails3 events

前端統一在 `useWailsEvents.ts` 註冊事件：

```ts
Events.On("log", (event) => {
  logs.addLog(event.data);
});
```

即時 log 由 `stores/logs.ts` 保存，頁面與元件只讀取 store 狀態。

## 指令

## 安裝
安裝 `Wails3 CLI` 及整理 `Go` 相依：
```bash
go install github.com/wailsapp/wails/v3/cmd/wails3@latest
go mod tidy
```

### 開發

啟動 Wails3 開發模式，會啟動桌面程式並監聽前後端變更：

```bash
wails3 dev
```

### 打包

建置正式版本：

```bash
wails3 build
```

### 更新建置資源

若有更換主系統 icon，或修改 `build/config.yml` 內的 `info`、`fileAssociations` 等建置資源設定，打包前需執行以下指令更新 build assets：

```bash
wails3 task common:update:build-assets
```

## 常見問題

### 啟動後找不到專案

請確認：

- `.env` 的 `PROJECT_BASE_PATH` 是否正確。
- 該目錄是否存在。
- 每個專案資料夾底下是否有 `project.json`。
- `project.json` 是否為陣列格式。
- 是否有符合目前 `runtime.GOOS` 的 `os` 設定。
- 必要欄位 `title`、`key`、`type`、`start`、`stop` 是否都有值。

### `pm2` 狀態一直判斷為未執行

請確認：

- `project.json` 的 `key` 是否與 PM2 process name 一致。
- 執行環境的 PATH 是否找得到 `pm2`、`node`。
- Linux 桌面啟動程式取得的 PATH 可能比 terminal 少，本專案會在 Unix-like 系統補強常見 PATH，但若使用特殊安裝路徑，仍應自行調整環境。

### 啟動指令 timeout

`start` 與 `stop` 指令應該在完成啟動或停止動作後結束。若指令本身會常駐，請自行背景化。

Windows 範例：

```cmd
start /B app.exe
```

Linux 範例：

```bash
nohup ./app > app.log 2>&1 &
```

