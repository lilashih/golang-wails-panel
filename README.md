# golang-wails-panel


## 主要功能
這個專案是一個桌面應用程式，主要提供兩大功能：

1. 啟動面板（Project Panel）
可將多個專案放在 [projects](release/projects/) 目錄下，程式會自動讀取該目錄底下的所有專案，並依照每個專案資料夾內的 `project.json` 設定檔來顯示專案資訊與操作按鈕，使用者可以直接在面板上執行啟動、停止、安裝等操作，方便管理多個專案。

`project.json` 是每個專案資料夾下的設定檔，用來描述專案的基本資訊與操作指令，主要欄位說明如下：

**pm2 專案設定範例：**
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

**exe 專案設定範例：**
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

| 欄位    | 說明 |
|---------|---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| title   | 專案名稱，會顯示在啟動面板上。 |
| type    | 專案類型，目前支援兩種：<br>• pm2：適用於 Node.js 等需用 pm2 管理的專案，相關指令（start、stop、install）通常是 pnpm/yarn/npm 等腳本。<br>• exe：適用於 Windows 執行檔（.exe）專案，start、stop 會是 Windows 指令，例如啟動 exe 或用 taskkill 關閉。 |
| key     | 專案唯一識別碼，通常用於內部識別。<br>• pm2：pm2 時需為 pam2 的識別名稱 name。<br>• exe：exe 時需為進程名稱，通常是該執行檔檔名。 |
| start   | 啟動專案的指令。 |
| stop    | 停止專案的指令。 |
| install | 安裝或初始化專案的指令（可選）。 |

> 不同 type 會對應不同的啟動、停止、安裝方式，讓面板能自動適應各種專案型態，並提供一鍵操作。

2. Log 檢視器
可自動讀取 [log](storage/log/) 目錄下的所有日誌檔案，並支援選擇其他額外檔案。即使是內容非常龐大的檔案，也能流暢閱讀與檢索，方便進行日誌分析與除錯。


## 目錄架構
本專案目錄結構如下：
```
├── app.go                  # Go 主程式進入點
├── go.mod                  # Go modules 設定
├── go.sum                  # Go modules 相依性雜湊
├── main.go                 # Go 主程式進入點
├── wails.json              # Wails 設定檔
├── README.md               # 專案說明文件
├── frontend/               # 前端程式，詳細說明請看[README.md](frontend\README.md)
├── src/                    # 後端 Go 程式
│   ├── core/               # 核心功能（cmd, config, helper, logger, pm2）
│   ├── def/                # 定義檔
│   └── service/            # 服務層（log_viewer, panel 等）
├── release/                # 發佈相關資源 
│   ├── bin/                # Go 打包後的執行檔
│   └── projects/           # 專案捷徑目錄
├── storage/                # 儲存資料用
│   └── log/                # 系統 log 檔案
```

## 開發指令

### 安裝
安裝 Wails CLI 及整理 Go 相依：

```bash
go install github.com/wailsapp/wails/v2/cmd/wails@latest
go mod tidy
```

### 執行
啟動開發模式（前後端程式碼變動會自動重載）：

```bash
wails dev
```

### 打包
編譯專案為可執行檔：

```bash
wails build
```

## 其他說明

### projects
若你在 `projects` 目錄下開發 Go 專案，建議將 build 或 release 目錄以「捷徑」方式同步到 `projects` 目錄下，方便啟動面板自動載入。

> ⚠️ 請務必使用指令建立資料夾捷徑（不可用右鍵產生超連結），僅建議於本機開發時使用，正式部署請勿保留捷徑。

指令範例如下，必須以 `projects` 目錄為基點產生捷徑路徑：

```bash
mklink /D your_link target_folder
mklink /D aaa path1\path2\release
```
