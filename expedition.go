package main

import (
	"github.com/alaingilbert/ogame"
	"log"
	"time"
)

func sendExpedition(bot *ogame.OGame, coord ogame.Coordinate, log log.Logger) {
	ships := ogame.ShipsInfos{SmallCargo: 0, LightFighter: 0, LargeCargo: 1666, EspionageProbe: 1, Pathfinder: 1,
		Battleship: 1, Reaper: 1, Battlecruiser: 1}
	log.Println("expedition.ank: lets start")
	time.Sleep(time.Duration(5+random(10)) * time.Second)

	for i := 0; true; i++ {
		log.Println("expedition.ank: loop start: ", i)

		slots := bot.GetSlots()
		reserved := reserverslot
		if slots.Total-slots.InUse <= reserved {
			log.Println("no avaible slot(1):", slots.InUse, "/", slots.Total, " - ", slots.ExpInUse, "/", slots.ExpTotal, " - reserved: ", reserverslot)
			time.Sleep((6 * time.Minute))
			continue
		}
		// Find next expedition fleet that will come back
		bigNum := 99999
		minSecs := int64(bigNum)

		// Sends new expeditions
		expeditionsPossible := slots.ExpTotal - slots.ExpInUse
		for expeditionsPossible > 0 {
			slots = bot.GetSlots()
			if slots.Total-slots.InUse <= reserved {
				log.Println("no avaible slot (2):", slots.InUse, "/", slots.Total, " - ", slots.ExpInUse, "/", slots.ExpTotal, " - reserved: ", reserverslot)
				time.Sleep((6 * time.Minute))
				continue
			}
			//checkTime()
			//return
			if time.Now().Hour() >= stoptimeExpedition {
				log.Println("time to stop exp: ", time.Now())
				return
			}
			//return
			newFleet, err := sendExpeditionfleet(bot, coord, ships)
			if err != nil {
				log.Println(err)
				break
			} else {
				if minSecs > newFleet.BackIn {
					minSecs = newFleet.BackIn
				}
				expeditionsPossible--
			}
			time.Sleep(time.Duration(15+random(20)) * time.Second)
		}
		fleets, slots := bot.GetFleets()
		for _, fleet := range fleets {
			if fleet.Mission == ogame.Expedition {
				if minSecs > fleet.BackIn {
					minSecs = fleet.BackIn
				}
			}
		}
		// If we didn't found any expedition fleet and didn't create any, let's wait 5min
		if minSecs == int64(bigNum) {
			minSecs = 15 * 60
		}

		time.Sleep(time.Duration(minSecs+60) * time.Second) // Sleep until one of the expedition fleet come back
	}
}
func sendExpeditionfleet(bot *ogame.OGame, coord ogame.Coordinate, ships ogame.ShipsInfos) (ogame.Fleet, error) {
	origin, err := bot.GetCelestial(coord)

	if err != nil {
		log.Println(err)
	}
	destination := ogame.Coordinate{coord.Galaxy, coord.System, 16, ogame.PlanetType}

	//ran := random(5)
	//log.Println(ran)
	//// if ran == 6{
	////     destination = NewCoordinate(destination.Galaxy, destination.System +2, 16, PLANET_TYPE)
	//// }
	//// if ran == 1{
	////     destination = NewCoordinate(destination.Galaxy, destination.System-2, 16, PLANET_TYPE)
	//// }
	//if ran == 5{
	//	destination = NewCoordinate(destination.Galaxy, destination.System+1, 16, PLANET_TYPE)
	//}
	//if ran == 2{
	//	destination = NewCoordinate(destination.Galaxy, destination.System-1, 16, PLANET_TYPE)
	//}
	// Print("random:", ran)

	fleet := ogame.NewFleetBuilder(bot)
	fleet.SetOrigin(origin.GetCoordinate())
	fleet.SetDestination(destination)
	fleet.SetSpeed(ogame.HundredPercent)
	fleet.SetMission(ogame.Expedition)
	fleet.SetShips(ships)
	fleet.SetDuration(1)
	f, err2 := fleet.SendNow()
	if err2 != nil {
		log.Println(err2)
	}
	return f, err2
}
