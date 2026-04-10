param(
    [ValidateSet("amd64", "arm64")]
    [string]$Arch = "amd64",

    [string]$OutputDir = "",

    [switch]$UseWebkit241,

    [switch]$NoCache
)

$ErrorActionPreference = "Stop"

$repoRoot = Split-Path -Parent $PSScriptRoot
$dockerfilePath = Join-Path $PSScriptRoot "Dockerfile.linux-build"

if ([string]::IsNullOrWhiteSpace($OutputDir)) {
    $OutputDir = "release"
}

$resolvedOutputDir = [System.IO.Path]::GetFullPath((Join-Path $repoRoot $OutputDir))
$stagingDir = Join-Path $repoRoot "tmp/docker-linux-build-$Arch"

if (-not (Test-Path -LiteralPath $dockerfilePath)) {
    throw "找不到 Dockerfile：$dockerfilePath"
}

if (Test-Path -LiteralPath $stagingDir) {
    Remove-Item -LiteralPath $stagingDir -Recurse -Force
}

$webkitPkg = "libwebkit2gtk-4.0-dev"
$wailsTags = ""

if ($UseWebkit241) {
    $webkitPkg = "libwebkit2gtk-4.1-dev"
    $wailsTags = "webkit2_41"
}

$arguments = @(
    "buildx", "build",
    "--pull",
    "--build-arg", "TARGETARCH=$Arch",
    "--build-arg", "WEBKIT_PKG=$webkitPkg",
    "--build-arg", "WAILS_TAGS=$wailsTags",
    "-f", $dockerfilePath,
    "--output", "type=local,dest=$stagingDir",
    "."
)

if ($NoCache) {
    $arguments = @("buildx", "build", "--pull", "--no-cache") + $arguments[3..($arguments.Length - 1)]
}

try {
    Push-Location $repoRoot
    & docker @arguments
    if ($LASTEXITCODE -ne 0) {
        throw "docker buildx build 執行失敗，結束碼：$LASTEXITCODE"
    }

    if (-not (Test-Path -LiteralPath $resolvedOutputDir)) {
        New-Item -ItemType Directory -Path $resolvedOutputDir -Force | Out-Null
    }

    Get-ChildItem -LiteralPath $stagingDir -Force | ForEach-Object {
        $targetPath = Join-Path $resolvedOutputDir $_.Name
        if (Test-Path -LiteralPath $targetPath) {
            Remove-Item -LiteralPath $targetPath -Recurse -Force
        }
        Copy-Item -LiteralPath $_.FullName -Destination $targetPath -Recurse -Force
    }
}
finally {
    Pop-Location
    if (Test-Path -LiteralPath $stagingDir) {
        Remove-Item -LiteralPath $stagingDir -Recurse -Force
    }
}

& docker builder prune -af | Out-Null
if ($LASTEXITCODE -ne 0) {
    throw "docker builder prune 執行失敗，結束碼：$LASTEXITCODE"
}

Write-Host "Linux 執行檔已同步至：$resolvedOutputDir"
