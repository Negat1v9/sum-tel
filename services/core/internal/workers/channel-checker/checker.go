package channelchecker

import (
	"context"
	"time"

	channelservice "github.com/Negat1v9/sum-tel/services/core/internal/channel/service"
	"github.com/Negat1v9/sum-tel/services/core/internal/model"
	"github.com/Negat1v9/sum-tel/shared/logger"
)

type CatcherNews struct {
	log       *logger.Logger
	chService *channelservice.ChannelService
}

func NewCatcherNews(log *logger.Logger, chService *channelservice.ChannelService) *CatcherNews {
	return &CatcherNews{
		log:       log,
		chService: chService,
	}
}

// panic on error loadAllChannels
func (c *CatcherNews) Start(ctx context.Context) {
	c.log.Infof("CatcherNews.Start: News catcher worker started")
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			c.log.Infof("CatcherNews.Start: News catcher worker received stop signal")
			return
		case <-ticker.C:
			ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
			channelsToParse, err := c.chService.ChannelsToParse(ctx, 100, 0) // todo: pagination
			cancel()
			c.log.Debugf("CatcherNews.Start: len(channelToParse): %d", len(channelsToParse))
			if err != nil {
				c.log.Errorf("CatcherNews.Start: Error fetching channels to parse: %v", err)
				continue
			}
			for _, ch := range channelsToParse {
				go c.parseChannel(&ch)
			}

		}
	}

}

func (c *CatcherNews) parseChannel(ch *model.Channel) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()

	// Fetch and parse the channel's messages
	err := c.chService.ParseChannel(ctx, ch.ID, ch.Username)
	if err != nil {
		c.log.Errorf("CatcherNews.parseChannel: Error parsing channel %d: %v", ch.ID, err)
		return
	}
	c.log.Debugf("CatcherNews.parseChannel: Successfully parsed channel %d", ch.ID)
}
