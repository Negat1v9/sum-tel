package tgparser

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

type ParsedMessage struct {
	Type      string
	Text      string
	HtmlText  string
	Link      string
	MsgId     int64
	Date      time.Time
	PhotoUrls []string
}

type ChannelInfo struct {
	Username    string
	Name        string
	Description string
	MsgInterval int32
	Messages    []ParsedMessage
}

const (
	baseUrl = "https://t.me/"
)

type TgParser struct {
	httpClient *http.Client
}

func NewTgParser() *TgParser {
	return &TgParser{
		httpClient: http.DefaultClient,
	}
}

// ParseChannel parses channel info and messages from telegram web page
func (p *TgParser) ParseChannel(ctx context.Context, username string) (*ChannelInfo, error) {
	url := baseUrl + "/s/" + username
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Accept", "application/json, text/javascript, */*; q=0.0")
	req.Header.Add("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/26.0.1 Safari/605.1.15")
	req.Header.Add("Origin", "https://t.me")
	req.Header.Add("Priority", "u=3, i")
	req.Header.Add("Referer", fmt.Sprintf("https://t.me/s/%s", username))
	req.Header.Add("Sec-Fetch-Dest", "empty")
	req.Header.Add("Sec-Fetch-Mode", "cors")
	req.Header.Add("Sec-Fetch-Site", "same-origin")
	res, err := p.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return nil, err
	}

	channelData := doc.Find(".tgme_header_right_column")

	channelInfo := &ChannelInfo{
		Username:    channelData.Find(".tgme_channel_info_header_username").Text(), // @username @ is included
		Name:        channelData.Find(".tgme_channel_info_header_title").Text(),
		Description: channelData.Find(".tgme_channel_info_description").Text(),
		Messages:    p.parseMsg(doc),
	}

	return channelInfo, nil
}

// ParseMessages parses new messages from telegram web page
func (p *TgParser) ParseMessages(ctx context.Context, username string, afterMsgId int64) ([]ParsedMessage, error) {
	url := baseUrl + "/s/" + username + fmt.Sprintf("?after=%d", afterMsgId)
	req, err := http.NewRequestWithContext(ctx, "POST", url, nil)
	if err != nil {
		return nil, err
	}

	// set headers needed for telegram server
	req.Header.Add("Accept", "application/json, text/javascript, */*; q=0.0")
	req.Header.Add("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/26.0.1 Safari/605.1.15")
	req.Header.Add("Origin", "https://t.me")
	req.Header.Add("Priority", "u=3, i")
	req.Header.Add("Referer", fmt.Sprintf("https://t.me/s/%s", username))
	req.Header.Add("Sec-Fetch-Dest", "empty")
	req.Header.Add("Sec-Fetch-Mode", "cors")
	req.Header.Add("Sec-Fetch-Site", "same-origin")
	req.Header.Add("X-Requested-With", "XMLHttpRequest")

	res, err := p.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	// the response looks like a list of divs and is presented as text, with line breaks and " at the end.
	// we remove all unnecessary characters so that the parser works.
	clear := bytes.ReplaceAll(body, []byte{'\\'}, []byte{' '})

	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(clear[1 : len(clear)-1]))
	if err != nil {
		return nil, err
	}

	// doesn't matter whether we received messages, even if there are 0
	// the request was sent successfully.
	return p.parseMsg(doc), nil
}

// get all messages from the document
func (p *TgParser) parseMsg(doc *goquery.Document) []ParsedMessage {
	msgs := make([]ParsedMessage, 0)
	doc.Find(".tgme_widget_message_wrap").Each(func(i int, s *goquery.Selection) {
		msg := ParsedMessage{}

		msgContent := s.Find(".tgme_widget_message_text")
		msgContent.Find(".emoji ").RemoveFiltered("i") // remove only <i> tags with emoji class
		msgContent.Find("tg-emoji").Remove()           // remove all <tg-emoji> tags
		msgContent.Find("a").RemoveAttr("onclick")     // remove onclick attributes from links

		msg.Text = msgContent.Text()
		if msg.Text == "" {
			// skip messages without text content
			return
		}
		msg.HtmlText, _ = msgContent.Html()
		msgLink, ok := s.Find(".tgme_widget_message_footer").Find("a").Attr("href")
		if ok {
			msg.Link = msgLink
			// extract msg id from link
			msg.MsgId, _ = strconv.ParseInt(msgLink[strings.LastIndex(msgLink, "/")+1:], 10, 64)
		}

		msgPublicationDateTime := s.Find(".tgme_widget_message_footer").Find("time").AttrOr("datetime", "")
		dateTime, err := time.Parse(time.RFC3339, msgPublicationDateTime)
		if err != nil {
			msg.Date = time.Now() // FIXME: what???
		}

		msg.Date = dateTime
		msg.Type = "text" // TODO: receive actual type

		msgs = append(msgs, msg)
	})

	return msgs
}
