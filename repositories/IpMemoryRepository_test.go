package repositories

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
	"wooc/models"
)

func TestIpMemoryRepository_Create(t *testing.T) {
	r := NewIpMemoryRepository()
	r.Create(models.IpModel{
		IpAddress: "8.8.8.8",
		ASN:       1234,
		Domains:   []string{"alpha.com", "beta.com"},
	})

	assert.Equal(t, &IpMemoryRepository{IpModels: []*models.IpModel{&models.IpModel{
		IpAddress: "8.8.8.8",
		ASN:       1234,
		Domains:   []string{"alpha.com", "beta.com"},
	}}}, r)
}

func TestIpMemoryRepository_FindAll(t *testing.T) {
	r := NewIpMemoryRepository()
	r.Create(models.IpModel{
		IpAddress: "8.8.8.8",
		ASN:       1,
		Domains:   []string{"alpha.com", "beta.com"},
	})
	r.Create(models.IpModel{
		IpAddress: "9.9.9.9",
		ASN:       2,
		Domains:   []string{"alpha.com", "beta.com"},
	})
	r.Create(models.IpModel{
		IpAddress: "7.7.7.7",
		ASN:       3,
		Domains:   []string{"alpha.com", "beta.com"},
	})
	r.Create(models.IpModel{
		IpAddress: "6.6.6.6",
		ASN:       4,
		Domains:   []string{"alpha.com", "beta.com"},
	})
	r.Create(models.IpModel{
		IpAddress: "0.0.0.0",
		ASN:       5,
		Domains:   []string{"alpha.com", "beta.com"},
	})

	result, _ := r.FindAll(1, 2)

	assert.Equal(t, []*models.IpModel{&models.IpModel{
		IpAddress: "8.8.8.8",
		ASN:       1,
		Domains:   []string{"alpha.com", "beta.com"},
	}, &models.IpModel{
		IpAddress: "9.9.9.9",
		ASN:       2,
		Domains:   []string{"alpha.com", "beta.com"},
	}}, result)
}

func TestIpMemoryRepository_FindAllStartOverLens(t *testing.T) {
	r := NewIpMemoryRepository()
	r.Create(models.IpModel{
		IpAddress: "8.8.8.8",
		ASN:       1,
		Domains:   []string{"alpha.com", "beta.com"},
	})
	r.Create(models.IpModel{
		IpAddress: "9.9.9.9",
		ASN:       2,
		Domains:   []string{"alpha.com", "beta.com"},
	})

	result, _ := r.FindAll(2, 3)

	assert.Equal(t, []*models.IpModel{}, result)
}

func TestIpMemoryRepository_FindAllZeroPage(t *testing.T) {
	r := NewIpMemoryRepository()
	_, err := r.FindAll(0, 2)
	assert.Equal(t, errors.New("page must be more than 0"), err)
}
