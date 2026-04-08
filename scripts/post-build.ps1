param(
    [Parameter(Mandatory = $true, Position = 0)]
    [string]$SourceBinary
)

$repoRoot = Split-Path -Parent $PSScriptRoot

if ([System.IO.Path]::IsPathRooted($SourceBinary)) {
    $sourcePath = [System.IO.Path]::GetFullPath($SourceBinary)
} else {
    $sourcePath = [System.IO.Path]::GetFullPath((Join-Path $repoRoot $SourceBinary))
}

if (-not (Test-Path -LiteralPath $sourcePath)) {
    Write-Error "Build output not found: $sourcePath"
    exit 1
}

$targetDir = Join-Path $repoRoot "release"
$targetPath = Join-Path $targetDir ([System.IO.Path]::GetFileName($sourcePath))

if (-not (Test-Path -LiteralPath $targetDir)) {
    New-Item -ItemType Directory -Path $targetDir | Out-Null
}

Copy-Item -LiteralPath $sourcePath -Destination $targetPath -Force
Write-Host "Copied build output to $targetPath"
