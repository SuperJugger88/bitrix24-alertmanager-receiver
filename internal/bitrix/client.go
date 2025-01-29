package bitrix

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"and.ivanov.go.bitrix24_receiver/internal/cache"
	"and.ivanov.go.bitrix24_receiver/internal/ctx"
	"and.ivanov.go.bitrix24_receiver/internal/worker"
)

// Client представляет клиент для работы с API Bitrix24
type Client struct {
	baseURL string
	pool    *worker.Pool
	cache   *cache.Cache
	ctxMgr  *ctx.ContextManager
	client  *http.Client
}

// NewClient создает новый клиент Bitrix24
func NewClient(baseURL string, workers int) *Client {
	return &Client{
		baseURL: baseURL,
		pool:    worker.NewPool(workers),
		cache:   cache.NewCache(5*time.Minute, 1000),
		ctxMgr:  ctx.NewContextManager(10 * time.Second),
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

type MessageJob struct {
	client  *Client
	message *Message
}

func (j *MessageJob) Execute() error {
	return j.client.SendMessageSync(j.message)
}

func (c *Client) SendMessage(msg *Message) {
	c.pool.Submit(&MessageJob{
		client:  c,
		message: msg,
	})
}

// SendMessage отправляет сообщение в Bitrix24
func (c *Client) SendMessageSync(msg *Message) error {
	jsonData, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("ошибка при формировании JSON: %w", err)
	}

	log.Printf("Отправка запроса в Bitrix24: %s", string(jsonData))

	resp, err := http.Post(c.baseURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("ошибка при отправке запроса: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	log.Printf("Ответ от Bitrix24: %s", string(body))

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("API вернул ошибку. Код: %d, Ответ: %s", resp.StatusCode, body)
	}

	return nil
}
