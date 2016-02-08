package main

import (
  "github.com/Unknwon/goconfig"
  "gopkg.in/redis.v3"
  "encoding/json"
  "fmt"
)

type Payload struct {
  Command string
  URI     string
}

func getClient() (*redis.Client, error) {
  c, _ := goconfig.LoadConfigFile(".config.ini")
  url, _ := c.GetValue("redis", "url")
  password, _ := c.GetValue("password", "url")
  client := redis.NewClient(&redis.Options{
      Addr:     url,
      Password: password, // no password set
      DB:       0,  // use default DB
  })

  pong, err := client.Ping().Result()
  if err != nil {
    fmt.Println(pong, err)
    return client, err
  } else {
    return client, nil
  }
}

func cacheKey(queue string) (string) {
  return "baccarat:"+ queue;
}


func Write( queue string, payload Payload ) (interface{}, error) {
  c, e := getClient()
  if e != nil {
    return "", e
  } else {
    key := cacheKey(queue)
    bytes, err := json.Marshal(payload)
    value := string(bytes)

    if err != nil {
      return "", err
    } else {
      fmt.Println("Pushing", value, "onto fifo", key)
      return c.LPush(key, string(bytes)).Result()
    }
  }
  return "", nil
}
func Read( queue string ) (Payload, error){
  c, e := getClient()
  var payload Payload
  if e != nil {
    return payload, e
  } else {

    key := cacheKey(queue)
    value, err := c.BRPop(0, key).Result()

    fmt.Println("Popping", value[1], "out of fifo", key)
    if err != nil {
      return payload, err
    }

    err = json.Unmarshal([]byte(value[1]), &payload)
    if err != nil {
      return payload, err
    } else {
      return payload, nil
    }

  }
}
