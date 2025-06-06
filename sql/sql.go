package sql;

import(
	Strings "strings"
);



type DriverDSN string;
const(
	DSN_LibSQL DriverDSN = "file:%s.sqlite?_pragma=journal_mode=MEMORY&_fk=true&_pragma=synchronous=NORMAL&mode=%s"
	DSN_DuckDB DriverDSN = "%s.duckdb?threads=2&memory_limit=1GB"
);



type DriverType string;
const(
	Driver_LibSQL DriverType = "libsql"
	Driver_DuckDB DriverType = "duckdb"
);



func IsDuplicateColumnError(err error) bool {
	if err == nil { return false; }
	if Strings.HasPrefix(err.Error(), "failed to execute query ALTER TABLE ") &&
	Strings.Contains(err.Error(), "Error executing statement: SQLite failure: `duplicate column name: ") {
		return true;
	}
	if Strings.HasPrefix(err.Error(), "Catalog Error: Column with name ") &&
	Strings.HasSuffix(err.Error(), " already exists!") {
		return true;
	}
	return false;
}
