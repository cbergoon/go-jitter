package main

import (
	"fmt"
	"time"

	jitter "github.com/cbergoon/go-jitter"
)

func main() {
	j, err := jitter.NewJitterer("github.com")
	if err != nil {
		fmt.Println(err)
	}
	j.OnFinish = func(s *jitter.Statistics) {
		fmt.Println(s.UncorrectedSD)
		fmt.Println(s.CorrectedSD)
		fmt.Println(s.RttRange)
		fmt.Println(s.RTTS)
	}
	j.SetBlockSampleSize(5)
	j.SetPingerPrivileged(true)
	j.SetPingerTimeout(time.Second * 10)

	j.Run()
}
