package algorithm

import "sort"

type CardType int

const (
	CardType_T CardType = iota
	CardType_S
	CardType_W
	CardType_F
)

func CheckWin(cards []int) (win bool) {
	//有个七对，特殊判断一下就可以了
	if len(cards) == 14 {
		for i := 0; i < 7; i++ {
			if cards[2*i] == cards[2*i+1] {
				win = true
			} else {
				win = false
				break
			}
		}
	}
	if win {
		return win
	}
	//用一个slice来装不同花色的牌
	diffColorCards := GetDiffColorCards(cards)
	for _, tempCards := range diffColorCards {
		switch len(tempCards) % 3 {
		case 0:
			frequency := findFrequencyOfArray(tempCards)
			win = CheckFor3(frequency)
			if !win {
				return win
			}
		case 1:
			win = false
			return win
		case 2:
			frequency := findFrequencyOfArray(tempCards)
			win = CheckFor2(frequency)
			if !win {
				return win
			}
		}
	}
	return win
}

func CheckFor2(frequency map[int]int) (result bool) {
	// 首先对map的key进行排序
	keys := make([]int, 0, len(frequency))

	for k := range frequency {
		keys = append(keys, k)
	}
	sort.Ints(keys)

	for _, item := range keys {
		//找到对子后，然后进行顺子或者刻子校验，一旦校验成功，即为可以胡牌，如果校验失败，则用下一个对子校验
		count := frequency[item]
		if count >= 2 {
			newFrequency := make(map[int]int)
			for k, v := range frequency {
				newFrequency[k] = v
			}
			newFrequency[item] -= 2
			result = CheckFor3(newFrequency)
			if result {
				return result
			}
		}
	}
	return result
}

func CheckFor3(frequency map[int]int) (result bool) {
	keys := make([]int, 0, len(frequency))

	for k := range frequency {
		keys = append(keys, k)
	}
	sort.Ints(keys)

	for _, item := range keys {
		count := frequency[item]
		switch count {
		case 4:
			frequency[item] = 1
			CheckFor3(frequency)
		case 3:
			delete(frequency, item)
			CheckFor3(frequency)
		case 2:
			if frequency[item+1] >= 2 && frequency[item+2] >= 2 {
				delete(frequency, item)
				frequency[item+1] -= 2
				frequency[item+2] -= 2
				CheckFor3(frequency)
			} else {
				result = false
				return result
			}
		case 1:
			if frequency[item+1] >= 1 && frequency[item+2] >= 1 {
				delete(frequency, item)
				frequency[item+1] -= 1
				frequency[item+2] -= 1
				CheckFor3(frequency)
			} else {
				result = false
				return result
			}
		case 0:
			result = true
		}
	}

	return result
}

// 统计每张牌出现的次数
func findFrequencyOfArray(arr []int) map[int]int {
	frequency := make(map[int]int)
	for _, item := range arr {
		if frequency[item] == 0 {
			frequency[item] = 1
		} else {
			frequency[item]++
		}
	}
	return frequency
}

//获取不同颜色的牌
func GetDiffColorCards(cards []int) map[CardType][]int {
	colorCards := make(map[CardType][]int)
	for _, card := range cards {
		switch card / 10 {
		case 0:
			colorCards[CardType_T] = append(colorCards[CardType_T], card)
		case 1:
			colorCards[CardType_S] = append(colorCards[CardType_S], card)
		case 2:
			colorCards[CardType_W] = append(colorCards[CardType_W], card)
		case 3:
			colorCards[CardType_F] = append(colorCards[CardType_F], card)
		}
	}
	return colorCards
}
