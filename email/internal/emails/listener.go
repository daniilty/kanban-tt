package emails

import (
	"context"
	"time"

	"github.com/daniilty/kanban-tt/email/internal/kafka"
	"github.com/daniilty/kanban-tt/schema"
	"go.uber.org/zap"
)

type EventsListener interface {
	Listen(ctx context.Context)
}

type eventsListener struct {
	logger        *zap.SugaredLogger
	timeout       time.Duration
	kafkaConsumer kafka.Consumer
	emailSender   Sender
}

func NewEventsListener(logger *zap.SugaredLogger, timeout time.Duration, consumer kafka.Consumer, emailSender Sender) EventsListener {
	return &eventsListener{
		logger:        logger,
		timeout:       timeout,
		kafkaConsumer: consumer,
		emailSender:   emailSender,
	}
}

func (e *eventsListener) Listen(ctx context.Context) {
	e.logger.Info("Listening for email events.")
	for {
		select {
		case <-ctx.Done():
			e.logger.Infow("Stopping users event handler.")
			e.kafkaConsumer.Close()

			return
		default:
			e.handleMessage(ctx)
		}
	}
}

func (e *eventsListener) handleMessage(ctx context.Context) {
	email := &schema.Email{}

	commit, err := e.kafkaConsumer.UnmarshalMessage(ctx, email)
	if err != nil {
		e.logger.Errorw("Unmarshal kafka message.", "err", err)

		return
	}

	err = e.emailSender.Send(email.To, "Welcome link", email.Msg)
	if err != nil {
		e.logger.Errorw("Send email.", "err", err)

		return
	}

	err = commit(ctx)
	if err != nil {
		e.logger.Errorw("Commit kafka message.", "err", err)
	}
}
