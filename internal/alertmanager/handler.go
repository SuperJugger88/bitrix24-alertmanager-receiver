package alertmanager

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"and.ivanov.go.bitrix24_receiver/internal/bitrix"
	"and.ivanov.go.bitrix24_receiver/internal/metrics"
	"and.ivanov.go.bitrix24_receiver/internal/template"
	"and.ivanov.go.bitrix24_receiver/internal/worker"

	"github.com/prometheus/alertmanager/notify/webhook"
)

// WebhookHandler обрабатывает входящие webhook от Alertmanager
type WebhookHandler struct {
	bitrixClient *bitrix.Client
	tmpl         *template.Processor
	pool         *worker.Pool
}

// NewWebhookHandler создает новый обработчик webhook
func NewWebhookHandler(bitrixClient *bitrix.Client, tmpl *template.Processor, workers int) *WebhookHandler {
	h := &WebhookHandler{
		bitrixClient: bitrixClient,
		tmpl:         tmpl,
		pool:         worker.NewPool(workers),
	}
	h.pool.Start()
	return h
}

type AlertJob struct {
	handler *WebhookHandler
	msg     *webhook.Message
	done    chan error
}

func (j *AlertJob) Execute() error {
	message, err := j.handler.tmpl.ProcessAlert(j.msg)
	if err != nil {
		j.done <- err
		return err
	}

	dialogID := os.Getenv("BITRIX_DIALOG_ID")

	bitrixMsg := &bitrix.Message{
		DialogID: dialogID,
		Message:  message,
	}

	err = j.handler.bitrixClient.SendMessageSync(bitrixMsg)
	j.done <- err
	return err
}

// Handle обрабатывает входящий webhook
func (h *WebhookHandler) Handle(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	defer func() {
		metrics.RequestDuration.WithLabelValues(r.Method, "/webhook").Observe(time.Since(start).Seconds())
	}()

	var msg webhook.Message
	if err := json.NewDecoder(r.Body).Decode(&msg); err != nil {
		log.Printf("Ошибка при разборе JSON: %v", err)
		http.Error(w, fmt.Sprintf("ошибка при разборе JSON: %v", err), http.StatusBadRequest)
		return
	}

	done := make(chan error)
	h.pool.Submit(&AlertJob{
		handler: h,
		msg:     &msg,
		done:    done,
	})

	// Ждем завершения обработки или таймаута
	select {
	case err := <-done:
		if err != nil {
			log.Printf("Ошибка обработки: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	case <-time.After(10 * time.Second):
		http.Error(w, "Timeout", http.StatusGatewayTimeout)
		return
	}

	w.WriteHeader(http.StatusOK)
}
