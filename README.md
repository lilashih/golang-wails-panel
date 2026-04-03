# golang-wails-panel

## ▍主要功能
這個專案是一個桌面應用程式，主要提供以下功能：

1. 啟動面板 (Project Panel)  
    可將多個專案放在 [projects](release/projects/) 目錄下，程式會自動讀取該目錄底下的所有專案，並依照每個專案資料夾內的 `project.json` 設定檔來顯示專案資訊與操作按鈕，使用者可以直接在面板上執行啟動、停止、安裝等操作，方便管理多個專案。

    `project.json` 是每個專案資料夾下的設定檔，用來描述專案的基本資訊與操作指令，主要欄位說明如下：  

    | 欄位    | 說明 |
    |---------|---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
    | title   | 專案名稱，會顯示在啟動面板上。 |
    | type    | 專案類型，目前支援兩種：<br>• pm2：適用於 Node.js 等需用 pm2 管理的專案，相關指令（start、stop、install）通常是 pnpm/yarn/npm 等腳本。<br>• exe：適用於 Windows 執行檔（.exe）專案，start、stop 會是 Windows 指令，例如啟動 exe 或用 taskkill 關閉。 |
    | key     | 專案唯一識別碼，通常用於內部識別。<br>• pm2：pm2 時需為 pam2 的識別名稱 name。<br>• exe：exe 時需為進程名稱，通常是該執行檔檔名。 |
    | start   | 啟動專案的指令。 |
    | stop    | 停止專案的指令。 |
    | install | 安裝或初始化專案的指令（可選）。 |


    > 不同 type 會對應不同的啟動、停止、安裝方式，讓面板能自動適應各種專案型態，並提供一鍵操作。

    - **`pm2` 專案設定範例：**
        ```json
        {
            "title": "Nodejs 專案設定範例",
            "type": "pm2",
            "key": "該專案在 pm2 的名稱",
            "start": "pnpm start", // 或 npm start，start 指令需在專案內的 package.json 定義
            "stop": "pnpm stop", // 或 npm stop，stop 指令需在專案內的 package.json 定義
            "install": "pnpm install" // 或 npm install
        }
        ```

    - **`exe` 專案設定範例：**
        ```json
        {
            "title": "Go 專案設定範例",
            "type": "exe",
            "key": "進程名稱，通常是該執行檔檔名",
            "start": "start /b app.exe",
            "stop": "taskkill /IM app.exe /F",
            "install": "" // 不須安裝留白即可
        }
        ```  
        

2. Log 檢視器  
可自動讀取 `storage/log/` 目錄下的所有日誌檔案，並支援選擇其他額外檔案。即使是內容非常龐大的檔案，也能流暢閱讀與檢索，方便進行日誌分析與除錯。

3. 系統托盤（Systray）  
支援最小化至系統托盤並在背景常駐，提供快速操作選單（如顯示/隱藏視窗、開啟日誌資料夾、結束程式等），方便在不占用工作列的情況下管理應用程式。

## ▍目錄架構
```
├── app.go                      # Go 主程式進入點
├── go.mod
├── go.sum
├── main.go                     # Go 主程式進入點
├── wails.json                  # Wails 設定檔
├── README.md                   # 專案說明文件
├── frontend/                   # 前端程式，詳細說明請看[README.md](frontend\README.md)
├── src/                        # 後端 Go 程式
│   ├── core/                   # 核心功能（cmd, config, helper, logger, pm2）
│   ├── def/                    # 定義檔
│   └── service/                # 服務層（log_viewer, panel 等）
├── release/                    # 發佈相關資源 
│   ├── golang-wails-panel.exe  # Go 打包後的執行檔
│   └── projects/               # 專案目錄
├── storage/                    # 儲存資料用
│   └── log/                    # 本地儲存（log/ 等）
```

## ▍套件
| 套件 | 簡短說明 |
|---|---|
| [wails](https://github.com/wailsapp/wails) | Wails 桌面應用框架 (v2版)。 |
| [fyne.io/systray](https://github.com/fyne-io/systray) | 建立系統托盤。 |
| [caarlos0/env](https://github.com/caarlos0/env) | 解析環境變數到結構體。 |
| [golobby/dotenv](https://github.com/golobby/dotenv) | 載入 `.env` 檔案設定。 |
| [lumberjack](https://github.com/natefinch/lumberjack) | 日誌檔案自動輪替 (log rotation)工具 。 |


## ▍環境變數
所有設定位於 [src/core/config](/src/core/config) ，由 [src/core/config/main.go](/src/core/config/main.go) 自動初始化，並載入 `.env`。 

- 使用範例： 
    ```go
    import (
    "fmt"
    "gbase/src/core/config"
    )

    fmt.Println(config.App.Name)
    ```

### 常用環境變數 (.env)

| 名稱 | 用途 | 預設值（正式環境） | 備註 |
|---|---|---|---|
| `APP_MODE` | 執行模式 (`release` / `debug` / `test`) | `release` | 正式環境請使用 `release` |
| `PROJECT_BASE_PATH` | 要掃描的專案目錄 | `./projects` | 直接填要掃描的目錄，相對路徑或絕對路徑皆可，例如 `./release/projects` 或 `C:\myDir\projects` |



### ▍日誌
如果只是想把 `log` 訊息寫入檔案，就用 `logger.Log`，如果想讓前端顯示就用 `runtime`。

| Logger | 行為 |
| ----------- | ------------------------------------------------------------ |
| `logger.Log` | `release` 模式僅寫入檔案，其餘模式同時輸出 `console`。 |
| `runtime` | 透過 `Wails` 的 `Logger Adapter` 改寫原有的 `Logger`，除了執行自定義的 `logger.Log` 外，還會透過 `runtime.EventsEmit` 將日誌事件即時傳送至前端。 |

- 產生的日誌檔案路徑：會在專案根目錄的 `storage` 底下自動建立 `log` 目錄（預設為 `./storage/log`）。
- 檔名格式：`log-YYYY-MM-DD.log`，每日產生新檔案（含 `UTF-8 BOM`）。
- 特性：
    - 使用 `lumberjack` 做檔案控管，寫入前自動監測系統日期，跨日就自動換檔。
    - 移除 `ANSI` 顏色控制碼，確保內容為純文字。
    - `logger.Log` 實際上為 `*log.Logger`（標準 `library`）的包裝，因此可以直接使用 `log` 的方法與行為：
        - 可直接使用 `Print/Printf/Println`、`Fatal/Fatalf`、`Panic/Panicf` 等方法。
        - 預設使用 `flags` `log.LstdFlags | log.Lshortfile`（含時間戳與呼叫來源檔案:行號）。  
    - 可於 `src/core/logger` 自訂格式、輸出與輪替策略。
- `Wails` 與前端日誌事件：
    - [wails.go](`src/core/logger/wails.go`) 裡定義了 `RegisterCtx` 與 `NewWailsLog()`，後端在啟動時會把這個 `adapter` 註冊給 `Wails`
    - `adapter` 會把每筆 `log` 送給 `logger.Log`，並同時觸發 `event` (`runtime.EventsEmit`) 給前端，讓前端可用來即時顯示 `log`。
- 使用範例： 
    ```go
    import (
        "gbase/src/core/logger"
        "github.com/wailsapp/wails/v2/pkg/runtime"
    )

    // 寫入一般 log
    logger.Log.Println("Server started")

    // wails 的 log，可顯示至前端 (需有 ctx)
    runtime.LogErrorf(s.ctx, "Server started")
    ```



## ▍開發指令

- 安裝  
    安裝 `Wails CLI` 及整理 `Go` 相依：
    ```bash
    go install github.com/wailsapp/wails/v2/cmd/wails@latest
    go mod tidy
    ```

- 執行  
    啟動開發模式（前後端皆會監聽程式碼變動自動重載）：
    ```bash
    wails dev
    ```

- 打包  
    ```bash
    wails build
    ```


## ▍其他說明

- projects  
    若你在 `projects` 目錄下開發 `Go` 專案，建議將該 `Go` 專案打包後的目錄以「捷徑」方式同步到 `projects` 目錄下，方便啟動面板自動載入。

    > ⚠️ 請務必使用指令建立資料夾捷徑（不可用右鍵產生超連結），僅建議於本機開發時使用，正式部署請勿保留捷徑。

    指令範例如下，必須以 `projects` 目錄為基點產生捷徑路徑：
    ```bash
    mklink /D your_link target_folder
    mklink /D aaa path1\path2\release
    ```
