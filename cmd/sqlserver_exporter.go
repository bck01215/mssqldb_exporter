package main

import (
	"database/sql"
	"errors"

	_ "github.com/denisenkom/go-mssqldb"
	"github.com/sirupsen/logrus"
)

func sql_exporter(name string) Config {

	conf, err := get_config(name, get_conns().Configs)
	if err != nil {
		logrus.Error(err)

	}
	return conf
}

func get_config(name string, configs []Config) (Config, error) {

	for _, config := range configs {
		if config.Connection == name {
			return config, nil
		}
	}
	var empty_conf Config
	return empty_conf, errors.New("Could not find the configuration file with name: " + name)

}
func get_metric_info(server_con string) Config {
	conf := sql_exporter(server_con)
	db := connect(conf)
	for i, metric := range conf.Metrics {
		values, err := run_query(db, metric.Statement, metric.Value)
		if err != nil {
			logrus.Error("Error collecting metrics:\n\t", err)
			continue
		}
		conf.Metrics[i].Values = values
	}
	db.Close()
	return conf
}

func connect(con Config) *sql.DB {
	connect := ""
	if con.Port != "" {
		connect = ("server=" + con.Host +
			";user id=" + con.Username +
			";password=" + con.Password +
			";port=" + con.Port + ";")
	} else {
		connect = ("server=" + con.Host +
			";user id=" + con.Username +
			";password=" + con.Password + ";")
	}
	db, err := sql.Open("mssql", connect)
	if err != nil {
		logrus.Error(err)
	}
	return db
}

func run_query(db *sql.DB, statement string, val_column string) (map[string]interface{}, error) {
	value := make(map[string]interface{})

	rows, err := db.Query(statement)
	if err != nil {
		return nil, err
	}
	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}
	count := len(columns)
	values := make([]interface{}, count)
	valuePtrs := make([]interface{}, count)
	for rows.Next() {
		for i := range columns {
			valuePtrs[i] = &values[i]
		}

		rows.Scan(valuePtrs...)

		for i, col := range columns {
			val := values[i]

			b, ok := val.([]byte)
			var v interface{}
			if ok {
				v = string(b)
			} else {
				v = val
			}

			value[col] = v

		}
	}
	return value, nil
}
