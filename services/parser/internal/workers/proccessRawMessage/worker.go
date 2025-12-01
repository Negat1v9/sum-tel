package proccessrawmessage

import (
	"context"
	"errors"
	"time"

	"github.com/Negat1v9/sum-tel/services/parser/internal/parser/service"
	"github.com/Negat1v9/sum-tel/shared/logger"
)

type Worker struct {
	log     *logger.Logger
	service *service.ParserService
}

func NewWorker(log *logger.Logger, service *service.ParserService) *Worker {
	return &Worker{
		log:     log,
		service: service,
	}
}

func (w *Worker) Start(ctx context.Context) {
	w.log.Infof("Worker.Start: proccess raw messages worker started")
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			w.log.Warnf("Worker.Start: proccess raw messages is stopping...")
			return
		case <-ticker.C:
			ctx, cancel := context.WithTimeout(ctx, 10*time.Minute)
			err := w.service.ProccessRawMessages(ctx, 25)
			switch {
			case errors.Is(err, service.ErrNoRawMessages):
				w.log.Debugf("Worker.Start: no raw messages to proccess")
				time.Sleep(30 * time.Second)
			case err != nil:
				w.log.Errorf("Worker.Start: error on proccess raw messages: %v", err)
			}
			cancel()
		}
	}
}
