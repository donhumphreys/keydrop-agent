package main

import (
  "encoding/json"
  "fmt"
  "io/ioutil"
  "log"
  "net"
  "net/http"
  "net/url"
)

type DockerContainer struct {
  ID string
  Image string
  Config DockerContainerConfig
}

type DockerContainerConfig struct {
  Labels map[string]string
}

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

// return an http client for the docker socket
func dockerClient() (*http.Client) {
  return &http.Client{ Transport: &http.Transport { Dial: dialDockerSocket, } }
}

// ping docker to make sure our connection is configured correctly
func pingDocker() (error) {

  // execute ping request to docker socket
  response, err := dockerClient().Get("http://docker/_ping")

  // fail if an error occurs during transport
  if err != nil {
    return fmt.Errorf("Failed to connect to %s: %s", config.DockerSock, err)
  }

  // fail if docker did not respond with 200 response code
  defer response.Body.Close()
  if response.StatusCode != 200 {
    if body, err := ioutil.ReadAll(response.Body); err != nil {
      return fmt.Errorf("Failed to ping %s: %s", config.DockerSock, err)
    } else {
      return fmt.Errorf("Failed to ping %s: %s", config.DockerSock, body)
    }
  }

  return nil
}

// fetch and decode container information for the container matching the id
func getDockerContainer(id string) (container DockerContainer, err error) {
  // open an http connection to the docker sock to request container info
  response, err := dockerClient().Get(fmt.Sprintf("http://docker/containers/%.12s/json", id))
  if err != nil {
    return container, fmt.Errorf("Failed to connect to %s: %s", config.DockerSock, err)
  }

  // close the response body when this method exits
  defer response.Body.Close()

  // make sure the http response code is 200, not anything else
  if response.StatusCode != 200 {
    if body, err := ioutil.ReadAll(response.Body); err != nil {
      return container, fmt.Errorf("Unexpected %d response. %s", response.StatusCode, err)
    } else {
      return container, fmt.Errorf("Unexpected %d response: %s", response.StatusCode, body)
    }
  }

  // decode container json into object describing container
  if body, err := ioutil.ReadAll(response.Body); err != nil {
    return container, fmt.Errorf("Could not read docker response for container %.12s. %s", id, err)
  } else if err = json.Unmarshal(body, &container); err != nil {
    return container, fmt.Errorf("Could not decode container %.12s json. %s", id, err)
  }

  return container, nil
}

// decode and stream start/die container events from docker sock to an event channel
func streamDockerEvents(events chan DockerEvent, errors chan error) {

  // close the event stream when this method exists
  defer close(events)

  // prepare list of filters for docker event stream
  filters, err := dockerEventFilters()
  if err != nil {
    errors <- fmt.Errorf("Failed to prepare docker event filters. %s", err)
    return
  }

  // open an http connection to the docker sock to listen to start/die events
  response, err := dockerClient().Get("http://docker/events?filters=" + url.QueryEscape(string(filters)))
  if err != nil {
    errors <- fmt.Errorf("Failed to connect to %s: %s", config.DockerSock, err)
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
  log.Println("Listening for container events...")

  // decode events into objects and send them to event stream
  for {
    var event DockerEvent
    if err = decoder.Decode(&event); err != nil {
      errors <- err
      return
    }
    events <- event
  }
}

// filter only container start/die events with keydrop.app-id label
func dockerEventFilters() ([]byte, error) {
  filters := make(map[string][]string)
  filters["type"] = []string{"container"}
  filters["event"] = []string{"start", "die"}
  filters["label"] = []string{"keydrop.app-id"}
  return json.Marshal(filters)
}

// provide a dialer for the http client to dial the docker unix socket
func dialDockerSocket(_, _ string) (net.Conn, error) {
  return net.Dial("unix", config.DockerSock)
}
