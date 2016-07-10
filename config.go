package main

import (
  "os"
)

type Config struct {
  DockerSock string
  VaultAddr string
  VaultSalt string
}

// load global configuration from environment variables
func loadConfig() (config Config) {
  var set bool
  if config.DockerSock, set = os.LookupEnv("DOCKER_SOCK"); !set {
    config.DockerSock = "/var/run/docker.sock"
  }
  config.VaultAddr, _ = os.LookupEnv("VAULT_ADDR")
  config.VaultSalt, _ = os.LookupEnv("VAULT_SALT")
  return
}
