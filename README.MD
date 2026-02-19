# Win Core Utils
Win Core Utils is a small collection of Windows ports of selected GNU/Linux core utilities, so they are easy to build and distribute on Windows.

## Goals:
- Provide familiar CLI tools for Windows users who prefer Unix-style utilities.
- Keep binaries small, dependency-free where practical, and easy to build.

## Included Utilities
- **file**: `file.exe` — Basic file type detector. Usage: `file.exe <path>`
- **ls**: `wls.exe` — Windows-friendly `ls` replacement with column output. Usage: `wls.exe [directory]`
  - **la**: `wla.exe` - Wrapper that invokes `wls.exe -l -a`; acts like `ls -l` or `la`. Usage: `wla.exe [directory]`
  - **ll**: `wll.exe` — Wrapper that invokes `wls.exe -l`; acts like `ls -l` or `ll`. Usage: `wll.exe [directory]`

## Notes & Limitations
These tools implement a focused subset of the upstream Unix utilities. They aim for convenience on Windows, not full feature parity.

