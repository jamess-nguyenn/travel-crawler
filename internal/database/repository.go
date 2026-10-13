package database

import (
	"context"
	"gorm.io/gorm"
)

type UpsertOptions struct {
	ConflictColumns []string
	UpdateColumns   []string
	DoNothing       bool
}

type BaseRepository[T any] interface {
	Create(
		ctx context.Context,
		data *T,
	) error

	Upsert(
		ctx context.Context,
		data *[]T,
		options UpsertOptions,
	) error

	Get(
		ctx context.Context,
		conditions map[string]any,
	) (*T, error)

	UpdateByID(
		ctx context.Context,
		ID uint64,
		updates map[string]any,
	) error

	List(
		ctx context.Context,
		listQuery ListQuery,
	) ([]T, int64, error)

	ListAll(
		ctx context.Context,
		listQuery ListQuery,
	) ([]T, error)
}

type Repository[T any] struct {
	db *gorm.DB
}

func NewRepository[T any](db *gorm.DB) *Repository[T] {
	return &Repository[T]{
		db: db,
	}
}

func (r *Repository[T]) Query(
	ctx context.Context,
) *gorm.DB {
	return DB(ctx, r.db).Model(new(T))
}

func (r *Repository[T]) Create(
	ctx context.Context,
	data *T,
) error {
	return DB(ctx, r.db).
		Create(data).
		Error
}

func (r *Repository[T]) Upsert(
	ctx context.Context,
	data *[]T,
	options UpsertOptions,
) error {
	return DB(ctx, r.db).
		Clauses(buildOnConflict(options)).
		Create(data).
		Error
}

func (r *Repository[T]) Get(
	ctx context.Context,
	conditions map[string]any,
) (*T, error) {
	var record T

	if err := DB(ctx, r.db).
		Where(conditions).
		First(&record).
		Error; err != nil {
		return nil, err
	}

	return &record, nil
}

func (r *Repository[T]) UpdateByID(
	ctx context.Context,
	ID uint64,
	updates map[string]any,
) error {
	return DB(ctx, r.db).
		Model(new(T)).
		Where("id = ?", ID).
		Updates(updates).
		Error
}

func (r *Repository[T]) List(
	ctx context.Context,
	listQuery ListQuery,
) ([]T, int64, error) {
	var (
		records []T
		total   int64
	)

	query := DB(ctx, r.db).
		Model(new(T)).
		Where(listQuery.Conditions).
		Scopes(listQuery.Filters...)

	if err := query.
		Count(&total).
		Error; err != nil {
		return nil, 0, err
	}

	if err := query.
		Scopes(listQuery.Scopes...).
		Select(listQuery.NormalizeColumns()).
		Find(&records).
		Error; err != nil {
		return nil, 0, err
	}

	return records, total, nil
}

func (r *Repository[T]) ListAll(
	ctx context.Context,
	listQuery ListQuery,
) ([]T, error) {
	var records []T

	if err := DB(ctx, r.db).
		Where(listQuery.Conditions).
		Scopes(listQuery.Filters...).
		Select(listQuery.NormalizeColumns()).
		Find(&records).
		Error; err != nil {
		return nil, err
	}

	return records, nil
}
