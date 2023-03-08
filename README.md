# EasyLog

A Simple & Easy Logging package for Go(golang). EasyLog is highly stable and thread-safe because it uses the [standard log package](https://pkg.go.dev/log).

## Installation

``` bash
$ go get github.com/Astera-org/easylog
```

## Usage

If you initialize `easylog` in `func main()`, it allows you to write log with same logger from any file in your project that is importing `easylog` package.

```go
import "github.com/Astera-org/easylog"

func main() {
    if err := easylog.Init("mydir/mylog.txt"); err != nil {
        panic(err)
    }
    // easylog.Debug calls will be ignored, since DEBUG < INFO
    easylog.SetLevel(easylog.INFO)
    // You can log strings as you would to fmt.Println
    easylog.Info("Hey", "dude")
    // You can log format strings as you would to fmt.Printf
    easylog.Errorf("I have %d computers", 12)
}
```

See [easylog_test.go](easylog_test.go) for more examples.

### Level (default: DEBUG)

Five levels defined below.

- DEBUG
- INFO
- WARN
- ERROR
- FATAL

You can set one of them for output level. Each function for levels allows you to output messages to level you want: `log.Debug()` `log.Info()` `log.Warn()` `log.Error()` `log.Fatal()`. The `log.Fatal()`, like built-in log package, outputs a message and then terminates the program.

### MaxSize (default: 1 MB)

You can set max size of log files in megabytes with `log.SetMaxSize()`. If log file exceed max size, the file will be changed to backup (`eg. sample.log.bak.20060102150405`) and create new one automatically.
