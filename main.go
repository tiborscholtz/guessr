package projects
package main

import (
	"encoding/hex"
	"fmt"
	"io"
	"os"
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
	"MP4":  "0000002066747970", // 'ftyp'
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

	fmt.Printf("Read bytes: % X\n\n", bytes)
	fmt.Println("Format likelihoods:")

	bestFormat := "Unknown"
	bestScore := 0

	for format, sig := range signatures {
		score := matchPercentage(bytes, sig)
		if score > bestScore {
			bestScore = score
			bestFormat = format
		}

		fmt.Printf("  %-6s : %3d%% match (sig: %s)\n", format, score, sig)
	}

	fmt.Printf("\nMost likely format: %s (%d%%)\n", bestFormat, bestScore)
}
