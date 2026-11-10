@echo off
setlocal
cd /d "%~dp0"
echo.
echo  SPECTRA v2.5 - Windows compiler
echo  ---------------------------------
echo.
where go >nul 2>&1
if errorlevel 1 (
    echo ERROR: Go is not installed or was not found on your PATH.
    echo Install the Windows Go installer from https://go.dev/dl/
    echo Then close this window, open a fresh Command Prompt,
    echo and double-click this build_windows.bat again.
    pause
    exit /b 1
)
echo Using:
go version
echo.
echo Building SPECTRA_v2.5.exe ...
go build -trimpath -ldflags "-H=windowsgui -s -w" -o "SPECTRA_v2.5.exe" .
if errorlevel 1 (
    echo.
    echo BUILD FAILED. Check error messages above.
    pause
    exit /b 1
)
echo.
echo SUCCESS. Your program is SPECTRA_v2.5.exe in this folder.
echo It launches the app in its own window using Edge or Chrome.
echo.
pause
