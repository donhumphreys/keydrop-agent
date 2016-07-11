package main

import (
  "log"
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
  if config.VaultAddr, set = os.LookupEnv("VAULT_ADDR"); !set {
    config.VaultAddr = "http://127.0.0.1:8200"
  }
  config.VaultSalt, _ = os.LookupEnv("VAULT_SALT")
  return
}

// verify docker and vault connections are configured correctly
func verifyConfiguration() {

  // check connection to docker
  if err := pingDocker(); err == nil {
    log.Printf("Connected to docker socket at: unix:///%s\n", config.DockerSock)
  } else {
    log.Println(err)
    os.Exit(1)
  }

  // check status of vault server
  if err := checkVaultHealth(); err == nil {
    log.Printf("Connected to vault server at: %s\n", config.VaultAddr)
  } else {
    log.Println(err)
    os.Exit(1)
  }
}
