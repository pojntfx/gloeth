package storage

import (
	"net"
	"sync"
	"testing"
	"time"
)

func TestDHT(t *testing.T) {
	laddrs := []string{"localhost:10001", "localhost:10002", "localhost:10003", "localhost:10004"}
	mladdrs := []string{"localhost:20001", "localhost:20002", "localhost:20003", "localhost:20004"}
	maladdrs := []string{"localhost:30001", "localhost:30002", "localhost:30003", "localhost:30004"}

	nodes := []*DHT{}
	for i, rawLaddr := range laddrs {
		peers := []*net.TCPAddr{}
		for ii, rawPeer := range mladdrs {
			if i != ii {
				peer, err := net.ResolveTCPAddr("tcp", rawPeer)
				if err != nil {
					t.Error(err)
				}

				peers = append(peers, peer)
			}
		}

		laddr, err := net.ResolveTCPAddr("tcp", rawLaddr)
		if err != nil {
			t.Error(err)
		}

		mladdr, err := net.ResolveTCPAddr("tcp", mladdrs[i])
		if err != nil {
			t.Error(err)
		}

		maladdr, err := net.ResolveTCPAddr("tcp", maladdrs[i])
		if err != nil {
			t.Error(err)
		}

		nodes = append(nodes, NewDHT(laddr, mladdr, maladdr, peers /* []*net.TCPAddr{peers[0]} would work here as well, only one node is required to discover the entire cluster */))
	}

	var wg sync.WaitGroup

	wg.Add(len(nodes))

	for _, node := range nodes {
		if err := node.Open(); err != nil {
			t.Error(err)
		}

		go func(inode *DHT, iwg *sync.WaitGroup) {
			if err := inode.Start(); err != nil {
				t.Error(err)
			}

			iwg.Done()
		}(node, &wg)
	}

	wg.Wait()

	time.Sleep(time.Second * 1)

	map1 := "map1"
	key1 := "key1"
	val1 := "val1"

	dmap1, err := nodes[0].GetMap(map1)
	if err != nil {
		t.Error(err)
	}

	if err := dmap1.Put(key1, val1); err != nil {
		t.Error(err)
	}

	wg.Add(len(nodes))

	for _, node := range nodes {
		go func(inode *DHT, iwg *sync.WaitGroup) {
			dmap, err := inode.GetMap(map1)
			if err != nil {
				t.Error(err)
			}

			outVal, err := dmap.Get(key1)
			if err != nil {
				t.Error(err)
			}

			t.Logf("node %v returned outVal %v", inode.laddr, outVal)

			if outVal != val1 {
				t.Errorf("outVal != val1: got %v, want %v", outVal, val1)
			}

			iwg.Done()
		}(node, &wg)
	}

	wg.Wait()

	wg.Add(len(nodes))

	for _, node := range nodes {
		go func(inode *DHT, iwg *sync.WaitGroup) {
			if err := inode.Stop(); err != nil {
				t.Error(err)
			}

			iwg.Done()
		}(node, &wg)
	}

	wg.Wait()
}
