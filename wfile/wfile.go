package main

import (
	"fmt"
	"io"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

var common = map[string]string{
	"txt":  "Text File",
	"md":   "Markdown",
	"jpg":  "JPEG Image",
	"jpeg": "JPEG Image",
	"png":  "PNG Image",
	"gif":  "GIF Image",
	"exe":  "Windows Executable",
	"dll":  "Windows DLL",
	"zip":  "ZIP Archive",
	"tar":  "TAR Archive",
	"gz":   "Gzip Archive",
	"pdf":  "PDF Document",
	"docx": "Word Document",
	"xlsx": "Excel Workbook",
	"pptx": "PowerPoint Presentation",
}

/**
 * This function attempts to determine the type of a file by first checking its extension against a common mapping,
 * then using the MIME type detection based on the extension, and finally falling back to content sniffing if necessary.
 *
 * @param path The file path to analyze
 * @return ext The file extension (without dot), typ The detected file type, and err if any error occurs
 */
func detectType(path string) (ext string, typ string, err error) {
	fi, err := os.Stat(path)
	if err != nil {
		return "", "", err
	}
	if fi.IsDir() {
		return "", "directory", nil
	}

	ext = strings.ToLower(filepath.Ext(path))
	ext = strings.TrimPrefix(ext, ".")

	// Try mapping common extensions
	if ext != "" {
		if v, ok := common[ext]; ok {
			return ext, v, nil
		}
		if m := mime.TypeByExtension("." + ext); m != "" {
			return ext, m, nil
		}
	}

	// Fallback: sniff content
	f, err := os.Open(path)
	if err != nil {
		return ext, "", err
	}
	defer f.Close()

	buf := make([]byte, 512)
	n, _ := io.ReadFull(f, buf)
	buf = buf[:n]
	ct := http.DetectContentType(buf)
	if ct == "application/octet-stream" && ext == "" {
		// Unknown
		return ext, "unknown", nil
	}
	return ext, ct, nil
}

/**
 * wfile - A simple file type detector
 *
 * This program takes one or more file paths as command-line arguments and attempts to determine their types.
 *
 * Usage:
 *  wfile <filename1> [filename2 ...]
 *
 * @author: Lemon
 */
func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "Usage: wfile <filename>")
		os.Exit(2)
	}

	for i := 1; i < len(os.Args); i++ {
		path := os.Args[i]
		ext, typ, err := detectType(path)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %s: %v\n", path, err)
			continue
		}
		if ext == "" {
			ext = "(none)"
		}
		fmt.Printf("%s: Extension: %s, File Type: %s\n", path, ext, typ)
	}
}
