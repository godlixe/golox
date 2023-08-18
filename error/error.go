package errorx

import (
	"fmt"
	"strconv"
)

func Error(line int, message string) {
	Report(line, "", message)
}

func Report(line int, where string, message string) {
	fmt.Println("[line " + strconv.Itoa(line) + "] Error" + where + ": " + message)
}
