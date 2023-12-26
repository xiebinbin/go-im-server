package server

import (
	"context"
	"imsdk/pkg/app"
	"imsdk/pkg/database/mongo"
	"imsdk/pkg/funcs"
	"imsdk/pkg/log"
	"imsdk/pkg/redis"
	"imsdk/pkg/validator"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

type domainConfig struct {
	Name     string `toml:"name"`
	Addr     string `toml:"addr"`
	CertFile string `toml:"cert_file"`
	KeyFile  string `toml:"key_file"`
}

var (
	Mode   string
	engine = gin.New()

	// domainConf
	domainConf domainConfig
)

// 启动各项服务
func start() {
	app.LoadConfig()
	log.Start()
	//orm.Start()
	mongo.Start()
	redis.Start()
	//elastic.Start()
	//kafka.Start()
	// 加载应用配置
	domainKey := "domain_" + Mode
	app.Config().Bind("domains", domainKey, &domainConf)
	staticUrl, err := app.Config().GetChildConf("global", "system", "static_url")
	if err != nil {
		log.Logger().Fatal(context.Background(), "failed to get system static_url")
	}
	funcs.SetStaticUrl(staticUrl.(string))
	binding.Validator = new(validator.Validator)
}

func Run(service func(engine *gin.Engine)) {
	start()
	runEnv, err := app.Config().GetChildConf("global", "system", "run_env")
	if err != nil {
		log.Logger().Fatal(context.Background(), "failed to get system config")
	}
	env := runEnv.(string)
	gin.SetMode(env)
	os.Setenv("RUN_ENV", env)

	engine.Use(logger, recovery)
	service(engine)
	engine.Run(domainConf.Addr)
}

func logger(ctx *gin.Context) {
	startTime := time.Now()
	path := ctx.Request.URL.Path
	raw := ctx.Request.URL.RawQuery
	ctx.Next()

	if raw != "" {
		path = path + "?" + raw
	}

	logFields := map[string]string{
		"action":     "ServerLogger",
		"client-ip":  ctx.ClientIP(),
		"path":       path,
		"latency":    time.Now().Sub(startTime).String(),
		"status":     strconv.Itoa(ctx.Writer.Status()),
		"body-size":  strconv.Itoa(ctx.Writer.Size()),
		"user-agent": ctx.Request.UserAgent(),
	}
	logCtx := log.WithFields(ctx, logFields)
	log.Logger().Info(logCtx, "index")
}

func recovery(ctx *gin.Context) {
	//defer func() {
	//	if err := recover(); err != nil {
	//		response.RespErr(ctx, errno.Add("sys-err", errno.SysErr))
	//		return
	//	}
	//}()
	ctx.Next()
}
