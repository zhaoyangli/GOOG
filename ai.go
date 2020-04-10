package main

import (
	"github.com/alaingilbert/ogame"
	"log"
	"time"
)

func ai(bot *ogame.OGame, log log.Logger) {
	maxenergy := int64(16)
	for {
		isreearch := true
		for _, cel := range bot.GetPlanets() {
			build1, build2, build3, build4 := cel.ConstructionsBeingBuilt()
			if !(build2 == 0) {
				log.Println(cel.Coordinate.String(), "there is sth being built: b1: ", build1, " b2: ", build2, " b3: ", build3, " b4: ", build4, " ")
				continue
			}
			if build4 == 0 {

				cel.BuildTechnology(ogame.EnergyTechnologyID)
				cel.BuildTechnology(ogame.ComputerTechnologyID)
				cel.BuildTechnology(ogame.ImpulseDriveID)
				isreearch = false
			}
			buildings, _ := cel.GetResourcesBuildings()
			res, _ := cel.GetResourcesDetails()
			log.Println(res)
			if res.Metal.Available >= res.Metal.StorageCapacity && buildings.MetalStorage < 10 {
				cel.BuildBuilding(ogame.MetalStorageID)
			}
			if res.Crystal.Available >= res.Crystal.StorageCapacity && buildings.MetalStorage < 9 {
				cel.BuildBuilding(ogame.CrystalStorageID)
			}
			production, num, _ := cel.GetProduction()
			log.Println(cel.Coordinate.String(), "production: ", production, " num: ", num)

			deu := buildings.DeuteriumSynthesizer
			metal := buildings.MetalMine
			cry := buildings.CrystalMine
			energy := buildings.SolarPlant
			missmetal := -(metal - deu - 4)
			misscry := -(cry - deu + 3)
			missenergy := -(energy - deu - 5)
			buildID := ogame.DeuteriumSynthesizerID
			if energy == 16 {
				missenergy = 0
			}
			if missenergy > 0 && energy < maxenergy {
				buildID = ogame.SolarPlantID
			}
			if missmetal > 0 && missmetal > missenergy {
				buildID = ogame.MetalMineID
			}
			if misscry > 0 && misscry > missmetal {
				buildID = ogame.CrystalMineID
			}
			cel.BuildBuilding(buildID)
			log.Println(cel.Coordinate.String(), " missmetal: ", missmetal, " misscry: ", misscry, " missenergy: ", missenergy, " buildID: ", buildID, " ")
			time.Sleep(time.Duration(5+random(10)) * time.Second)
		}
		if !isreearch {
			sendtele(bot.Universe + " no research")
		}
		time.Sleep(time.Duration(50+random(10)) * time.Minute)
	}
}
