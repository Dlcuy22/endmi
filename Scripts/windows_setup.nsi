; =========================================================
; Endmi Windows Installer Script (NSIS)
; =========================================================

!define APP_NAME "Endmi"
!define APP_VERSION "1.0.0"
!define EXE_NAME "endmi.exe"
!define INSTALL_DIR "$PROGRAMFILES64\${APP_NAME}"

; Set compressor
SetCompressor /SOLID lzma

Name "${APP_NAME}"
OutFile "endmi-setup.exe"
InstallDir "${INSTALL_DIR}"
RequestExecutionLevel admin

; ---------------------------------------------------------
; Pages
; ---------------------------------------------------------
Page directory
Page instfiles

UninstPage uninstConfirm
UninstPage instfiles

; ---------------------------------------------------------
; Installation Section
; ---------------------------------------------------------
Section "Install"
    SetOutPath "$INSTDIR"
    
    ; Add the executable
    ; Ensure you have built endmi.exe in the root folder before compiling
    File "..\${EXE_NAME}"
    
    ; Add the PATH setup script
    SetOutPath "$INSTDIR\Scripts"
    File "add_endmi_to_path.bat"
    
    ; Create uninstaller
    WriteUninstaller "$INSTDIR\uninstall.exe"
    
    ; Create Add/Remove Programs entries
    WriteRegStr HKLM "Software\Microsoft\Windows\CurrentVersion\Uninstall\${APP_NAME}" "DisplayName" "${APP_NAME}"
    WriteRegStr HKLM "Software\Microsoft\Windows\CurrentVersion\Uninstall\${APP_NAME}" "UninstallString" "$INSTDIR\uninstall.exe"
    WriteRegStr HKLM "Software\Microsoft\Windows\CurrentVersion\Uninstall\${APP_NAME}" "DisplayVersion" "${APP_VERSION}"
    WriteRegStr HKLM "Software\Microsoft\Windows\CurrentVersion\Uninstall\${APP_NAME}" "Publisher" "Endmi Team"

    DetailPrint "Running PATH setup script..."
    ExecWait '"$INSTDIR\Scripts\add_endmi_to_path.bat"'

SectionEnd

; ---------------------------------------------------------
; Uninstaller Section
; ---------------------------------------------------------
Section "Uninstall"
    Delete "$INSTDIR\${EXE_NAME}"
    Delete "$INSTDIR\Scripts\add_endmi_to_path.bat"
    RMDir "$INSTDIR\Scripts"
    Delete "$INSTDIR\uninstall.exe"
    RMDir "$INSTDIR"

    DeleteRegKey HKLM "Software\Microsoft\Windows\CurrentVersion\Uninstall\${APP_NAME}"
SectionEnd
