// Package memory используется, когда память системы выбрана хранилищем.
package memory

import (
	"github.com/levshindenis/sprint1/internal/app/models"
)

// Memory - хранилище.
// Arr - массив из данных.
type Memory struct {
	Arr []models.MSItem
}
