package storage

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/buraksezer/olric"
	"github.com/buraksezer/olric/config"
)

// DHT is a distributed hash table
type DHT struct {
	laddr, mladdr, maladdr *net.TCPAddr
	peers                  []*net.TCPAddr
	startedChan            chan bool
	db                     *olric.Olric
}

// NewDHT creates a new DHT
func NewDHT(laddr, mladdr, maladdr *net.TCPAddr, peers []*net.TCPAddr) *DHT {
	return &DHT{
		laddr,
		mladdr,
		maladdr,
		peers,
		make(chan bool),
		nil,
	}
}

// Open opens the DHT
func (d *DHT) Open() error {
	config := config.New("local")
	config.Name = fmt.Sprintf("%v:%v", d.laddr.IP.String(), d.laddr.Port)
	config.BindAddr = d.laddr.IP.String()
	config.BindPort = d.laddr.Port
	peers := []string{}
	for _, peer := range d.peers {
		peers = append(peers, peer.String())
	}
	config.Peers = peers
	config.Started = func() {
		d.startedChan <- true
	}

	config.MemberlistConfig.BindAddr = d.mladdr.IP.String()
	config.MemberlistConfig.BindPort = d.mladdr.Port
	config.MemberlistConfig.AdvertiseAddr = d.maladdr.IP.String()
	config.MemberlistConfig.AdvertisePort = d.maladdr.Port
	config.MemberlistConfig.Name = fmt.Sprintf("%v:%v", d.mladdr.IP.String(), d.mladdr.Port)

	db, err := olric.New(config)
	if err != nil {
		return err
	}

	d.db = db

	return nil
}

// Start starts the DHT
func (d *DHT) Start() error {
	errChan := make(chan error)

	go func() {
		if err := d.db.Start(); err != nil {
			errChan <- err
		}
	}()

	for {
		select {
		case <-d.startedChan:
			return nil
		case err := <-errChan:
			return err
		}
	}
}

// GetMap gets a distributed map (if map does not exist, it creates it)
func (d *DHT) GetMap(name string) (*olric.DMap, error) {
	return d.db.NewDMap(name)
}

// Stop stops the DHT
func (d *DHT) Stop() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	defer cancel()
	return d.db.Shutdown(ctx)
}
