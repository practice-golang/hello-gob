package main // import "gob-aes"

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/gob"
	"fmt"
	"io"
	"os"
	"slices"
)

type Person struct {
	Name string
	Age  int
}

func saveGOB(fileName string, data interface{}, key []byte) error {
	var gobBuffer bytes.Buffer
	iv := make([]byte, aes.BlockSize)

	file, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := gob.NewEncoder(&gobBuffer)
	err = encoder.Encode(data)
	if err != nil {
		return err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return err
	}

	_, err = io.ReadFull(rand.Reader, iv)
	if err != nil {
		return err
	}

	_, err = file.Write(iv)
	if err != nil {
		return err
	}

	stream := cipher.NewCFBEncrypter(block, iv)

	writer := &cipher.StreamWriter{S: stream, W: file}
	_, err = writer.Write(gobBuffer.Bytes())
	if err != nil {
		return err
	}

	return nil
}

func loadGOB(fileName string, key []byte, data interface{}) error {
	// encryptedData := make([]byte, 4096)

	file, err := os.Open(fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	// _, err = file.Read(encryptedData)
	// if err != nil && err != io.EOF {
	// 	return err
	// }
	fileInfo, err := file.Stat()
	if err != nil {
		return fmt.Errorf("loadHostData/Stat: %s", err)
	}

	encryptedData := make([]byte, fileInfo.Size())
	_, err = io.ReadFull(file, encryptedData)
	if err != nil {
		return fmt.Errorf("loadHostData/ReadFull: %s", err)
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return err
	}
	iv := encryptedData[:aes.BlockSize]
	stream := cipher.NewCFBDecrypter(block, iv)

	encryptedData = encryptedData[aes.BlockSize:]

	reader := cipher.StreamReader{S: stream, R: bytes.NewReader(encryptedData)}
	decoder := gob.NewDecoder(&reader)
	err = decoder.Decode(data)
	if err != nil {
		return err
	}

	return nil
}

func main() {
	persons := []Person{
		{Name: "James", Age: 32},
		{Name: "Alice", Age: 30},
		{Name: "John", Age: 12},
	}
	fileName := "person.gob"
	key := []byte("0123456789abcdef0123456789abcdef") // AES key (32byte = 256bit)

	err := saveGOB(fileName, persons, key)
	if err != nil {
		fmt.Println("Error saving struct:", err)
		return
	}

	fmt.Println("Struct saved to", fileName)

	var decodedPersons []Person
	err = loadGOB(fileName, key, &decodedPersons)
	if err != nil {
		fmt.Println("Error loading struct:", err)
		return
	}

	fmt.Println("Struct loaded as:", decodedPersons)

	// Find index of John
	idx := slices.IndexFunc(decodedPersons, func(p Person) bool {
		return p.Name == "John"
	})
	fmt.Println("idx of John's data:", idx)
}
