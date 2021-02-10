package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/BurntSushi/toml"
	_ "github.com/denisenkom/go-mssqldb"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type metric struct {
	Name      string
	Help      string
	Statement string
}
type connection struct {
	Database        string
	ID              string
	Password        string
	Port            string
	RepeatInSeconds time.Duration
	Metric          []metric
}
type connections struct {
	Connection []connection
}
type queryMetric struct {
	Gauge     prometheus.Gauge
	Statement string
}

func main() {
	var monitor connections

	if _, err := toml.DecodeFile("metrics.toml", &monitor); err != nil {
		fmt.Println(err)
		return
	}
	for _, con := range monitor.Connection {
		gauges := makeGauge(con)
		recordMetrics(con, gauges)
	}
	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(":2112", nil)
}

func recordMetrics(con connection, gauges []queryMetric) {
	connect := ""
	if con.Port != "" {
		connect = ("server=" + con.Database +
			";user id=" + con.ID +
			";password=" + con.Password +
			";port=" + con.Port + ";")
	} else {
		connect = ("server=" + con.Database +
			";user id=" + con.ID +
			";password=" + con.Password + ";")
	}

	for _, queryMetric := range gauges {
		postNew(con, connect, queryMetric)
	}
}

func postNew(con connection, connect string, metric queryMetric) {
	go func() {

		for {

			condb, errdb := sql.Open("mssql", connect)
			if errdb != nil {
				fmt.Println(" Error open db:", errdb.Error())
			}
			var response string = query(metric.Statement, condb)
			fmt.Println("Ran metrics on " + con.Database)
			i, err := strconv.Atoi(response)
			if err != nil {
				i = -1
			}
			metric.Gauge.Set(float64(i))
			time.Sleep(con.RepeatInSeconds * time.Second)

		}
	}()
}

func query(statement string, condb *sql.DB) string {

	var (
		response string
	)

	rows, err := condb.Query(statement)
	if err != nil {
		log.Fatal(err)
	}
	for rows.Next() {
		err := rows.Scan(&response)
		if err != nil {
			log.Fatal(err)
		}
	}

	defer condb.Close()
	return response
}

func makeGauge(con connection) []queryMetric {
	var gaugeStatements []queryMetric
	for _, metric := range con.Metric {
		var gaugeStatement queryMetric
		gaugeStatement.Gauge = promauto.NewGauge(prometheus.GaugeOpts{
			Name: metric.Name, Help: metric.Help})
		gaugeStatement.Statement = metric.Statement
		gaugeStatements = append(gaugeStatements, gaugeStatement)
	}
	return gaugeStatements

}
