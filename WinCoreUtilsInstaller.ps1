<#
WinCoreUtilsInstaller.ps1

Simple user-scope installer for the WinCoreUtils collection (file, wls, wla, wll).

Usage:
  - Install for current user (no admin required):
      .\WinCoreUtilsInstaller.ps1

  - Force reinstall:
      .\WinCoreUtilsInstaller.ps1 -Force

  - Uninstall for current user:
      .\WinCoreUtilsInstaller.ps1 -Uninstall

This script copies the executables located beside this script into
%LOCALAPPDATA%\Programs\WinCoreUtils and adds/removes that folder from the current user's PATH.
#>

param(
    [switch]$Uninstall,
    [switch]$Force,
    [switch]$Help
)

function Add-ToPath([string]$folder) {
    $targetScope = 'User'
    $folderNorm = $folder.TrimEnd('\')
    $currentPath = [Environment]::GetEnvironmentVariable('Path', $targetScope)
    if (-not $currentPath) { $currentPath = '' }
    $pathParts = $currentPath -split ';' | ForEach-Object { $_.Trim() } | Where-Object { $_ -ne '' }
    foreach ($p in $pathParts) {
        if ($p.TrimEnd('\').ToLower() -eq $folderNorm.ToLower()) {
            Write-Host "Path already contains: $folder"
            return
        }
    }
    if ($currentPath.Trim() -eq '') { $newPath = $folder } else { $newPath = $currentPath.TrimEnd(';') + ";$folder" }
    [Environment]::SetEnvironmentVariable('Path', $newPath, $targetScope)
    Write-Host "Added $folder to user PATH. You may need to open a new shell to see it."
}

function Remove-FromPath([string]$folder) {
    $targetScope = 'User'
    $folderNorm = $folder.TrimEnd('\')
    $currentPath = [Environment]::GetEnvironmentVariable('Path', $targetScope)
    if (-not $currentPath) { return }
    $parts = $currentPath -split ';' | ForEach-Object { $_.Trim() } | Where-Object { $_ -ne '' }
    $remaining = @()
    foreach ($p in $parts) {
        if ($p.TrimEnd('\').ToLower() -ne $folderNorm.ToLower()) { $remaining += $p }
    }
    $newPath = ($remaining -join ';')
    [Environment]::SetEnvironmentVariable('Path', $newPath, $targetScope)
    Write-Host "Removed $folder from user PATH (if present)."
}

function Is-AnyRunning {
    $names = @('wfile', 'wls', 'wla', 'wll')
    foreach ($n in $names) {
        try {
            $p = Get-Process -Name $n -ErrorAction SilentlyContinue
            if ($p) { return $true }
        } catch { }
    }
    return $false
}

function Perform-Install {
    if (Is-AnyRunning) {
        Write-Host "Error: one of the utilities appears to be running. Close them before installing." -ForegroundColor Red
        exit 1
    }
    $destBase = Join-Path $env:LOCALAPPDATA 'Programs\WinCoreUtils'

    # Map destination exe names to their relative source paths inside the repo
    $exeMap = @{
        'wfile.exe' = 'wfile\wfile.exe'
        'wls.exe' = 'wls\wls.exe'
        'wla.exe' = 'wls\wla.exe'
        'wll.exe' = 'wls\wll.exe'
    }

    $foundAny = $false
    foreach ($pair in $exeMap.GetEnumerator()) {
        $src = Join-Path $PSScriptRoot $pair.Value
        if (Test-Path $src) { $foundAny = $true; break }
    }
    if (-not $foundAny) {
        Write-Error "Could not find any of the expected executables in the project subfolders. Ensure the build output exists under 'file/' and 'wls/' and try again."
        exit 1
    }

    if (-not (Test-Path $destBase)) {
        New-Item -ItemType Directory -Path $destBase -Force | Out-Null
    }

    foreach ($pair in $exeMap.GetEnumerator()) {
        $exe = $pair.Key
        $rel = $pair.Value
        $sourceExePath = Join-Path $PSScriptROOT $rel
        if (-not (Test-Path $sourceExePath)) {
            Write-Host "Skipping: $exe (not found at $sourceExePath)"
            continue
        }
        $destExePath = Join-Path $destBase $exe
        if ((Test-Path $destExePath) -and -not $Force) {
            Write-Host "$exe already installed at $destExePath. Use -Force to overwrite."
            continue
        }
        Copy-Item -Path $sourceExePath -Destination $destExePath -Force:$Force
        Write-Host "Copied $exe -> $destExePath"
    }

    Add-ToPath -folder $destBase
    Write-Host "Installation finished. You can run the utilities from any new shell."
}

function Perform-Uninstall {
    if (Is-AnyRunning) {
        Write-Host "Error: one of the utilities appears to be running. Close them before uninstalling." -ForegroundColor Red
        exit 1
    }
    $destBase = Join-Path $env:LOCALAPPDATA 'Programs\WinCoreUtils'
    if (Test-Path $destBase) {
        try {
            Remove-Item -Path $destBase -Recurse -Force
            Write-Host "Removed folder: $destBase"
        } catch {
            Write-Warning ("Failed to remove {0}: {1}" -f $destBase, $_)
        }
    } else {
        Write-Host "No install found at $destBase"
    }
    Remove-FromPath -folder $destBase
    Write-Host "Uninstall complete. You may need to open a new shell to see PATH changes."
}

function Show-Help {
    Write-Host "WinCoreUtils Installer - user-scoped installer for wfile, wls, wla, wll"
    Write-Host ""
    Write-Host "Usage: .\WinCoreUtilsInstaller.ps1 [options]"
    Write-Host "  -Force       Overwrite existing installation when installing"
    Write-Host "  -Uninstall   Remove the current user install"
    Write-Host "  -Help        Show this help message"
    Write-Host ""
    Write-Host "Examples:";
    Write-Host "  .\WinCoreUtilsInstaller.ps1            # Install for current user (no admin required)"
    Write-Host "  .\WinCoreUtilsInstaller.ps1 -Force     # Force reinstall / overwrite"
    Write-Host "  .\WinCoreUtilsInstaller.ps1 -Uninstall # Uninstall for current user"
}

if ($Help) { Show-Help; exit }

if ($Uninstall) { Perform-Uninstall; exit }

Perform-Install
