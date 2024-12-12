package pkg

import (
	"fmt"
	"log"
	"time"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
)

func InitLogger() {
	rl, err := rotatelogs.New(
		"../log/server_%Y%m%d.log",
		rotatelogs.WithRotationTime(time.Hour*time.Duration(24)), //rotation time
		rotatelogs.WithRotationCount(7),                          //max backup count
	)
	if err != nil {
		fmt.Println("error: rotatelog.New()")
		return
	}
	log.SetFlags(log.Ldate | log.Ltime)
	log.SetOutput(rl)
}
