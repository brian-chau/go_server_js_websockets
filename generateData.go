package main

import (
	"fmt"
	"math/rand"
	"time"
)

func GenerateData(generatedDataChannel chan []byte){
	for {
		randomSleepTime := rand.Intn(5)
		time.Sleep(time.Duration(randomSleepTime) * time.Second)
		outStr := fmt.Sprintf("field=%d", randomSleepTime)
		generatedDataChannel <- []byte(outStr)
	}
}