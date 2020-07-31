package interfaces

import "wooc/models"

type RepositoryInterface interface {
	Create(ipModel models.IpModel) *models.IpModel
	FindAll(page uint64, pageSize uint64) ([]*models.IpModel, error)
}
