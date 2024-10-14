package ds

import (
	"bridge/models"
	"sync"
)

func (d *Discord) SendBridgeFunc(text, username string, channelID []string, extra []models.FileInfo, Avatar string, reply *models.BridgeMessageReply, resultChannel chan<- models.MessageIds, wg *sync.WaitGroup) {

	defer wg.Done()

	m := models.BridgeSendToMessenger{
		Text:      text,
		Sender:    username,
		ChannelId: channelID,
		Avatar:    Avatar,
		Extra:     extra,
		Reply:     reply,
	}

	mesarray, err := d.MarshalDataSendBridgeAsync(m)
	if err != nil {
		d.log.ErrorErr(err)
		d.log.InfoStruct("SendBridgeFunc m: ", m)
	}

	for _, ds := range mesarray {
		resultChannel <- ds
	}
}
