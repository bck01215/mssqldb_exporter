<!-- 

[logo]: https://github.com/adam-p/markdown-here/raw/master/src/common/images/icon48.png "Logo Title Text 2" -->
 <!-- ![alt text][logo] -->
# SQLServer_exporter
## Notes:
WIP: a Microsoft SQL database exporter written in go.

Currently only gauge metrics are used.

## How it works


### Config file

You can follow the patter in the [config.yaml.](https://github.com/bck01215/mssqldb_exporter/blob/main/config.yaml). Its location can be specified with the ```--config``` flag or it will default to ```$PWD/config.yaml```. The connection name can be anything you would like to be passed in as the target paramater in the web request (i.e. ```localhost:9101/probe?target=connection_name```). The metric files will be searched for in the ```metrics-folder```. It defaults to ```$PWD/metrics```. Port is not required if your connection use the default port.


### Metric files

If you will need to specify which database you need to connect to as part of the statement. The labels will be collected from the column of the same name. You will also want to specify the name of the column used to collect the value. ```value: value_column```. If the metric is incorrectly made, the logs should provide more info. All other metrics should continue to work as expected, but I do not claim to have made a safety catch for every crazy combination, and it may cause an http panic.


### How to use

```
go build sql-server-exporter

sql-server-exporter --config=/path/to/config.yaml --metrics-folder=/path/to/metrics

```
Or you can use the docker container by mounting the volume
```
docker run -p 9105:9101 -v $PWD:/app bkauffman7/sql-server-exporter:v1 --config=/app/config.yaml --metrics-folder=/app/metric
````