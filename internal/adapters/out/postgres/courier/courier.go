package courier

import (
	"context"

	"github.com/Haba1234/delivery/internal/adapters/out/postgres"
	"github.com/Haba1234/delivery/internal/core/domain/model/courier"
	"github.com/Haba1234/delivery/internal/core/ports"
	"github.com/Haba1234/delivery/internal/pkg/errs"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var _ ports.ICourierRepository = &Repository{}

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) (*Repository, error) {
	if db == nil {
		return nil, errs.NewValueIsRequiredError("db")
	}

	return &Repository{
		db: db,
	}, nil
}

func (r *Repository) Add(ctx context.Context, aggregate *courier.Courier) error {
	modelCourier := toModelDB(aggregate)

	tx := postgres.GetTxFromContext(ctx)
	if tx == nil {
		tx = r.db
	}

	err := tx.
		Session(&gorm.Session{FullSaveAssociations: true}).
		Create(&modelCourier).Error
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) Update(ctx context.Context, aggregate *courier.Courier) error {
	modelCourier := toModelDB(aggregate)

	tx := postgres.GetTxFromContext(ctx)
	if tx == nil {
		tx = r.db
	}

	err := tx.
		Session(&gorm.Session{FullSaveAssociations: true}).
		Save(&modelCourier).Error
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) Get(ctx context.Context, id uuid.UUID) (*courier.Courier, error) {
	modelCourier := ModelCourier{}

	tx := postgres.GetTxFromContext(ctx)
	if tx == nil {
		tx = r.db
	}

	result := tx.
		Preload(clause.Associations).
		Find(&modelCourier, id)
	if result.RowsAffected == 0 {
		return nil, postgres.ErrNotFound
	}

	return toDomain(modelCourier), nil
}

func (r *Repository) GetAllInFreeStatus(ctx context.Context) ([]*courier.Courier, error) {
	var modelCouriers []ModelCourier

	tx := postgres.GetTxFromContext(ctx)
	if tx == nil {
		tx = r.db
	}

	result := tx.
		Preload(clause.Associations).
		Where("status = ?", courier.StatusFree).
		Find(&modelCouriers)
	if result.Error != nil {
		return nil, result.Error
	}

	if result.RowsAffected == 0 {
		return nil, postgres.ErrNotFound
	}

	aggregates := make([]*courier.Courier, len(modelCouriers))
	for i, modelCourier := range modelCouriers {
		aggregates[i] = toDomain(modelCourier)
	}

	return aggregates, nil
}
