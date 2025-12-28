package ds

import (
	"compendium_s/models"
	"context"
)

func (c *Client) CheckRoleDs(guildId, memderId, roleid string) bool {
	req := &CheckRoleRequest{
		Guild:    guildId,
		Memberid: memderId,
		Roleid:   roleid,
	}
	fr, err := c.client.CheckRole(context.Background(), req)
	if err != nil {
		c.log.ErrorErr(err)
		return false
	}
	return fr.Flag
}
func (c *Client) GetMembersRoles(guildId string) ([]models.DsMembersRoles, error) {
	req := &GuildRequest{Guild: guildId}
	roles, err := c.client.GetMembersRoles(context.Background(), req)
	if err != nil {
		c.log.ErrorErr(err)
		return nil, err
	}
	var roles2 []models.DsMembersRoles
	for _, memberRole := range roles.Memberroles {
		roles2 = append(roles2, models.DsMembersRoles{
			Userid:  memberRole.Userid,
			RolesId: memberRole.RolesId,
		})
	}
	return roles2, nil
}
func (c *Client) GetRoles(guildId string) ([]models.CorpRole, error) {
	req := &GuildRequest{Guild: guildId}
	roles, err := c.client.GetRoles(context.Background(), req)
	if err != nil {
		c.log.ErrorErr(err)
		return nil, err
	}
	var roles2 []models.CorpRole
	for _, role := range roles.Roles {
		roles2 = append(roles2, models.CorpRole{
			Id:   role.Id,
			Name: role.Name,
		})
	}
	return roles2, nil
}
