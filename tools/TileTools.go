package tools

import (
	"fmt"
	"math/rand"
	"time"
)

func RandomTiles(originalArray []int) []int {
	copies := 4
	copiedArrays := make([][]int, copies)

	for i := 0; i < copies; i++ {
		copiedArrays[i] = make([]int, len(originalArray))
		copy(copiedArrays[i], originalArray)
	}

	// 打乱顺序
	rand.Seed(time.Now().UnixNano())
	for i := range copiedArrays {
		rand.Shuffle(len(copiedArrays[i]), func(j, k int) {
			copiedArrays[i][j], copiedArrays[i][k] = copiedArrays[i][k], copiedArrays[i][j]
		})
	}

	// 合并打乱后的数组
	var mergedArray []int
	for _, arr := range copiedArrays {
		mergedArray = append(mergedArray, arr...)
	}
	fmt.Println("洗牌结果:", mergedArray)

	return mergedArray
}
