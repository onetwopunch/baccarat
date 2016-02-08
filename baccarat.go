package main

import (
  "fmt"
  "os"
  "strconv"
  "runtime"
)

func main() {
  args := os.Args[1:]
  isQ, queue := namedCommand(args, "-q")
  isPool, pool := namedCommand(args, "--pool")
  if isQ {
    if isPool{
      var err error
      threads, err := strconv.ParseInt(pool, 10, 32)
      if err != nil {
        fmt.Println(err)
      }
      listen( queue, int(threads) )
    } else {
      threads := maxThreads()
      listen( queue, threads )
    }
  } else {
    // baccarat write queue 'ls -la' http://google.com

    cmd := args[0]
    if cmd == "write" {
      if len(args) != 4 {
        fmt.Println("Incorrect number of arguments")
        return
      }
      payload := Payload{args[2], args[3]}
      val, e := Write(args[1], payload)
      fmt.Println("Result:", val, e)
    } else if cmd == "read" {
      if len(args) != 2 {
        fmt.Println("Incorrect number of arguments")
        return
      }
      payload, e := Read(args[1])
      fmt.Println("Result:", payload, e)
    }
  }
}

func maxThreads() int {
    maxProcs := runtime.GOMAXPROCS(0)
    numCPU := runtime.NumCPU()
    if maxProcs < numCPU {
        return maxProcs
    }
    return numCPU
}

func namedCommand(args []string, item string) (bool, string) {
  for i, arg := range args {
    if arg == item {
      return true, args[i + 1]
    }
  }
  return false, ""
}
