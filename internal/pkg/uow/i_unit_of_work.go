package uow

import "context"

//go:generate mockgen -destination=./../../mocks/uow_mock.go -package=mocks . IUnitOfWork
type IUnitOfWork interface {
	Begin(ctx context.Context) context.Context
	Commit(ctx context.Context) error
	Rollback(ctx context.Context) error
}
