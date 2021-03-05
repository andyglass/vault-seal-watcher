package main

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"time"
	"vault-seal-watcher/logger"

	"github.com/hashicorp/vault/api"
	"github.com/sirupsen/logrus"
)

var vaultHealthInit = true

type (
	// VaultHealth structure
	VaultHealth struct {
		Address     string
		Version     string
		ClusterName string
		Initialized bool
		Standby     bool
		Sealed      bool
	}
)

func getVaultClient(vaultAddr string, vaultTimeout time.Duration, insecureSkipVerify bool) (*api.Client, error) {
	vault, err := api.NewClient(&api.Config{
		Address: vaultAddr,
		HttpClient: &http.Client{
			Timeout: vaultTimeout,
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: insecureSkipVerify,
				},
			},
		}})
	if err != nil {
		return nil, err
	}

	return vault, nil
}

// GetVaultHealth ...
func (s *Server) GetVaultHealth() (*VaultHealth, error) {
	health, err := s.Vault.Sys().Health()
	if err != nil {
		return nil, err
	}

	vault := &VaultHealth{
		Address:     s.Cfg.VaultAddr,
		Version:     health.Version,
		ClusterName: health.ClusterName,
		Initialized: health.Initialized,
		Standby:     health.Standby,
		Sealed:      health.Sealed,
	}

	if vaultHealthInit {
		vaultHealthInit = false

		logger.Log.Infof("Vault address: %s", vault.Address)
		logger.Log.Infof("Vault cluster: %s", vault.ClusterName)
		logger.Log.Infof("Vault version: %s", vault.Version)
		logger.Log.Infof("Vault initialized: %v", vault.Initialized)
		logger.Log.Infof("Vault standby: %v", vault.Standby)
		logger.Log.Infof("Vault sealed: %v", vault.Sealed)
	}

	return vault, nil
}

//UnsealVault ...
func (s *Server) UnsealVault() error {
	type VaultUnseal struct {
		KeyShares    int
		KeyThreshold int
		KeyProgress  int
	}

	seal, err := s.Vault.Sys().ResetUnsealProcess()
	if err != nil {
		return err
	}

	vault := &VaultUnseal{
		KeyShares:    seal.N,
		KeyThreshold: seal.T,
		KeyProgress:  seal.Progress,
	}

	if len(s.Cfg.VaultUnsealKeys) < vault.KeyThreshold {
		return fmt.Errorf(
			"Provided count of unseal key-shares (%d) less than key-threshold (%d)",
			len(s.Cfg.VaultUnsealKeys),
			vault.KeyThreshold,
		)
	}

	tick := time.NewTicker(s.Cfg.VaultUnsealDelay)

	for i := vault.KeyProgress; i < vault.KeyThreshold; i++ {
		unseal, err := s.Vault.Sys().Unseal(s.Cfg.VaultUnsealKeys[i])
		if err != nil {
			return err
		}
		logger.Log.Debugf("Unseal key-shares progress: %d/%d", i+1, unseal.T)
		<-tick.C
	}

	logger.Log.Info("Vault is unsealed")

	return nil
}

// RunVaultWatcher ...
func (s *Server) RunVaultWatcher() {
	tick := time.NewTicker(s.Cfg.VaultWatchPeriod)

	for {
		vault, err := s.GetVaultHealth()
		if err != nil {
			logger.Log.Error(err)
			<-tick.C
			continue
		}

		logger.Log.WithFields(
			logrus.Fields{
				"address":     vault.Address,
				"cluster":     vault.ClusterName,
				"initialized": vault.Initialized,
				"standby":     vault.Standby,
				"sealed":      vault.Sealed,
			}).Infof("")

		if !vault.Initialized {
			logger.Log.Warn("Vault initialization required")
		}

		if vault.Sealed && vault.Initialized {
			logger.Log.Warn("Vault is in sealed state")
			err := s.UnsealVault()
			if err != nil {
				logger.Log.Error(err)
			}
		}

		<-tick.C
	}
}
