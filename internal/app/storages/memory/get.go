package memory

import "errors"

// GetData - нужна для получение каких-либо данных из Memory.
// Если param == key, то будет возвращен длинный URL и параметр deleted.
// Если param == value, то будет возвращен короткий URL и параметр deleted.
// Если param == all, то будут возвращены все записи по полученному UserID.
func (ms *Memory) GetData(value string, param string, userid string) (string, []bool, error) {
	if param == "key" {
		for _, elem := range ms.Arr {
			if elem.Key == value {
				return elem.Value, []bool{elem.Deleted}, nil
			}
		}
		return "", nil, nil
	}
	if param == "value" {
		for _, elem := range ms.Arr {
			if elem.Value == value {
				return elem.Key, []bool{elem.Deleted}, nil
			}
		}
		return "", nil, nil
	}
	if param == "all" {
		mystr := ""
		var mybool []bool
		for _, elem := range ms.Arr {
			if elem.UserID == userid {
				mystr += elem.Key + "*" + elem.Value + "*"
				mybool = append(mybool, elem.Deleted)
			}
		}
		if mystr != "" {
			return mystr[:len(mystr)-1], mybool, nil
		}
		return "", nil, nil
	}
	return "", nil, errors.New("unknown param")
}
