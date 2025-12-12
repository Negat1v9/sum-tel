package newsrepository

import (
	"context"
	"fmt"
	"strings"

	"github.com/Negat1v9/sum-tel/services/core/internal/model"
	"github.com/Negat1v9/sum-tel/shared/sqltransaction"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type NewsRepository struct {
	db *sqlx.DB
}

func NewNewsRepository(db *sqlx.DB) *NewsRepository {
	return &NewsRepository{db: db}
}

func (r *NewsRepository) Create(ctx context.Context, tx sqltransaction.Txx, news *model.News) error {
	_, err := tx.ExecContext(
		ctx,
		createNewsQuery,
		news.ID,
		news.Title,
		news.Summary,
		news.Language,
		news.Category,
	)

	return err
}

func (r *NewsRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.News, error) {
	news := &model.News{}
	err := r.db.GetContext(ctx, news, getNewsByIDQuery, id)
	if err != nil {
		return nil, err
	}

	return news, nil
}

func (r *NewsRepository) GetAll(ctx context.Context, limit, offset int) ([]model.News, error) {
	news := []model.News{}
	err := r.db.SelectContext(ctx, &news, getAllNewsQuery, limit, offset)
	if err != nil {
		return nil, err
	}

	return news, nil
}

func (r *NewsRepository) GetByUserSubscription(ctx context.Context, userID int, limit, offset int) (*model.NewsList, error) {
	var total int
	err := r.db.GetContext(ctx, &total, countNewsByUserSourcesQuary, userID)
	if err != nil {
		return nil, err
	}

	if total == 0 {
		return &model.NewsList{TotalRecords: total}, nil
	}

	rows, err := r.db.QueryxContext(ctx, getNewsByUserSourcesQuary, userID, limit, offset)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	newsList := make([]model.News, 0, limit)
	for rows.Next() {
		var news model.News
		err = rows.Scan(&news.ID, &news.Title, &news.Summary, &news.Language, &news.Category, &news.CreatedAt, &news.NumberOfSources)
		if err != nil {
			return nil, err
		}
		newsList = append(newsList, news)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return &model.NewsList{
		TotalRecords: total,
		News:         newsList,
	}, nil
}
func (r *NewsRepository) Delete(ctx context.Context, id uuid.UUID) (*model.News, error) {
	news := &model.News{}
	err := r.db.GetContext(ctx, news, deleteNewsQuery, id)
	if err != nil {
		return nil, err
	}

	return news, nil
}

func (r *NewsRepository) CreateNewsSource(ctx context.Context, tx sqltransaction.Txx, source *model.NewsSource) error {
	_, err := tx.ExecContext(
		ctx,
		createNewsSourceQuery,
		source.NewsID,
		source.MessageID,
		source.ChannelID,
	)

	return err
}

func (r *NewsRepository) CreateNewsSources(ctx context.Context, tx sqltransaction.Txx, sources []model.NewsSource) error {
	if len(sources) == 0 {
		return nil
	}

	// build dynamic query with VALUES for each source
	query, args := buildCreateNewsSourcesBatchQuery(sources)

	_, err := tx.ExecContext(ctx, query, args...)
	return err
}
func (r *NewsRepository) GetNewsSourcesByNewsID(ctx context.Context, newsID uuid.UUID) ([]model.NewsSource, error) {
	rows, err := r.db.QueryxContext(ctx, getNewsSourcesByNewsIDQuery, newsID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	sources := make([]model.NewsSource, 0, 5)
	for rows.Next() {
		var source model.NewsSource
		err = rows.Scan(&source.ID, &source.MessageID, &source.NewsID, &source.ChannelID, &source.ChannelName)
		if err != nil {
			return nil, err
		}
		sources = append(sources, source)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return sources, nil
}

// build dynamic batch insert query for news sources
func buildCreateNewsSourcesBatchQuery(sources []model.NewsSource) (string, []any) {

	var valuesClauses []string
	var args []any

	for i, source := range sources {
		paramIndex := i * 3
		valuesClauses = append(valuesClauses, fmt.Sprintf("($%d, $%d, $%d)", paramIndex+1, paramIndex+2, paramIndex+3))
		args = append(args, source.NewsID, source.MessageID, source.ChannelID)
	}

	return fmt.Sprintf(createNewsSourcesQueryPrefix, strings.Join(valuesClauses, ", ")), args
}
