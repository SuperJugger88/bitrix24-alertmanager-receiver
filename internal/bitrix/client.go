package bitrix

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

// Client представляет клиент для работы с API Bitrix24
type Client struct {
	baseURL string
	client  *http.Client
}

// NewClient создает новый клиент Bitrix24
func NewClient(baseURL string) *Client {
	return &Client{
		baseURL: baseURL,
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// SendMessage отправляет сообщение в Bitrix24
func (c *Client) SendMessage(msg *Message) error {
	jsonData, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("ошибка при формировании JSON: %w", err)
	}

	log.Printf("Отправка запроса в Bitrix24: %s", string(jsonData))

	resp, err := http.Post(c.baseURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("ошибка при отправке запроса: %w", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Printf("Ошибка при закрытии тела ответа: %v", err)
		}
	}(resp.Body)

	body, _ := io.ReadAll(resp.Body)
	log.Printf("Ответ от Bitrix24: %s", string(body))

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("API вернул ошибку. Код: %d, Ответ: %s", resp.StatusCode, body)
	}

	return nil
}
