package models

import (
	"github.com/Unknwon/com"
	"strings"
)

func getIdFromSlug(s string) int64 {
	arr := strings.Split(s, "_")
	strId := arr[len(arr)-1]
	//todo handle error
	id := com.StrTo(strId).MustInt64()
	return id
}
