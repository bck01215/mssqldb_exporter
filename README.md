# mssqldb_exporter
WIP: a Microsoft SQL database exporter written in go

The example metrics.toml shows the needed information to connect to the database.
Add as many connections as needed, and as many metrics as needed to each connection

Currently only gauge metrics are used, more may be provided in the future.

MSSQL driver is not included in this repo

No default MSSQL metrics are included in this repo.
