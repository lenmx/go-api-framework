package db

import (
	"project-name/pkg/xxorm"

	//"github.com/arthurkiller/rollingwriter"
	_ "github.com/go-sql-driver/mysql"
	"runtime"
	"project-name/config"
	"project-name/pkg/xlogger"
	"xorm.io/xorm"
	"xorm.io/xorm/names"
)

type Adapter struct {
	driverName     string
	dataSourceName string
	dbName         string
	*xorm.Engine
}

var Pop *Adapter

func InitAdapter() {
	Pop = NewAdapter("mysql", config.G_config.Db.DataSourceName, config.G_config.Db.Name)
}

func NewAdapter(driverName string, dataSourceName string, dbname string) (adapter *Adapter) {
	adapter = &Adapter{
		driverName:     driverName,
		dataSourceName: dataSourceName,
		dbName:         dbname,
	}

	err := adapter.open()
	if err != nil {
		xlogger.ErrorLogger.Errorf("connect to db err: %s", err)
		panic(err)
	}
	runtime.SetFinalizer(adapter, finalizer)

	return adapter
}

func (adapter *Adapter) open() (err error) {
	var engine *xorm.Engine
	engine, err = xorm.NewEngine(adapter.driverName, adapter.dataSourceName+adapter.dbName);
	if err != nil {
		return
	}

	engine.AddHook(xxorm.NewLogHook())
	engine.SetMapper(names.SameMapper{})

	adapter.Engine = engine
	return
}

func (adapter *Adapter) close() (err error) {
	if err = adapter.Engine.Close(); err != nil {
		return
	}

	adapter.Engine = nil
	return
}

func finalizer(adapter *Adapter) {
	if err := adapter.close(); err != nil {
		panic(err)
	}
}
