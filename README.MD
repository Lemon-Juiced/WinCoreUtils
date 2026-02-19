# Win Core Utils
Win Core Utils is a small collection of Windows ports of selected GNU/Linux core utilities, so they are easy to build and distribute on Windows.

## Goals:
- Provide familiar CLI tools for Windows users who prefer Unix-style utilities.
- Keep binaries small, dependency-free where practical, and easy to build.

## Included Utilities
- **file**: `wfile.exe` — Basic file type detector. Usage: `wfile.exe <path>`
- **ls**: `wls.exe` — Windows-friendly `ls` replacement with column output. Usage: `wls.exe [directory]`
  - **la**: `wla.exe` - Wrapper that invokes `wls.exe -l -a`; acts like `ls -l` or `la`. Usage: `wla.exe [directory]`
  - **ll**: `wll.exe` — Wrapper that invokes `wls.exe -l`; acts like `ls -l` or `ll`. Usage: `wll.exe [directory]`

## Notes & Limitations
These tools implement a focused subset of the upstream Unix utilities. They aim for convenience on Windows, not full feature parity.

## Quick Installer
A simple user-scoped PowerShell installer is included: `WinCoreUtilsInstaller.ps1`.
The installer copies the executables into `%LOCALAPPDATA%\Programs\WinCoreUtils` and adds that folder to the current user's `PATH` so you can run the utilities from any new terminal.

Usage (run from the folder containing the executables and `WinCoreUtilsInstaller.ps1`):

```powershell
.\WinCoreUtilsInstaller.ps1            # Install for current user (no admin required)
.\WinCoreUtilsInstaller.ps1 -Force     # Force reinstall / overwrite
.\WinCoreUtilsInstaller.ps1 -Uninstall # Remove install and PATH entry for current user
.\WinCoreUtilsInstaller.ps1 -Help      # Shows help message for the installer
```

### Notes:
- The script expects the .exe files in the correct subdirectory relative to `WinCoreUtilsInstaller.ps1` when run.
  - /wfile/wfile.exe
  - /wls/wla.exe
  - /wls/wll.exe
  - /wls/wls.exe
- Since, the .exe files do not come prebundled with the installer they need to be built via the `make` command in their subdirectories.
- PATH changes take effect in new shells only (reopen the terminal to use the utilities).