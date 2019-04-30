package main

import (
	"github.com/jasonlvhit/gocron"
)

func main() {
	gocron.Every(1).Hour().Do(scheduleTask)
	<-gocron.Start()
}
