package model

import (
	"math/rand"
	"strconv"
	"strings"
	"time"
)

// The functions in this file generate random value.
// They are called by Column, because it is Column that is responsible for generating random data.

var chars = []rune("")
var lowerChars = []rune("abcdefghijklmnopqrstuvwxyz")
var capitalChars = []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ")
var numChars = []rune("0123456789")

const layout = "2006-01-02 15:04:05"

// mysql int range
var intRangeMap = map[ColumnTypeBase][]int{
	Tinyint:   {-128, 127},
	Smallint:  {-32768, 32767},
	Mediumint: {-8388608, 8388607},
	Int:       {-2147483648, 2147483647},
	Bigint:    {-9223372036854775808, 9223372036854775807},
}

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

func generateRandomInt(t ColumnTypeBase, unsigned bool) string {
	var m int
	switch unsigned {
	case true:
		m = rand.Intn(intRangeMap[t][1])
	case false:
		m = rand.Intn(intRangeMap[t][1]-intRangeMap[t][0]) + intRangeMap[t][0]
	}
	return strconv.Itoa(m)
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
	return strings.Join([]string{`{"json":"`, string(str), `"}`}, "")
}
