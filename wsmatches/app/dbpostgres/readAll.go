package dbpostgres

import (
	"ws/models"
)

//func (d *Db) ReadCorpLevelAll() ([]models.LevelCorp, error) {
//	var n []models.LevelCorp
//	sel := "SELECT * FROM kzbot.corplevel"
//	results, err := d.pool.Query(context.Background(), sel)
//	if err != nil {
//		d.log.ErrorErr(err)
//		return nil, err
//	}
//	for results.Next() {
//		var t models.LevelCorp
//		err = results.Scan(&t.CorpName, &t.Level, &t.EndDate, &t.HCorp, &t.Percent)
//		if err != nil {
//			d.log.ErrorErr(err)
//		}
//		n = append(n, t)
//	}
//	return n, nil
//}
//func (d *Db) ReadCorpLevel(CorpName string) (models.LevelCorp, error) {
//	var n models.LevelCorp
//	err := d.pool.QueryRow(context.Background(), "SELECT * FROM kzbot.corplevel WHERE corpname = $1", CorpName).Scan(
//		&n.CorpName, &n.Level, &n.EndDate, &n.HCorp, &n.Percent)
//	if err != nil {
//		return models.LevelCorp{}, err
//	}
//	return n, nil
//}

func (d *Db) ReadCorpsLevelAll() ([]models.LevelCorps, error) {
	ctx, cancel := d.GetContext()
	defer cancel()
	var nn []models.LevelCorps
	sel := "SELECT * FROM kzbot.corpslevel"
	results, err := d.pool.Query(ctx, sel)
	defer results.Close()
	if err != nil {
		d.log.ErrorErr(err)
		return nil, err
	}
	for results.Next() {
		var n models.LevelCorps
		err = results.Scan(&n.CorpName, &n.Level, &n.EndDate, &n.HCorp, &n.Percent, &n.LastUpdate, &n.Relic)
		if err != nil {
			d.log.ErrorErr(err)
		}
		nn = append(nn, n)
	}
	return nn, nil
}

func (d *Db) ReadCorpsLevel(hCorp string) (models.LevelCorps, error) {
	ctx, cancel := d.GetContext()
	defer cancel()
	var n models.LevelCorps
	err := d.pool.QueryRow(ctx, "SELECT * FROM kzbot.corpslevel WHERE hcorp = $1", hCorp).Scan(
		&n.CorpName, &n.Level, &n.EndDate, &n.HCorp, &n.Percent, &n.LastUpdate, &n.Relic)
	if err != nil {
		return models.LevelCorps{}, err
	}
	return n, nil
}
