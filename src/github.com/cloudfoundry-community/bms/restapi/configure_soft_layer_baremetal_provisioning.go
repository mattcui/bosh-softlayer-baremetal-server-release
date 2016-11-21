package restapi

import (
	"crypto/tls"
	"net/http"

	"fmt"

	errors "github.com/go-openapi/errors"
	runtime "github.com/go-openapi/runtime"
	middleware "github.com/go-openapi/runtime/middleware"

	"github.com/cloudfoundry-community/bms/handlers"
	"github.com/cloudfoundry-community/bms/restapi/operations"

	"log"
)

// This file is safe to edit. Once it exists it will not be overwritten

//go:generate swagger generate server --target .. --name  --spec ../swagger.json

func configureFlags(api *operations.SoftLayerBaremetalProvisioningAPI) {
	// api.CommandLineOptionsGroups = []swag.CommandLineOptionsGroup{ ... }
}

func configureAPI(api *operations.SoftLayerBaremetalProvisioningAPI) http.Handler {
	fmt.Println("configuring the API")

	// configure the api here
	api.ServeError = errors.ServeError

	// Set your custom logger if needed. Default one is log.Printf
	// Expected interface func(string, ...interface{})
	//
	// Example:
	api.Logger = log.Printf

	api.JSONConsumer = runtime.JSONConsumer()

	api.JSONProducer = runtime.JSONProducer()

	api.GetBmsDeploymentNameHandler = operations.GetBmsDeploymentNameHandlerFunc(func(params operations.GetBmsDeploymentNameParams) middleware.Responder {
		return middleware.NotImplemented("operation .GetBmsDeploymentName has not yet been implemented")
	})

	api.GetInfoHandler = operations.GetInfoHandlerFunc(handlers.GetInfoHandlerFunc)
	api.GetLoginUsernamePasswordHandler = operations.GetLoginUsernamePasswordHandlerFunc(handlers.GetLoginUsernamePasswordHandlerFunc)
	api.GetStemcellsHandler = operations.GetStemcellsHandlerFunc(handlers.GetStemcellsHandlerFunc)

	api.GetTasksHandler = operations.GetTasksHandlerFunc(func(params operations.GetTasksParams) middleware.Responder {
		return middleware.NotImplemented("operation .GetTasks has not yet been implemented")
	})
	api.GetTasksTaskIDTxtLevelHandler = operations.GetTasksTaskIDTxtLevelHandlerFunc(func(params operations.GetTasksTaskIDTxtLevelParams) middleware.Responder {
		return middleware.NotImplemented("operation .GetTasksTaskIDTxtLevel has not yet been implemented")
	})
	api.PutBaremetalServerIDStateHandler = operations.PutBaremetalServerIDStateHandlerFunc(func(params operations.PutBaremetalServerIDStateParams) middleware.Responder {
		return middleware.NotImplemented("operation .PutBaremetalServerIDState has not yet been implemented")
	})
	api.PutBaremetalsHandler = operations.PutBaremetalsHandlerFunc(func(params operations.PutBaremetalsParams) middleware.Responder {
		return middleware.NotImplemented("operation .PutBaremetals has not yet been implemented")
	})

	api.ServerShutdown = func() {}

	return setupGlobalMiddleware(api.Serve(setupMiddlewares))
}

// The TLS configuration before HTTPS server starts.
func configureTLS(tlsConfig *tls.Config) {
	// Make all necessary changes to the TLS configuration here.
}

// The middleware configuration is for the handler executors. These do not apply to the swagger.json document.
// The middleware executes after routing but before authentication, binding and validation
func setupMiddlewares(handler http.Handler) http.Handler {
	return handler
}

// The middleware configuration happens before anything, this middleware also applies to serving the swagger.json document.
// So this is a good place to plug in a panic handling middleware, logging and metrics
func setupGlobalMiddleware(handler http.Handler) http.Handler {
	return handler
}
