package main

import (
	"fmt"
	"log"
	"sync"

	"github.com/james-wilder/betdaq-bot/recorder"
	"github.com/james-wilder/betdaq/client"
	"github.com/james-wilder/betdaq/config"
	"github.com/james-wilder/betdaq/model"
)

const configFilename = "config.json"

func main() {
	fmt.Println("Hello world")

	conf, err := config.ReadConfig(configFilename)
	if err != nil {
		log.Fatal(err)
		panic("Couldn't load config file" + configFilename)
	}

	c := betdaq.NewClient(conf.Username, conf.Password)

	// getOddsLadder(c)
	getEventSubTreeNoSelections(c, 100004) // Horse Racing
}

func getOddsLadder(c *betdaq.BetdaqClient) *model.GetOddsLadderResponse {
	getLadderResponse, err := c.GetOddsLadder(model.GetOddsLadder{
		GetOddsLadderRequest: model.GetOddsLadderRequest{
			PriceFormat: 1,
		},
	})
	if err != nil {
		log.Fatal(err)
		panic("Couldn't get the odds ladder")
	}
	for _, price := range getLadderResponse.GetOddsLadderResult.Ladder {
		fmt.Println(price.Price, price.Representation)
	}

	return getLadderResponse
}

func getEventSubTreeNoSelections(c *betdaq.BetdaqClient, id int64) *model.GetEventSubTreeNoSelectionsResponse {
	getEventSubTreeNoSelections, err := c.GetEventSubTreeNoSelections(model.GetEventSubTreeNoSelections{
		GetEventSubTreeNoSelectionsRequest: model.GetEventSubTreeNoSelectionsRequest{
			EventClassifierIds: []int64{
				id,
			},
		},
	})
	if err != nil {
		log.Fatal(err)
		panic("Couldn't do GetEventSubTreeNoSelections")
	}

	startRecorders(c, getEventSubTreeNoSelections.GetEventSubTreeNoSelectionsResult.EventClassifiers)

	return getEventSubTreeNoSelections
}

func startRecorders(c *betdaq.BetdaqClient, eventClassifiers []model.EventClassifierType) {
	var wg sync.WaitGroup

	for _, eventClassifier := range eventClassifiers {
		//fmt.Println("Event", eventClassifier.Id, eventClassifier.Name)

		for _, market := range eventClassifier.Markets {
			//fmt.Println("  Market", market.Id, market.Name, market.Type, market.StartTime)

			wg.Add(1)
			go recorder.Recorder(c, eventClassifier, market, wg)
		}
		startRecorders(c, eventClassifier.EventClassifiers)
	}

	wg.Wait()
}
