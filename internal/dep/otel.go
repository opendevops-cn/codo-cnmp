package dep

import (
	"context"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"net/url"
	"os"
	"regexp"
	"strings"

	"codo-cnmp/internal/conf"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/metrics"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"

	"go.opentelemetry.io/otel/exporters/prometheus"
	"go.opentelemetry.io/otel/metric"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"

	sdktrace "go.opentelemetry.io/otel/sdk/trace"

	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/trace"

	kratoszap "github.com/go-kratos/kratos/contrib/log/zap/v2"
)

func NewMeterProvider(bc *conf.Bootstrap) (metric.MeterProvider, error) {
	md := bc.APP
	metricConf := bc.OTEL.METRIC
	exporter, err := prometheus.New()
	if err != nil {
		return nil, err
	}

	if metricConf.ENABLE_EXEMPLAR {
		err = metrics.EnableOTELExemplar()
		if err != nil {
			return nil, err
		}
	}

	provider := sdkmetric.NewMeterProvider(
		sdkmetric.WithResource(
			resource.NewWithAttributes(
				semconv.SchemaURL,
				semconv.ServiceNameKey.String(md.NAME),
				attribute.String("environment", md.ENV),
			),
		),
		sdkmetric.WithReader(exporter),
		sdkmetric.WithView(
			metrics.DefaultSecondsHistogramView(metrics.DefaultServerSecondsHistogramName),
		),
	)
	otel.SetMeterProvider(provider)
	return provider, nil
}

func NewTracerProvider(_ context.Context, bc *conf.Bootstrap, textMapPropagator propagation.TextMapPropagator, logger log.Logger) (trace.TracerProvider, error) {
	const (
		protocolUdp = "udp"
	)

	md := bc.APP
	traceConf := bc.OTEL.TRACE

	var exp sdktrace.SpanExporter

	u, err := url.Parse(traceConf.ENDPOINT)
	if err == nil && u.Scheme == protocolUdp {
		exp, err = jaeger.New(jaeger.WithAgentEndpoint(
			jaeger.WithAgentHost(u.Hostname()), jaeger.WithAgentPort(u.Port()),
		))
	} else {
		exp, err = jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(traceConf.ENDPOINT)))
	}
	if err != nil {
		return nil, err
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exp),
		sdktrace.WithResource(
			resource.NewWithAttributes(
				semconv.SchemaURL,
				semconv.ServiceNameKey.String(md.NAME),
				attribute.String("environment", md.ENV),
			),
		),
	)

	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(textMapPropagator)
	otel.SetErrorHandler(otel.ErrorHandlerFunc(func(err error) {
		helper := log.NewHelper(logger)
		helper.Errorf("[otel] error: %v", err)
	}))

	return tp, nil
}

func NewTextMapPropagator() propagation.TextMapPropagator {
	return propagation.NewCompositeTextMapPropagator(
		tracing.Metadata{},
		propagation.Baggage{},
		propagation.TraceContext{},
	)
}

type logWrapper func(level log.Level, keyvals ...interface{}) error

func (f logWrapper) Log(level log.Level, keyvals ...interface{}) error {
	return f(level, keyvals...)
}

func NewLogger(bc *conf.Bootstrap) (log.Logger, error) {
	md := bc.APP

	logConf := bc.OTEL.LOG
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "t",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stack",
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.FullCallerEncoder,
	}
	writeSyncer := zapcore.Lock(os.Stdout)
	encoder := zapcore.NewJSONEncoder(encoderConfig)
	core := zapcore.NewCore(encoder, writeSyncer, zapConvertLevel(logConf.LEVEL))
	zapLogger := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(3))
	logger := kratoszap.NewLogger(zapLogger)
	hostname, err := os.Hostname()
	if err != nil {
		return nil, err
	}
	// 定义脱敏过滤器
	filteredLogger := log.NewFilter(logger, log.FilterFunc(func(level log.Level, keyvals ...interface{}) bool {
		for i := 0; i < len(keyvals); i++ {
			if keyvals[i] == "args" {
				if val, ok := keyvals[i+1].(string); ok {
					sanitized := sanitizeArgs(val)
					keyvals[i+1] = sanitized
				}
			}
		}
		return false
	}))
	return log.With(filteredLogger,
		"service.id", hostname,
		"service.name", md.NAME,
		"service.version", "v0.0.1",
		"trace", tracing.TraceID(),
		"span", tracing.SpanID(),
	), nil
}

// sanitizeArgs 处理 args 字符串中需要脱敏的字段
func sanitizeArgs(s string) string {
	sensitiveRegex := regexp.MustCompile(`(import_detail):\{[^}].*\}|(app_secret):\\?"[^"]*\\?"`)

	// 使用回调函数处理每个匹配
	result := sensitiveRegex.ReplaceAllStringFunc(s, func(match string) string {
		// 判断是哪个字段被匹配到
		if regexp.MustCompile(`^import_detail:`).MatchString(match) {
			return "import_detail:\"******\""
		} else if regexp.MustCompile(`^app_secret:`).MatchString(match) {
			return "app_secret:\"******\""
		}
		return match
	})
	return result
}

func zapConvertLevel(level string) zapcore.Level {
	switch strings.ToLower(level) {
	case "debug":
		return zapcore.DebugLevel
	case "info":
		return zapcore.InfoLevel
	case "warning", "warn":
		return zapcore.WarnLevel
	case "error":
		return zapcore.ErrorLevel
	case "fatal":
		return zapcore.FatalLevel
	}
	return zapcore.DebugLevel
}
