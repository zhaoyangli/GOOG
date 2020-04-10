package main

import (
	"fmt"
	_ "go.uber.org/automaxprocs"
	"io"
	"log"
	"os"
	"time"
)
import _ "go4.org/lock"
import "github.com/alaingilbert/ogame"

var quit chan int = make(chan int)
var reserverslot int64
var logger log.Logger
var stoptime1 = 21
var stoptime2 = 22
var esp_probes_for_scans int64 = 17

func main() {
	var c conf
	conf := c.confget()
	fmt.Println(conf.Host)
	params := ogame.Params{conf.User, conf.Pwd, conf.Uni, conf.Lang,
		0, true, "", "", "", "", false,
		"lobby", "", ""}

	bot, _ := ogame.NewWithParams(params)
	fileName := "ll.log"
	f, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	fmt.Println(err)
	defer f.Close()
	writers := []io.Writer{f, os.Stdout}
	fileAndStdoutWriter := io.MultiWriter(writers...)
	logger.SetOutput(fileAndStdoutWriter)
	logger.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	logger.SetPrefix("[" + bot.Universe + "]")
	logger.Println(bot)
	attacked, _ := bot.IsUnderAttack()

	logger.Println(bot.IsLoggedIn())   // False
	logger.Println(attacked)           // False
	logger.Println(bot.GetUserInfos()) // False

	////////////////////normal////////////////////

	sendtele(bot.Universe + " origin start")
	time.Sleep(5 * time.Second)
	reserverslot = GetFleetSlotsReserved()
	go start(bot, conf)
	go defender(bot, 180, logger)
	coord := ogame.Coordinate{5, 403, 9, ogame.MoonType}
	go sendExpedition(bot, coord, logger)
	go ai(bot, logger)

	logger.Println("to debug")
	for i := 0; i < 3; i++ {
		<-quit
	}
}
