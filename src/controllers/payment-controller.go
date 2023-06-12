package controllers

import (
	"net/http"
	"payment/src/security"

	trace "github.com/JohnSalazar/microservices-go-common/trace/otel"
	"github.com/gin-gonic/gin"
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
