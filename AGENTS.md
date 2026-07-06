# AGENTS.md

## 語言規範

* 所有回覆、說明、文件、註解、skill內容，若無特別要求一律使用「繁體中文」
* 程式碼內的註解也必須使用繁體中文
* API 文件（Swagger）、README、說明文件，若無特別要求皆需為繁體中文

---

## 命名規則

### Go 檔案命名

* 使用 snake_case
* 範例：panel_service.go

### Go 測試檔案

* 統一放在 test 目錄下

---

## AI 行為規範（非常重要）

### 必須遵守

* 必須遵循本專案分層架構
* 必須使用既有 helper
* 必須維持風格一致
* 必須寫可維護、可擴展的程式碼

### 禁止行為

* 不要隨意新增新的架構風格
* 不要破壞既有命名規則

---

## 文件與 Skill 撰寫規範

* skill 文件一律放置 .agents/skills/ 中
* 所有 skill 必須使用「繁體中文」撰寫
* 所有說明文件必須使用「繁體中文」
* 語氣需清楚、結構化、易讀
* 可使用條列式、標題分段

---

## 開發規範

### 後端

- 業務服務放在 `src/service/<feature>`。
- 共用能力放在 `src/core`。
- Wails3 service 統一由 `src/app/services.go` 註冊。
- 不要在 `main.go` 堆疊業務邏輯。
- 不要在 `src/core` 直接依賴特定頁面。
- 需要跨平台行為時，優先使用既有 `core/cmd`、`helper`、runner 分層。
- log 請優先使用 `logger` package，不要直接散落 `fmt.Println`。

### 前端

- 頁面放在 `frontend/src/pages`。
- 可重用元件放在 `frontend/src/components`。
- 跨頁狀態放在 `frontend/src/stores`。
- Wails event 統一由 `useWailsEvents.ts` 註冊。
- 後端 bindings 呼叫優先集中在 store 或 API wrapper。
- UI 風格以管理工具為主，保持清楚、穩定、可長時間使用。
- 常見操作圖示優先使用 `lucide-vue-next`。
- 樣式色彩優先使用 `frontend/src/styles/style.css` 中的語意 token。

### 文件

- 文件預設使用繁體中文。
- README、API 文件、開發文件若無特別需求都使用繁體中文。
- 新增規畫或設計說明放在 `docs`。

---

## 原則總結

AI 在本專案中的角色應該是：

> 嚴格遵守規範的工程師，而不是自由發揮的生成器

所有產出的內容，必須：

* 可讀
* 可維護
* 一致
* 符合專案既有設計
