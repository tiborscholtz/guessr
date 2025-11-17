package main

import (
	"encoding/hex"
	"fmt"
	"io"
	"os"

	"github.com/fatih/color"
	"github.com/rodaine/table"
)

// (magic numbers)
var signatures = map[string]string{
	"JPEG": "FFD8FF",
	"PNG":  "89504E47",
	"PDF":  "25504446",
	"ZIP":  "504B0304",
	"RAR":  "526172211A07",
	"7Z":   "377ABCAF271C",
	"GZIP": "1F8B08",
	"BMP":  "424D",
	"MP4":  "0000002066747970",
	"FLAC": "664C6143",
	"ICO":  "00000100",
}

func readFirstBytes(path string, n int) ([]byte, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	buf := make([]byte, n)
	_, err = io.ReadFull(f, buf)
	if err != nil {
		return nil, err
	}
	return buf, nil
}

// Return a percentage based result
func matchPercentage(fileBytes []byte, sigHex string) int {
	sigBytes, _ := hex.DecodeString(sigHex)

	if len(fileBytes) < len(sigBytes) {
		return 0
	}

	matches := 0
	for i := 0; i < len(sigBytes); i++ {
		if fileBytes[i] == sigBytes[i] {
			matches++
		}
	}

	return (matches * 100) / len(sigBytes)
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: fileguess <filename>")
		return
	}

	filename := os.Args[1]

	// Read first 8 bytes for now
	bytes, err := readFirstBytes(filename, 8)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Printf("Read bytes: ")
	color.Green("% X\n\n", bytes)

	bestFormat := "Unknown"
	bestScore := 0
	headerFmt := color.New(color.FgHiWhite).SprintfFunc()
	columnFmt := color.New(color.FgCyan).SprintfFunc()

	tbl := table.New("Format", "Similarity")
	tbl.WithHeaderFormatter(headerFmt).WithFirstColumnFormatter(columnFmt)
	for format, sig := range signatures {
		score := matchPercentage(bytes, sig)
		if score > bestScore {
			bestScore = score
			bestFormat = format
		}
		tbl.AddRow(format, score)
	}
	tbl.Print()
	fmt.Printf("\nMost likely format: ")
	color.Green("%s (%d%%)\n", bestFormat, bestScore)
}
