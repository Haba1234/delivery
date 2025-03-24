package postgres

import (
	"context"
	"errors"

	"github.com/Haba1234/delivery/internal/pkg/errs"
	"github.com/Haba1234/delivery/internal/pkg/uow"

	"gorm.io/gorm"
)

var _ uow.IUnitOfWork = &UnitOfWork{}

var ErrNotFound = errors.New("not found")

type txKey struct{}

type UnitOfWork struct {
	db *gorm.DB
}

func NewUnitOfWork(db *gorm.DB) (*UnitOfWork, error) {
	if db == nil {
		return nil, errs.NewValueIsRequiredError("db")
	}
	return &UnitOfWork{db: db}, nil
}

func (u *UnitOfWork) Begin(ctx context.Context) context.Context {
	tx := u.db.Begin()
	return context.WithValue(ctx, txKey{}, tx)
}

func (*UnitOfWork) Commit(ctx context.Context) error {
	tx := GetTxFromContext(ctx)
	if tx != nil {
		return tx.Commit().Error
	}
	return nil
}

func (*UnitOfWork) Rollback(ctx context.Context) error {
	tx := GetTxFromContext(ctx)
	if tx != nil {
		return tx.Rollback().Error
	}
	return nil
}

func GetTxFromContext(ctx context.Context) *gorm.DB {
	tx, ok := ctx.Value(txKey{}).(*gorm.DB)
	if ok {
		return tx
	}
	return nil
}
