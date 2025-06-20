# PowerShell script for development with hot reload
# Usage: .\scripts\dev.ps1

Write-Host "Starting AI Knowledge Base Middleware in development mode..." -ForegroundColor Green

# Check if air is installed
$airPath = Get-Command air -ErrorAction SilentlyContinue
if (-not $airPath) {
    Write-Host "Air not found. Installing air..." -ForegroundColor Yellow
    go install github.com/air-verse/air@latest
    if ($LASTEXITCODE -ne 0) {
        Write-Host "Failed to install air. Please install it manually." -ForegroundColor Red
        exit 1
    }
    Write-Host "Air installed successfully!" -ForegroundColor Green
}

# Check if tmp directory exists, create if not
if (-not (Test-Path "tmp")) {
    New-Item -ItemType Directory -Path "tmp" | Out-Null
    Write-Host "Created tmp directory" -ForegroundColor Yellow
}

# Start air for hot reload
Write-Host "Starting hot reload server..." -ForegroundColor Cyan
air 