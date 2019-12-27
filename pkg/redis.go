package pkg

import (
	"github.com/go-redis/redis/v7"
	"strings"
)

type Redis struct {
	client redis.Client
}

func (d *Redis) Connect(hostPort, password string) {
	d.client = *redis.NewClient(&redis.Options{
		Addr:     hostPort,
		Password: password,
		DB:       0,
	})
}

func (d *Redis) RegisterNode(macAddress, tcpReadHostPort string) error {
	return d.client.Set("node+"+macAddress, tcpReadHostPort, 0).Err()
}

func (d *Redis) getTcpReadHostPortForMacAddress(macAddress string) (string, error) {
	return d.client.Get("node+" + macAddress).Result()
}

func (d *Redis) GetTcpReadHostPortsForMacAddress(macAddress string) ([]string, error) {
	if strings.TrimSpace(macAddress) == "ff:ff:ff:ff:ff:ff" {
		macAddresses, err := d.client.Keys("node+*").Result()
		if err != nil {
			return []string{}, err
		}

		var tcpReadHostPorts []string
		for _, macAddress := range macAddresses {
			tcpReadHostPort, err := d.getTcpReadHostPortForMacAddress(strings.Replace(macAddress, "node+", "", -1))
			if err != nil {
				return []string{}, err
			}

			tcpReadHostPorts = append(tcpReadHostPorts, tcpReadHostPort)
		}

		return tcpReadHostPorts, nil
	}

	tcpReadHostPort, err := d.getTcpReadHostPortForMacAddress(macAddress)
	if err != nil {
		return []string{}, err
	}

	return []string{tcpReadHostPort}, nil
}
