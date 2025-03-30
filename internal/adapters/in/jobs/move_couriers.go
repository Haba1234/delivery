package jobs

import (
	"context"
	"log"

	"github.com/Haba1234/delivery/internal/core/application/usecases/commands"
	"github.com/Haba1234/delivery/internal/pkg/errs"

	"github.com/robfig/cron/v3"
)

var _ cron.Job = &MoveCouriers{}

type MoveCouriers struct {
	moveCouriersCommandHandler commands.IMoveCouriersHandler
}

func NewMoveCouriers(
	moveCouriersCommandHandler commands.IMoveCouriersHandler,
) (*MoveCouriers, error) {
	if moveCouriersCommandHandler == nil {
		return nil, errs.NewValueIsRequiredError("moveCouriersCommandHandler")
	}

	return &MoveCouriers{
		moveCouriersCommandHandler: moveCouriersCommandHandler,
	}, nil
}

func (j *MoveCouriers) Run() {
	ctx := context.Background()

	command, err := commands.NewMoveCouriers()
	if err != nil {
		log.Printf("failed to commands.NewMoveCouriers(): %v", err)
	}

	err = j.moveCouriersCommandHandler.Handle(ctx, command)
	if err != nil {
		log.Printf("failed to moveCouriersCommandHandler.Handle(): %v", err)
	}
}
