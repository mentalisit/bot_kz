package TgApi

import "kz_bot/pkg/utils"

func (t *TgApi) CheckAdminTg(ChatId string, userName string) bool {
	ch := utils.WaitForMessage("CheckAdminTg")
	defer close(ch)
	m := apiRs{
		FuncApi:  funcCheckAdmin,
		Channel:  ChatId,
		UserName: userName,
	}

	a, err := convertAndSend(m)
	if err != nil {
		t.log.ErrorErr(err)
		return false
	}
	return a.ArrBool
}
