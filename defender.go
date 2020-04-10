package main

import (
	"github.com/alaingilbert/ogame"
	"log"
	"strconv"
	"time"
)

func defender(bot *ogame.OGame, interval int, log log.Logger) {
	colist := []string{}
	messagel := []string{}
	reserverslot_o := reserverslot
	for {
		underattack, err := bot.IsUnderAttack()
		if err != nil {
			panic(err)
		}
		//log.Println("underattack: ",underattack)
		if underattack {
			log.Print("under attack\n")
			reserverslot = bot.GetSlots().Total / 3
			atts, _ := bot.GetAttacks()
			for _, eachatt := range atts {
				message := bot.Universe + "\n" + eachatt.String()
				coor := eachatt.DestinationName
				if Contains(messagel, message) == -1 {
					go handle(bot, eachatt.Destination)
					colist = append(colist, coor)
					messagel = append(messagel, message)
					/////////////////////
					out := ""
					for _, ship := range ogame.Ships {
						shipID := ship.GetID()
						nbr := eachatt.Ships.ByID(shipID)
						if !(nbr == 0) {
							out = out + shipID.String() + strconv.FormatInt(nbr, 10) + "\n"
						}
					}
					////////////////////////////

					message = message + "\n" +
						"    ArrivalIn: " + strconv.FormatInt(eachatt.ArriveIn, 10) + "\n" +
						"       shipinfo: " + out + "\n" +
						"       shipcost: " + eachatt.Ships.FleetCost().String() + "\n" +
						"       shipvalue: " + strconv.FormatInt(eachatt.Ships.FleetValue(), 10) + "\n"
					sendtele(message)
					playsound(10)
					bot.GetFleetsFromEventList()
				}
			}
		} else {
			log.Println("defender good\n")
			reserverslot = reserverslot_o
		}
		time.Sleep(time.Duration(int64(interval)+random(20)) * time.Second)
	}
}

func handle(bot *ogame.OGame, coord ogame.Coordinate) {
	remaintime := caltime(bot, coord)
	log.Println("time to fs: ", remaintime)
	for {
		if remaintime == -1 {
			return
		}
		if remaintime < 300 {
			fsfleet, err := fs(bot, coord)
			log.Println(fsfleet)
			if err != nil {
				log.Println(err)
			}
			for {
				recalltime := caltime(bot, coord)
				if recalltime == -1 {
					log.Println("time to recall")
					time.Sleep(60 * time.Second)
					bot.CancelFleet(fsfleet.ID)
					return
				}
				time.Sleep(time.Duration(60+random(10)) * time.Second)
			}
		}
		time.Sleep(time.Duration(60+random(10)) * time.Second)
		remaintime = caltime(bot, coord)
	}

}

func caltime(bot *ogame.OGame, coord ogame.Coordinate) int64 {
	cel, _ := bot.GetCelestial(coord)
	atts, _ := bot.GetAttacksUsing(cel.GetID())
	if len(atts) == 0 {
		return -1
	}
	timereturn := int64(99999)
	for _, eachatt := range atts {
		if timereturn > eachatt.ArriveIn {
			timereturn = eachatt.ArriveIn
		}
	}
	return timereturn
}

func fs(bot *ogame.OGame, coord ogame.Coordinate) (ogame.Fleet, error) {
	log.Println("fs start: ", coord)
	origin, err := bot.GetCelestial(coord)
	if err != nil {
		log.Println(err)
	}
	//destination := ogame.Coordinate{coord.Galaxy, coord.System, 16, ogame.PlanetType}
	destinations := bot.GetCachedCelestials()
	destination := destinations[0]

	distance := int64(999999)
	for _, cel := range destinations {
		log.Println(cel.GetCoordinate(), "dis: ", bot.Distance(coord, cel.GetCoordinate()))

		if coord.Equal(cel.GetCoordinate()) {
			log.Println("origin")
			continue
		}
		if distance > bot.Distance(coord, cel.GetCoordinate()) {
			distance = bot.Distance(coord, cel.GetCoordinate())
			destination = cel
			log.Println("Gottcha")
		}
	}
	log.Println("fs to: ", destination.GetCoordinate())
	if destination.GetCoordinate().Equal(origin.GetCoordinate()) {
		destination := ogame.Coordinate{coord.Galaxy, coord.System, 16, ogame.PlanetType}
		fleet := ogame.NewFleetBuilder(bot)
		fleet.SetOrigin(origin.GetCoordinate())
		fleet.SetDestination(destination)
		fleet.SetSpeed(ogame.TenPercent)
		fleet.SetMission(ogame.Spy)
		fleet.SetAllShips()
		fleet.SetAllResources()
		f, err2 := fleet.SendNow()
		return f, err2

	} else {
		fleet := ogame.NewFleetBuilder(bot)
		fleet.SetOrigin(origin.GetCoordinate())
		fleet.SetDestination(destination)
		fleet.SetSpeed(ogame.TenPercent)
		fleet.SetMission(ogame.Park)
		fleet.SetAllShips()
		fleet.SetAllResources()
		f, err2 := fleet.SendNow()
		return f, err2
	}

}
