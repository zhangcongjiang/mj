package main

import (
	"mj/database"
	"mj/models"
	"time"
)

func main() {
	database.Setup()
	g := models.Game{}
	game := g.Init()
	game.Start()
	go game.ReceiveMsg()
	time.Sleep(time.Second * 1000 * 60)

}
