package main

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/fatih/color"
	"github.com/rodaine/table"
)

type Signatures struct {
	Signatures []Signature `json:"signatures"`
}

type Signature struct {
	Extension    string `json:"extension"`
	Magicnumbers string `json:"magicnumbers"`
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
	for i := range sigBytes {
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
	bytes, err := readFirstBytes(filename, 50)
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
	jsonFile, err := os.Open("./signatures.json")
	if err != nil {
		fmt.Println(err)
	}
	defer jsonFile.Close()
	byteValue, _ := io.ReadAll(jsonFile)
	var signatures Signatures
	json.Unmarshal(byteValue, &signatures)
	for i := 0; i < len(signatures.Signatures); i++ {
		slice := append([]byte(nil), bytes[0:len(signatures.Signatures[i].Magicnumbers)]...)
		score := matchPercentage(slice, signatures.Signatures[i].Magicnumbers)
		if score > bestScore {
			bestScore = score
			bestFormat = signatures.Signatures[i].Extension
		}
		tbl.AddRow(signatures.Signatures[i].Extension, score)
	}
	tbl.Print()
	fmt.Printf("\nMost likely format: ")
	color.Green("%s (%d%%)\n", bestFormat, bestScore)
}
