package main

import (
  "fmt"
  "net/http"
)

// ping vault to make sure we can connect and that the vault is unsealed
func checkVaultHealth() (error) {

  // execute ping request to docker socket
  response, err := http.Head(vaultUrl("sys/health"))

  // fail if an error occurs during transport
  if err != nil {
    return fmt.Errorf("Failed to connect to vault at %s", err)
  }

  // fail if vault did not respond with 200 response code
  if response.StatusCode != 200 {
    return fmt.Errorf("Found unhealthy or sealed vault at %s", config.VaultAddr)
  }

  return nil
}

// generate url for vault api
func vaultUrl(path string) (string) {
  return fmt.Sprintf("%s/v1/%s", config.VaultAddr, path)
}
