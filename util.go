package main

import (
	"fmt"
	"github.com/alaingilbert/ogame"
	"github.com/faiface/beep"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/valyala/fastrand"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
	"reflect"
	"time"
)

type conf struct {
	Host               string `yaml:"host"`
	User               string `yaml:"user"`
	Pwd                string `yaml:"pwd"`
	Uni                string `yaml:"uni"`
	Uni2               string `yaml:"uni2"`
	Lang               string `yaml:"lang"`
	Dbname             string `yaml:"dbname"`
	FleetSlotsReserved int64  `yaml:"FleetSlotsReserved"`
	Teletoken          string `yaml:"teletoken"`
	Teleid             int64  `yaml:"teleid"`
}

func findhomeplanet(bot *ogame.OGame) ogame.Moon {
	moons := bot.GetMoons()
	result := moons[0]
	shipini, _ := result.GetShips()
	value := shipini.FleetValue()
	for _, cel := range bot.GetMoons() {
		ships, _ := cel.GetShips()
		valuen := ships.FleetValue()
		if valuen > value {
			value = valuen
			result = cel
		}
	}
	return result
}

func random(x int) int64 {
	return int64(fastrand.Uint32n(uint32(x)) + 1)

}
func Contains(array interface{}, val interface{}) (index int) {
	index = -1
	switch reflect.TypeOf(array).Kind() {
	case reflect.Slice:
		{
			s := reflect.ValueOf(array)
			for i := 0; i < s.Len(); i++ {
				if reflect.DeepEqual(val, s.Index(i).Interface()) {
					index = i
					return
				}
			}
		}
	}
	return
}
func sendtele(message string) {
	var cutil conf
	confutil := cutil.confget()
	token := confutil.Teletoken
	id := confutil.Teleid
	telegrambot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(id, token)
	msg := tgbotapi.NewMessage(id, message)
	telegrambot.Send(msg)
}

func (c *conf) confget() *conf {

	yamlFile, err := ioutil.ReadFile("conf.yaml")
	if err != nil {
		fmt.Println(err.Error())
	}
	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		fmt.Println(err.Error())
	}
	return c
}

func GetFleetSlotsReserved() int64 {
	var c conf
	conf := c.confget()
	return conf.FleetSlotsReserved
}

func printTime() {
	fmt.Printf("%v\n", time.Now())
}

func playsound(times int) {
	f, err := os.Open("12310.mp3")
	if err != nil {
		log.Fatal(err)
	}
	streamer, format, err := mp3.Decode(f)
	if err != nil {
		log.Fatal(err)
	}
	defer streamer.Close()
	speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))
	loop := beep.Loop(times, streamer)

	done := make(chan bool)
	speaker.Play(beep.Seq(loop, beep.Callback(func() {
		done <- true
	})))

	for {
		select {
		case <-done:
			return
		case <-time.After(5 * time.Second):
			speaker.Lock()
			fmt.Println("Being attack", format.SampleRate.D(streamer.Position()).Round(time.Second))
			speaker.Unlock()
		}
	}
}

func loginyourself(params ogame.Params) (*ogame.OGame, error) {
	bot, err := ogame.NewWithParams(params)
	for err != nil {
		printTime()
		fmt.Print("login failed")
		time.Sleep(2 * time.Second)
		bot, err := ogame.NewWithParams(params)
		err = err
		bot = bot
	}
	return bot, err
}

func start(bot *ogame.OGame, conf *conf) {
	e := echo.New()
	e.Use(middleware.CORS())
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			ctx.Set("bot", bot)
			//ctx.Set("version", version)
			//ctx.Set("commit", commit)
			//ctx.Set("date", date)
			return next(ctx)
		}
	})
	e.HideBanner = true
	e.HidePort = true
	e.Debug = false
	e.GET("/", ogame.HomeHandler)
	e.GET("/tasks", ogame.TasksHandler)
	e.GET("/bot/server", ogame.GetServerHandler)
	e.POST("/bot/set-user-agent", ogame.SetUserAgentHandler)
	e.GET("/bot/server-url", ogame.ServerURLHandler)
	e.GET("/bot/language", ogame.GetLanguageHandler)
	e.GET("/bot/empire/type/:typeID", ogame.GetEmpireHandler)
	e.POST("/bot/page-content", ogame.PageContentHandler)
	e.GET("/bot/login", ogame.LoginHandler)
	e.GET("/bot/logout", ogame.LogoutHandler)
	e.GET("/bot/username", ogame.GetUsernameHandler)
	e.GET("/bot/universe-name", ogame.GetUniverseNameHandler)
	e.GET("/bot/server/speed", ogame.GetUniverseSpeedHandler)
	e.GET("/bot/server/speed-fleet", ogame.GetUniverseSpeedFleetHandler)
	e.GET("/bot/server/version", ogame.ServerVersionHandler)
	e.GET("/bot/server/time", ogame.ServerTimeHandler)
	e.GET("/bot/is-under-attack", ogame.IsUnderAttackHandler)
	e.GET("/bot/user-infos", ogame.GetUserInfosHandler)
	e.POST("/bot/send-message", ogame.SendMessageHandler)
	e.GET("/bot/fleets", ogame.GetFleetsHandler)
	e.GET("/bot/fleets/slots", ogame.GetSlotsHandler)
	e.POST("/bot/fleets/:fleetID/cancel", ogame.CancelFleetHandler)
	e.GET("/bot/espionage-report/:msgid", ogame.GetEspionageReportHandler)
	e.GET("/bot/espionage-report/:galaxy/:system/:position", ogame.GetEspionageReportForHandler)
	e.GET("/bot/espionage-report", ogame.GetEspionageReportMessagesHandler)
	e.POST("/bot/delete-report/:messageID", ogame.DeleteMessageHandler)
	e.POST("/bot/delete-all-espionage-reports", ogame.DeleteEspionageMessagesHandler)
	e.POST("/bot/delete-all-reports/:tabIndex", ogame.DeleteMessagesFromTabHandler)
	e.GET("/bot/attacks", ogame.GetAttacksHandler)
	e.GET("/bot/get-auction", ogame.GetAuctionHandler)
	e.POST("/bot/do-auction", ogame.DoAuctionHandler)
	e.GET("/bot/galaxy-infos/:galaxy/:system", ogame.GalaxyInfosHandler)
	e.GET("/bot/get-research", ogame.GetResearchHandler)
	e.GET("/bot/buy-offer-of-the-day", ogame.BuyOfferOfTheDayHandler)
	e.GET("/bot/price/:ogameID/:nbr", ogame.GetPriceHandler)
	e.GET("/bot/moons", ogame.GetMoonsHandler)
	e.GET("/bot/moons/:moonID", ogame.GetMoonHandler)
	e.GET("/bot/moons/:galaxy/:system/:position", ogame.GetMoonByCoordHandler)
	e.GET("/bot/celestials/:celestialID/items", ogame.GetCelestialItemsHandler)
	e.GET("/bot/celestials/:celestialID/items/:itemRef/activate", ogame.ActivateCelestialItemHandler)
	e.GET("/bot/planets", ogame.GetPlanetsHandler)
	e.GET("/bot/planets/:planetID", ogame.GetPlanetHandler)
	e.GET("/bot/planets/:galaxy/:system/:position", ogame.GetPlanetByCoordHandler)
	e.GET("/bot/planets/:planetID/resource-settings", ogame.GetResourceSettingsHandler)
	e.POST("/bot/planets/:planetID/resource-settings", ogame.SetResourceSettingsHandler)
	e.GET("/bot/planets/:planetID/resources-buildings", ogame.GetResourcesBuildingsHandler)
	e.GET("/bot/planets/:planetID/defence", ogame.GetDefenseHandler)
	e.GET("/bot/planets/:planetID/ships", ogame.GetShipsHandler)
	e.GET("/bot/planets/:planetID/facilities", ogame.GetFacilitiesHandler)
	e.POST("/bot/planets/:planetID/build/:ogameID/:nbr", ogame.BuildHandler)
	e.POST("/bot/planets/:planetID/build/cancelable/:ogameID", ogame.BuildCancelableHandler)
	e.POST("/bot/planets/:planetID/build/production/:ogameID/:nbr", ogame.BuildProductionHandler)
	e.POST("/bot/planets/:planetID/build/building/:ogameID", ogame.BuildBuildingHandler)
	e.POST("/bot/planets/:planetID/build/technology/:ogameID", ogame.BuildTechnologyHandler)
	e.POST("/bot/planets/:planetID/build/defence/:ogameID/:nbr", ogame.BuildDefenseHandler)
	e.POST("/bot/planets/:planetID/build/ships/:ogameID/:nbr", ogame.BuildShipsHandler)
	e.POST("/bot/planets/:planetID/teardown/:ogameID", ogame.TeardownHandler)
	e.GET("/bot/planets/:planetID/production", ogame.GetProductionHandler)
	e.GET("/bot/planets/:planetID/constructions", ogame.ConstructionsBeingBuiltHandler)
	e.POST("/bot/planets/:planetID/cancel-building", ogame.CancelBuildingHandler)
	e.POST("/bot/planets/:planetID/cancel-research", ogame.CancelResearchHandler)
	e.GET("/bot/planets/:planetID/resources", ogame.GetResourcesHandler)
	e.POST("/bot/planets/:planetID/send-fleet", ogame.SendFleetHandler)
	e.POST("/bot/planets/:planetID/send-ipm", ogame.SendIPMHandler)
	e.GET("/bot/moons/:moonID/phalanx/:galaxy/:system/:position", ogame.PhalanxHandler)
	e.GET("/game/allianceInfo.php", ogame.GetAlliancePageContentHandler) // Example: //game/allianceInfo.php?allianceId=500127

	// Get/Post Page Content
	e.GET("/game/index.php", ogame.GetFromGameHandler)
	e.POST("/game/index.php", ogame.PostToGameHandler)

	// For AntiGame plugin
	// Static content
	e.GET("/cdn/*", ogame.GetStaticHandler)
	e.GET("/headerCache/*", ogame.GetStaticHandler)
	e.GET("/favicon.ico", ogame.GetStaticHandler)
	e.GET("/game/sw.js", ogame.GetStaticHandler)

	// JSON API
	/*
		/api/serverData.xml
		/api/localization.xml
		/api/players.xml
		/api/universe.xml
	*/
	e.GET("/api/*", ogame.GetStaticHandler)
	e.HEAD("/api/*", ogame.GetStaticHEADHandler) // AntiGame uses this to check if the cached XML files need to be refreshed

	//if enableTLS {
	//	log.Println("Enable TLS Support")
	//	return e.StartTLS(host+":"+strconv.Itoa(port), tlsCertFile, tlsKeyFile)
	//}
	log.Println("Disable TLS Support")
	e.Start(conf.Host)
}
