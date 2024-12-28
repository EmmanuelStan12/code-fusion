package service

import (
	"github.com/EmmanuelStan12/code-fusion/internal/common/utils"
	"github.com/EmmanuelStan12/code-fusion/internal/db"
)

type BaseService struct {
	Manager *db.PersistenceManager
	Jwt     utils.JwtUtils
}
