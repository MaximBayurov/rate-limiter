package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/MaximBayurov/rate-limiter/internal/app"
	"github.com/MaximBayurov/rate-limiter/internal/iplists"
	"github.com/MaximBayurov/rate-limiter/internal/logger"
)

const contentType = "application/json"

type Handler func(http.ResponseWriter, *http.Request)

func AddIPHandler( //nolint:dupl
	app app.App,
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
			IP        string `json:"ip"`
			Type      string `json:"type"`
			Overwrite bool   `json:"overwrite,omitempty"`
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

		err = app.AddIP(
			request.IP,
			request.Type,
			request.Overwrite,
		)
		if err != nil && !errors.Is(err, iplists.ErrBase) {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		var message string
		if err != nil {
			message = err.Error()
		} else {
			message = "IP успешно добавлен в список"
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

func DeleteIPHandler(
	app app.App,
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
			IP   string `json:"ip"`
			Type string `json:"type"`
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

		err = app.DeleteIP(
			request.IP,
			request.Type,
		)
		if err != nil && !errors.Is(err, iplists.ErrBase) {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		var message string
		if err != nil {
			message = err.Error()
		} else {
			message = "IP успешно удален из списка"
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
