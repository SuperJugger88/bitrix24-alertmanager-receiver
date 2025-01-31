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

	"github.com/prometheus/alertmanager/notify/webhook"
)

// WebhookHandler обрабатывает входящие webhook от Alertmanager
type WebhookHandler struct {
	bitrixClient *bitrix.Client
	tmpl         *template.Processor
}

// NewWebhookHandler создает новый обработчик webhook
func NewWebhookHandler(bitrixClient *bitrix.Client, tmpl *template.Processor) *WebhookHandler {
	return &WebhookHandler{
		bitrixClient: bitrixClient,
		tmpl:         tmpl,
	}
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

	if err := h.processAlert(&msg); err != nil {
		log.Printf("Ошибка обработки: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *WebhookHandler) processAlert(msg *webhook.Message) error {
	message, err := h.tmpl.ProcessAlert(msg)
	if err != nil {
		return err
	}

	dialogID := os.Getenv("BITRIX_DIALOG_ID")
	bitrixMsg := &bitrix.Message{
		DialogID: dialogID,
		Message:  message,
	}

	return h.bitrixClient.SendMessage(bitrixMsg)
}
