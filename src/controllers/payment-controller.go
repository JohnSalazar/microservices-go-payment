package controllers

import (
	"net/http"
	"payment/src/security"

	"github.com/gin-gonic/gin"
	trace "github.com/oceano-dev/microservices-go-common/trace/otel"
)

type PaymentController struct {
	managerSecurityKeys security.ManagerSecurityRSAKeys
}

func NewPaymentController(
	managerSecurityKeys security.ManagerSecurityRSAKeys,
) *PaymentController {
	return &PaymentController{
		managerSecurityKeys: managerSecurityKeys,
	}
}

func (payment *PaymentController) RSAPublicKey(c *gin.Context) {
	_, span := trace.NewSpan(c.Request.Context(), "PaymentController.RSAPublicKey")
	defer span.End()

	keys := payment.managerSecurityKeys.GetAllRSAPublicKeys()

	c.JSON(http.StatusOK, keys)
}
