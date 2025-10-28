# Frontend Project Guide

This document describes the `frontend` directory, including project structure, main features, used packages, and how to install and start the project.

## ▍Directory Structure

```
frontend/
├── index.html                # Entry HTML file
├── package.json              # Frontend dependency management
├── package.json.md5          # MD5 checksum for package.json
├── README.md                 # This documentation file
├── tailwind.config.js        # Tailwind CSS configuration
├── vite.config.js            # Vite configuration
├── src/                      # Main frontend code
│   ├── App.vue               # Vue entry component
│   ├── main.js               # Frontend entry point
│   ├── asset/                # Static assets (fonts, images)
│   ├── component/            # Vue components
│   │   ├── layout/           # Layout-related components
│   │   ├── page/             # Page components
│   │   └── widget/           # Widget components
│   ├── css/                  # Style files
│   ├── router/               # Router configuration
│   └── store/                # State management
├── wailsjs/                  # Wails auto-generated JS files
│   ├── go/                   # Go API bindings
│   └── runtime/              # Wails API
```

## ▍Features
- Uses Vue 3 as the frontend framework
- Uses Vite for development and build
- Tailwind CSS for rapid UI design
- Interacts with the Wails backend

## ▍Packages
- [Vue 3](https://vuejs.org/)
- [Vite](https://vitejs.dev/)
- [Tailwind CSS](https://tailwindcss.com/)
- [Wails](https://wails.io/)

## ▍Installation & Startup
You can run these commands in the frontend directory, or simply run `wails dev` in the project root.

1. Enter the frontend directory:
    ```bash
    cd frontend
    ```
2. Install dependencies:
    ```bash
    npm install
    ```
3. Start the development server:
    ```bash
    npm run dev
    ```
