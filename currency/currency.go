package currency

import (
	"math/rand"
	"strconv"
)

var currency = 50
var history *CNode

func InitCurrency() {
	history = new(CNode)
}

func UpdateCurrency() string {
	if history.Size() == 20 {
		history.RemoveFirst()
	}

	history.Add(currency)

	var val int
	if rand.Intn(2) == 1 {
		val = 1
	} else {
		val = -1
	}

	currency += val * rand.Intn(7)

	currencyString := strconv.Itoa(currency)
	return currencyString
}

func GetHistory(size int) []int {
	return history.ToSlice(size)
}

func GetCurrency() int {
	return currency
}
