package main

import (
  "fmt"
  "os"
  "runtime"
  "github.com/Unknwon/goconfig"
)
type Config struct {
  pool int
  redis string
  headers map[string]string
}

func main() {
  args := os.Args[1:]
  hasQ, queue := namedCommand(args, "-q")
  hasConfig, path := namedCommand(args, "-c")
  config, err := NewConfig(path)
  if hasQ && hasConfig {
    if err != nil{
      fmt.Println("ERROR (config)", err)
      return
    }
    listen( queue, config )

  } else if hasConfig {
    // baccarat write queue 'ls -la' http://google.com
    cmd := args[0]
    if cmd == "write" {
      if len(args) < 4 {
        fmt.Println("Incorrect number of arguments")
        return
      }
      payload := Payload{args[2], args[3]}
      val, e := Write(args[1], payload, config)
      fmt.Println("Result:", val, e)
    } else if cmd == "read" {
      if len(args) < 2 {
        fmt.Println("Incorrect number of arguments")
        return
      }
      payload, e := Read(args[1], config)
      fmt.Println("Result:", payload, e)
    } else {
      usage()
    }
  }
}
func usage()  {
  fmt.Println("baccarat - Concurrent background job scheduler")
  fmt.Println("    Listen: baccarat -q queueName -c path/to/config.ini")
  fmt.Println("    Write:  baccarat write queueName 'echo $PATH' http://example.com/notifications -c path/to/config.ini")
  fmt.Println("    Read:   baccarat read queueName -c path/to/config.ini")
}
func maxThreads() int {
    maxProcs := runtime.GOMAXPROCS(0)
    numCPU := runtime.NumCPU()
    if maxProcs < numCPU {
        return maxProcs
    }
    return numCPU
}
func NewConfig(path string) (*Config, error) {
  var config Config
  c, err := goconfig.LoadConfigFile(path)
  if err != nil {
    return &config, err
  }
  config.headers, err = c.GetSection("headers")
  config.redis, err = c.GetValue("redis", "url")
  pool, pErr := c.Int("default", "pool")

  if pErr == nil {
    config.pool = pool
  } else {
    config.pool = maxThreads()
  }
  return &config, err
}
func namedCommand(args []string, item string) (bool, string) {
  for i, arg := range args {
    if arg == item {
      return true, args[i + 1]
    }
  }
  return false, ""
}
