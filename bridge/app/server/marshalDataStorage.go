package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

func MarshalDataStorage(message any, t string) {
	data, err := json.Marshal(message)
	if err != nil {
		fmt.Println("Error marshalling JSON:", err)
		return
	}
	url := fmt.Sprintf("http://storage/storage/bridge/%s", t)

	aa, err := http.Post(url, "application/json", bytes.NewBuffer(data))
	if err != nil {
		fmt.Println("Error sending message to storage:", err)
		return
	}
	fmt.Println(aa)
	fmt.Printf("send to storage %+v\n", message)
}
