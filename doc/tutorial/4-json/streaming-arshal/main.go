package main

import (
	"encoding/json/v2"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"time"
)

type Message struct {
	Body string
	Time int
}

func main() {
	msgFilePath := getFilePath()
	read(msgFilePath)
	write(msgFilePath)
}

func read(msgFilePath string) {
	var readOnlyfile = openFile(msgFilePath, false)
	var msg Message

	fmt.Printf("msg before UnmarshalRead: %+v\n", msg)
	json.UnmarshalRead(readOnlyfile, &msg)
	fmt.Printf("msg after UnmarshalRead: %+v\n", msg)
}

func write(msgFilePath string) {
	var writablefile = openFile(msgFilePath, true)
	var newMsg = Message{Body: "Hello", Time: int(time.Now().Unix())}

	fmt.Println("file content before MarshalWrite:", getFileContent(writablefile))
	err := json.MarshalWrite(writablefile, newMsg)
	if err != nil {
		panic("Failed to MarshalWrite:\n" + err.Error())
	}
	fmt.Println("file content after MarshalWrite:", getFileContent(writablefile))
}

func getFilePath() string {
	// _, filename, _, _ gets the absolute path of this current source code file
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		panic("No caller information")
	}

	name := "message.json"
	// Get the directory containing this .go file
	srcDir := filepath.Dir(filename)

	// Build the path relative to this source folder
	path := filepath.Join(srcDir, name)

	return path
}

func getFileContent(file *os.File) string {
	resetCursor(file)
	bytes, err := io.ReadAll(file)

	if err != nil {
		panic("Failed to getFileContent: " + err.Error())
	}

	resetCursor(file)

	return string(bytes)
}

func openFile(name string, writable bool) *os.File {
	var file *os.File
	var err error

	if writable == true {
		file, err = os.OpenFile(name, os.O_RDWR, 0666)
	} else {
		file, err = os.Open(name)
	}

	if err != nil {
		panic("Failed to open file " + name)
	}

	return file
}

func resetCursor(file *os.File) {
	_, err := file.Seek(0, io.SeekStart)
	if err != nil {
		panic(err)
	}
}
