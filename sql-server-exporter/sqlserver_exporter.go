package main

import (
	"database/sql"
	"errors"
	"sync"

	_ "github.com/denisenkom/go-mssqldb"
	"github.com/sirupsen/logrus"
)

var waitGroup sync.WaitGroup

func sql_exporter(name string) Config {

	conf, err := get_config(name, Get_Conns().Configs)
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
	db, err := connect(conf)
	if err == nil {
		waitGroup.Add(len(conf.Metrics))
		for i, metric := range conf.Metrics {
			go func(i int, metric Metric) {
				values, err := run_query(db, metric.Statement, metric.Value)
				if err != nil {
					logrus.Error("Error collecting metrics:\n\t", err)

				} else {
					conf.Metrics[i].Values = values
				}
				defer waitGroup.Done()
			}(i, metric)
		}
		waitGroup.Wait()
	} else {
		conf.Metrics = make([]Metric, 0)
	}
	db.Close()
	return conf
}

func connect(con Config) (*sql.DB, error) {
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
	db, _ := sql.Open("mssql", connect)
	sqlServerUp.Reset()
	Err := db.Ping()
	if Err != nil {
		logrus.Error(Err)
		sqlServerUp.WithLabelValues(Err.Error()).Set(float64(0))

	} else {
		sqlServerUp.WithLabelValues("").Set(float64(1))
	}

	return db, Err
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

		if err := rows.Scan(valuePtrs...); err != nil {
			logrus.Error(err)
		}

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
