package service

import (
	"github.com/EmmanuelStan12/code-fusion/client"
	"github.com/EmmanuelStan12/code-fusion/internal/db"
)

type BaseService struct {
	Manager *db.PersistenceManager
	Jwt     client.JwtClient
}
