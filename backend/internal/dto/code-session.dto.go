package dto

import "github.com/EmmanuelStan12/code-fusion/configs"

type CreateCodeSessionDTO struct {
	Title       string
	Language    configs.Language
	MemoryLimit configs.MemoryLimit
	Timeout     configs.Timeout
}
