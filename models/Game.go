package models

import (
	"fmt"
	uuid "github.com/satori/go.uuid"
	"mj/algorithm"
	"sort"
	"strconv"
	"strings"
	"time"
)

var (
	gameCache         = map[string]Game{}
	newTileChanCache  = map[string]ChanContainer{}
	dropTileChanCache = map[string]ChanContainer{}
	playerHandTile    = map[string][]Tile{}
	playerDropTile    = map[string][]Tile{}
	playerShowTile    = map[string][]Tile{}
)

type Game struct {
	Id        string   `json:"id"`
	Players   []Player `json:"players"`
	TileQueue Wall     `json:"tile-queue"`
	DropTiles []Tile   `json:"drop-tiles"`
}
type ChanContainer struct {
	Ch chan TileMessage
}

func (g *Game) Init() *Game {
	queue := Wall{}
	tiles := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 11, 12, 13, 14, 15, 16, 17, 18, 19, 31, 32, 33}
	wall := queue.Init(tiles)

	var players []Player
	gameId := uuid.NewV4().String()
	dropTileChanCache[gameId] = ChanContainer{Ch: make(chan TileMessage)}
	for i := 1; i <= 3; i++ {
		playerId := strconv.Itoa(i)
		var handTiles []Tile
		startTiles := wall.Pop13()
		for _, t := range startTiles {
			handTiles = append(handTiles, t)
		}
		playerHandTile[playerId] = g.SortTiles(handTiles)
		playerDropTile[playerId] = []Tile{}
		playerShowTile[playerId] = []Tile{}
		fmt.Println(fmt.Sprintf("用户%s,起手手牌：%s", playerId, g.TilesToString(handTiles)))

		player := &Player{
			Id:          playerId,
			CurrentGame: gameId,
		}

		newTileChanCache[playerId] = ChanContainer{Ch: make(chan TileMessage)}
		go player.ReceiveMsg()

		players = append(players, *player)

	}

	game := &Game{
		Id:        gameId,
		Players:   players,
		TileQueue: *wall,
	}

	gameCache[gameId] = *game
	return game
}

func (g *Game) Start() {
	fmt.Println("游戏开始，管理员准备发牌")
	tileQueue := g.TileQueue
	players := g.Players
	time.Sleep(1)
	player1 := players[0]
	msg := &TileMessage{
		MessageType:    4,
		TargetPlayerId: player1.Id,
		Tile:           tileQueue.Pop(),
	}
	g.TileQueue = tileQueue
	channel := newTileChanCache[player1.Id].Ch
	fmt.Println(fmt.Sprintf("游戏管理员给用户 %s 发了一张牌：%d", player1.Id, msg.Tile.Number))
	channel <- *msg

}

func (g *Game) Shuffler(target string, tile Tile) {
	// 处理给player发牌逻辑
}

func (g *Game) ReceiveMsg() {
	for msg := range dropTileChanCache[g.Id].Ch {

		switch msg.MessageType {
		case 5:
			//接收到弃牌信息
			currentPlayer := msg.TargetPlayerId

			next := g.getNextPlayerId(currentPlayer)

			for {
				if next != currentPlayer {
					checkTiles := append(playerHandTile[next], msg.Tile)
					checkTiles = g.SortTiles(checkTiles)

					if algorithm.CheckWin(g.TilesToInt(checkTiles)) {
						fmt.Println(fmt.Sprintf("用户 %s 可以胡牌，他的手牌是：%s", next, g.TilesToString(checkTiles)))
						winMsg := TileMessage{
							MessageType:    3,
							SourcePlayerId: currentPlayer,
							TargetPlayerId: next,
							Tile:           msg.Tile,
						}
						channel := newTileChanCache[next].Ch
						channel <- winMsg
						time.Sleep(1 * time.Second)
					}

				} else {
					break
				}
				next = g.getNextPlayerId(next)
			}
			tileQueue := g.TileQueue
			if tileQueue.IsEmpty() {
				fmt.Println("游戏结束，平局")
				g.Stop()
			}

			nextPlayerId := g.getNextPlayerId(currentPlayer)
			newTile := tileQueue.Pop()
			newTileMsg := TileMessage{
				MessageType:    4,
				TargetPlayerId: nextPlayerId,
				Tile:           newTile,
			}
			channel := newTileChanCache[nextPlayerId].Ch
			fmt.Println(fmt.Sprintf("游戏管理员给用户 %s 发了一张牌：%d", nextPlayerId, newTile.Number))
			time.Sleep(time.Second * 1)
			channel <- newTileMsg
			g.TileQueue = tileQueue
		}
	}
	// 处理player的操作信息，例如用户碰牌后，下一张牌则发给碰牌用户的下家，用户杠牌后，从牌堆末尾给他发一张牌，
	// 用户弃牌后，给下家发一张牌，用户胡牌后，游戏结束
}

func (g *Game) Stop() {
	delete(gameCache, g.Id)
	close(dropTileChanCache[g.Id].Ch)
	delete(dropTileChanCache, g.Id)
	for _, player := range g.Players {
		ch := newTileChanCache[player.Id].Ch
		close(ch)
		delete(newTileChanCache, player.Id)
	}

}

func (g *Game) getNextPlayerId(playerId string) string {
	players := g.Players
	index := 0
	for i := 0; i < len(players); i++ {
		if players[i].Id == playerId {
			index = i
			break
		}
	}
	if index == len(players)-1 {
		return players[0].Id
	} else {
		return players[index+1].Id
	}
}

func (g *Game) SortTiles(tiles []Tile) []Tile {
	sort.Slice(tiles, func(i, j int) bool {
		return tiles[i].Number < tiles[j].Number
	})
	return tiles

}
func (g *Game) TilesToString(tiles []Tile) string {
	var tilesString []string

	for _, h := range tiles {
		tilesString = append(tilesString, strconv.Itoa(h.Number))
	}

	handString := "手牌：" + strings.Join(tilesString, ", ")

	return handString
}

func (g *Game) TilesToInt(tiles []Tile) []int {
	var arr []int

	for _, h := range tiles {
		arr = append(arr, h.Number)
	}

	return arr
}

func (g *Game) GetGameById(gameId string) Game {
	return gameCache[gameId]
}
func (g *Game) GetNewChanById(playerId string) ChanContainer {
	return newTileChanCache[playerId]
}

func (g *Game) GetDropTileByPlayerId(playerId string) []Tile {
	return playerDropTile[playerId]
}

func (g *Game) GetShowTileByPlayerId(playerId string) []Tile {
	return playerShowTile[playerId]
}

func (g *Game) GetHandTileByPlayerId(playerId string) []Tile {
	return playerHandTile[playerId]
}

func (g *Game) SetShowTileByPlayerId(playerId string, tiles []Tile) {
	playerShowTile[playerId] = tiles
}
func (g *Game) SetDropTileByPlayerId(playerId string, tiles []Tile) {
	playerDropTile[playerId] = tiles
}

func (g *Game) SetHandTileByPlayerId(playerId string, tiles []Tile) {
	playerHandTile[playerId] = tiles
}

func (g *Game) GetDropChanById(gameId string) ChanContainer {
	return dropTileChanCache[gameId]
}
