package recorder

import (
	"fmt"
	"sync"
	"time"

	"github.com/james-wilder/betdaq/client"
	"github.com/james-wilder/betdaq/model"
)

// Example from API:
//                  2018-04-18T12:50:00.0000000+00:00
const TimeFormat = "2006-01-02T15:04:05.0000000+00:00"

func Recorder(c *betdaq.BetdaqClient, marketType model.MarketType, wg sync.WaitGroup) {
	defer wg.Done()

	t, err := time.Parse(TimeFormat, marketType.StartTime)
	if err != nil {
		fmt.Println(err)
		return
	}

	duration := time.Until(t)
	if duration < time.Minute*15 {
		fmt.Println("    Skip market", duration)
		return
	}

	fmt.Println("    Do something in", duration)

	time.Sleep(duration)

	// TODO: get market prices every second until it starts
	//getMarketInformation, err := c.GetMarketInformation(model.GetMarketInformation{
	//	GetMarketInformationRequest: model.GetMarketInformationRequest{
	//		MarketIds: []int64{
	//			marketType.Id,
	//		},
	//	},
	//})
	//if err != nil {
	//	log.Println("Couldn't do GetMarketInformation")
	//	log.Fatal(err)
	//}
	//
	//for _, market := range getMarketInformation.GetMarketInformationResult.Markets {
	//	fmt.Println(market.Id, market.Name, market.Type, market.Status, market.StartTime)
	//	for _, selection := range market.Selections {
	//		fmt.Println("  ", selection.Id, selection.Name, selection.Status)
	//	}
	//}

	// TODO: save it to disk
}
