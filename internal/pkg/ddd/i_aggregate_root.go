package ddd

type IAggregateRoot interface {
	GetDomainEvents() []IDomainEvent
	ClearDomainEvents()
}
