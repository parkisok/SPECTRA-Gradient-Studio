SPECTRA v2.5 - Build your own Windows application

WHAT IS INCLUDED
  index.html         - SPECTRA v2.5, all UI and editor code
  main.go            - Go desktop launcher; embeds index.html into the EXE
  go.mod             - Go project declaration
  build_windows.bat  - double-click this to compile

ONE-TIME SETUP
  1. Install Go for Windows from https://go.dev/dl/
  2. Make sure Microsoft Edge or Google Chrome is installed.
  3. Unzip this folder to a location where you can edit files.

COMPILE
  Double-click build_windows.bat.
  Or open Command Prompt in this folder and run:
    go build -trimpath -ldflags "-H=windowsgui -s -w" -o SPECTRA_v2.5.exe .

RUN / SHARE
  Double-click SPECTRA_v2.5.exe. You can send that single EXE to
  another Windows 10/11 x64 PC with Edge or Chrome installed.
  You do NOT need to send the .html or Go source along with the EXE.

MAKE YOUR OWN CHANGES
  1. Edit index.html and save it.
  2. Double-click build_windows.bat again to recompile.
  3. Share the new EXE.

WHAT KIND OF DESKTOP APP IS THIS?
  The EXE embeds the full HTML studio, extracts it to a per-user
  AppData folder, and opens an Edge/Chrome desktop app window.
  It doesn't launch a normal web browser tab. The editor works offline.
  It uses the browser's rendering engine, so Edge or Chrome is required.

USER DATA
  The app profile/cache and any local sprite-library data are stored
  in each user's local Windows profile; they are not included in the EXE
  you share. Replacing the EXE in the same version uses the same profile.

TRUST & WINDOWS SMARTSCREEN
  This executable will not be digitally signed by Go. Windows may warn
  users about an unknown publisher. Signing/distribution is separate from
  compilation.
