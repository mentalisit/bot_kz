package bot

import (
	"fmt"
	"rs/models"
	"rs/storage"
	"sort"
	"strconv"
	"time"
)

func (b *Bot) OptimizationSborkz() {
	time.Sleep(5 * time.Second)
	go NewOpti(b.storage.DbFunc)
}

type Opti struct {
	st storage.DbFunc
}

func NewOpti(st storage.DbFunc) *Opti {
	o := &Opti{st: st}
	o.readAll()
	return o
}

func (o *Opti) readAll() {
	sb := o.st.ReadAllActive()
	fmt.Printf("ReadAllActive %d \n", len(sb))

	time.Sleep(2 * time.Second)
	o.sCorp(sb)
	fmt.Println("Optimization Done")
}
func (o *Opti) sCorp(sb []models.Sborkz) {
	sort.Slice(sb, func(i, j int) bool {
		return sb[i].Corpname < sb[j].Corpname
	})
	corp := make(map[string][]models.Sborkz)
	for _, sborkz := range sb {
		corp[sborkz.Corpname] = append(corp[sborkz.Corpname], sborkz)
	}

	for _, sborkzs := range corp {
		time.Sleep(5 * time.Second)
		o.sNames(sborkzs)
	}
}

func (o *Opti) sNames(sb []models.Sborkz) {
	sort.Slice(sb, func(i, j int) bool {
		return sb[i].Name < sb[j].Name
	})
	name := make(map[string][]models.Sborkz)
	for _, sborkz := range sb {
		name[sborkz.Name] = append(name[sborkz.Name], sborkz)
	}

	for _, sborkzs := range name {
		o.sLvl(sborkzs)
	}
}

func (o *Opti) sLvl(sb []models.Sborkz) {
	sort.Slice(sb, func(i, j int) bool {
		return sb[i].Lvlkz < sb[j].Lvlkz
	})
	lvl := make(map[string][]models.Sborkz)
	for _, sborkz := range sb {
		lvl[sborkz.Lvlkz] = append(lvl[sborkz.Lvlkz], sborkz)
	}

	for _, sborkzs := range lvl {
		time.Sleep(1 * time.Second)
		o.sNumEvent(sborkzs)
	}
}
func (o *Opti) sNumEvent(sb []models.Sborkz) {
	sort.Slice(sb, func(i, j int) bool {
		return sb[i].Numberevent < sb[j].Numberevent
	})
	nEvent := make(map[int][]models.Sborkz)
	for _, sborkz := range sb {
		nEvent[sborkz.Numberevent] = append(nEvent[sborkz.Numberevent], sborkz)
	}

	for s, sborkzs := range nEvent {
		fmt.Printf("sort CorpName: %s Name: %s level: %s NumEvent: %d",
			sborkzs[0].Corpname, sborkzs[0].Name, sborkzs[0].Lvlkz, s)

		o.sEventPoints(sborkzs)
	}
}

func (o *Opti) sEventPoints(sb []models.Sborkz) {
	points := make(map[string]int)
	srs := make(map[string]models.Sborkz)
	count := make(map[string]int)

	for _, sborkz := range sb {
		points[sborkz.Name] += sborkz.Eventpoints
		count[sborkz.Name]++
		fmt.Printf(".")

		value, exists := srs[sborkz.Name]
		if !exists {
			srs[sborkz.Name] = sborkz
		}
		if value.Id > sborkz.Id {
			o.st.DeleteSborkzId(sborkz.Id)
		} else {
			o.st.DeleteSborkzId(srs[sborkz.Name].Id)
			srs[sborkz.Name] = sborkz
		}
	}
	fmt.Printf("\n")

	for Name, sborkz := range srs {
		if count[Name] > 1 {
			if points[sborkz.Name] != 0 {
				fmt.Printf("updateNamesEvent Name: %s Count: %d Points:%d\n\n",
					Name, count[Name], points[Name])
			} else {
				fmt.Printf("updateNames Name: %s Count: %d \n\n",
					Name, count[Name])
			}

			active := strconv.Itoa(count[Name])
			Points := points[Name]
			o.st.UpdateSborkzPoints(active, sborkz.Id, Points)
		}
	}
}
