package main

import (
  "log"
  "os"
)

var config = loadConfig()

func main() {
  // create channels for passing docker events and errors
  events := make(chan DockerEvent)
  errors := make(chan error)

  // create new thread to react to docker events
  go printEvents(events)

  // create new thread to listen to docker socket and pass events to listen
  go streamEvents(events, errors)

  // wait for errors to occur, exit unsuccessfully if they do
  if err := <- errors; err != nil {
    log.Println(err)
    os.Exit(1)
  }

  // exit successfully
  os.Exit(0)
}

// listen to event stream and print out event information
func printEvents(events chan DockerEvent) {
  for event := range events {
    log.Printf("Type: %s, ID: %s, Action: %s\n", event.Type, event.Actor.ID, event.Action)
  }
}
