package template

import (
	"fmt"
	"strings"
	"text/template"

	"github.com/prometheus/alertmanager/notify/webhook"
)

// Processor обрабатывает шаблоны сообщений
type Processor struct {
	tmpl *template.Template
}

// NewProcessor создает новый обработчик шаблонов
func NewProcessor(templatePath string) (*Processor, error) {
	tmpl, err := template.ParseFiles(templatePath)
	if err != nil {
		return nil, fmt.Errorf("ошибка при загрузке шаблона: %w", err)
	}
	return &Processor{tmpl: tmpl}, nil
}

// ProcessAlert обрабатывает алерт и возвращает отформатированное сообщение
func (p *Processor) ProcessAlert(msg *webhook.Message) (string, error) {
	var message strings.Builder
	if err := p.tmpl.ExecuteTemplate(&message, "bitrix24.message", msg); err != nil {
		return "", fmt.Errorf("ошибка при выполнении шаблона: %w", err)
	}

	result := strings.ReplaceAll(message.String(), "\n", "")
	return result, nil
}
