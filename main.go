// SPECTRA Desktop — self-contained Windows desktop launcher.
// The entire offline editor is embedded in this executable; no local server,
// separate HTML download, Python, npm, or installation is necessary.
//
// Rendering is provided by a separately profiled local Microsoft Edge (or Chrome)
// in --app mode, which displays only the desktop editor in its own window.
// Edge is preinstalled on supported Windows 10/11 machines.
package main

import (
	"bytes"
	_ "embed"
	"errors"
	"fmt"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
	"unsafe"
)

//go:embed index.html
var editorHTML []byte

const appTitle = "SPECTRA — Dual Sprite Cutout Lab v2.5"
const alreadyExists = syscall.Errno(183)

var (
	kernel32       = syscall.NewLazyDLL("kernel32.dll")
	user32         = syscall.NewLazyDLL("user32.dll")
	createMutexW   = kernel32.NewProc("CreateMutexW")
	closeHandle    = kernel32.NewProc("CloseHandle")
	messageBoxW    = user32.NewProc("MessageBoxW")
	enumWindows    = user32.NewProc("EnumWindows")
	getWindowTextW = user32.NewProc("GetWindowTextW")
	isVisible      = user32.NewProc("IsWindowVisible")
	showWindow     = user32.NewProc("ShowWindow")
	setForeground  = user32.NewProc("SetForegroundWindow")
)

func main() {
	// Mutex remains owned while our dedicated browser window is open.
	mutexName, _ := syscall.UTF16PtrFromString("Local\\SPECTRA_Dual_Sprite_v25")
	mutex, _, lastErr := createMutexW.Call(0, 0, uintptr(unsafe.Pointer(mutexName)))
	if mutex == 0 {
		fail(fmt.Errorf("could not initialize the application: %v", lastErr))
		return
	}
	defer closeHandle.Call(mutex)
	if errors.Is(lastErr, alreadyExists) {
		focusExistingWindow()
		return
	}

	if err := startStudio(); err != nil {
		fail(err)
	}
}

func fail(err error) {
	textPtr, _ := syscall.UTF16PtrFromString(err.Error())
	titlePtr, _ := syscall.UTF16PtrFromString("SPECTRA Desktop")
	messageBoxW.Call(0, uintptr(unsafe.Pointer(textPtr)), uintptr(unsafe.Pointer(titlePtr)), 0x10)
}

func focusExistingWindow() {
	cb := syscall.NewCallback(func(hwnd uintptr, lparam uintptr) uintptr {
		visible, _, _ := isVisible.Call(hwnd)
		if visible == 0 {
			return 1
		}
		buf := make([]uint16, 512)
		n, _, _ := getWindowTextW.Call(hwnd, uintptr(unsafe.Pointer(&buf[0])), uintptr(len(buf)))
		if n > 0 {
			title := syscall.UTF16ToString(buf[:n])
			if strings.Contains(title, appTitle) {
				showWindow.Call(hwnd, 9) // SW_RESTORE
				setForeground.Call(hwnd)
				return 0
			}
		}
		return 1
	})
	enumWindows.Call(cb, 0)
}

func startStudio() error {
	local := os.Getenv("LOCALAPPDATA")
	if local == "" {
		local = os.TempDir()
	}
	appDir := filepath.Join(local, "SPECTRA Dual Sprite v2.5")
	if err := os.MkdirAll(appDir, 0700); err != nil {
		return fmt.Errorf("could not create application data folder: %w", err)
	}
	htmlPath := filepath.Join(appDir, "studio.html")
	current, readErr := os.ReadFile(htmlPath)
	if readErr != nil || !bytes.Equal(current, editorHTML) {
		tmp, err := os.CreateTemp(appDir, "studio-*.html")
		if err != nil {
			return fmt.Errorf("could not prepare editor: %w", err)
		}
		tmpName := tmp.Name()
		defer os.Remove(tmpName)
		if _, err = tmp.Write(editorHTML); err != nil {
			tmp.Close()
			return err
		}
		if err = tmp.Close(); err != nil {
			return err
		}
		_ = os.Remove(htmlPath)
		if err = os.Rename(tmpName, htmlPath); err != nil {
			return fmt.Errorf("could not update editor file: %w", err)
		}
	}

	browser, err := findBrowser()
	if err != nil {
		return err
	}
	htmlURL := url.URL{Scheme: "file", Path: "/" + filepath.ToSlash(htmlPath)}
	profileDir := filepath.Join(appDir, "Browser Profile")
	// Isolated profile prevents other Edge tabs/windows from being reused.
	args := []string{
		"--user-data-dir=" + profileDir,
		"--no-first-run",
		"--no-default-browser-check",
		"--disable-extensions",
		"--window-size=1430,960",
		"--app=" + htmlURL.String(),
	}
	cmd := exec.Command(browser, args...)
	// No browser output is needed; no console window is created at startup.
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("unable to open SPECTRA window: %w", err)
	}
	// Hold the single-instance mutex until the app window closes.
	if err := cmd.Wait(); err != nil {
		return fmt.Errorf("SPECTRA window closed unexpectedly: %w", err)
	}
	return nil
}

func findBrowser() (string, error) {
	var candidates []string
	for _, base := range []string{
		os.Getenv("ProgramFiles(x86)"),
		os.Getenv("ProgramFiles"),
		os.Getenv("LOCALAPPDATA"),
	} {
		if base == "" {
			continue
		}
		candidates = append(candidates,
			filepath.Join(base, "Microsoft", "Edge", "Application", "msedge.exe"),
			filepath.Join(base, "Google", "Chrome", "Application", "chrome.exe"),
		)
	}
	for _, p := range candidates {
		if stat, err := os.Stat(p); err == nil && !stat.IsDir() {
			return p, nil
		}
	}
	for _, n := range []string{"msedge.exe", "chrome.exe"} {
		if p, err := exec.LookPath(n); err == nil {
			return p, nil
		}
	}
	return "", errors.New("SPECTRA Desktop needs Microsoft Edge or Google Chrome to display its window. Install or enable Edge, then start SPECTRA again. The app itself works completely offline.")
}
