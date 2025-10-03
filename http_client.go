package inpu

import "net/http"

func getDefaultClient() *http.Client {
	return &http.Client{}
}
