package main

import (
	"github.com/jasonlvhit/gocron"
)

func main() {
	gocron.Every(15).Minutes().Do(scheduleTask)
	<-gocron.Start()
}
