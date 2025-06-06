package sql;

import(
	Log     "log"
	Fmt     "fmt"
	SQL     "database/sql"
	Strings "strings"
	UtilsFS "github.com/PxnPub/pxnGoUtils/utils/fs"
//	_ "github.com/marcboeker/go-duckdb/v2"
//	_ "github.com/tursodatabase/go-libsql"
);



type Builder struct {
	Driver DriverType
	PathDB string
	FileDB string
	RW     bool
}

type PxDb struct {
	Driver DriverType
	DB     *SQL.DB
	Tables []PxTb
}

type PxTb struct {
	PXDB      *PxDb
	TableName string
	Fields    []PxFd
}

type PxFd struct {
	FieldName string
	FieldType string
}



func NewBuilder(driver DriverType) *Builder {
	return &Builder {
		Driver: driver,
		RW:     true,
	};
}

func (build *Builder) Build() *PxDb {
	driver := build.Driver; if driver == "" { driver = Driver_LibSQL; }
	pathdb := build.PathDB; if pathdb == "" { pathdb = "db/"; }
	filedb := build.FileDB; if filedb == "" { filedb = "db";  }
	mkdir, err := UtilsFS.CreateDIR(pathdb);
	if err != nil { panic(err); }
	if mkdir { Log.Printf("Created database dir: %s", pathdb); }
	var dsn string;
	switch driver {
	case Driver_LibSQL:
		var mode_rw string;
		if build.RW { mode_rw = "rwc";
		} else {      mode_rw = "ro"; }
		dsn = Fmt.Sprintf(string(DSN_LibSQL), pathdb+filedb, mode_rw);
	case Driver_DuckDB: dsn = Fmt.Sprintf(string(DSN_DuckDB), pathdb+filedb);
	default: panic(Fmt.Errorf("Invalid database driver: %s", driver));
	}
	// open database
	db, err := SQL.Open(string(driver), dsn);
	if err != nil { panic(err); }
	if build.RW { db.SetMaxOpenConns( 1); db.SetMaxIdleConns(1);
	} else {      db.SetMaxOpenConns(10); db.SetMaxIdleConns(5); }
	if err := db.Ping(); err != nil { panic(err); }
	return &PxDb{
		Driver: driver,
		DB:     db,
	};
}



func  (build *Builder) WithPath(pathdb string) *Builder { build.PathDB = pathdb; return build; }
func  (build *Builder) WithFile(filedb string) *Builder { build.FileDB = filedb; return build; }

func (build *Builder) ReadWrite() *Builder { build.RW = true;  return build; }
func (build *Builder) ReadOnly()  *Builder { build.RW = false; return build; }



func (pxdb *PxDb) Table(table_name string) *PxTb {
	table := PxTb{
		PXDB:      pxdb,
		TableName: table_name,
	};
	pxdb.Tables = append(pxdb.Tables, table);
	return &table;
}



func (pxtb *PxTb) PrimaryKey(field_name string) *PxTb {
	return pxtb.Field(field_name,
		Fmt.Sprintf(`INTEGER NOT NULL, PRIMARY KEY ("%s")`, field_name));
}

func (pxtb *PxTb) AutoIncrement(field_name string) *PxTb {
	return pxtb.Field(field_name, `INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT`);
}

func (pxtb *PxTb) Field(field_name string, field_type string) *PxTb {
	switch field_type {
	case "TIMESTAMP":
		switch pxtb.PXDB.Driver {
		case Driver_LibSQL: field_type = `INTEGER NOT NULL`;
		case Driver_DuckDB: field_type = `TIMESTAMP`;
		default: panic(Fmt.Errorf(
			"Invalid driver: %s; cannot detect timestamp type: %s %s",
			pxtb.PXDB.Driver, pxtb.TableName, field_name));
		}
	default:
	}
	var sql string;
	// first field
	if len(pxtb.Fields) == 0 {
		sql = Fmt.Sprintf(
			`CREATE TABLE IF NOT EXISTS "%s" ("%s" %s);`,
			pxtb.TableName,
			field_name,
			field_type,
		);
	// append table
	} else {
		sql = Fmt.Sprintf(
			`ALTER TABLE "%s" ADD COLUMN "%s" %s;`,
			pxtb.TableName,
			field_name,
			field_type,
		);
	}
	if _, err := pxtb.PXDB.DB.Exec(sql); err != nil && !IsDuplicateColumnError(err) {
		panic(Fmt.Errorf("SQL: %s\nERR: %s\n", sql, err)); }
	pxtb.Fields = append(pxtb.Fields, PxFd{ field_name, field_type });
	return pxtb;
}

func (pxtb *PxTb) Unique(fields string) *PxTb {
	unique_key := Strings.ReplaceAll(Strings.ReplaceAll(fields, ",", ""), " ", "_");
	sql := Fmt.Sprintf(
		`CREATE UNIQUE INDEX IF NOT EXISTS idx_unique_%s ON "%s" (%s);`,
		unique_key,
		pxtb.TableName,
		fields,
	);
	if _, err := pxtb.PXDB.DB.Exec(sql); err != nil {
		panic(Fmt.Errorf("SQL: %s\nERR: %s\n", sql, err)); }
	return pxtb;
}
