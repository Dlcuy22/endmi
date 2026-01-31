# Makefile
# Build and development automation for Endmi.
#
# Targets:
#   - run: runs the application locally
#   - build-amd64: builds the binary for windows amd64
#   - build-all: builds binaries for Linux, Windows, and macOS (multiple architectures)

BINARY_NAME=endmi

# Detect OS
ifeq ($(OS),Windows_NT)
	# Windows (PowerShell is used for cross-platform environment variable setting)
	RM = powershell -Command "if (Test-Path bin) { Remove-Item -Recurse -Force bin }; if (Test-Path $(BINARY_NAME).exe) { Remove-Item -Force $(BINARY_NAME).exe }"
	MKDIR = powershell -Command "if (!(Test-Path bin)) { New-Item -ItemType Directory bin }"
	MV_INSTALLER = powershell -Command "Move-Item -Path Scripts/endmi-setup.exe -Destination bin/endmi-setup.exe -Force"
	MAKENSIS = makensis
	# Helper for cross-compiling on Windows
	GO_BUILD = powershell -Command "$$env:GOOS='$(1)'; $$env:GOARCH='$(2)'; go build -o $(3) ."
else
	# Linux / macOS
	RM = rm -rf bin/ $(BINARY_NAME).exe
	MKDIR = mkdir -p bin
	MV_INSTALLER = mv Scripts/endmi-setup.exe bin/
	MAKENSIS = makensis
	GO_BUILD = GOOS=$(1) GOARCH=$(2) go build -o $(3) .
endif

run:
	@echo Running Endmi...
	go run .

build-amd64:
	@echo Building local Windows binary...
	go build -o $(BINARY_NAME).exe .

build-all: 
	@echo =====================================
	@echo Starting full cross-platform build...
	@echo =====================================
	@$(MAKE) build-windows
	@$(MAKE) build-linux
	@$(MAKE) build-macos
	@echo =====================================
	@echo Build complete. Check the bin/ folder.
	@echo =====================================

build-windows:
	@echo [Windows] Preparing directory...
	@$(MKDIR)
	@echo [Windows] Building standalone binary (amd64)...
	@$(call GO_BUILD,windows,amd64,bin/$(BINARY_NAME)-windows-amd64.exe)
	@echo [Windows] Building installer prerequisite...
	@go build -o $(BINARY_NAME).exe .
	@echo [Windows] Compiling NSIS installer...
	@$(MAKENSIS) Scripts/windows_setup.nsi
	@echo [Windows] Moving installer to bin/...
	@$(MV_INSTALLER)

build-linux:
	@echo [Linux] Preparing directory...
	@$(MKDIR)
	@echo [Linux] Building amd64...
	@$(call GO_BUILD,linux,amd64,bin/$(BINARY_NAME)-linux-amd64)
	@echo [Linux] Building arm64...
	@$(call GO_BUILD,linux,arm64,bin/$(BINARY_NAME)-linux-arm64)
	@echo [Linux] Building arm...
	@$(call GO_BUILD,linux,arm,bin/$(BINARY_NAME)-linux-arm)

build-macos:
	@echo [MacOS] Preparing directory...
	@$(MKDIR)
	@echo [MacOS] Building amd64...
	@$(call GO_BUILD,darwin,amd64,bin/$(BINARY_NAME)-darwin-amd64)
	@echo [MacOS] Building arm64...
	@$(call GO_BUILD,darwin,arm64,bin/$(BINARY_NAME)-darwin-arm64)

clean:
	@echo Cleaning up build artifacts...
	@$(RM)



