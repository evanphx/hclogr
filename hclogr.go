package hclogr

import (
	"github.com/go-logr/logr"
	"github.com/hashicorp/go-hclog"
)

// Wrap returns a new logr.Logger that sends it's logs to the given hclog.Logger
func Wrap(log hclog.Logger) logr.Logger {
	return logr.New(&logSink{
		log: log,
	})
}

func New(cfg *hclog.LoggerOptions) logr.Logger {
	cp := *cfg

	return logr.New(&logSink{
		log: hclog.New(cfg),
		cfg: &cp,
	})
}

type logSink struct {
	log hclog.Logger
	cfg *hclog.LoggerOptions
}

var _ logr.LogSink = &logSink{}

func (l *logSink) Enabled(level int) bool {
	return l.log.GetLevel() <= fromLevel(level)
}

func (l *logSink) Error(err error, msg string, kvs ...interface{}) {
	vals := append([]interface{}{"error", err}, kvs...)
	l.log.Error(msg, vals...)
}

func (l *logSink) Info(level int, msg string, kvs ...interface{}) {
	l.log.Log(fromLevel(level), msg, kvs...)
}

func (l *logSink) Init(info logr.RuntimeInfo) {
	if l.cfg != nil {
		l.cfg.AdditionalLocationOffset = info.CallDepth + 1
		l.log = hclog.New(l.cfg)
	}
}

func (l *logSink) WithName(name string) logr.LogSink {
	return &logSink{
		log: l.log.Named(name),
	}
}

func (l *logSink) WithValues(kvs ...interface{}) logr.LogSink {
	return &logSink{
		log: l.log.With(kvs),
	}
}

func toLevel(hl hclog.Level) int {
	switch hl {
	case hclog.Trace:
		return 4
	case hclog.Debug:
		return 3
	case hclog.Info:
		return 2
	case hclog.Warn:
		return 1
	case hclog.Error:
		return 0
	default:
		return 2
	}
}

var standardLevels = map[int]hclog.Level{
	0: hclog.Error,
	1: hclog.Warn,
	2: hclog.Info,
	3: hclog.Debug,
	4: hclog.Trace,
}

func fromLevel(level int) hclog.Level {
	hl, ok := standardLevels[level]
	if ok {
		return hl
	}

	if level > 4 {
		return hclog.Trace
	}

	return hclog.Error
}
