package main

import (
  "fmt"
  "time"
  "os/exec"
  "bytes"
  "encoding/json"
  "net/http"
)
var openSlots = make(chan int)

type Response struct {
  status int
  message string
}
func listen(queue string, threads int) {
  openSlots <- threads
  fmt.Println("Listening for redis messages on queue", queue, "...")

  for true {
    open := <- openSlots
    if open < threads {
      howMuch := threads - open
      for i := 0; i < howMuch; i++ {
        go doWork(queue)
      }
    } else {
      time.Sleep(time.Sleep(1 * time.Second))
    }
  }
}

func doWork( queue string ) {
  var err error
  var jsonStr string
  var response Response
  var payload Payload
  var httpResponse *http.Response
  var req *http.Request
  var out []byte
  var buff []byte

  // 1. pull work off queue
  payload, err = Read(args[1])

  // 2. execute
  out, err = exec.Command(payload.Command).Output()

  // 3. send stdout as json array in post
  if err != nil{
    response.status = 1
    response.message = err
  } else {
    response.status = 0
    response.message = string(out)
  }
  buff, err = json.Marshal(payload)
  httpResponse, err = http.Post(payload.URI, "application/json", &buff)
  req, err = http.NewRequest("POST", url, bytes.NewBuffer(buff))
  req.Header.Set("Content-Type", "application/json")
  // TODO: Loop through a config of headers that the user defines

  client := &http.Client{}
  httpResponse, err = client.Do(req)
  if err != nil {
      fmt.Println("ERROR (request):", err)
  }
  defer httpResponse.Body.Close()
  fmt.Println("response Status:", httpResponse.Status)


  // 4. decrement openSlots
  open := <-openSlots
  open--
  openSlots <- open
}
