## Synopsis

Defines an interface for loggers and stats reporters that can be passed to Uber Go libraries.  
Provides implementations which wrap a common logging module, [logrus](https://github.com/Sirupsen/logrus), 
and a common stats reporting module [go-statsd-client](https://github.com/cactus/go-statsd-client).  
Clients may also choose to implement these interfaces themselves.

## Key Interfaces

### Logging

```go
/*
 * Interface for loggers accepted by Uber's libraries.
 */
type Logger interface {
	Debug(args ...interface{})
	Debugf(format string, args ...interface{})
	Info(args ...interface{})
	Infof(format string, args ...interface{})
	Warn(args ...interface{})
	Warnf(format string, args ...interface{})
	Error(args ...interface{})
	Errorf(format string, args ...interface{})
	Fatal(args ...interface{})
	Fatalf(format string, args ...interface{})
	Panic(args ...interface{})
	Panicf(format string, args ...interface{})
	WithField(key string, value interface{}) Logger
	WithFields(keyValues LogFields) Logger
}
```

### Stats Reporting

```go
/*
 * Interface for statsd-like stats reporters accepted by Uber's libraries.
 */
type StatsReporter interface {
	IncCounter(name string, tags map[string]string, value int64)
	UpdateGauge(name string, tags map[string]string, value int64)
	RecordTimer(name string, tags map[string]string, d time.Duration)
}
```

## Basic Usage

```go
logger := logrus.New()
barkLogger := bark.NewLoggerFromLogrus(logger)
barkLogger.WithFields(bark.Fields{"someField":"someValue"}).Info("Message")

statsd, err := statsd.New("127.0.0.1:8125", "barktest")
if err != nil {
    logger.Fatal("Example code failed")
}

barkStatsReporter := bark.NewStatsReporterFromCactus(statsd)  
barkStatsReporter.IncCounter("foo", map[string]string{"tag":"val"}, 1)
 
ubermodule.New(ubermodule.Config{
    logger: barkLogger
    statsd: barkStatsReporter
})
```

## Contributors

dh

## License

bark is available under the MIT license. See the LICENSE file for more info.
