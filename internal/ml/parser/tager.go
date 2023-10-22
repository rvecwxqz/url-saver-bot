package parser

import (
	"context"
	"errors"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"sync"
	"time"
	pb "url-saver-bot/internal/proto"
	"url-saver-bot/internal/storage"
)

const (
	noDataTag = "page have no data"
	errorTag  = "error while parsing page"
)

type TagWorker struct {
	buff        []storage.Page
	maxBuffSize int
	ch          chan storage.Page
	errChan     chan error
	ticker      *time.Ticker
	parser      parser
	storage     storage.Storage
	ctx         context.Context
}

func NewTagWorker(ctx context.Context, s storage.Storage, maxBufferSize int) *TagWorker {
	w := &TagWorker{
		buff:        make([]storage.Page, 0, maxBufferSize),
		maxBuffSize: maxBufferSize,
		ch:          make(chan storage.Page),
		errChan:     make(chan error),
		parser:      NewParser(),
		storage:     s,
		ticker:      time.NewTicker(3 * time.Second),
		ctx:         ctx,
	}

	go func() {
		for {
			select {
			case val := <-w.ch:
				w.buff = append(w.buff, val)
				if len(w.buff) == w.maxBuffSize {
					go w.processPages()
				}
			case <-w.ticker.C:
				if len(w.buff) > 0 {
					go w.processPages()
				}
			}
		}
	}()
	go func() {
		for err := range w.errChan {
			log.Println(err)
		}
	}()

	return w
}

func (w *TagWorker) AppendPage(page storage.Page) {
	w.ch <- page
}

func (w *TagWorker) processPages() {
	pages := w.buff
	w.buff = w.buff[:0]

	conn, err := grpc.Dial(":3233", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Println(err)
	}
	defer conn.Close()
	client := pb.NewBertClassifierClient(conn)

	var wg sync.WaitGroup
	wg.Add(len(pages))
	for i := 0; i < len(pages); i++ {
		go func(page *storage.Page) {
			tag := ""
			text, err := w.parser.parse(page.URL)
			var e *NoDataError
			if errors.As(err, &e) {
				tag = noDataTag
			} else if err != nil {
				tag = errorTag
				w.errChan <- err
			}

			req := &pb.PredictRequest{Text: text}
			resp, err := client.Predict(w.ctx, req)
			if err != nil {
				tag = noDataTag
				w.errChan <- err
			}
			if tag == "" {
				tag = resp.Prediction
			}
			page.Tags = tag

			wg.Done()
		}(&pages[i])
	}
	wg.Wait()

	go func(pages []storage.Page) {
		err := w.storage.BatchUpdate(w.ctx, pages)
		if err != nil {
			w.errChan <- fmt.Errorf("tag worker update error: %w", err)
		}
	}(pages)
}
