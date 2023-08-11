package models

import (
	"fmt"
	"mj/tools"
)

type Wall struct {
	Items  []Tile // 队列元素
	Length int    // 队列元素个数
}

// NewSequentialQueue 初始化队列
func (queue *Wall) Init(tiles []int) *Wall {
	initTiles := []Tile{}
	for _, t := range tools.RandomTiles(tiles) {
		initTiles = append(initTiles, Tile{Number: t})
	}

	return &Wall{
		Items:  initTiles,
		Length: len(initTiles),
	}
}

// IsEmpty 判断队列是否为空
func (queue *Wall) IsEmpty() bool {
	return queue.Length == 0
}

// Pop 将该队列首元素弹出并返回
func (queue *Wall) Pop() Tile {
	if queue.IsEmpty() {
		fmt.Println("空队列")
	}
	fmt.Println("牌池剩余牌数：", queue.Length)
	item := queue.Items[0]
	queue.Items = queue.Items[1:]
	queue.Length--
	return item
}

func (queue *Wall) Pop13() []Tile {
	if queue.IsEmpty() {
		fmt.Println("空队列")
	}
	item := queue.Items[0:13]
	queue.Items = queue.Items[13:]
	queue.Length -= 13
	return item
}

// Back 获取该队列尾元素
func (queue *Wall) Back() Tile {
	if queue.IsEmpty() {
		fmt.Println("空队列")
	}
	return queue.Items[queue.Length-1]
}
