package logic

//func (c *Hs) updateAvatars() {
//	fmt.Printf("updateAvatars")
//	guilds, err := c.guilds.GuildGetAll()
//	if err != nil {
//		c.log.ErrorErr(err)
//		return
//	}
//	for _, guild := range guilds {
//		fmt.Printf(" Guild: %s", guild.Name)
//		membersRead, errm := c.corpMember.CorpMembersRead(guild.ID)
//		if errm != nil {
//			c.log.ErrorErr(errm)
//			return
//		}
//		for _, member := range membersRead {
//			var avatarURL string
//
//			if guild.Type == "ds" {
//				avatarURL = c.ds.GetAvatarUrl(member.UserId)
//			} else if guild.Type == "tg" {
//				avatarURL = c.tg.GetAvatarUrl(member.UserId)
//			}
//
//			if avatarURL != "" && avatarURL != member.AvatarUrl {
//				erru := c.corpMember.CorpMemberAvatarUpdate(member.UserId, guild.ID, avatarURL)
//				if erru != nil {
//					c.log.ErrorErr(erru)
//				}
//				fmt.Printf("Avatar update %s %s %s\n", guild.Name, member.Name, avatarURL)
//			}
//			time.Sleep(1 * time.Second)
//			fmt.Printf(".")
//		}
//	}
//	fmt.Println("updateAvatars() DONE")
//}
