package logger

import (
	"context"
	"fmt"
	"io"
	"os"
	"runtime"
	"strings"
	"sync"

	"github.com/aws/aws-lambda-go/lambdacontext"
	"github.com/goccy/go-json"
	reflect "github.com/goccy/go-reflect"
	"go.uber.org/zap"
	"go.uber.org/zap/buffer"
	"go.uber.org/zap/zapcore"
)

var instance ZapLogger
var lock = &sync.Mutex{}

type zap_logger struct {
	log *zap.Logger
}

type customEncoder struct {
	zapcore.Encoder
}

type Logger interface {
	ZapLogger
}

type ZapLogger interface {
	AnyField(key string, value interface{}) zap.Field
	StringField(key, value string) zap.Field
	Any(val interface{}) zap.Field
	Struct(val interface{}) zap.Field

	Debugf(format string, v ...interface{})
	Debug(msg string, keyvals ...zap.Field)
	Errorf(format string, v ...interface{})
	Error(msg string, keyvals ...zap.Field)
	Infof(format string, v ...interface{})
	Info(msg string, keyvals ...zap.Field)
}

func (l *zap_logger) Struct(val interface{}) zap.Field {
	return customField(val)
}

func (l *zap_logger) Any(value interface{}) zap.Field {
	return customField(value)
}

func (l *zap_logger) StringField(key, value string) zap.Field {
	return zap.String(key, value)
}

func (l *zap_logger) AnyField(key string, value interface{}) zap.Field {
	return zap.Any("data."+key, value)
}

func (l *zap_logger) Debugf(format string, v ...interface{}) {
	l.Debug(fmt.Sprint(format, v))
}

func (l *zap_logger) Debug(msg string, keyvals ...zap.Field) {
	l.log.Debug(msg, keyvals...)
}

func (l *zap_logger) Errorf(format string, v ...interface{}) {
	l.Error(fmt.Sprint(format, v))
}

func (l *zap_logger) Error(msg string, keyvals ...zap.Field) {
	l.log.Error(msg, keyvals...)
}

func (l *zap_logger) Infof(format string, v ...interface{}) {
	l.Info(fmt.Sprint(format, v))
}

func (l *zap_logger) Info(msg string, keyvals ...zap.Field) {
	l.log.Info(msg, keyvals...)
}

// NewNopLambdaLogger returns a new LambdaLogger with a nop writer (no operations).
func NewNopLambdaLogger() ZapLogger {
	return &zap_logger{log: zap.NewNop()}
}

// NewDefaultLogger configures logger.
func NewDefaultLogger(cfg *LoggerConfiguration, ctx context.Context, businessId string) ZapLogger {
	return &zap_logger{log: newZapLogger(cfg, ctx, businessId)}
}

func newZapLogger(cfg *LoggerConfiguration, ctx context.Context, businessId string) *zap.Logger {
	encoderCfg := zapcore.EncoderConfig{
		MessageKey:          "data.message",
		LevelKey:            "level",
		TimeKey:             "date",
		NameKey:             "logger",
		CallerKey:           "caller",
		FunctionKey:         "function",
		StacktraceKey:       "error.stackTrace",
		SkipLineEnding:      false,
		EncodeLevel:         zapcore.CapitalLevelEncoder,
		EncodeTime:          zapcore.RFC3339TimeEncoder,
		EncodeDuration:      zapcore.StringDurationEncoder,
		NewReflectedEncoder: defaultReflectedEncoder,
	}

	encoder := &customEncoder{Encoder: zapcore.NewJSONEncoder(encoderCfg)}

	addDefaultFields(encoder, cfg, ctx, businessId)

	zapLevel := zap.NewAtomicLevelAt(zap.InfoLevel)

	if strings.ToLower(cfg.General.LogLevel) == "debug" {
		zapLevel = zap.NewAtomicLevelAt(zap.DebugLevel)
	}

	core := zapcore.NewCore(encoder, os.Stderr, zapLevel)

	return zap.New(core)
}

func addDefaultFields(encoder zapcore.Encoder, cfg *LoggerConfiguration, ctx context.Context, businessId string) {
	encoder.AddString("data.businessId", businessId)
	lc, _ := lambdacontext.FromContext(ctx)
	encoder.AddString("data.requestId", lc.AwsRequestID)
	encoder.AddString("func.country", cfg.General.Country)
	encoder.AddString("func.env", cfg.General.Environment)
	encoder.AddString("func.version", cfg.General.Version)
	encoder.AddString("func.name", cfg.General.Name)
	encoder.AddString("func.level", cfg.General.LogLevel)
	encoder.AddString("func.runtime", runtime.Version())
}

func defaultReflectedEncoder(w io.Writer) zapcore.ReflectedEncoder {
	enc := json.NewEncoder(w)

	enc.SetEscapeHTML(false)

	return enc
}

func (e *customEncoder) EncodeEntry(entry zapcore.Entry, fields []zapcore.Field) (*buffer.Buffer, error) {
	return e.Encoder.EncodeEntry(entry, fields) //nolint: wrapcheck
}

func customField(val interface{}) zap.Field {
	reflectType := reflect.TypeOf(val)

	var key string

	if reflectType.Kind() == reflect.Ptr {
		key = reflectType.Elem().Name()
	} else {
		key = reflectType.Name()
	}

	reflectTypeValue := reflect.ValueOf(val)

	switch reflectTypeValue.Kind() { //nolint: nolintlint
	case reflect.Struct:
		return zap.Any(key, maskSensitiveFields(reflectType, reflectTypeValue))
	default:
		return zap.Any(key, val)
	}
}

func maskSensitiveFields(reflectType reflect.Type, reflectTypeValue reflect.Value) map[string]interface{} {
	out := map[string]interface{}{}

	for index := 0; index < reflectType.NumField(); index++ {
		reflectTypeField := reflectType.Field(index)

		jsonTag := reflectTypeField.Tag.Get("json")

		if reflectTypeField.Tag.Get("zap") == "sensitive" {
			out[jsonTag] = "[filtered]"
		} else {
			if jsonTag != "" {
				out[strings.Split(jsonTag, ",")[0]] = reflectTypeValue.Field(index).Interface()
			} else {
				out[reflectTypeField.Name] = reflectTypeValue.Field(index).Interface()
			}
		}
	}

	return out
}

func GetLogger(ctx context.Context) ZapLogger {
	return SetLogger(ctx, "")
}

func SetLogger(ctx context.Context, businessId string) ZapLogger {
	if instance != nil {
		return instance
	}
	cfg := NewLoggerConfiguration()
	instance = NewDefaultLogger(cfg, ctx, businessId)
	return instance
}
