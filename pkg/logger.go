package pkg

import (
	"fmt"
	"io"
	"log"
	"os"
	"time"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
)

var (
	InputLocationLogger  *log.Logger
	OutputLocationLogger *log.Logger
)

func InitLogger() {
	if err := createLogFolder("../log"); err != nil {
		log.Fatal(err.Error())
	}
	rl, err := rotatelogs.New(
		"../log/server_%Y%m%d.log",
		rotatelogs.WithRotationTime(time.Hour*time.Duration(24)), //rotation time
		rotatelogs.WithRotationCount(7),                          //max backup count
	)
	fmt.Println(os.Getwd())
	if err != nil {
		fmt.Println("error: rotatelog.New()")
		return
	}
	log.SetFlags(log.Ldate | log.Ltime)
	log.SetOutput(rl)
}

func InitTimestampLogger() {
	if err := createLogFolder("../gps-in"); err != nil {
		log.Fatal(err.Error())
	}
	inputLocationLogFile, err := os.OpenFile("../gps-in/gps-in.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatalf("Failed to open access log file: %s", err.Error())
	}
	if err := createLogFolder("../gps-out"); err != nil {
		log.Fatal(err.Error())
	}
	outputLocationLogFile, err := os.OpenFile("../gps-out/gps-out.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatalf("Failed to open access log file: %s", err.Error())
	}

	// ログ出力を分ける
	InputLocationLogger = log.New(io.MultiWriter(inputLocationLogFile, os.Stdout), "", log.LstdFlags)
	OutputLocationLogger = log.New(io.MultiWriter(outputLocationLogFile, os.Stderr), "", log.LstdFlags)
}

func createLogFolder(folderName string) error {
	_, err := os.Stat(folderName)
	if os.IsNotExist(err) {
		err := os.Mkdir(folderName, os.ModePerm)
		if err != nil {
			return fmt.Errorf("フォルダの作成に失敗しました: %w", err)
		}
		fmt.Println("フォルダを作成しました:", folderName)
	} else if err != nil {
		return fmt.Errorf("フォルダの状態確認に失敗しました: %w", err)
	} else {
		fmt.Println("フォルダは既に存在します:", folderName)
	}
	return nil
}

func CreateLogFile(logFilePath string) (*log.Logger, func(), error) {
	logFile, err := os.Create(logFilePath)
	if err != nil {
		log.Printf("Error creating log file: %v\n", err.Error())
		return nil, nil, err
	}

	cleanUp := func() {
		logFile.Close()
		// ファイルサイズを確認
		fileInfo, err := os.Stat(logFilePath)
		if err != nil {
			log.Printf("Error checking log file: %v\n", err)
			return
		}
		fmt.Println(fileInfo.Size())
		if fileInfo.Size() == 0 {
			if err := os.Remove(logFilePath); err != nil {
				log.Printf("Error deleting empty log file: %v\n", err)
			} else {
				log.Printf("Deleted empty log file: %s\n", logFilePath)
			}
		}
	}

	return log.New(logFile, "", log.LstdFlags), cleanUp, nil
}
