package TgApi

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"kz_bot/pkg/utils"
	"net/http"
	"time"
)

func (t *TgApi) CheckAdminTg(ChatId string, userName string) (bool, error) {
	ch := utils.WaitForMessage("CheckAdminTg")
	defer close(ch)
	m := apiRs{
		//FuncApi:  funcCheckAdmin,
		Channel:  ChatId,
		UserName: userName,
	}

	data, err := json.Marshal(m)
	if err != nil {
		return false, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "POST", "http://telegram/check_admin", bytes.NewBuffer(data))
	if err != nil {
		return false, err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	var a answer
	err = json.NewDecoder(resp.Body).Decode(&a)
	if err != nil {
		body, _ := ioutil.ReadAll(resp.Body)
		t.log.ErrorErr(err)
		t.log.Info(string(body))
	}

	return a.ArrBool, a.ArrError
}
