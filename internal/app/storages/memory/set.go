package memory

import "github.com/levshindenis/sprint1/internal/app/models"

func (ms *Memory) SetData(key string, value string, userid string) error {
	ms.Arr = append(ms.Arr, models.MSItem{Key: key, Value: value, UserId: userid, Deleted: false})
	return nil
}
