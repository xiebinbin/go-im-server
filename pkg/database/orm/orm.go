package orm

import (
	"context"
	"fmt"
	"github.com/fatih/structs"
	"gorm.io/gorm"
	"imsdk/pkg/log"

	//_ "github.com/jinzhu/gorm/dialects/mysql" // 使用MySQL
	"gorm.io/driver/postgres" // 使用postgres
	"imsdk/pkg/app"
	"sync/atomic"
)

type (
	// Orm gorm 连接对象, 包含Master和Slaves, 由配置决定, Slaves 使用 atomic 包进行循环获取
	Orm struct {
		Master    *gorm.DB
		OldMaster *gorm.DB
		Slaves    []*gorm.DB
	}

	connInfo struct {
		Host     string `toml:"host"`
		Port     string `toml:"port"`
		Username string `toml:"username"`
		Password string `toml:"password"`
		DbName   string `toml:"dbname"`
		MaxIdle  int    `toml:"max_idle"`
		MaxOpen  int    `toml:"max_open"`
	}

	config struct {
		Master connInfo   `toml:"master"`
		Slaves []connInfo `toml:"slave"`
	}
)

var (
	orm       = &Orm{}
	slavesLen int
	err       error
	cursor    int64
	conf      config
)

func createMysqlConnectionURL(username, password, addr, dbName string) string {
	url := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local", username, password, addr, dbName)
	log.Logger().Debug(context.Background(), "mysql connect url: ", url)
	return url
}

func createPostgreConnectionUrl(conf map[string]interface{}) string {
	url := fmt.Sprintf("host=%s port=%s user=%s password=%s  dbname=%s sslmode=disable", conf["Host"], conf["Port"], conf["Username"], conf["Password"], conf["DbName"])
	log.Logger().Debug(context.Background(), "pgsql connect url: ", url)
	return url
}

// Start 启动数据库
func Start() {
	err = app.Config().Bind("db", "database", &conf)
	if err == app.ErrNodeNotExists {
		return
	}

	//var confInfo map[string]string
	confInfo := structs.New(conf.Master).Map()
	masterUrl := createPostgreConnectionUrl(confInfo)
	//orm.Master, err = gorm.Open("postgres", masterUrl)
	logCtx := log.WithFields(context.Background(), map[string]string{"action": "startPgsql"})
	orm.Master, err = gorm.Open(postgres.Open(masterUrl), &gorm.Config{})
	if err != nil {
		log.Logger().Warn(logCtx, "master database connect error: ", err, "url: ", masterUrl)
	}
	//orm.Master.LogMode(true)
	//orm.Master.DB().SetMaxIdleConns(configs.Master.MaxIdle)
	//orm.Master.DB().SetMaxOpenConns(configs.Master.MaxOpen)

	oldConfInfo := confInfo
	oldConfInfo["DbName"] = "im"
	oldUrl := createPostgreConnectionUrl(oldConfInfo)
	//connect, err := gorm.Open("postgres", slaveUrl)
	orm.OldMaster, err = gorm.Open(postgres.Open(oldUrl), &gorm.Config{})
	if err != nil {
		log.Logger().Warn(logCtx, "oldmaster database connect error: ", err, "url: ", oldUrl)
	}

	if conf.Slaves != nil {
		for _, slave := range conf.Slaves {
			slaveConfInfo := structs.New(slave).Map()
			slaveUrl := createPostgreConnectionUrl(slaveConfInfo)
			//connect, err := gorm.Open("postgres", slaveUrl)
			connect, err := gorm.Open(postgres.Open(slaveUrl), &gorm.Config{})
			if err != nil {
				log.Logger().Warn(logCtx, "slave database connect error: ", err, "url: ", slaveUrl)
			}
			orm.Slaves = append(orm.Slaves, connect)
		}
		slavesLen = len(orm.Slaves)
	}
}

// Slave 获得一个从库连接对象, 使用 atomic.AddInt64 计算调用次数，然后按 Slave 连接个数和次数进行取模操作之后获取指定index的Slave
func Slave() *gorm.DB {
	rs := atomic.AddInt64(&cursor, 1)
	return orm.Slaves[rs%int64(slavesLen)]
}

// Master 获得主库连接
func Master() *gorm.DB {
	return orm.Master
}

func OldMaster() *gorm.DB {
	return orm.OldMaster
}
