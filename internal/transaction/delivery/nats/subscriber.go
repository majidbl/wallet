package nats

import (
	"context"
	"encoding/json"
	"github.com/avast/retry-go"
	"github.com/go-playground/validator/v10"
	"github.com/nats-io/stan.go"
	"github.com/opentracing/opentracing-go"
	"sync"
	"time"

	"github.com/pkg/errors"

	"github.com/majidbl/wallet/internal/models"
	"github.com/majidbl/wallet/internal/transaction"
	"github.com/majidbl/wallet/pkg/logger"
)

const (
	retryAttempts = 3
	retryDelay    = 1 * time.Second
)

type transactionSubscriber struct {
	stanConn      stan.Conn
	log           logger.Logger
	transactionUC transaction.UseCase
	validator     *validator.Validate
}

// NewTransactionSubscriber transaction subscriber constructor
func NewTransactionSubscriber(stanConn stan.Conn, log logger.Logger, transactionUC transaction.UseCase, validator *validator.Validate) *transactionSubscriber {
	return &transactionSubscriber{stanConn: stanConn, log: log, transactionUC: transactionUC, validator: validator}
}

// Subscribe  to subject and run workers with given callback for handling messages
func (s *transactionSubscriber) Subscribe(subject, qgroup string, workersNum int, cb stan.MsgHandler) {
	s.log.Infof("Subscribing to Subject: %v, group: %v", subject, qgroup)
	wg := &sync.WaitGroup{}

	for i := 0; i <= workersNum; i++ {
		wg.Add(1)
		go s.runWorker(
			wg,
			i,
			s.stanConn,
			subject,
			qgroup,
			cb,
			stan.SetManualAckMode(),
			stan.AckWait(ackWait),
			stan.DurableName(durableName),
			stan.MaxInflight(maxInflight),
			stan.DeliverAllAvailable(),
		)
	}
	wg.Wait()
}

func (s *transactionSubscriber) runWorker(
	wg *sync.WaitGroup,
	workerID int,
	conn stan.Conn,
	subject string,
	qgroup string,
	cb stan.MsgHandler,
	opts ...stan.SubscriptionOption,
) {
	s.log.Infof("Subscribing worker: %v, subject: %v, qgroup: %v", workerID, subject, qgroup)
	defer wg.Done()

	_, err := conn.QueueSubscribe(subject, qgroup, cb, opts...)
	if err != nil {
		s.log.Errorf("WorkerID: %v, QueueSubscribe: %v", workerID, err)
		if err := conn.Close(); err != nil {
			s.log.Errorf("WorkerID: %v, conn.Close error: %v", workerID, err)
		}
	}
}

// Run start subscribers
func (s *transactionSubscriber) Run(ctx context.Context) {
	go s.Subscribe(createTransactionSubject, transactionGroupName, createTransactionWorkers, s.processCreateTransaction(ctx))
}

func (s *transactionSubscriber) processCreateTransaction(ctx context.Context) stan.MsgHandler {
	return func(msg *stan.Msg) {
		span, ctx := opentracing.StartSpanFromContext(ctx, "transactionSubscriber.processCreateTransaction")
		defer span.Finish()

		s.log.Infof("subscriber process Create Transaction: %s", msg.String())
		totalSubscribeMessages.Inc()

		var m models.Transaction
		if err := json.Unmarshal(msg.Data, &m); err != nil {
			errorSubscribeMessages.Inc()
			s.log.Errorf("json.Unmarshal : %v", err)
			return
		}

		if err := retry.Do(func() error {
			return s.transactionUC.Create(ctx, &m)
		},
			retry.Attempts(retryAttempts),
			retry.Delay(retryDelay),
			retry.Context(ctx),
		); err != nil {
			errorSubscribeMessages.Inc()
			s.log.Errorf("transactionUC.Create : %v", err)

			if msg.Redelivered && msg.RedeliveryCount > maxRedeliveryCount {
				if err := s.publishErrorMessage(ctx, msg, err); err != nil {
					s.log.Errorf("publishErrorMessage : %v", err)
					return
				}
				if err := msg.Ack(); err != nil {
					s.log.Errorf("msg.Ack: %v", err)
					return
				}
			}
			return
		}

		if err := msg.Ack(); err != nil {
			s.log.Errorf("msg.Ack: %v", err)
		}
		successSubscribeMessages.Inc()
	}
}

func (s *transactionSubscriber) publishErrorMessage(ctx context.Context, msg *stan.Msg, err error) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "transactionSubscriber.publishErrorMessage")
	defer span.Finish()

	s.log.Infof("publish dead letter queue message: %v", msg)

	errMsg := &models.TransactionErrorMsg{
		Subject:   msg.Subject,
		Sequence:  msg.Sequence,
		Data:      msg.Data,
		Timestamp: msg.Timestamp,
		Error:     err.Error(),
		Time:      time.Now().UTC(),
	}

	errMsgBytes, err := json.Marshal(&errMsg)
	if err != nil {
		return errors.Wrap(err, "json.Marshal")
	}

	return s.stanConn.Publish(deadLetterQueueSubject, errMsgBytes)
}
