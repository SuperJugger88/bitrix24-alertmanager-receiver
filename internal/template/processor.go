package template

import (
	"fmt"
	"strings"
	"sync"
	"text/template"

	"github.com/prometheus/alertmanager/notify/webhook"
)

// Processor обрабатывает шаблоны сообщений
type Processor struct {
	tmpl  *template.Template
	cache sync.Map // Кэш для обработанных шаблонов
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
	// Пробуем получить из кэша
	cacheKey := fmt.Sprintf("%s_%s", msg.Status, msg.Receiver)
	if cached, ok := p.cache.Load(cacheKey); ok {
		return cached.(string), nil
	}

	var message strings.Builder
	if err := p.tmpl.ExecuteTemplate(&message, "bitrix24.message", msg); err != nil {
		return "", fmt.Errorf("ошибка при выполнении шаблона: %w", err)
	}

	result := message.String()
	p.cache.Store(cacheKey, result)
	return result, nil
}
