package main

import (
	"fmt"
	"net/http"

	"github.com/sirupsen/logrus"

	Config "mini-accounting/config"
	Constants "mini-accounting/constants"
	Library "mini-accounting/library"
	CustomErrorPackage "mini-accounting/pkg/custom_error"
	LoggerPackage "mini-accounting/pkg/logger"
	Wire "mini-accounting/wire"
)

func main() {
	path := "main"
	// INIT LIBRARY
	library := Library.New()
	// EMBED LOGGER PACKAGE INTO THIS PROJECT
	LoggerPackage.New(library)
	// INIT CONFIGURATION
	config := Config.New(library)
	// SETUP CONFIGURATION
	err := config.Setup()
	// WHEN SETUP CONFIGURATION RETURNS AN ERROR
	if err != nil {
		err = err.(*CustomErrorPackage.CustomError).UnshiftPath(path)
		LoggerPackage.WriteLog(logrus.Fields{
			"path":  err.(*CustomErrorPackage.CustomError).GetPath(),
			"title": err.(*CustomErrorPackage.CustomError).GetDisplay().Error(),
		}).Panic(err.(*CustomErrorPackage.CustomError).GetPlain())
	}
	// SET PORT THAT WE GET FROM .env FILE
	port := fmt.Sprintf(":%s", config.GetConfig().App.Port)
	// INIT ROUTER
	router := Wire.InjectRoute(config, library)
	// SETUP ROUTER
	router.Setup()
	// WHEN RUNNING SERVER RETURNS AN ERROR
	if err := http.ListenAndServe(port, router.GetEngine()); err != nil {
		err = CustomErrorPackage.New(Constants.ErrServeFailed, err, path, library)
		LoggerPackage.WriteLog(logrus.Fields{
			"path":  err.(*CustomErrorPackage.CustomError).GetPath(),
			"title": err.(*CustomErrorPackage.CustomError).GetDisplay().Error(),
		}).Fatal(err.(*CustomErrorPackage.CustomError).GetPlain())
	}
}
