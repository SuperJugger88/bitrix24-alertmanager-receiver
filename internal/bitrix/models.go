package bitrix

// Message представляет сообщение для отправки в Bitrix24
type Message struct {
	DialogID string `json:"DIALOG_ID"` // ID чата в Bitrix24
	Message  string `json:"MESSAGE"`   // Текст сообщения
}
