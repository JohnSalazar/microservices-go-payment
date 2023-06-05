package tasks

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math/rand"
	command "payment/src/application/commands"
	"payment/src/dtos"
	"payment/src/repositories/interfaces"
	"payment/src/security"
	"time"

	common_models "github.com/oceano-dev/microservices-go-common/models"
	common_nats "github.com/oceano-dev/microservices-go-common/nats"
	common_security "github.com/oceano-dev/microservices-go-common/security"
	common_service "github.com/oceano-dev/microservices-go-common/services"
	trace "github.com/oceano-dev/microservices-go-common/trace/otel"
)

type verifyPaymentTask struct {
	manager           security.ManagerSecurityRSAKeys
	common_manager    common_security.ManagerSecurityRSAKeys
	paymentRepository interfaces.PaymentRepository
	email             common_service.EmailService
	publisher         common_nats.Publisher
}

func NewVerifyPaymentTask(
	manager security.ManagerSecurityRSAKeys,
	common_manager common_security.ManagerSecurityRSAKeys,
	paymentRepository interfaces.PaymentRepository,
	email common_service.EmailService,
	publisher common_nats.Publisher,
) *verifyPaymentTask {
	return &verifyPaymentTask{
		manager:           manager,
		common_manager:    common_manager,
		paymentRepository: paymentRepository,
		email:             email,
		publisher:         publisher,
	}
}

var payments []*dtos.ProcessPayment

func (task *verifyPaymentTask) AddPayment(payment *dtos.ProcessPayment) {
	payments = append(payments, payment)
}

func (task *verifyPaymentTask) Run() {
	ticker := time.NewTicker(2 * time.Second)
	quit := make(chan struct{})
	go func() {
		for {
			select {
			case <-ticker.C:
				for i := range payments {
					if payments[i].VerifiedAt.IsZero() {
						payments[i] = nil

					}
					if payments[i] != nil && payments[i].VerifiedAt.Before(time.Now().UTC()) {

						success, err := task.paymentGetway(payments[i])
						if err != nil {
							payments[i].VerifiedAt = time.Now().UTC().Add(30 * time.Second)

							ctx := context.Background()
							_, span := trace.NewSpan(ctx, "tasks.VerifyPaymentTask")
							defer span.End()
							msg := fmt.Sprintf("error update payment %s: %s", payments[i].ID, err.Error())
							trace.FailSpan(span, msg)
							log.Print(msg)
							go task.email.SendSupportMessage(msg)
							ticker.Reset(15 * time.Second)
							break
						}

						paymentCommand := &command.UpdateStatusPaymentCommand{
							ID:       payments[i].ID,
							Status:   uint(common_models.PaymentConfirmed),
							StatusAt: time.Now().UTC(),
						}

						if !success {
							paymentCommand.Status = uint(common_models.PaymentRejected)
							fmt.Println("unauthorized payment!!!")
						} else {
							fmt.Println("authorized payment!")
						}

						data, _ := json.Marshal(paymentCommand)
						_ = task.publisher.Publish(string(common_nats.PaymentUpdate), data)

						payments[i] = nil
					}
				}

				if len(payments) > 0 {
					task.clearPayment()
				}

				//fmt.Printf("payment success checked %s\n", time.Now().UTC())
				ticker.Reset(5 * time.Second)
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()
}

func (task *verifyPaymentTask) clearPayment() {
	newPayment := []*dtos.ProcessPayment{}

	for i := range payments {
		if payments[i] != nil && !payments[i].VerifiedAt.IsZero() {
			newPayment = append(newPayment, payments[i])
		}
	}

	payments = newPayment
}

func (task *verifyPaymentTask) paymentGetway(payment *dtos.ProcessPayment) (bool, error) {
	cardNumber, err := task.decryptCardNumber(payment)
	if err != nil {
		return false, err
	}

	fmt.Println(cardNumber)

	time.Sleep(8 * time.Second)

	errOcurred := rand.Float32() < 0.5
	if errOcurred {
		return false, errors.New("communication error")
	}

	success := rand.Float32() < 0.5

	return success, nil
}

func (task *verifyPaymentTask) decryptCardNumber(payment *dtos.ProcessPayment) (string, error) {
	privateKey := task.manager.GetKeyFromKid(payment.Kid)
	if privateKey == nil {
		return "", errors.New("key to decrypt information not found")
	}

	// data, err := json.Marshal(payment.CardNumber)
	// if err != nil {
	// 	return "", err
	// }

	cardNumber, err := task.common_manager.Decrypt(payment.CardNumber, privateKey)
	if err != nil {
		return "", err
	}

	return cardNumber, nil
}
