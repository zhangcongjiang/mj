package models

import (
	"fmt"
)

type Player struct {
	Id          string `json:"id"`
	CurrentGame string `json:"current-game"`
}

func (p *Player) ReceiveDropTileMsg() {
	// 接收其他用户的弃牌信息
}

func (p *Player) ReceiveMsg() {
	g := &Game{}
	currentGame := p.CurrentGame
	game := g.GetGameById(currentGame)
	channel := game.GetNewChanById(p.Id).Ch

	for msg := range channel {
		if msg.MessageType == 4 {
			handTiles := game.GetHandTileByPlayerId(p.Id)
			fmt.Println(fmt.Sprintf("抓牌前，用户 %s 当前手牌：%s", p.Id, game.TilesToString(handTiles)))
			newHandTiles := append(handTiles, msg.Tile)
			//todo check win
			game.SetHandTileByPlayerId(p.Id, newHandTiles)
			fmt.Println(fmt.Sprintf("弃牌前，用户 %s 当前手牌：%s", p.Id, game.TilesToString(newHandTiles)))

			p.DropTile()
		} else if msg.MessageType == 5 {
			fmt.Println(fmt.Sprintf("我是用户 %s ,接收到接收到用户 %s 的弃牌信息，弃的牌是 %d", p.Id, msg.SourcePlayerId, msg.Tile.Number))
		
		}

	}

}

func (p *Player) DropTile() {
	g := &Game{}
	gameId := p.CurrentGame
	game := g.GetGameById(gameId)
	handTiles := game.GetHandTileByPlayerId(p.Id)
	dropTiles := game.GetDropTileByPlayerId(p.Id)
	sortedHandTiles := game.SortTiles(handTiles)

	newHandTiles := sortedHandTiles[:len(sortedHandTiles)-1]
	dropTile := sortedHandTiles[len(sortedHandTiles)-1]
	fmt.Println(fmt.Sprintf("用户 %s 弃牌：%d", p.Id, dropTile.Number))
	msg := TileMessage{
		MessageType:    5,
		TargetPlayerId: p.Id,
		Tile:           dropTile,
	}
	dropChannel := game.GetDropChanById(gameId).Ch
	dropChannel <- msg

	dropTiles = append(dropTiles, dropTile)
	game.SetDropTileByPlayerId(p.Id, dropTiles)
	game.SetHandTileByPlayerId(p.Id, newHandTiles)
	fmt.Println(fmt.Sprintf("弃牌后，用户 %s 当前手牌：%s", p.Id, game.TilesToString(newHandTiles)))
}
