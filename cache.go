package main

import (
  "gopkg.in/redis.v3"
  "encoding/json"
  "fmt"
)

type Payload struct {
  Command string
  URI     string
}

func getClient(config *Config) (*redis.Client, error) {

  client := redis.NewClient(&redis.Options{
      Addr:     config.redis,
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


func Write( queue string, payload Payload, config *Config ) (interface{}, error) {
  c, e := getClient(config)
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
func Read( queue string, config *Config ) (Payload, error){
  c, e := getClient(config)
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
