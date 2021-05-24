package main

import (
	"errors"
	"fmt"
	"os"
	"runtime/debug"
	"time"
)

type ErrorWithTrace struct {
	message string
	when    time.Time
	where   string
}

func New(text string) error {
	return &ErrorWithTrace{
		message: text,
		when:    time.Now(),
		where:   string(debug.Stack()),
	}
}

func (e *ErrorWithTrace) Error() string {
	return fmt.Sprintf("error: %s\nraised at %v\nwith trace:\n%s", e.message, e.when, e.where)
}

func mainRecover() {
	if v := recover(); v != nil {
		err := New(v.(string))
		fmt.Println(err)
	}
}

func main() {
	defer mainRecover()
	_, err := os.OpenFile("_____")
	if err != nil {
		panic(errors.New("file not found"))
	}
}
