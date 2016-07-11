package main

import (
  "log"
  "os"
)

var config = loadConfig()

func main() {
  // validate configuration before starting work
  verifyConfiguration()

  // create channels for passing docker events and errors
  events := make(chan DockerEvent)
  errors := make(chan error)

  // create new thread to react to docker events
  go printEvents(events, errors)

  // create new thread to listen to docker socket and pass events to listen
  go streamDockerEvents(events, errors)

  // wait for errors to occur, exit unsuccessfully if they do
  if err := <- errors; err != nil {
    log.Println(err)
    os.Exit(1)
  }

  // exit successfully
  os.Exit(0)
}

// listen to event stream and print out event information
func printEvents(events chan DockerEvent, errors chan error) {
  for event := range events {
    switch event.Action {
      case "start":
        handleStartEvent(event, errors)
      case "die":
        handleStopEvent(event, errors)
      default:
        log.Printf("Unhandled event \"%s\" for container %.12s.\n", event.Action, event.Actor.ID)
    }
  }
}

// fetch and print detailed container information for container start event
func handleStartEvent(event DockerEvent, errors chan error) {

  // read container information from docker api
  container, err := getDockerContainer(event.Actor.ID)
  if err != nil {
    errors <- err
  }

  // print event and container information
  log.Printf("ID: %.12s, Action: %s, SHA: %.19s, App-ID: %s, IP: %s", event.Actor.ID, event.Action, container.Image, container.Config.Labels["keydrop.app-id"], container.NetworkSettings.IPAddress)
}

// print event information for container stop event
func handleStopEvent(event DockerEvent, errors chan error) {
  log.Printf("ID: %.12s, Action: %s\n", event.Actor.ID, event.Action)
}
