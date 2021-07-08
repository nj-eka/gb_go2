package client

import (
	"errors"
	"net/http"
	"testing"

	check "gopkg.in/check.v1"
)

func Test(t *testing.T) { check.TestingT(t) }

type TestSuite struct{}

var _ = check.Suite(&TestSuite{})

// Проверка на отправление невалидного JSON-объекта {{"data": 1}.
// Мы ожидаем возвращения ошибки ErrorIncorrectBodyFormat.
func (s *TestSuite) TestErrorIncorrectJson(c *check.C) {
	var err error = PostJson(http.DefaultClient, "http://google.com", `{{"data": 1}`)
	c.Assert(err, check.NotNil)
	c.Assert(errors.Is(err, ErrorIncorrectBodyFormat), check.Equals, true)
}

// Проверка на отправление данных по некорректному адресу URL,
// ожидаем ошибку ErrorSendRequest.
func (s *TestSuite) TestErrorSendRequest(c *check.C) {
	var err error = PostJson(http.DefaultClient, "http://icorrect.url", `{}`)
	c.Assert(err, check.NotNil)
	c.Assert(errors.Is(err, ErrorSendRequest), check.Equals, true)
}

// Проверка на создание собственной ошибки типа HttpStatusError с использованием вспомогательного сервиса httpstat.us.
// Сервис возвращает код ответа, который мы указываем в запросе.
func (s *TestSuite) TestIncorrectStatus(c *check.C) {
	var err error = PostJson(http.DefaultClient, "http://httpstat.us/500", `{}`)
	c.Assert(err, check.NotNil)
	e, ok := err.(*HttpStatusError)
	c.Assert(ok, check.Equals, true)
	c.Assert(e.Status(), check.Equals, http.StatusInternalServerError)
}

// Проверка корректного выполнения функции.
func (s *TestSuite) TestOk(c *check.C) {
	var err error = PostJson(http.DefaultClient, "http://httpstat.us/200", `{}`)
	c.Assert(err, check.IsNil)
}

func (s *TestSuite) TestUnknownError(c *check.C) {
	needPanic = true
	var err error = PostJson(http.DefaultClient, "http://httpstat.us/200", `{}`)
	c.Assert(err, check.NotNil)
	c.Assert(errors.Is(err, ErrorUnknown), check.Equals, true)
	needPanic = false
}
