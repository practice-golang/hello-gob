// ref - https://medium.com/@kpbird/golang-serialize-struct-using-gob-part-2-f6134dd4f22c

package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
)

// Book - 학생
type Book struct {
	Title  string
	Author string
}

func main() {
	fmt.Println("Gob Practice")

	bookEncode := Book{Title: "책 이름", Author: "저자 이름"}

	var buf bytes.Buffer
	e := gob.NewEncoder(&buf)
	if err := e.Encode(bookEncode); err != nil {
		panic(err)
	}
	fmt.Println("Encoded Struct ", buf)

	var bookDecode Book
	d := gob.NewDecoder(&buf)
	if err := d.Decode(&bookDecode); err != nil {
		panic(err)
	}

	fmt.Println("Decoded Struct ", bookDecode.Title, "\t", bookDecode.Author)
}
