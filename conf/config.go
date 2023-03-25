package conf

import (
	"os"
)

var Cfg Config

type Config struct {
	Basic struct {

	}
	RabbitMQ struct {
		Host string
		Port int64
		Username string
		Password string
		Exchange string
	}
	Redis struct {
		Host string
		Port int64
		Password string
		DB int64
	}
	PostgresDB struct {
		Host string
		Port int64
		Username string
		Password string
		DBNAME string
		SCHEMA string
	}
}

func InitConfig()  {
	if len(os.Getenv("DEBUG")) > 0{
		// RabbitMQ
		Cfg.RabbitMQ.Host = "localhost"
		Cfg.RabbitMQ.Port = 5672
		Cfg.RabbitMQ.Username = "admin"
		Cfg.RabbitMQ.Password = "admin"
		Cfg.RabbitMQ.Exchange = "super_ex"

		// Redis
		Cfg.Redis.Host = "localhost"
		Cfg.Redis.Port = 6379
		Cfg.Redis.Password = "testadmin123"
		Cfg.Redis.DB = 0

		// Postgres
		Cfg.PostgresDB.Host = "localhost"
		Cfg.PostgresDB.Port = 5432
		Cfg.PostgresDB.Username = "admin"
		Cfg.PostgresDB.Password = "admin"
		Cfg.PostgresDB.DBNAME = "admin"
		Cfg.PostgresDB.SCHEMA = "superarch"
	}else{
		// rabbitMQ
		Cfg.RabbitMQ.Host = "localhost"
		Cfg.RabbitMQ.Port = 5672
		Cfg.RabbitMQ.Username = "admin"
		Cfg.RabbitMQ.Password = "testadmin123ashore"
		Cfg.RabbitMQ.Exchange = "super_ex"

		// Redis
		Cfg.Redis.Host = "localhost"
		Cfg.Redis.Port = 6379
		Cfg.Redis.Password = "testadmin123ashore"
		Cfg.Redis.DB = 0

		// Postgres
		Cfg.PostgresDB.Host = "localhost"
		Cfg.PostgresDB.Port = 5432
		Cfg.PostgresDB.Username = "admin"
		Cfg.PostgresDB.Password = "testadmin123ashore"
		Cfg.PostgresDB.DBNAME = "admin"
		Cfg.PostgresDB.SCHEMA = "superarch"

	}

}