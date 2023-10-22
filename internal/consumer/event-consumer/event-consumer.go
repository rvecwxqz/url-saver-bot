package event_consumer

import (
	"log"
	"sync"
	"time"
	"url-saver-bot/internal/events"
)

type Consumer struct {
	fetcher   events.Fetcher
	processor events.Processor
	batchSize int
}

func New(f events.Fetcher, p events.Processor, b int) Consumer {
	return Consumer{
		fetcher:   f,
		processor: p,
		batchSize: b,
	}
}

func (c *Consumer) Start() error {
	for {
		gotEvents, err := c.fetcher.Fetch(c.batchSize)
		if err != nil {
			log.Printf("[ERR] consumer: %v\n", err.Error())

			time.Sleep(1 * time.Second)
			continue
		}

		if len(gotEvents) == 0 {
			time.Sleep(1 * time.Second)
			continue
		}

		if err = c.handleEvents(gotEvents); err != nil {
			log.Printf("[ERR] consumer handling events: %v\n", err.Error())
			continue
		}
	}
}

func (c *Consumer) handleEvents(e []events.Event) error {
	var wg sync.WaitGroup
	wg.Add(len(e))
	for _, event := range e {
		log.Printf("got event %v", event.Text)

		go func(e events.Event) {
			if err := c.processor.Process(event); err != nil {
				log.Printf("[ERR] error during proccessing event: %v", err.Error())
			}
			wg.Done()
		}(event)
	}
	return nil
}
