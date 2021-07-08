// @ gb go 2 manual lsn 1
// HTTP-клиент для отправления JSON-объектов.
// HTTP-клиент должен содержать собственные ошибки, соответствующие статусам кодов-ответов HTTP.

package client

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

const (
	jsonContentType = "application/json"
)

var (
	ErrorIncorrectBodyFormat = errors.New("incorrect body format")
	ErrorSendRequest         = errors.New("couldn't send request")
	ErrorUnknown             = errors.New("unknown error")
	needPanic                bool
)

type HttpStatusError struct {
	status int
}

func NewHttpStatusError(status int) error {
	return &HttpStatusError{status}
}

func (s *HttpStatusError) Error() string {
	return fmt.Sprintf("status code: %d", s.status)
}

func (s *HttpStatusError) Status() int {
	return s.status
}

// PostJson - основная функция пакета, выполняет отправку строки на указанный адрес URL, но только если она представляет собой JSON.
// В противном случае возвращается ошибка ErrorIncorrectBodyFormat.
// Если возникнут ошибки при установлении HTTP-соединения, то возвращается ошибка ErrorSendRequest.
// При некорректном коде ответа будет возвращена собственная ошибка типа HttpStatusError.
// Некорректным кодом ответа считаем любой HTTP-код, отличный от 200.
func PostJson(client *http.Client, url, body string) (err error) {
	var (
		js map[string]interface{}
	)
	// be aware of encoding/json is panicist library
	defer func() {
		if v := recover(); v != nil {
			err = fmt.Errorf("%w: %v", ErrorUnknown, v)
		}
	}()
	err = json.Unmarshal([]byte(body), &js)
	// условия возникновения панической ситуации при вызове json.Unmarshal весьма специфичны:
	// это должна быть ошибка в реализации декодера либо конкурентная модификация исходного слайса байт.
	// Очевидно, стабильно воспроизвести такую ситуацию мы не можем.
	// Но зато можем эмулировать возникновение паники:
	if needPanic {
		panic("JSON decoder out of sync - data changing underfoot?")
	}
	if err != nil {
		return fmt.Errorf("%w: %s", ErrorIncorrectBodyFormat, err.Error())
	}
	res, err := client.Post(url, jsonContentType, bytes.NewReader([]byte(body)))
	if err != nil {
		return fmt.Errorf("%w: %s", ErrorSendRequest, err.Error())
	}
	if res.StatusCode != http.StatusOK {
		return NewHttpStatusError(res.StatusCode)
	}
	return nil
}
