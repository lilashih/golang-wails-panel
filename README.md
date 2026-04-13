# golang-wails-panel

## ▍主要功能
這個專案是一個桌面應用程式，主要提供以下功能：

1. 啟動面板 (Project Panel)  
    可將多個專案放在 [projects](release/projects/) 目錄下，程式會自動讀取該目錄底下的所有專案，並依照每個專案資料夾內的 `project.json` 設定檔來顯示專案資訊與操作按鈕，使用者可以直接在面板上執行啟動、停止、安裝等操作，方便管理多個專案。

    `project.json` 是每個專案資料夾下的設定檔，用來描述專案的基本資訊與操作指令，主要欄位說明如下：  

    | 欄位    | 說明 |
    |---------|---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
    | os      | 作業系統名稱，windows 或 linux。 |
    | title   | 專案名稱，會顯示在啟動面板上。 |
    | type    | 專案類型，目前支援兩種：<br>• pm2：適用於 Node.js 等需用 pm2 管理的專案，相關指令（start、stop、install）通常是 pnpm/yarn/npm 等腳本。<br>• exe：適用於 Windows/Linux 執行檔（.exe）專案，start、stop 會是 Windows 指令，例如啟動 exe 或用 taskkill 關閉。 |
    | key     | 專案唯一識別碼，通常用於內部識別。<br>• pm2：pm2 時需為 pm2 的識別名稱 name。<br>• exe：exe 時需為進程名稱，通常是該執行檔檔名。 |
    | start   | 啟動專案的指令。 |
    | stop    | 停止專案的指令。 |
    | install | 安裝或初始化專案的指令（可選）。 |


    > 不同 type 會對應不同的啟動、停止、安裝方式，讓面板能自動適應各種專案型態，並提供一鍵操作。

    > 注意：`project.json` 的設定是一個包含多筆物件的陣列。若需要針對不同作業系統（Windows / Linux）指定不同的安裝、啟動或停止指令，請在每筆物件中加入 `os` 欄位（可填 `windows` 或 `linux`）。啟動面板會依目前執行的作業系統選取對應的設定並執行該設定內的 `install` / `start` / `stop` 指令；同一個 `project.json` 可同時包含 Windows 與 Linux 的設定範本。

    > 範例如下：可在同一個 `project.json` 以陣列方式同時提供 Windows 與 Linux 的設定，面板會選用與系統相符的那一筆。

    - **`pm2` 專案設定範例：**
        ```json
        [{
            "os": "windows",
            "title": "Nodejs 專案設定範例",
            "type": "pm2",
            "key": "該專案在 pm2 的名稱",
            "start": "pnpm start", // 或 npm start，start 指令需在專案內的 package.json 定義
            "stop": "pnpm stop", // 或 npm stop，stop 指令需在專案內的 package.json 定義
            "install": "pnpm install" // 或 npm install
        }]
        ```

    - **`exe` 專案設定範例：**
        ```json
        [{
            "os": "windows",
            "title": "執行檔專案設定範例",
            "type": "exe",
            "key": "進程名稱，通常是該執行檔檔名",
            "start": "start /b app_windows.exe",
            "stop": "taskkill /IM app_windows.exe /F",
            "install": "app_windows.exe migrate" // 不須安裝留白即可
         }, {
            "os": "linux",
            "title": "執行檔專案設定範例",
            "type": "exe",
            "key": "app_linux",
            "start": "nohup ./app_linux > ./app.log 2>&1 &",
            "stop": "pkill -f app_linux",
            "install": "./app_linux migrate"
        }] 
        ```  
        

2. Log 檢視器  
可自動讀取 `storage/log/` 目錄下的所有日誌檔案，並支援選擇其他額外檔案。即使是內容非常龐大的檔案，也能流暢閱讀與檢索，方便進行日誌分析與除錯。

3. 系統托盤（`Systray`）  
支援最小化至系統托盤並在背景常駐，提供快速操作選單（如顯示/隱藏視窗、開啟日誌資料夾、結束程式等），方便在不占用工作列的情況下管理應用程式。

## ▍目錄架構
```
├── app.go                      # Wails 註冊點
├── go.mod
├── go.sum
├── main.go                     # Go 主程式進入點
├── wails.json                  # Wails 設定檔
├── README.md                   # 專案說明文件
├── scripts/                    # 執行腳本
│   ├── build-linux-appimage.ps1 # 在 Windows 透過 Docker 打包 Linux AppImage
│   ├── build-linux-docker.ps1  # 在 Windows 透過 Docker 打包 Linux 執行檔
│   ├── Dockerfile.linux-appimage # Linux AppImage 打包用 Dockerfile
│   ├── Dockerfile.linux-build  # Linux 打包用 Dockerfile
│   ├── package-appimage.sh     # Linux AppImage 封裝腳本
│   ├── post-build.ps1          # Windows 打包完成後同步至 release
│   └── post-build.sh           # Linux 打包完成後同步至 release
├── frontend/                   # 前端程式，詳細說明請看[README.md](frontend\README.md)
├── windows/                    # Wails 在 Windows 打包時附加的資源與設定檔
├── src/                        # 後端 Go 程式
│   ├── core/                   # 核心功能（cmd, config, def, helper, logger, pm2）
│   └── service/                # 服務層（log_viewer, panel 等）
├── release/                    # 發佈相關資源 
│   ├── golang-wails-panel.exe  # Go 打包後的執行檔
│   ├── golang-wails-panel.AppImage # Linux AppImage
│   └── projects/               # 專案目錄
└── storage/                    # 儲存資料用（log/ 等）
```

## ▍套件
| 套件 | 簡短說明 |
|---|---|
| [wails](https://github.com/wailsapp/wails) | Wails 桌面應用框架 (v2版)。 |
| [fyne.io/systray](https://github.com/fyne-io/systray) | 建立系統托盤。 |
| [caarlos0/env](https://github.com/caarlos0/env) | 解析環境變數到結構體。 |
| [golobby/dotenv](https://github.com/golobby/dotenv) | 載入 `.env` 檔案設定。 |


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
| `PROJECT_BASE_PATH` | 要掃描的專案目錄 | `./projects` | 相對路徑或絕對路徑皆可，例如 `./release/projects` 或 `C:\myDir\projects` |



### ▍日誌
如果只是想把 `log` 訊息寫入檔案，就用 `logger.Log`，如果想讓前端顯示就用 `runtime`。

| Logger | 行為 |
| ----------- | ------------------------------------------------------------ |
| `logger.Log` | `release` 模式僅寫入檔案，其餘模式同時輸出 `console`。 |
| `runtime` | 透過 `Wails` 的 `Logger Adapter` 改寫原有的 `Logger`，除了執行自定義的 `logger.Log` 外，還會透過 `runtime.EventsEmit` 將日誌事件即時傳送至前端。 |

- 產生的日誌檔案路徑：會在專案根目錄的 `storage` 底下自動建立 `log` 目錄（預設為 `./storage/log`）。
- 檔名格式：`log-YYYY-MM-DD.log`，每日產生新檔案（含 `UTF-8 BOM`）。
- 特性：
    - 寫入前自動監測系統日期，跨日就自動換檔。
    - 移除 `ANSI` 顏色控制碼，確保內容為純文字。
    - `logger.Log` 實際上為 `*log.Logger`（標準 `library`）的包裝，因此可以直接使用 `log` 的方法與行為：
        - 可直接使用 `Print/Printf/Println`、`Fatal/Fatalf`、`Panic/Panicf` 等方法。
        - 預設使用 `flags` `log.LstdFlags | log.Lshortfile`（含時間戳與呼叫來源檔案:行號）。  
    - 過期日誌會自動壓縮為 `.gz` 以節省磁碟空間。  
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


## ▍打包  
執行以下指令進行打包：  
 ```bash
wails build
```

`Wails` 的打包流程分為兩個階段：  
 1. 預設輸出目錄 `bin`  
    執行 `wails build` 後，編譯完成的執行檔會先產生在目錄 `bin/` 底下（ `Wails` 的打包檔一定會包一層 `bin/` ）。
    - Windows：bin/golang-wails-panel.exe  
    - Linux：bin/golang-wails-panel  

2. 自動同步至 `release` 目錄  
    專案透過 `wails.json` 中設定的 `postBuildHooks`，在打包完成後自動執行腳本，將檔案複製到 `release/` 目錄，方便發佈使用。  

    不同作業系統會呼叫對應腳本：
    - Windows：scripts/post-build.ps1
    - Linux：scripts/post-build.sh

### 使用 Docker 在 Windows 打包 Linux 執行檔
如果你是在 `Windows` 主機上，但想快速產生 `Linux` 執行檔，專案已提供 `Docker` 版建置流程。此方式會：  
- 在 `Linux` 容器內執行 `wails build`  
- 將產物直接匯出到專案的 `release/` 目錄  
- 建置完成後自動執行 `docker builder prune -af` 清除 `build cache  `
- 不保留最終 `image`  

#### 前置需求
- 已安裝 `Docker Desktop`
- `Docker` 可正常執行 `docker buildx`

#### 基本用法
在專案根目錄執行：
```powershell
.\scripts\build-linux-docker.ps1
```

預設會輸出到：
```text
release/golang-wails-panel
```

若需要同時保留不同架構的產物，請自行指定 `-OutputDir`，例如 `release/linux-arm64`。

#### 可用參數
```powershell
.\scripts\build-linux-docker.ps1 -Arch amd64
.\scripts\build-linux-docker.ps1 -Arch arm64
.\scripts\build-linux-docker.ps1 -Arch amd64 -NoCache
.\scripts\build-linux-docker.ps1 -Arch amd64 -UseWebkit241
```

參數說明：
| 參數 | 說明 |
|---|---|
| `-Arch` | 目標架構，可用 `amd64` 或 `arm64`。 |
| `-OutputDir` | 自訂輸出目錄，預設為 `release`。若要同時保留不同架構產物，請自行指定不同目錄。 |
| `-NoCache` | 不使用 `Docker build cache`。 |
| `-UseWebkit241` | 針對較新的 `Linux` 環境改用 `libwebkit2gtk-4.1-dev`，並加上 `-tags webkit2_41`。 |

#### 直接使用 Dockerfile
如果你不想透過 PowerShell 腳本，也可以直接執行：
```powershell
docker buildx build --pull `
  --build-arg TARGETARCH=amd64 `
  -f .\scripts\Dockerfile.linux-build `
  --output type=local,dest=.\release `
  .
docker builder prune -af
```

若目標 `Linux` 發行版需要 `webkit2gtk-4.1`，可改用：
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

### 使用 Docker 在 Windows 打包 Linux AppImage
如果你的目標是讓 Linux 使用者更接近「下載後直接雙擊執行」，建議改用 `AppImage`。專案已提供獨立的 `Docker` 打包流程，會：

- 在 `Linux` 容器內執行 `wails build`
- 透過 `linuxdeploy` 與 `gtk plugin` 將 `GTK / WebKitGTK` 相關相依一併封裝進 `AppImage`
- 將產物直接匯出到專案的 `release/` 目錄
- 建置完成後自動執行 `docker builder prune -af` 清除 `build cache`

#### 基本用法
在專案根目錄執行：
```powershell
.\scripts\build-linux-appimage.ps1
```

預設會輸出到：
```text
release/golang-wails-panel.AppImage
```

若需要同時保留不同架構的產物，請自行指定 `-OutputDir`，例如 `release/appimage-arm64`。

#### 可用參數
```powershell
.\scripts\build-linux-appimage.ps1 -Arch amd64
.\scripts\build-linux-appimage.ps1 -Arch arm64
.\scripts\build-linux-appimage.ps1 -Arch amd64 -NoCache
.\scripts\build-linux-appimage.ps1 -Arch amd64 -UseWebkit241
```

參數說明：
| 參數 | 說明 |
|---|---|
| `-Arch` | 目標架構，可用 `amd64` 或 `arm64`。 |
| `-OutputDir` | 自訂輸出目錄，預設為 `release`。若要同時保留不同架構產物，請自行指定不同目錄。 |
| `-NoCache` | 不使用 `Docker build cache`。 |
| `-UseWebkit241` | 針對較新的 `Linux` 環境改用 `libwebkit2gtk-4.1-dev`，並加上 `-tags webkit2_41`。 |

#### 直接使用 Dockerfile
如果你不想透過 PowerShell 腳本，也可以直接執行：
```powershell
docker buildx build --pull `
  --build-arg TARGETARCH=amd64 `
  -f .\scripts\Dockerfile.linux-appimage `
  --output type=local,dest=.\release `
  .
docker builder prune -af
```

若目標 `Linux` 發行版需要 `webkit2gtk-4.1`，可改用：
```powershell
docker buildx build --pull `
  --build-arg TARGETARCH=amd64 `
  --build-arg WEBKIT_PKG=libwebkit2gtk-4.1-dev `
  --build-arg WAILS_TAGS=webkit2_41 `
  -f .\scripts\Dockerfile.linux-appimage `
  --output type=local,dest=.\release `
  .
docker builder prune -af
```

#### AppImage 補充說明
- `AppImage` 的目的，是降低目標機器手動安裝 `libwebkit2gtk-4.0.so.37` 這類相依套件的機率，比單一 `ELF` 執行檔更適合發佈給一般 Linux 使用者。
- 某些 Linux 發行版若缺少 `FUSE` 執行環境，雙擊 `AppImage` 仍可能無法啟動；常見需安裝 `libfuse2` 或新版系統的 `libfuse2t64`。
- 專案目前在 `release` 模式會優先以執行檔所在目錄作為相對路徑基準，因此雙擊 `AppImage` 時，`projects` 與 `storage` 會比原本單純依賴目前工作目錄更穩定。



## ▍測試
測試檔案攤平放在 ``test`` 目錄下，不要建立子目錄，然後執行
```
go test -v ./test/...
```


## ▍其他說明
- projects  
    若你在 `projects` 目錄下開發 `Go` 專案，建議將該 `Go` 專案打包後的目錄以「捷徑」方式同步到 `projects` 目錄下，方便啟動面板自動載入。

    > ⚠️ 請務必使用指令建立資料夾捷徑（不可用右鍵產生超連結），僅建議於本機開發時使用，正式部署請勿保留捷徑。

    指令範例如下，必須以 `projects` 目錄為基點產生捷徑路徑：
    ```bash
    mklink /D your_link target_folder
    mklink /D aaa path1\path2\release

