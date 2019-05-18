package main

import "github.com/jasonlvhit/gocron"

func main() {
	gocron.Every(60).Minutes().Do(scheduleTask)
	<-gocron.Start()
}
