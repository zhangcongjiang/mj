package models

type TileMessage struct {
	MessageType    int    `json:"message-type"`
	TargetPlayerId string `json:"target-player-id"`
	SourcePlayerId string `json:"source-player-id"`
	Tile           Tile   `json:"tile"`
}
