package main

import (
  "fmt"
  "time"
  "os/exec"
  "bytes"
  "encoding/json"
  "net/http"
  "strings"
)

type Response struct {
  Status int
  Message string
  Error error
}

var slots chan bool

func listen(queue string, conf *Config) {
  fmt.Println("Listening for redis messages on queue", queue, "...")
  open := conf.pool
  slots = make(chan bool, open)

  for {
    used := conf.pool - open
    // fmt.Println("Used", used)
    // fmt.Println("Open", open)
    if used < conf.pool  {
      if (open < 0){ open = 0 }
      for i := 0; i < open; i++ {
        go doWork(queue, conf)
        open--
        // fmt.Println("Open after work queued", open)
      }
    } else {
      time.Sleep(1 * time.Second)
    }
    // Wait for a slot to open
    // If there's more than one done in the queue, it'll get picked off
    // the next time around
    <- slots
    open++
    // fmt.Println("Open after work done", open)
  }
}

func doWork( queue string, config *Config ) {
  var err error
  var payload Payload
  var httpResponse *http.Response
  var req *http.Request
  var out []byte
  var buff []byte

  // 1. pull work off queue
  payload, err = Read(queue, config)
  fmt.Println("Payload", payload)

  // 2. execute
  cmd := strings.Split(payload.Command, " ")
  out, err = exec.Command(cmd[0], cmd[1:]...).Output()

  // fmt.Println("Output", out, err)

  // 3. send stdout as json array in post
  var response = new(Response)
  if err != nil{
    response.Status = 1
    response.Error = err
  } else {
    response.Status = 0
    response.Message = string(out)
  }
  buff, err = json.Marshal(response)

  // fmt.Println("JSON", string(buff))
  // fmt.Println("URI", payload.URI)
  req, err = http.NewRequest("POST", payload.URI, bytes.NewBuffer(buff))
  req.Header.Set("Content-Type", "application/json")

  for key := range config.headers {
    req.Header.Set(key, config.headers[ key ])
  }

  fmt.Println("Request", req)
  client := &http.Client{}
  httpResponse, err = client.Do(req)
  if err != nil {
    fmt.Println("ERROR (request):", err)
  }
  fmt.Println("response Status:", httpResponse.Status)
  defer httpResponse.Body.Close()

  // 4. decrement open slots
  slots <- true
}
