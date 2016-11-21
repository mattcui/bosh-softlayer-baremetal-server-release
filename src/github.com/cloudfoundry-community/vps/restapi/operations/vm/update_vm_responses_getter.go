package vm

import (
	"github.com/cloudfoundry-community/vps/models"
)

func (o *UpdateVMOK) GetPayload() string {
	return o.Payload
}

func (o *UpdateVMNotFound) GetStatusCode() int {
	return 404
}

func (o *UpdateVMDefault) GetStatusCode() int {
	return o._statusCode
}

func (o *UpdateVMDefault) GetPayload() *models.Error {
	return o.Payload
}
