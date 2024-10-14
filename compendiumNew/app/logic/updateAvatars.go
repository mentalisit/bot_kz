package logic

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

func (c *Hs) updateAvatars() {
	fmt.Printf("updateAvatars")
	guilds, err := c.guilds.GuildGetAll()
	if err != nil {
		c.log.ErrorErr(err)
		return
	}
	for _, guild := range guilds {
		fmt.Printf(" Guild: %s", guild.Name)
		membersRead, errm := c.corpMember.CorpMembersRead(guild.ID)
		if errm != nil {
			c.log.ErrorErr(errm)
			return
		}
		for _, member := range membersRead {
			url := fmt.Sprintf("http://telegram/GetAvatarUrl?userid=%s", member.UserId)
			if guild.Type == "ds" {
				url = fmt.Sprintf("http://kz_bot/discord/GetAvatarUrl?userid=%s", member.UserId)
			}
			var avatarURL string

			getAvatarUrl(url, &avatarURL)

			if avatarURL != "" && avatarURL != member.AvatarUrl {
				erru := c.corpMember.CorpMemberAvatarUpdate(member.UserId, guild.ID, avatarURL)
				if erru != nil {
					c.log.ErrorErr(erru)
				}
				fmt.Printf("Avatar update %s %s %s\n", guild.Name, member.Name, avatarURL)
			}
			time.Sleep(1 * time.Second)
			fmt.Printf(".")
		}
	}
	fmt.Println("updateAvatars() DONE")
}

func getAvatarUrl(url string, result *string) {
	// Создаем контекст с тайм-аутом 3 секунды
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel() // Освобождаем ресурсы контекста после завершения функции

	// Выполнение GET-запроса с контекстом
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		fmt.Println("Ошибка при создании запроса:", err)
		return
	}

	response, err := http.DefaultClient.Do(req)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			fmt.Println("Время ожидания запроса истекло")
		} else {
			fmt.Println("Ошибка при выполнении запроса:", err)
		}
		return
	}
	defer response.Body.Close()

	// Проверка кода ответа
	if response.StatusCode != http.StatusOK {
		fmt.Printf("Неправильный статус код: %d\n", response.StatusCode)
		return
	}

	// Чтение тела ответа
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println("Ошибка при чтении ответа:", err)
		return
	}

	err = json.Unmarshal(body, &result)
	if err != nil {
		fmt.Println("Ошибка при декодировании JSON:", err)
		return
	}
}
