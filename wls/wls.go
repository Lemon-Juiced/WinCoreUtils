package main

import (
	"C" // Required for cgo to compile, even though we don't actually use it in this file

	"flag"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"golang.org/x/sys/windows"
	"golang.org/x/term"
)

/**
 * wls - Windows version of ls
 *
 * This is a simple implementation of the Unix `ls` command for Windows.
 * It lists the contents of a directory in columns, similar to the default behavior of `ls` on Unix systems.
 * Additionally, it supports a long listing format with the `-l` flag that shows permissions, size, and modification time.
 * It also colors the output based on file types:
 * - Directories: Blue
 * - Executables: Green
 * - Archives: Red
 * - Images/Videos: Magenta
 * - Audio files: Cyan
 *
 * Usage:
 *  wls [directory]
 *  wls -a [directory]  # Include hidden files
 *  wls -l [directory]  # Long listing format
 *
 * If no directory is specified, it lists the contents of the current directory.
 *
 * @author: Lemon
 */
func main() {
	longFlag := flag.Bool("l", false, "long listing")
	allFlag := flag.Bool("a", false, "include hidden files")
	// Expand combined short flags (e.g. -la -> -l -a) so `-la` works like many shells
	if len(os.Args) > 1 {
		os.Args = append([]string{os.Args[0]}, expandCombinedFlags(os.Args[1:])...)
	}
	flag.Parse()

	// Try to enable ANSI escape processing on Windows consoles
	_ = enableANSI()

	dir := "."
	if args := flag.Args(); len(args) > 0 {
		dir = args[0]
	}

	entries, err := os.ReadDir(dir)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}

	names := make([]string, 0, len(entries))
	for _, e := range entries {
		name := e.Name()
		// skip dotfiles and Windows hidden files unless -a provided
		if !*allFlag {
			if strings.HasPrefix(name, ".") {
				continue
			}
			full := filepath.Join(dir, name)
			if isHidden(full) {
				continue
			}
		}
		names = append(names, name)
	}

	// If -a requested, include the special entries "." and ".." (like Unix `ls -a`).
	if *allFlag {
		hasDot := false
		hasDotDot := false
		for _, n := range names {
			if n == "." {
				hasDot = true
			}
			if n == ".." {
				hasDotDot = true
			}
		}
		if !hasDot {
			names = append(names, ".")
		}
		if !hasDotDot {
			names = append(names, "..")
		}
	}

	sort.Strings(names)

	// Determine terminal width
	width := 80
	if w, _, err := term.GetSize(int(os.Stdout.Fd())); err == nil && w > 0 {
		width = w
	}

	// If long listing requested, print one entry per line with details
	if *longFlag {
		for _, name := range names {
			path := filepath.Join(dir, name)
			fi, err := os.Stat(path)
			if err != nil {
				fmt.Fprintln(os.Stderr, "error:", err)
				continue
			}
			perms := fi.Mode().String()
			size := fi.Size()
			mtime := fi.ModTime().Format("Jan _2 15:04")
			fmt.Printf("%s %8d %s %s\n", perms, size, mtime, colorName(name, path, fi))
		}
		return
	}

	// Layout in columns like unix `ls` (vertical filling)
	maxLen := 0
	for _, n := range names {
		if l := len(n); l > maxLen {
			maxLen = l
		}
	}
	if maxLen == 0 {
		return
	}

	colWidth := maxLen + 2
	cols := width / colWidth
	if cols < 1 {
		cols = 1
	}
	rows := (len(names) + cols - 1) / cols

	for r := 0; r < rows; r++ {
		for c := 0; c < cols; c++ {
			idx := c*rows + r
			if idx >= len(names) {
				continue
			}
			name := names[idx]
			path := filepath.Join(dir, name)
			fi, _ := os.Stat(path)
			// Pad all but the last printed column. When padding, print colored name then spaces to maintain alignment.
			if c == cols-1 || c*rows+r+rows >= len(names) {
				fmt.Print(colorName(name, path, fi))
			} else {
				padded := colWidth - len(name)
				if padded < 0 {
					padded = 0
				}
				fmt.Print(colorName(name, path, fi))
				fmt.Print(strings.Repeat(" ", padded))
			}
		}
		fmt.Println()
	}
}

// expandCombinedFlags turns combined short flags like `-la` into `-l -a`.
// It leaves long flags (`--foo`) and non-flag arguments unchanged.
func expandCombinedFlags(args []string) []string {
	out := make([]string, 0, len(args))
	for _, a := range args {
		if !strings.HasPrefix(a, "-") || strings.HasPrefix(a, "--") || len(a) == 2 {
			out = append(out, a)
			continue
		}
		// short combined flags, split into separate `-x` entries
		for i := 1; i < len(a); i++ {
			out = append(out, "-"+string(a[i]))
		}
	}
	return out
}

/**
 * Enables virtual terminal processing so ANSI colors work on Windows consoles.
 */
func enableANSI() error {
	h := windows.Handle(os.Stdout.Fd())
	var mode uint32
	if err := windows.GetConsoleMode(h, &mode); err != nil {
		return err
	}
	const enableVirtualTerminalProcessing = 0x0004
	mode |= enableVirtualTerminalProcessing
	return windows.SetConsoleMode(h, mode)
}

// isHidden returns true if the file has the Windows hidden attribute.
func isHidden(path string) bool {
	p, err := windows.UTF16PtrFromString(path)
	if err != nil {
		return false
	}
	attrs, err := windows.GetFileAttributes(p)
	if err != nil {
		return false
	}
	return attrs&windows.FILE_ATTRIBUTE_HIDDEN != 0
}

/**
 * Returns the name wrapped in ANSI color codes based on file type
 *
 * @param name The file or directory name to color
 * @param path The full path to the file or directory (used for stat)
 * @param fi os.FileInfo for the file or directory (can be nil if not available)
 * @return The name wrapped in ANSI color codes if applicable, otherwise the original name
 */
func colorName(name, path string, fi os.FileInfo) string {
	if fi == nil {
		return name
	}
	if fi.IsDir() {
		return blue(name)
	}
	ext := strings.ToLower(strings.TrimPrefix(filepath.Ext(name), "."))
	if isExecutable(ext) {
		return green(name)
	}
	if isArchive(ext) {
		return red(name)
	}
	if isImageOrVideo(ext) {
		return magenta(name)
	}
	if isAudio(ext) {
		return cyan(name)
	}
	return name
}

/**
 * Helper functions to determine file types based on extensions for coloring
 *
 * @param ext The file extension (without dot)
 * @return bool indicating if the file is of a certain type
 */
func isExecutable(ext string) bool {
	switch ext {
	case "exe", "bat", "cmd", "com", "ps1":
		return true
	}
	return false
}

/**
 * Checks if the file extension corresponds to a common archive format
 *
 * @param ext The file extension (without dot)
 * @return bool indicating if the file is an archive
 */
func isArchive(ext string) bool {
	switch ext {
	case "zip", "tar", "gz", "tgz", "7z", "rar":
		return true
	}
	return false
}

/**
 * Checks if the file extension corresponds to a common image or video format
 *
 * @param ext The file extension (without dot)
 * @return bool indicating if the file is an image or video
 */
func isImageOrVideo(ext string) bool {
	switch ext {
	case "jpg", "jpeg", "png", "gif", "bmp", "webp", "mp4", "mkv", "mov", "avi":
		return true
	}
	return false
}

/**
 * Checks if the file extension corresponds to a common audio format
 *
 * @param ext The file extension (without dot)
 * @return bool indicating if the file is an audio file
 */
func isAudio(ext string) bool {
	switch ext {
	case "mp3", "wav", "flac", "aac", "ogg":
		return true
	}
	return false
}

/**
 * Wraps a string in ANSI color codes for the given color code
 *
 * @param code The ANSI color code (e.g. 31 for red, 32 for green)
 * @param s The string to wrap in color codes
 * @return The input string wrapped in ANSI color codes
 */
func colorWrap(code int, s string) string {
	return fmt.Sprintf("\x1b[%dm%s\x1b[0m", code, s)
}

func blue(s string) string    { return colorWrap(34, s) }
func green(s string) string   { return colorWrap(32, s) }
func red(s string) string     { return colorWrap(31, s) }
func magenta(s string) string { return colorWrap(35, s) }
func cyan(s string) string    { return colorWrap(36, s) }
