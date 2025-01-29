package main

import (
	"log"
	"net/http"
	"os"
	"runtime"
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
	numWorkers := runtime.NumCPU()
	bitrixClient := bitrix.NewClient(bitrixURL, numWorkers)
	tmplProcessor, err := template.NewProcessor(templatePath)
	if err != nil {
		log.Fatalf("Ошибка при инициализации обработчика шаблонов: %v", err)
	}

	handler := alertmanager.NewWebhookHandler(bitrixClient, tmplProcessor, numWorkers)

	// Регистрируем метрики
	prometheus.MustRegister(metrics.RequestDuration)

	// Добавляем обработчики
	http.HandleFunc("/webhook", handler.Handle)
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

	log.Printf("Сервер запущен на порту %s с %d воркерами", port, numWorkers)
	log.Printf("Метрики доступны по адресу http://localhost:%s/metrics", port)
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Ошибка при запуске сервера: %v", err)
	}
}
