package model

import (
	"math/rand"
	"strings"
	"time"
)

var chars = []rune("")
var lowerChars = []rune("abcdefghijklmnopqrstuvwxyz")
var capitalChars = []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ")
var numChars = []rune("0123456789")

var layout = "2006-01-02 15:04:05"

//mysql int range
var intMin = -2147483648
var intMax = 2147483647

func init() {
	rand.Seed(time.Now().UnixNano())
	chars = append(chars, lowerChars...)
	chars = append(chars, capitalChars...)
	chars = append(chars, numChars...)
}

func generateRandomString(n int) string {
	str := make([]rune, n)
	for i := range str {
		str[i] = chars[rand.Intn(len(chars))]
	}
	return string(str)
}

func generateRandomInt(n int) string {
	//intMaxとnの桁数で場合分け
	rand.Seed(time.Now().UnixNano())
	rand.Intn(intMax-intMin) + intMin

	str := make([]rune, n)
	for i := range str {
		str[i] = numChars[rand.Intn(len(numChars))]
	}
	return string(str)
}

func generateRandomTinyint() string {
	str := make([]rune, 1)
	for i := range str {
		str[i] = numChars[rand.Intn(len(numChars))%2]
	}
	return string(str)
}

func generateRandomDate() string {
	min := time.Date(1971, 1, 0, 0, 0, 0, 0, time.UTC).Unix() //the min of timestamp in mysql is 1970-01-01
	max := time.Date(2037, 1, 0, 0, 0, 0, 0, time.UTC).Unix() //2038 problem for mysql timestamp
	delta := max - min

	sec := rand.Int63n(delta) + min
	return time.Unix(sec, 0).Format(layout)
}

// TODO: mod random data
func generateRandomJson() string {
	str := make([]rune, 10)
	for i := range str {
		str[i] = numChars[rand.Intn(len(numChars))]
	}
	return strings.Join([]string{`{\\"json\\":"`, string(str), `"}`}, "")
}
