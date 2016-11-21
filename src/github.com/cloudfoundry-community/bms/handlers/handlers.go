package handlers

import (
	"fmt"

	middleware "github.com/go-openapi/runtime/middleware"

	"github.com/cloudfoundry-community/bms/models"
	"github.com/cloudfoundry-community/bms/restapi/operations"
)

func GetInfoHandlerFunc(params operations.GetInfoParams) middleware.Responder {
	fmt.Printf("GetInfoHandlerFunc params: %#v\n", params)

	infoOK := operations.NewGetInfoOK()
	infoOK.SetPayload(&models.Info{
		Status: 200,
		Data: &models.InfoData{
			Name:    "Baremetal Server Golang version",
			Version: "0.2.0",
		},
	})

	return infoOK
}

func GetLoginUsernamePasswordHandlerFunc(params operations.GetLoginUsernamePasswordParams) middleware.Responder {
	fmt.Printf("GetLoginUsernamePasswordHandlerFunc params: %#v\n", params)

	if params.Username == "admin" && params.Password == "admin" {
		loginUsernamePasswordOK := operations.NewGetLoginUsernamePasswordOK()
		loginUsernamePasswordOK.SetPayload(&models.Login{
			Status: 200,
			Data:   "Login Successful",
		})
		return loginUsernamePasswordOK
	}

	loginUsernamePasswordDefault := operations.NewGetLoginUsernamePasswordDefault(403)
	loginUsernamePasswordDefault.SetPayload(&models.Error{
		Code:    403,
		Fields:  "fake-fields",
		Message: "Wrong username or password!",
	})

	return loginUsernamePasswordDefault
}

func GetStemcellsHandlerFunc(params operations.GetStemcellsParams) middleware.Responder {
	stemcellsOK := operations.NewGetStemcellsOK()
	stemcellsOK.SetPayload(&models.Stemcells{
		Status: 200,
		Data:   []string{"fake-stemcell-0", "fake-stemcell-1"},
	})

	return stemcellsOK
}
