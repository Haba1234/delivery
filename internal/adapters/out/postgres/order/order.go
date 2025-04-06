package order

import (
	"context"
	"errors"

	"github.com/Haba1234/delivery/internal/adapters/out/postgres"
	"github.com/Haba1234/delivery/internal/core/domain/model/order"
	"github.com/Haba1234/delivery/internal/core/ports"
	"github.com/Haba1234/delivery/internal/pkg/ddd"
	"github.com/Haba1234/delivery/internal/pkg/errs"

	"github.com/google/uuid"
	"github.com/mehdihadeli/go-mediatr"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var _ ports.IOrderRepository = &Repository{}

var (
	ErrOrderNotFound          = errors.New("order not found")
	ErrOrderCreatedNotFound   = errors.New("order with the status was Created not found")
	ErrAssignedOrdersNotFound = errors.New("assigned orders not found")
)

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

func (r *Repository) Add(ctx context.Context, aggregate *order.Order) error {
	modelOrder := toModelDB(aggregate)

	tx := postgres.GetTxFromContext(ctx)
	if tx == nil {
		tx = r.db
	}

	err := tx.
		Session(&gorm.Session{FullSaveAssociations: true}).
		Create(&modelOrder).Error
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) Update(ctx context.Context, aggregate *order.Order) error {
	modelOrder := toModelDB(aggregate)

	tx := postgres.GetTxFromContext(ctx)
	if tx == nil {
		tx = r.db
	}

	err := tx.
		Session(&gorm.Session{FullSaveAssociations: true}).
		Save(&modelOrder).Error
	if err != nil {
		return err
	}

	err = r.PublishDomainEvents(ctx, aggregate)
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) Get(ctx context.Context, id uuid.UUID) (*order.Order, error) {
	modelOrder := ModelOrder{}

	tx := postgres.GetTxFromContext(ctx)
	if tx == nil {
		tx = r.db
	}

	result := tx.
		Preload(clause.Associations).
		Find(&modelOrder, id)
	if result.RowsAffected == 0 {
		return nil, ErrOrderNotFound
	}

	aggregate := toDomain(modelOrder)
	return aggregate, nil
}

func (r *Repository) GetFirstInCreatedStatus(ctx context.Context) (*order.Order, error) {
	modelOrder := ModelOrder{}

	tx := postgres.GetTxFromContext(ctx)
	if tx == nil {
		tx = r.db
	}

	result := tx.
		Preload(clause.Associations).
		Where("status = ?", order.StatusCreated).
		First(&modelOrder)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, ErrOrderCreatedNotFound
		}
		return nil, result.Error
	}

	aggregate := toDomain(modelOrder)
	return aggregate, nil
}

func (r *Repository) GetAllInAssignedStatus(ctx context.Context) ([]*order.Order, error) {
	var modelOrders []ModelOrder

	tx := postgres.GetTxFromContext(ctx)
	if tx == nil {
		tx = r.db
	}

	result := tx.
		Preload(clause.Associations).
		Where("status = ?", order.StatusAssigned).
		Find(&modelOrders)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, ErrAssignedOrdersNotFound
	}

	aggregates := make([]*order.Order, len(modelOrders))
	for i, modelOrder := range modelOrders {
		aggregates[i] = toDomain(modelOrder)
	}

	return aggregates, nil
}

func (*Repository) PublishDomainEvents(ctx context.Context, aggregate ddd.IAggregateRoot) error {
	for _, event := range aggregate.GetDomainEvents() {
		switch event.(type) {
		case order.CompletedDomainEvent:
			err := mediatr.Publish[order.CompletedDomainEvent](ctx, event.(order.CompletedDomainEvent))
			if err != nil {
				return err
			}
		}
	}

	aggregate.ClearDomainEvents()

	return nil
}
