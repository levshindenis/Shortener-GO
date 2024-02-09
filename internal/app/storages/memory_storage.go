package storages

import "errors"

type MemoryStorage struct {
	arr []MSItem
}

type MSItem struct {
	key     string
	value   string
	userid  string
	deleted bool
}

func (st *MemoryStorage) GetArr() []MSItem {
	return st.arr
}

func (st *MemoryStorage) SetData(key string, value string, userid string) error {
	st.arr = append(st.GetArr(), MSItem{key: key, value: value, userid: userid, deleted: false})
	return nil
}

func (st *MemoryStorage) GetData(value string, param string, userid string) (string, []bool, error) {
	if param == "key" {
		for _, elem := range st.GetArr() {
			if elem.key == value {
				return elem.value, []bool{elem.deleted}, nil
			}
		}
		return "", nil, nil
	}
	if param == "value" {
		for _, elem := range st.GetArr() {
			if elem.value == value {
				return elem.key, []bool{elem.deleted}, nil
			}
		}
		return "", nil, nil
	}
	if param == "all" {
		mystr := ""
		var mybool []bool
		for _, elem := range st.GetArr() {
			if elem.userid == userid {
				mystr += elem.key + "*" + elem.value + "*"
				mybool = append(mybool, elem.deleted)
			}
		}
		if mystr != "" {
			return mystr[:len(mystr)-1], mybool, nil
		}
		return "", nil, nil
	}
	return "", nil, errors.New("unknown param")
}
