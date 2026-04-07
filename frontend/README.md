# 前端專案說明

本文件為 `frontend` 目錄的說明，包含專案結構、主要功能、使用到的套件與安裝啟動方式。

## ▍目錄結構

```
frontend/
├── index.html                # 入口 HTML 文件
├── package.json              # 前端依賴管理
├── package.json.md5          # package.json 的 MD5 校驗
├── README.md                 # 本說明文件
├── tailwind.config.js        # Tailwind CSS 設定
├── vite.config.js            # Vite 設定
├── src/                      # 前端主要程式
│   ├── App.vue               # Vue 入口組件
│   ├── main.js               # 前端進入點
│   ├── asset/                # 靜態資源（字型、圖片）
│   ├── component/            # Vue 組件
│   │   ├── layout/           # 佈局相關組件
│   │   ├── page/             # 頁面組件
│   │   └── widget/           # 小元件
│   ├── css/                  # 樣式檔案
│   ├── router/               # 路由設定
│   └── store/                # 狀態管理
├── wailsjs/                  # Wails 自動生成的 JS 檔案
│   ├── go/                   # Go 端 API 對應
│   └── runtime/              # Wails API
```

## ▍主要功能
- 使用 `Vue 3` 作為前端框架
- 使用 `Vite` 作為開發與建構工具
- `Tailwind CSS` 快速設計 UI
- 與 `Wails` 後端進行互動

## ▍使用到的主要套件
- [Vue 3](https://vuejs.org/)
- [Vite](https://vitejs.dev/)
- [Tailwind CSS](https://tailwindcss.com/)
- [Wails](https://wails.io/)


## ▍安裝與啟動
可在 `frontend` 目錄執行，或直接在專案目錄執行 `wails dev`
1. 進入 frontend 目錄：
    ```bash
    cd frontend
    ```
2. 安裝依賴：
    ```bash
    npm install
    ```
3. 啟動開發伺服器：
    ```bash
    npm run dev
    ```
