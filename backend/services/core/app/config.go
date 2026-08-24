package main

import (
	"github.com/scalent.io/scalent-hrms/internal/middleware"
	"github.com/scalent.io/scalent-hrms/pkg/casbin"
	"github.com/scalent.io/scalent-hrms/pkg/db/cache"
	"github.com/scalent.io/scalent-hrms/pkg/db/sqlx"
	"github.com/scalent.io/scalent-hrms/pkg/reacher"
	"github.com/scalent.io/scalent-hrms/pkg/server"
	"github.com/scalent.io/scalent-hrms/services/core/service"
	"github.com/spf13/viper"
)

type CoreConfig struct {
	Server     *server.Config
	Middleware middleware.Middleware
	DB         *sqlx.DbConfig
	CasbinImpl casbin.CasbinImplConfig
	Cache      cache.Config
	Service    service.Config
	Reacher    reacher.ReacherConfig
}

var config CoreConfig

func initConfig() error {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.SetEnvPrefix("core")

	err := viper.ReadInConfig()
	if err != nil {
		return err
	}

	viper.Unmarshal(&config)
	return nil
}

func NewCacheConfig() *cache.Config {
	return &config.Cache
}

func NewServiceConfig() service.Config {
	return config.Service
}
