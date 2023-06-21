package manager

import (
	"github.com/ip2location/ip2location-go/v9"
	"ssp/lib/mongo"
)

//go:generate go run github.com/golang/mock/mockgen -destination=./mock_manager.go -package=manager . IDSPPubModel
type Manager struct {
	Client   *mongo.Client
	IPClient *ip2location.DB
}

func (m *Manager) GetMongoClient() *mongo.Client {
	return m.Client
}

func (m *Manager) GetIPClient() *ip2location.DB {
	return m.IPClient
}
