package authorize

import (
	"os"
	"runtime"

	// mysql driver
	_ "github.com/go-sql-driver/mysql"
	"xorm.io/xorm"
)

var ormManager *OrmManager

// InitOrmManager init
func InitOrmManager() {
	dbSource, set := os.LookupEnv("CASBIN_DB_STRING")
	if !set {
		dbSource = "test:test@tcp(localhost:3306)/casbin_metadata_dev"
	}
	ormManager = NewOrmManager("mysql", dbSource)
}

// OrmManager represents the MySQL ormManager for policy storage.
type OrmManager struct {
	driverName     string
	dataSourceName string
	engine         *xorm.Engine
}

// finalizer is the destructor for OrmManager.
func finalizer(a *OrmManager) {
	err := a.engine.Close()
	if err != nil {
		panic(err)
	}
}

// NewOrmManager is the constructor for OrmManager.
func NewOrmManager(driverName string, dataSourceName string) *OrmManager {
	a := &OrmManager{}
	a.driverName = driverName
	a.dataSourceName = dataSourceName

	// Open the DB, create it if not existed.
	a.open()

	// Call the destructor when the object is released.
	runtime.SetFinalizer(a, finalizer)
	return a
}

func (a *OrmManager) open() {
	engine, err := xorm.NewEngine(a.driverName, a.dataSourceName)
	if err != nil {
		panic(err)
	}
	a.engine = engine
	a.createTables()
}

func (a *OrmManager) close() {
	a.engine.Close()
	a.engine = nil
}

func (a *OrmManager) createTables() {
	err := a.engine.Sync2(new(Model))
	if err != nil {
		panic(err)
	}

	err = a.engine.Sync2(new(Adapter))
	if err != nil {
		panic(err)
	}

	err = a.engine.Sync2(new(Enforcer))
	if err != nil {
		panic(err)
	}

	err = a.engine.Sync2(new(PolicyList))
	if err != nil {
		panic(err)
	}
}
