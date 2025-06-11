package config

import (
	"fmt"
	"os"
	"time"

	Environments "mini-accounting/config/environments"
	Constants "mini-accounting/constants"
	Library "mini-accounting/library"
	CustomErrorPackage "mini-accounting/pkg/custom_error"
)

type Config interface {
	Setup() error
	GetConfig() *Environments.ConfigModel
}

type ConfigImpl struct {
	model   Environments.ConfigModel
	library Library.Library
}

func New(
	library Library.Library,
) Config {
	return &ConfigImpl{
		library: library,
	}
}

func (o *ConfigImpl) Setup() error {
	path := "Config:Setup"
	// GET PROJECT DIRECTORY
	directory, err := os.Getwd()
	// WHEN GET DIRECTORY RETURNS ERROR
	if err != nil {
		return CustomErrorPackage.New(Constants.ErrConfiguration, err, path, o.library)
	}
	// INIT TEMPORARY CONFIGURATION MODEL
	var model Environments.ConfigModel
	// GET ".env" FILE (FULL PATH DIRECTORY)
	envFile := fmt.Sprintf("%s/.env", directory)
	// BINDING ".env" FILE INTO "model" POINTER
	err = o.library.ReadConfig(envFile, &model)
	// WHEN BINDING PROCESS RETURNS ERROR
	if err != nil {
		return CustomErrorPackage.New(Constants.ErrNoConfiguration, err, path, o.library)
	}

	// SET UP GLOBAL TIMEZONE
	location, err := time.LoadLocation(model.App.LocationTimezone)
	if err != nil {
		err := CustomErrorPackage.New(Constants.ErrNoConfiguration, Constants.ErrNoConfiguration, path, o.library)
		return err
	}

	time.Local = location

	// SET BINDED MODEL INTO THIS CONFIG PROPERTY
	o.model = model
	return nil
}

func (o *ConfigImpl) GetConfig() *Environments.ConfigModel {
	return &o.model
}
