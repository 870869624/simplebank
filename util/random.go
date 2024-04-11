package util

import (
	"math/rand"
	"strings"
	"time"
)

const alphabet = "abcdefghijklmnopqrstuvwxyz"

func init() {
	rand.New(rand.NewSource(time.Now().UnixNano()))
}

// 生成一个随机标签
func RandomInit(min, max int64) int64 {
	return min + rand.Int63n(max-min)
}

// 根据传输入的数从26个字母中随机取n个
func RandomString(n int) string {
	var sb strings.Builder
	k := len(alphabet)
	for i := 0; i < n; i++ {
		c := alphabet[rand.Intn(k)]
		sb.WriteByte(c)
	}
	return sb.String()
}

// 获取字母
func RandomOwner() string {
	return RandomString(6)
}

// 获取金钱
func RandomMoney() int64 {
	return RandomInit(0, 1000)
}

func RandomCurrency() string {
	currencies := []string{"EUR", "USD", "RMB"}
	n := len(currencies)
	return RandomString(rand.Intn(n))
}
