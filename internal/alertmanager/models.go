package alertmanager

import (
	"github.com/prometheus/alertmanager/template"
)

// AlertGroups группирует алерты по статусу
type AlertGroups struct {
	Firing   []*template.Alert // Активные алерты
	Resolved []*template.Alert // Разрешенные алерты
}
