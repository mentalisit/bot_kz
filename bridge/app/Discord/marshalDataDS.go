package ds

import (
	"bridge/models"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

const endpoint = "kz_bot"

func (d *Discord) MarshalDataSendBridgeAsync(message models.BridgeSendToMessenger) []models.MessageIds {
	data, err := json.Marshal(message)
	if err != nil {
		d.log.ErrorErr(err)
		return nil
	}

	resp, err := http.Post("http://"+endpoint+"/discord/send/bridge", "application/json", bytes.NewBuffer(data))
	if err != nil {
		resp, err = http.Post("http://192.168.100.131:802/send/bridge", "application/json", bytes.NewBuffer(data))
		if err != nil {
			d.log.ErrorErr(err)
			return nil
		}
	}
	var dataReply []models.MessageIds
	err = json.NewDecoder(resp.Body).Decode(&dataReply)
	if err != nil {
		d.log.Info(fmt.Sprintf("err resp.Body %+v\n", resp.Body))
		d.log.ErrorErr(err)
		return nil
	}
	return dataReply
}

func (d *Discord) MarshalDataDiscordDel(message models.DeleteMessageStruct) {
	data, err := json.Marshal(message)
	if err != nil {
		d.log.ErrorErr(err)
		return
	}

	_, err = http.Post("http://"+endpoint+"/discord/del", "application/json", bytes.NewBuffer(data))
	if err != nil {
		_, err = http.Post("http://192.168.100.131:802/data", "application/json", bytes.NewBuffer(data))
		d.log.ErrorErr(err)
		return
	}
}

func (d *Discord) MarshalDataDiscordSendDel(message models.SendTextDeleteSeconds) {
	data, err := json.Marshal(message)
	if err != nil {
		d.log.ErrorErr(err)
		return
	}

	_, err = http.Post("http://"+endpoint+"/discord/SendDel", "application/json", bytes.NewBuffer(data))
	if err != nil {
		_, err = http.Post("http://192.168.100.131:802/data", "application/json", bytes.NewBuffer(data))
		d.log.ErrorErr(err)
		return
	}
}
