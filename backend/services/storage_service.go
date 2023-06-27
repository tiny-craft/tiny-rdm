package services

import "sync"

type storageService struct {
}

var storage *storageService
var onceStorage sync.Once

func Storage() *storageService {
	if storage == nil {
		onceStorage.Do(func() {
			storage = &storageService{}
		})
	}
	return storage
}
