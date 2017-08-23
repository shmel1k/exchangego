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
	if history.Size() == 10 {
		history.RemoveFirst()
	}

	history.Add(currency)

	currency += rand.Intn(10) - 5
	currencyString := strconv.Itoa(currency)
	return currencyString
}

func GetHistory(size int) []int {
	return history.ToSlice(size)
}
