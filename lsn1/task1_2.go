// 1. ... напишите программу, содержащую вызов функции, которая будет создавать паническую ситуацию неявно.
// Затем создайте отложенный вызов, который будет обрабатывать эту паническую ситуацию и, в частности, печатать предупреждение в консоль.
// Критерием успешного выполнения задания является то, что программа не завершается аварийно ни при каких условиях.
// 2. Дополните функцию из п.1 возвратом собственной ошибки в случае возникновения панической ситуации.
// Собственная ошибка должна хранить время обнаружения панической ситуации.
// Критерием успешного выполнения задания является наличие обработки созданной ошибки в функции main и вывод ее состояния в консоль

package main

import (
	"bytes"
	"errors"
	"fmt"
	"runtime"
	"time"
)

type TracedError struct {
	Err        error
	When       time.Time
	StackTrace []byte
}

func (e *TracedError) Error() string {
	var buf bytes.Buffer
	fmt.Fprintf(&buf, "Error: %s\n", e.Err)
	fmt.Fprintf(&buf, "Raised at %v\n", e.When)
	fmt.Fprintf(&buf, "Trace: %s\n", e.StackTrace)
	return buf.String()
}

func (e *TracedError) Unwrap() error {
	return e.Err
}

func NewTracedError(value interface{}) error {
	var err error
	switch v := value.(type) {
	case error:
		err = v
	case string:
		err = errors.New(v)
	default:
		err = errors.New(v.(string))
	}
	buf := make([]byte, 1024)
	for {
		n := runtime.Stack(buf, true)
		if n < len(buf) {
			break
		}
		buf = make([]byte, 2*len(buf))
	}
	return &TracedError{
		Err:        err,
		When:       time.Now(),
		StackTrace: buf,
	}
}

func HandlePanic() {
	defer func() {
		if v := recover(); v != nil {
			tracedErr := NewTracedError(v)
			fmt.Println(tracedErr)
		}
	}()
	panic(errors.New("Paaaanic"))
}

func TestTracedError() {
	var someErr = errors.New("some error")

	err := fmt.Errorf("found error: %w", someErr)
	fmt.Println(errors.Is(err, someErr))

	tracedErr := NewTracedError(someErr)
	fmt.Println(errors.Is(tracedErr, someErr))
}
