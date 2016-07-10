package main

import (
  "encoding/json"
  "fmt"
  "io/ioutil"
  "net"
  "net/http"
  "net/url"
)

type DockerEvent struct {
  Action string
  Actor DockerActor
  Time int64
  TimeNano int64
  Type string
}

type DockerActor struct {
  ID string
  Attributes map[string]string
}

// provide a dialer for the http client to dial the docker unix socket
func dialUnix(_, _ string) (net.Conn, error) {
  return net.Dial("unix", config.DockerSock)
}

// decode and stream start/die container events from docker sock to an event channel
func streamEvents(events chan DockerEvent, errors chan error) {

  // close the event stream when this method exists
  defer close(events)

  // open an http connection to the docker sock to listen to start/die events
  client := &http.Client{ Transport: &http.Transport { Dial: dialUnix, } }
  response, err := client.Get("http://docker/events?filters=" + url.QueryEscape("{\"event\": [\"start\", \"die\"]}"))
  if err != nil {
    errors <- fmt.Errorf("Failed to connect to docker. %s", err)
    return
  }

  // close the response body when this method exits
  defer response.Body.Close()

  // make sure the http response code is 200, not anything else
  if response.StatusCode != 200 {
    if body, err := ioutil.ReadAll(response.Body); err != nil {
      errors <- fmt.Errorf("Unexpected %d response. %s", response.StatusCode, err)
    } else {
      errors <- fmt.Errorf("Unexpected %d response: %s", response.StatusCode, body)
    }
    return
  }

  // attach a json decoder to the response stream
  decoder := json.NewDecoder(response.Body)

  // decode events into objects and send them to event stream
  for {
    var event DockerEvent
    if err = decoder.Decode(&event); err != nil {
      errors <- err
      break
    }
    events <- event
  }
}
