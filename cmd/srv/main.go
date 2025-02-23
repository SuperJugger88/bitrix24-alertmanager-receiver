package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"and.ivanov.go.bitrix24_receiver/internal/alertmanager"
	"and.ivanov.go.bitrix24_receiver/internal/bitrix"
	"and.ivanov.go.bitrix24_receiver/internal/metrics"
	"and.ivanov.go.bitrix24_receiver/internal/template"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	bitrixURL    = os.Getenv("BITRIX_WEBHOOK_URL")
	templatePath = os.Getenv("MESSAGE_TEMPLATE_PATH")
)

func main() {
	bitrixClient := bitrix.NewClient(bitrixURL)
	tmplProcessor, err := template.NewProcessor(templatePath)
	if err != nil {
		log.Fatalf("Ошибка при инициализации обработчика шаблонов: %v", err)
	}

	handler := alertmanager.NewWebhookHandler(bitrixClient, tmplProcessor)

	// Регистрируем метрики
	prometheus.MustRegister(metrics.RequestDuration)

	// Добавляем обработчики
	http.HandleFunc("/webhook", func(w http.ResponseWriter, r *http.Request) {
		dialogID := r.URL.Query().Get("dialog_id")
		if dialogID == "" {
			http.Error(w, "dialog_id parameter is required", http.StatusBadRequest)
			return
		}
		if !strings.HasPrefix(dialogID, "chat") {
			dialogID = "chat" + dialogID
		}
		ctx := context.WithValue(r.Context(), "dialogID", dialogID)
		handler.Handle(w, r.WithContext(ctx))
	})
	http.Handle("/metrics", promhttp.Handler()) // Эндпоинт для метрик
	port := os.Getenv("APP_PORT")

	location, err := time.LoadLocation("Europe/Moscow")
	if err != nil {
		log.Printf("Ошибка установки часового пояса: %v, используется UTC", err)
		location = time.UTC
	}
	time.Local = location

	log.Printf("Используется часовой пояс: %s", time.Local.String())

	server := &http.Server{
		Addr:         ":" + port,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	log.Printf("Сервер запущен на порту %s", port)
	log.Printf("Метрики доступны по адресу http://localhost:%s/metrics", port)
	if err := server.ListenAndServe(); err != nil {
		log.Panicf("Ошибка при запуске сервера: %v", err)
	}
}
