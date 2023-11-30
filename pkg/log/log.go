package log

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"
)

type Log struct {
	logger *zap.Logger
}

const loggerKey = iota
const (
	FieldReqId = "req-id"
)

var log = &Log{}

func Start() {
	level := zap.DebugLevel
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	encoderConfig.MessageKey = "msg"
	encoderConfig.TimeKey = "ts"
	encoderConfig.LevelKey = "level"
	encoder := zapcore.NewConsoleEncoder(encoderConfig)

	core := zapcore.NewCore(
		encoder,
		zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), zapcore.AddSync(getWriter())),
		level,
	)
	log.logger = zap.New(core, zap.AddStacktrace(zap.ErrorLevel))
}

func getWriter() io.Writer {
	logDir := GetRoot()
	logWriter, _ := rotatelogs.New(logDir+"/%Y%m%d.log",
		rotatelogs.WithMaxAge(time.Hour*24*7),
		rotatelogs.WithRotationTime(time.Hour*24),
	)
	return logWriter
}

func Logger() *Log {
	return log
}

func WithFields(ctx context.Context, fields map[string]string) context.Context {
	fieldArr := make([]zap.Field, 0)
	if _, ok := fields[FieldReqId]; !ok {
		reqId := ctx.Value(FieldReqId)
		fieldArr = append(fieldArr, zap.Any(FieldReqId, reqId))
	}
	uid := ctx.Value("address")
	if uid != nil {
		fieldArr = append(fieldArr, zap.Any("uid", uid))
	}

	for k, v := range fields {
		f := zap.Any(k, v)
		fieldArr = append(fieldArr, f)
	}
	l := WithCtx(ctx)
	return context.WithValue(ctx, loggerKey, l.With(fieldArr...))
}

func NewContext(ctx context.Context, fields ...zapcore.Field) context.Context {
	return context.WithValue(ctx, loggerKey, WithCtx(ctx).With(fields...))
}

func WithCtx(ctx context.Context) *zap.Logger {
	if ctx == nil {
		return log.logger
	}
	if ctxLogger, ok := ctx.Value(loggerKey).(*zap.Logger); ok {
		return ctxLogger
	}
	return log.logger
}

func (l *Log) Info(ctx context.Context, args ...interface{}) {
	pc, _, line, _ := runtime.Caller(1)
	f := runtime.FuncForPC(pc)
	WithCtx(ctx).Info(fmt.Sprint(args), zap.String("line", strconv.Itoa(line)), zap.String("func", f.Name()))
}

func (l *Log) Error(ctx context.Context, args ...interface{}) {
	pc, _, line, _ := runtime.Caller(1)
	f := runtime.FuncForPC(pc)
	WithCtx(ctx).Error(fmt.Sprint(args), zap.String("line", strconv.Itoa(line)), zap.String("func", f.Name()))
}

func (l *Log) Debug(ctx context.Context, args ...interface{}) {
	pc, _, line, _ := runtime.Caller(1)
	f := runtime.FuncForPC(pc)
	WithCtx(ctx).Debug(fmt.Sprint(args), zap.String("line", strconv.Itoa(line)), zap.String("func", f.Name()))
}

func (l *Log) Warn(ctx context.Context, args ...interface{}) {
	pc, _, line, _ := runtime.Caller(1)
	f := runtime.FuncForPC(pc)
	WithCtx(ctx).Warn(fmt.Sprint(args), zap.String("line", strconv.Itoa(line)), zap.String("func", f.Name()))
}

func (l *Log) Fatal(ctx context.Context, args ...interface{}) {
	pc, _, line, _ := runtime.Caller(1)
	f := runtime.FuncForPC(pc)
	WithCtx(ctx).Fatal(fmt.Sprint(args), zap.String("line", strconv.Itoa(line)), zap.String("func", f.Name()))
}

func GrpcUnaryServerInterceptor(l *Log) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		startTime := time.Now()
		reqId := ""
		headers, ok := metadata.FromIncomingContext(ctx)
		if ok {
			reqIdArr := headers.Get("Req-Id")
			if len(reqIdArr) > 0 {
				reqId = reqIdArr[0]
			}
		}
		if reqId == "" {
			reqId = uuid.New().String()
		}
		items := map[string]string{
			"method": info.FullMethod,
			"req-id": reqId,
		}
		newCtx := WithFields(ctx, items)
		resp, err = handler(newCtx, req)
		code := status.Code(err)
		duration := time.Since(startTime)

		items1 := map[string]string{
			"code":     strconv.Itoa(int(code)),
			"duration": duration.String(),
		}
		fields := make([]zap.Field, 0)
		for k, v := range items1 {
			f := zap.String(k, v)
			fields = append(fields, f)
		}
		WithCtx(newCtx).Info("serverLogger", fields...)
		return resp, err
	}
}

func GetRoot() string {
	dir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	return strings.Replace(dir+"/logs", "\\", "/", -1)
}
