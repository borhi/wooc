package repositories

import (
	"errors"
	"sync"
	"wooc/models"
)

type IpMemoryRepository struct {
	IpModels []*models.IpModel
	sync.Mutex
}

func NewIpMemoryRepository() *IpMemoryRepository {
	return &IpMemoryRepository{}
}

func (r *IpMemoryRepository) Create(ipModel models.IpModel) *models.IpModel {
	r.Lock()
	defer r.Unlock()

	r.IpModels = append(r.IpModels, &ipModel)

	return &ipModel
}

func (r *IpMemoryRepository) FindAll(page uint64, pageSize uint64) ([]*models.IpModel, error) {
	if page == 0 {
		return nil, errors.New("page must be more than 0")
	}
	page -= 1

	r.Lock()
	defer r.Unlock()

	lens := uint64(len(r.IpModels))

	start := page * pageSize
	if start > lens {
		start = lens
	}

	end := start + pageSize

	if end > lens {
		end = lens
	}

	return r.IpModels[start:end], nil
}
