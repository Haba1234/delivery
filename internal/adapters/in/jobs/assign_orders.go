package jobs

import (
	"context"
	"log"

	"github.com/Haba1234/delivery/internal/core/application/usecases/commands"
	"github.com/Haba1234/delivery/internal/pkg/errs"

	"github.com/robfig/cron/v3"
)

var _ cron.Job = &AssignOrders{}

type AssignOrders struct {
	assignOrdersCommandHandler commands.IAssignOrderHandler
}

func NewAssignOrders(
	assignOrdersCommandHandler commands.IAssignOrderHandler,
) (*AssignOrders, error) {
	if assignOrdersCommandHandler == nil {
		return nil, errs.NewValueIsRequiredError("assignOrdersCommandHandler")
	}

	return &AssignOrders{
		assignOrdersCommandHandler: assignOrdersCommandHandler,
	}, nil
}

func (j *AssignOrders) Run() {
	ctx := context.Background()

	command, err := commands.NewAssignOrder()
	if err != nil {
		log.Printf("failed to commands.NewAssignOrder(): %v", err)
	}

	err = j.assignOrdersCommandHandler.Handle(ctx, command)
	if err != nil {
		log.Printf("failed to assignOrdersCommandHandler.Handle(): %v", err)
	}
}
