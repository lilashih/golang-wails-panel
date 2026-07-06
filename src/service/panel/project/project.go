package project

// Project 是 Panel 實際管理的一個專案實例，有些資料會傳給前端顯示。
//
// Install、Start、Stop、CheckRunning 是 builder 依 Config.Type 綁定的後端操作函式，不會輸出到前端 JSON。
type Project struct {
	Config ProjectConfig // 來自 project.json 中符合目前 OS 的設定

	Path string `json:"path"` // 專案所在的實體目錄路徑，程式會自動解析絕對路徑、相對路徑或捷徑

	Running bool `json:"running"` // 目前 runner 判斷出的執行狀態

	CheckRunning func() `json:"-"` // 更新 Running 狀態的流程

	Install func() error `json:"-"` // 執行專案安裝或重建依賴的流程

	Start func() error `json:"-"` // 執行專案啟動流程

	Stop func() error `json:"-"` // 執行專案停止流程
}

// ProjectConfig 是 project.json 中的資料格式。
//
// 同一個 project.json 可以放多筆 ProjectConfig，例如 windows、linux 各一筆。
// 載入時會依目前 runtime.GOOS 選出對應設定。
// Type 決定要使用哪一種 runner，例如 exe、pm2。
type ProjectConfig struct {
	OS string `json:"os"` // OS 指定這筆設定適用的作業系統，例如 windows、linux、darwin

	Title string `json:"title"` // 前端顯示與 log 使用的專案名稱

	Key string `json:"key"` // runner 用來辨識 process 或服務的識別值

	Type string `json:"type"` // Type 指定 runner 類型，例如 exe、pm2

	Start string `json:"start"` // 動專案時要執行的 shell 指令，不可為空

	Stop string `json:"stop"` // 停止專案時要執行的 shell 指令，不可為空

	Install string `json:"install"` // 安裝或重建專案依賴時要執行的 shell 指令，可為空
}
