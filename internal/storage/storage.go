package storage

import (
	"fmt"
	"log"
	"sync"

	"github.com/nxvrlxv/hosts_monitor/internal/checker"
)

type Storage struct {
	Data map[int]checker.Result
	Mu   *sync.RWMutex
}

func CreateNewStorage() *Storage {
	var mu sync.RWMutex
	return &Storage{
		Data: make(map[int]checker.Result),
		Mu:   &mu,
	}
}

func (s *Storage) FillData(data []checker.Result) error {
	if len(data) == 0 {
		return fmt.Errorf("nothing to fill: empty data")
	}

	s.Mu.Lock()
	defer s.Mu.Unlock()
	for _, res := range data {
		if v, ok := s.Data[res.ID]; ok {
			if v.IsUp != res.IsUp {
				log.Printf("WARN: [%d] site status changed from %v to %v\n", res.ID, v.IsUp, res.IsUp)
			}
		}

		s.Data[res.ID] = res
	}

	return nil
}

func (s *Storage) GetResultById(id int) (*checker.Result, error) {
	s.Mu.RLock()
	defer s.Mu.RUnlock()

	if res, ok := s.Data[id]; !ok {
		return nil, fmt.Errorf("Not Found")
	} else {
		return &res, nil
	}
}

func (s *Storage) GetAllData() []checker.Result {
	s.Mu.RLock()
	defer s.Mu.RUnlock()
	if len(s.Data) == 0 {
		return []checker.Result{}
	}
	var res []checker.Result

	for _, v := range s.Data {
		res = append(res, v)
	}

	return res
}
