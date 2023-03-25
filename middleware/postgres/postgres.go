package postgres

import (
	"SuperArch/conf"
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
)


func GetPostgresClient() (*sql.DB){
	pdb, err := sql.Open("postgres", fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		conf.Cfg.PostgresDB.Username, conf.Cfg.PostgresDB.Password, conf.Cfg.PostgresDB.Host, conf.Cfg.PostgresDB.Port, conf.Cfg.PostgresDB.DBNAME))
	if err != nil{
		logrus.Errorf("[Postgres Module][InitPostgresClient] %s", err)
		return nil
	}
	pdb.Exec(fmt.Sprintf("set search_path=%s", conf.Cfg.PostgresDB.SCHEMA))

	return pdb
}
