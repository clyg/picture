package main

import (
	"fmt"
	"os"
	"time"
	"tk_calib/tkcore"
)

var argsCount int = 2

func main() {
	start := time.Now()
	tkPath := input()
	fmt.Println("tkPath: ", tkPath)

	tkData := tkcore.MakeTestDataAndWrite()
	tkData, err := tkData.GetTimeRangeData(tkData.GetTimestampForIndex(2), int64(4*tkcore.Microsecond))
	if err != nil {
		fmt.Println(err)
		return
	}
	err = tkData.Write(".\\timestampRangeData.csv")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("耗时: 	", time.Since(start))
}

func input() []string {
	if len(os.Args) < argsCount {
		fmt.Println("参数不足")
		os.Exit(1)
	}
	fmt.Println("START!")
	return os.Args[1:]
}
