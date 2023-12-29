package server

import (
	"fmt"
	"hash/fnv"
	"sort"
)

// type ServerPool struct {
// 	Servers []*Server

// 	ServersCount int
// }

// func Hash(name string) uint64 {
// 	hash := fnv.New64()
// 	hash.Write([]byte(name))
// 	return hash.Sum64()
// }

// func (pool *ServerPool) AddToPool(host, port string) {
// 	server := &Server{Host: host, Port: port}
// 	pool.Servers = append(pool.Servers, server)

// 	pool.ServersCount = len(pool.Servers)
// }

// func (pool *ServerPool) ListPool() {
// 	for _, s := range pool.Servers {
// 		fmt.Printf("pool { Host : %v, Port : %v}\n", (*s).Host, (*s).Port)
// 	}
// }

// func (pool *ServerPool) FindAvailableServer(name string) *Server {
// 	h := Hash(name)
// 	nextServerIndex := int(h) % pool.ServersCount
// 	return pool.Servers[nextServerIndex]
// }

// type Server struct {
// 	Host string
// 	Port string
// }

// func (s *Server) Address() string {
// 	return fmt.Sprintf("%s:%s", s.Host, s.Port)
// }

type ServerPool struct {
	Servers      []*Server
	ServersCount int
	HashMap      map[uint64]*Server
	//Mmap         *InMemoryFileMap
}

func Hash(name string) uint64 {
	hash := fnv.New64()
	hash.Write([]byte(name))
	return hash.Sum64()
}

func (pool *ServerPool) AddToPool(host, port string) error {
	server := &Server{Host: host, Port: port}
	pool.Servers = append(pool.Servers, server)
	pool.updateHashMap(server.Address(), server)
	pool.ServersCount = len(pool.Servers)
	return nil
}

func (pool *ServerPool) updateHashMap(address string, server *Server) error {
	h := Hash(address)
	if _, ok := pool.HashMap[h]; ok {
		return fmt.Errorf("already in map")
	}
	pool.HashMap[h] = server
	return nil
}

func (pool *ServerPool) ListPool() {
	for _, s := range pool.Servers {
		fmt.Printf("pool { Host : %v, Port : %v}\n", (*s).Host, (*s).Port)
	}
}

func (pool *ServerPool) FindAvailableServer(name string) *Server {

	// case : if filename exists in mmap
	h := Hash(name)
	keys := make([]uint64, 0, pool.ServersCount)
	for k := range pool.HashMap {
		keys = append(keys, k)
	}
	sort.Slice(keys, func(i, j int) bool {
		return keys[i] < keys[j]
	})
	i := sort.Search(len(keys), func(i int) bool {
		return keys[i] >= h
	})

	// If i is out of bounds, wrap around to the first server
	if i == len(keys) {
		i = 0
	}

	fmt.Println(i)

	return pool.HashMap[keys[i]]
}

type Server struct {
	Host string
	Port string
}

func (s *Server) Address() string {
	return fmt.Sprintf("%s:%s", s.Host, s.Port)
}
