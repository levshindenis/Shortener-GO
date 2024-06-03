package tools

import "github.com/3th1nk/cidr"

// MakeIps - возвращает список верных ip
func MakeIps(str string) []string {
	if str == "" {
		return nil
	}

	var arr []string

	c, _ := cidr.Parse(str)

	c.Each(func(ip string) bool {
		arr = append(arr, ip)
		return true
	})

	return arr
}
