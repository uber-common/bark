## Synopsis

Defines an interface for loggers and stats reporters that can be passed to Uber Go libraries.  
Provides implementations which wrap a common logging module, [logrus](https://github.com/Sirupsen/logrus), 
and a common stats reporting module [go-statsd-client](https://github.com/cactus/go-statsd-client).  
Clients may also choose to implement these interfaces themselves.

## Basic Usage

```go
logger := logrus.New()
statsd, err := statsd.New("127.0.0.1:8125", "barktest")
if err != nil {
    logger.Fatal("Example code failed")
}

ubermodule.New(ubermodule.Config{
    logger: bark.NewLoggerFromLogrus(logger)
    statsd: bark.NewStatsReporterFromCactus(statsd)
})
```

## Contributors

dh

## License

bark is available under the MIT license. See the LICENSE file for more info.
