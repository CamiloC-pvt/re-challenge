package main

import (
	"flag"
	"log"
	"os"
	"strconv"

	"github.com/CamiloC-pvt/re-challenge/app"
)

func main() {
	flagPort := flag.String("port", "3001", "Define the port to use")
	flag.Parse()

	strPort := os.Getenv("PORT")
	if strPort == "" {
		strPort = *flagPort
	}

	port, err := strconv.ParseInt(strPort, 10, 32)
	if err != nil {
		log.Fatalf("Cannot read specified port: %s\n", strPort)
	}

	reChallengeApp := app.InitReChallenge(int(port))
	reChallengeApp.StartReChallenge()
}
