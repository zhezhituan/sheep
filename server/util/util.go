package util

import (
	"net/http"
)

//WWrite 向用户返回信息,返回resp和一个错误信息
func WWrite(w http.ResponseWriter, msg string) error {
	_, err := w.Write([]byte(msg))
	return err
}
