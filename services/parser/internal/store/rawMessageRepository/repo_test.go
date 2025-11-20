package rawMessageRepository

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/Negat1v9/sum-tel/services/parser/internal/model"
	sqltransaction "github.com/Negat1v9/sum-tel/services/parser/internal/store/sqlTransaction"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/stretchr/testify/require"
)

func TestCreateMessages(t *testing.T) {
	sqlmockDB, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	require.NoError(t, err)

	sqlxDB := sqlx.NewDb(sqlmockDB, "sqlmock")

	repo := NewRawMessageRepository(sqlxDB)
	sqlTx := sqltransaction.NewSqlTransaction(sqlxDB)

	t.Run("CreateMessagesSuccess", func(t *testing.T) {
		created1 := time.Now()
		created2 := time.Now()
		receivedAt1 := time.Now()
		receivedAt2 := time.Now()
		msgs := []model.RawMessage{{ChannelID: "ch-1", ContentType: "text", TelegramMessageID: 1, HTMLText: "<a>text</a>", Status: "new", MediaURLs: pq.StringArray{"t1"}, MessageDate: created1, ReceivedAt: receivedAt1},
			{ChannelID: "ch-2", ContentType: "text", TelegramMessageID: 2, HTMLText: "<a>text</a>", Status: "new", MessageDate: created2, ReceivedAt: receivedAt2}}
		// rows := sqlmock.NewRows([]string{"channel_id", "content_type", "telegram_message_id", "html_text", "status", "media_urls", "message_date", "received_at"}).
		// 	AddRow("ch-1", "text", 1, "<a>text</a>", "new", pq.StringArray{}, time.Now(), time.Now()).
		// 	AddRow("ch-2", "text", 2, "<a>text</a>", "new", pq.StringArray{}, time.Now(), time.Now()).
		// 	AddRow("ch-3", "text", 3, "<a>text</a>", "new", pq.StringArray{}, time.Now(), time.Now())
		var placeholder []string
		for i := range 2 {
			offset := i * 8
			placeholder = append(placeholder, fmt.Sprintf("($%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d)", offset+1, offset+2, offset+3, offset+4, offset+5, offset+6, offset+7, offset+8))
			// args = append(args, msg.ChannelID, msg.ContentType, msg.TelegramMessageID, msg.HTMLText, msg.Status, msg.MediaURLs, msg.MessageDate, msg.ReceivedAt)
		}
		quary := fmt.Sprintf(createMessagesQuery, strings.Join(placeholder, ","))
		mock.ExpectBegin()
		tx, err := sqlTx.StartTx(context.Background())
		require.NoError(t, err)

		mock.ExpectExec(quary).WithArgs("ch-1", "text", 1, "<a>text</a>", "new", pq.StringArray{"t1"}, created1, receivedAt1, "ch-2", "text", 2, "<a>text</a>", "new", nil, created2, receivedAt2).WillReturnResult(sqlmock.NewResult(2, 2))
		err = repo.CreateMessages(context.Background(), tx, msgs)
		require.NoError(t, err)

	})
}

func TestGetChannelMessages(t *testing.T) {
	sqlmockDB, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	require.NoError(t, err)

	sqlxDB := sqlx.NewDb(sqlmockDB, "sqlmock")

	repo := NewRawMessageRepository(sqlxDB)

	t.Run("GetChannelMessagesSuccess", func(t *testing.T) {
		created1 := time.Now()
		created2 := time.Now()
		created3 := time.Now()
		receivedAt1 := time.Now()
		receivedAt2 := time.Now()
		receivedAt3 := time.Now()
		rows := sqlmock.NewRows([]string{"id", "channel_id", "content_type", "telegram_message_id", "html_text", "status", "media_urls", "message_date", "received_at"}).
			AddRow(1, "ch-1", "text", 1, "<a>text</a>", "new", pq.StringArray{"t1", "t2"}, created1, receivedAt1).
			AddRow(2, "ch-1", "text", 2, "<a>text</a>", "new", pq.StringArray{}, created2, receivedAt2).
			AddRow(2, "ch-1", "text", 3, "<a>text</a>", "new", nil, created3, receivedAt3)
		mock.ExpectQuery(getChannelMessagesQuery).WillReturnRows(rows)

		msgs, err := repo.GetChannelMessages(context.Background(), "ch-1", 3, 0)
		require.NoError(t, err)
		require.Equal(t, 3, len(msgs))
	})
}

// TODO: write more tests !
