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

func (r *NewsRepository) DeleteNewsSource(ctx context.Context, id int) (*model.NewsSource, error) {
	source := &model.NewsSource{}
	err := r.db.GetContext(ctx, source, deleteNewsSourceQuery, id)
	if err != nil {
		return nil, err
	}

	return source, nil
}

func (r *NewsRepository) DeleteNewsSourcesByNewsID(ctx context.Context, newsID uuid.UUID) error {
	_, err := r.db.ExecContext(ctx, deleteNewsSourcesByNewsIDQuery, newsID)
	return err
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
