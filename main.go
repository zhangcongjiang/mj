package main

import (
	"mj/models"
	"time"
)

func main() {
	//database.Setup()
	g := models.Game{}
	game := g.Init()
	game.Start()
	go game.ReceiveMsg()
	time.Sleep(time.Second * 1000 * 60)

}

//func main() {
//
//	list := []int{1, 1, 4, 5, 6, 7, 9, 11, 12, 13, 14, 14, 14, 33}
//	result := algorithm.CheckWin(list)
//	fmt.Println(result)
//}
