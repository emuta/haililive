package uniqueid

import (
	"time"
	"fmt"
	"strconv"
)

func New(t time.Time, code string) int64 {
	str := fmt.Sprintf("%d%s", t.Unix(), code)
	val, _ := strconv.ParseInt(str, 10, 64)
	return val
}