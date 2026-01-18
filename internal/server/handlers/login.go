package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	application "github.com/MaximBayurov/rate-limiter/internal/app"
	"github.com/MaximBayurov/rate-limiter/internal/logger"
)

func TryLoginHandler( //nolint:dupl
	app application.App,
	logger logger.Logger,
) Handler {
	return func(w http.ResponseWriter, r *http.Request) {
		// Проверяем метод запроса
		if r.Method != http.MethodPut {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Проверяем Content-Type
		if contentType != r.Header.Get("Content-Type") {
			http.Error(w, fmt.Sprintf("Content-Type must be %s", contentType), http.StatusUnsupportedMediaType)
			return
		}

		var request struct {
			Login    string `json:"login"`
			Password string `json:"password"`
			IP       string `json:"ip"`
		}

		// Декодируем JSON тело запроса в структуру
		decoder := json.NewDecoder(r.Body)
		decoder.DisallowUnknownFields() // запрещаем неизвестные поля

		err := decoder.Decode(&request)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Закрываем тело запроса
		defer func() {
			if err := r.Body.Close(); err != nil {
				logger.Error(err.Error())
			}
		}()

		err = app.TryLogin(
			request.Login,
			request.Password,
			request.IP,
		)
		if err != nil && !errors.Is(err, application.ErrBase) {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		var message string
		if err != nil {
			message = err.Error()
		} else {
			message = "авторизация разрешена"
		}
		// Отправляем ответ
		w.Header().Set("Content-Type", "application/json")
		if err = json.NewEncoder(w).Encode(map[string]interface{}{
			"success": err == nil,
			"message": message,
		}); err != nil {
			logger.Error(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func ClearBucketHandler(
	app application.App,
	logger logger.Logger,
) Handler {
	return func(w http.ResponseWriter, r *http.Request) {
		// Проверяем метод запроса
		if r.Method != http.MethodDelete {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Проверяем Content-Type
		if contentType != r.Header.Get("Content-Type") {
			http.Error(w, fmt.Sprintf("Content-Type must be %s", contentType), http.StatusUnsupportedMediaType)
			return
		}

		var request struct {
			Login string `json:"login"`
			IP    string `json:"ip"`
		}

		// Декодируем JSON тело запроса в структуру
		decoder := json.NewDecoder(r.Body)
		decoder.DisallowUnknownFields() // запрещаем неизвестные поля

		err := decoder.Decode(&request)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Закрываем тело запроса
		defer func() {
			if err := r.Body.Close(); err != nil {
				logger.Error(err.Error())
			}
		}()

		err = app.ClearLoginAttempts(
			request.Login,
			request.IP,
		)
		if err != nil && !errors.Is(err, application.ErrBase) {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		var message string
		if err != nil {
			message = err.Error()
		} else {
			message = "количество попыток успешно сброшено"
		}
		// Отправляем ответ
		w.Header().Set("Content-Type", "application/json")
		if err = json.NewEncoder(w).Encode(map[string]interface{}{
			"success": err == nil,
			"message": message,
		}); err != nil {
			logger.Error(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}
