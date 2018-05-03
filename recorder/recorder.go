package recorder

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"sync"
	"time"

	"github.com/james-wilder/betdaq/client"
	"github.com/james-wilder/betdaq/model"
)

// Example from API:
//                  2018-04-18T12:50:00.0000000+00:00
const TimeFormat = "2006-01-02T15:04:05.0000000+00:00"

type HistoricEvent struct {
	Event  *model.EventClassifierType
	Market *model.MarketType
	Prices []*model.GetPricesResponse
}

func Recorder(c *betdaq.BetdaqClient, event model.EventClassifierType, market model.MarketType, wg *sync.WaitGroup) {
	defer wg.Done()

	if market.Type != 1 {
		log.Println("Not a win market, dropping market", market.Id, event.Name)
		return
	}

	if market.NumberOfWinningSelections != 1 {
		log.Println("Winning selections not 1, dropping market", market.Id, event.Name)
		return
	}

	startTime, err := time.Parse(TimeFormat, market.StartTime)
	if err != nil {
		fmt.Println(err)
		return
	}

	fifteenBeforeStart := startTime.Add(-15 * time.Minute)
	duration := time.Until(fifteenBeforeStart)
	if duration < 0 {
		fmt.Println("Skip market - not got > 15 minutes til start", market.Id, event.Name)
		return
	}

	// TODO: order by start time more sensibly
	// believe it or not this is enough to order by start time
	time.Sleep(duration / 1000000)
	fmt.Println("Recording market in", duration, market.Id, event.Name)

	time.Sleep(duration)

	fmt.Println("Waking up...", market.Id, event.Name)

	var prices []*model.GetPricesResponse
	for !time.Now().After(startTime) {
		log.Println("Requesting prices", market.Id, event.Name)
		getPrices, err := c.GetPrices(model.GetPrices{
			GetPricesRequest: model.GetPricesRequest{
				MarketIds: []int64{
					market.Id,
				},
				ThresholdAmount:              "0",
				NumberAgainstPricesRequired:  3,
				NumberForPricesRequired:      3,
				WantMarketMatchedAmount:      true,
				WantSelectionMatchedDetails:  true,
				WantSelectionsMatchedAmounts: true,
			},
		})
		if err != nil {
			log.Println("Couldn't do GetMarketInformation")
			continue
		} else {
			prices = append(prices, getPrices)
		}

		time.Sleep(1 * time.Second)
	}

	log.Println("Save historic data to disk", market.Id, event.Name)
	ev := HistoricEvent{
		Event:  &event,
		Market: &market,
		Prices: prices,
	}

	b, err := json.MarshalIndent(ev, "", "  ")
	if err != nil {
		log.Println("JSON problem")
		log.Println(err)
		return
	}

	filename := fmt.Sprintf("data/%d.json", market.Id)
	err = ioutil.WriteFile(filename, b, 0644)
	if err != nil {
		log.Println("writing JSON to file failed")
		log.Println(err)
		return
	}
}
