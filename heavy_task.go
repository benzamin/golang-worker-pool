package main

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"
)

type HeavyTask struct{}

func NewHeavyTask() HeavyTask {
	return HeavyTask{}
}

func (e *HeavyTask) Run(payload Payload) JobResult {

	var sleepFor = RandInt(20, 50)

	if payload.Params["sleep"] != "" {
		val, _ := strconv.Atoi(payload.Params["sleep"])
		sleepFor = val
	}
	//simulating heavy task by Sleeping for X ms
	time.Sleep(time.Millisecond * time.Duration(sleepFor))

	//send some random failed cases
	if sleepFor%10 == 0 {
		return NewJobResult(nil, fmt.Errorf("error occured while executing heavy task"))
	}

	return NewJobResult(`{"task":"done", "time":"`+strconv.Itoa(sleepFor)+`"}`, nil)

}

func RandInt(min, max int) int {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return r.Intn(max-min) + min
}
