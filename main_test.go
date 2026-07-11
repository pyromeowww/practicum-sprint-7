package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCafeNegative(t *testing.T) {
	handler := http.HandlerFunc(mainHandle)

	requests := []struct {
		request string
		status  int
		message string
	}{
		{"/cafe", http.StatusBadRequest, "unknown city"},
		{"/cafe?city=omsk", http.StatusBadRequest, "unknown city"},
		{"/cafe?city=tula&count=na", http.StatusBadRequest, "incorrect count"},
	}
	for _, v := range requests {
		response := httptest.NewRecorder()
		req := httptest.NewRequest("GET", v.request, nil)
		handler.ServeHTTP(response, req)

		assert.Equal(t, v.status, response.Code)
		assert.Equal(t, v.message, strings.TrimSpace(response.Body.String()))
	}
}

func TestCafeWhenOk(t *testing.T) {
	handler := http.HandlerFunc(mainHandle)

	requests := []string{
		"/cafe?count=2&city=moscow",
		"/cafe?city=tula",
		"/cafe?city=moscow&search=ложка",
	}
	for _, v := range requests {
		response := httptest.NewRecorder()
		req := httptest.NewRequest("GET", v, nil)

		handler.ServeHTTP(response, req)

		assert.Equal(t, http.StatusOK, response.Code)
	}
}

func TestCafeCount(t *testing.T) {
	handler := http.HandlerFunc(mainHandle)
	moscowCafeCount := len(cafeList["moscow"])
	request := []struct {
		count int // передаваемое значение count
		want  int // ожидаемое количество кафе в ответе
	}{
		{0, 0},
		{1, 1},
		{2, 2},
		{100, moscowCafeCount},
	}
	for _, v := range request {
		url := fmt.Sprintf("/cafe?city=moscow&count=%d", v.count) // делаем URL с параметром count
		req := httptest.NewRequest("GET", url, nil)               // создаём HTTP запрос
		response := httptest.NewRecorder()                        // запись ответа
		handler.ServeHTTP(response, req)                          // вызов обработчика
		if response.Code != http.StatusOK {                       // проверка на статус кода
			t.Errorf("для count=%d ожидается статус 200, получено %d", v.count, response.Code)
			continue
		}
		// проверка на пустоту
		body := strings.TrimSpace(response.Body.String())
		if body == "" {
			if v.want != 0 {
				t.Errorf("для count=%d ожидается %d кафе, получен пустой ответ", v.count, v.want)
			}
			continue
		}
		// разбитие ответа на слайс и сравниваем длину с ожидаемой
		cafes := strings.Split(body, ",")
		if len(cafes) != v.want {
			t.Errorf("для count=%d ожидается %d кафе, получено %d", v.count, v.want, len(cafes))
		}
	}
}
