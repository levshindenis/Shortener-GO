package memory

import "github.com/levshindenis/sprint1/internal/app/models"

// SetData - нужна для записи значений в Memory.
func (ms *Memory) SetData(key string, value string, userid string) error {
	ms.Arr = append(ms.Arr, models.MSItem{Key: key, Value: value, UserID: userid, Deleted: false})
	return nil
}
