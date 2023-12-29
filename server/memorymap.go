package server

import "fmt"

type InMemoryFileMap struct {
	FileMap map[string]uint64
}

func NewInMemoryFileMap() *InMemoryFileMap {
	return &InMemoryFileMap{FileMap: map[string]uint64{}}
}

func (mmap *InMemoryFileMap) Store(fName string, hash uint64) error {
	if _, ok := mmap.FileMap[fName]; ok {
		return fmt.Errorf("exists file: [%s]", fName)
	}
	mmap.FileMap[fName] = hash
	return nil
}

func (mmap *InMemoryFileMap) GetFileHash(fName string) (uint64, error) {
	h, ok := mmap.FileMap[fName]
	if !ok {
		return 0, fmt.Errorf("does not exist : [%s]", fName)
	}
	return h, nil
}
