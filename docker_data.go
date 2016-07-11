package main

type DockerActor struct {
  ID string
  Attributes map[string]string
}

type DockerContainer struct {
  ID string
  Image string
  Config DockerContainerConfig
  NetworkSettings DockerContainerNetworkSettings
}

type DockerContainerConfig struct {
  Image string
  Labels map[string]string
}

type DockerContainerNetworkSettings struct {
  IPAddress string
}

type DockerEvent struct {
  Action string
  Actor DockerActor
  Time int64
  TimeNano int64
  Type string
}
