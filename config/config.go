package config

import (
	"github.com/codegangsta/cli"
	"github.com/rancher/go-rancher/v2"
)

type Config struct {
	KubernetesURL   string
	CattleURL       string
	CattleAccessKey string
	CattleSecretKey string
	WorkerCount     int
	HealthCheckPort int
}

func Conf(context *cli.Context) Config {
	config := Config{
		KubernetesURL:   context.String("kubernetes-url"),
		CattleURL:       context.String("cattle-url"),
		CattleAccessKey: context.String("cattle-access-key"),
		CattleSecretKey: context.String("cattle-secret-key"),
		WorkerCount:     context.Int("worker-count"),
		HealthCheckPort: context.Int("health-check-port"),
	}

	return config
}

func GetRancherClient(conf Config) (*client.RancherClient, error) {
	return client.NewRancherClient(&client.ClientOpts{
		Url:       conf.CattleURL,
		AccessKey: conf.CattleAccessKey,
		SecretKey: conf.CattleSecretKey,
	})
}
