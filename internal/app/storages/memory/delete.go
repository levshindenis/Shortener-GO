package memory

import "github.com/levshindenis/sprint1/internal/app/models"

func (ms *Memory) DeleteData(delValues []models.DeleteValue) error {
	for _, elem := range delValues {
		for ind, msi := range ms.Arr {
			if msi.Key == elem.Value && msi.UserID == elem.Userid {
				ms.Arr[ind].Deleted = true
			}
		}
	}
	return nil
}
