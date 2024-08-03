# Confo

Golang command line tool with environment supproted, Command line arguments have higher priority than environment variables.


## Examples

```go
logFile    = confo.NewString("log.file", "", "the log file")
logLevel   = confo.NewBool("log.level", false, "the log level, default is bool")
serverHost = confo.NewString("server.addr", ":7070", "the address http will serve on")
serverPort = confo.NewInt("server.port", 8080, "the tcp port will listen on")
secs = confo.NewDuration("secs", "30s", "the duration flag, 30 secodns")
secs2 = confo.NewDuration("days", "30d", "the duration flag, 30 days")
secs3 = confo.NewDuration("weeks", "2w", "the duration flag, 2 weeks")
bytes = confo.NewBytes("bytes", 128, "the bytes flag, 128 Byte")
kbs = confo.NewBytes("kbytes", 64, "the bytes flag, 64KB")
mbs = confo.NewBytes("mbytes", 64, "the bytes flag, 64KB")
gbs = confo.NewBytes("gbytes", 64, "the bytes flag, 64KB")
strs = confo.NewArrayString("array.str", "the string array flag, default is empty")
ints = confo.NewArrayInt("array.int", "the string array flag, default is empty")
```