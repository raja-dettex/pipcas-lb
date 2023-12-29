package server

import (
	"fmt"
	"testing"
)

func TestSeverPool(t *testing.T) {
	pool := &ServerPool{HashMap: make(map[uint64]*Server)}
	for i := 1; i <= 4; i++ {
		pool.AddToPool(fmt.Sprintf("testHost%d", i), fmt.Sprintf("testPort%d", i))
	}
	for i := 1; i < 7; i++ {
		s := pool.FindAvailableServer(fmt.Sprintf("hello%d.txt", i))
		fmt.Println(s.Address())
	}
	fmt.Println("-----------")
	pool.AddToPool("testHost5", "testPort5")
	for i := 1; i < 7; i++ {
		s := pool.FindAvailableServer(fmt.Sprintf("hello%d.txt", i))
		fmt.Println(s.Address())
	}
}
