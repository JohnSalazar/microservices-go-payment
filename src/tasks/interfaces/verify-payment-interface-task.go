package interfaces

import "payment/src/dtos"

type VerifyPaymentTask interface {
	AddPayment(payment *dtos.ProcessPayment)
	Run()
}
