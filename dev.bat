@echo off
echo Starting AI Knowledge Base Middleware in development mode...

REM Check if air is installed
where air >nul 2>&1
if %errorlevel% neq 0 (
    echo Air not found. Installing air...
    go install github.com/air-verse/air@latest
    if %errorlevel% neq 0 (
        echo Failed to install air. Please install it manually.
        pause
        exit /b 1
    )
    echo Air installed successfully!
)

REM Check if tmp directory exists, create if not
if not exist "tmp" (
    mkdir tmp
    echo Created tmp directory
)

REM Start air for hot reload
echo Starting hot reload server...
set AIR_RUNNER=cmd
cmd /C "set AIR_RUNNER=cmd && air"