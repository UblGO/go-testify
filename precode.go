package main

import (
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var cafeList = map[string][]string{
	"moscow": []string{"Мир кофе", "Сладкоежка", "Кофе и завтраки", "Сытый студент"},
}

func mainHandle(w http.ResponseWriter, req *http.Request) {
	countStr := req.URL.Query().Get("count")
	if countStr == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("count missing"))
		return
	}

	count, err := strconv.Atoi(countStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("wrong count value"))
		return
	}

	city := req.URL.Query().Get("city")

	cafe, ok := cafeList[city]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("wrong city value"))
		return
	}

	if count > len(cafe) {
		count = len(cafe)
	}

	answer := strings.Join(cafe[:count], ",")

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(answer))
}

// Функция получения ответа
func getResponseRecorder(url string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(http.MethodGet, url, nil)

	res := httptest.NewRecorder()
	handler := http.HandlerFunc(mainHandle)
	handler.ServeHTTP(res, req)

	return res
}

func TestMainHandlerWhenCountMoreThanTotal(t *testing.T) {
	totalCount := 4
	responseRecorder := getResponseRecorder("/cafe?count=20&city=moscow")
	// проверка кода ответа
	assert.Equal(t, http.StatusOK, responseRecorder.Code)

	// проверка на непустое тело ответа
	require.NotEmpty(t, responseRecorder.Body)
	cafeSlice := strings.Split(responseRecorder.Body.String(), ",")

	// проверка количества кафе
	assert.Len(t, cafeSlice, totalCount)
}

func TestMainHandlerWhenOk(t *testing.T) {
	responseRecorder := getResponseRecorder("/cafe?count=3&city=moscow")
	count := 3

	// проверка кода ответа
	assert.Equal(t, http.StatusOK, responseRecorder.Code)

	// проверка на непустое тело ответа
	require.NotEmpty(t, responseRecorder.Body)
	cafeSlice := strings.Split(responseRecorder.Body.String(), ",")

	// проверка количества кафе
	assert.Len(t, cafeSlice, count)
}

func TestMainHandlerWhenWrongCity(t *testing.T) {
	responseRecorder := getResponseRecorder("/cafe?count=3&city=tomsk")
	expBody := `wrong city value`

	// проверка кода ответа
	assert.Equal(t, http.StatusBadRequest, responseRecorder.Code)

	// проверка на непустое тело ответа
	require.NotEmpty(t, responseRecorder.Body)
	body := responseRecorder.Body.String()

	// проверка тела ответа
	assert.Equal(t, expBody, body)
}
