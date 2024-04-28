package main // import "gob-plain"

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"io/fs"
	"os"
	"slices"
)

type Book struct {
	Title  string
	Author string
}

func saveGOB(fileName string, data interface{}) error {
	var buf bytes.Buffer
	encoder := gob.NewEncoder(&buf)
	if err := encoder.Encode(data); err != nil {
		return err
	}
	return os.WriteFile(fileName, buf.Bytes(), fs.FileMode(0644))
}

func loadGOB(fileName string, data interface{}) error {
	fileData, err := os.ReadFile(fileName)
	if err != nil {
		return err
	}

	buf := bytes.NewBuffer(fileData)
	decoder := gob.NewDecoder(buf)
	if err := decoder.Decode(data); err != nil {
		return err
	}

	return nil
}

func main() {
	bookEncode := []Book{
		{Title: "Holy bible genesis", Author: "Mosses"},
		{Title: "Sutta nipāta", Author: "Buddha"},
		{Title: "القرآن(Qur'an)", Author: "محمد(Muhammad)"},
	}
	fileName := "book.gob"

	if err := saveGOB(fileName, bookEncode); err != nil {
		fmt.Println("Error saving struct:", err)
		return
	}

	fmt.Println("Struct saved to", fileName)

	var bookDecode []Book
	if err := loadGOB(fileName, &bookDecode); err != nil {
		fmt.Println("Error loading struct:", err)
		return
	}

	fmt.Println("Struct loaded as:", bookDecode)

	// Find index of Buddha
	idx := slices.IndexFunc(bookDecode, func(p Book) bool {
		return p.Author == "Buddha"
	})
	fmt.Println("idx of Buddha's data:", idx)
}
