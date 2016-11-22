package vm

import (
	"github.com/cloudfoundry-community/vps/models"
)

func (o *OrderVMByFilterOK) GetPayload() *models.VMResponse {
	return o.Payload
}

func (o *OrderVMByFilterNotFound) GetStatusCode() int {
	return 404
}

func (o *OrderVMByFilterDefault) GetStatusCode() int {
	return o._statusCode
}

func (o *OrderVMByFilterDefault) GetPayload() *models.Error {
	return o.Payload
}
