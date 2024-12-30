package configs

import "github.com/EmmanuelStan12/code-fusion/internal/common/utils"

type DockerConfig struct {
	Timeouts     []Timeout     `json:"timeouts"`
	Languages    []Language    `json:"languages"`
	MemoryLimits []MemoryLimit `json:"memoryLimits"`
}

type Timeout int

type Language string

type MemoryLimit int

const (
	TSec30  = 30_000
	TSec60  = 60_000
	TSec120 = 120_000

	LanguageJavaScript = "JavaScript"
	LanguageTypeScript = "TypeScript"

	MLMegaByte8  = 8
	MLMegaByte16 = 16
	MLMegaByte32 = 32
	MLMegaByte64 = 64
)

func NewDockerConfig() *DockerConfig {
	return &DockerConfig{
		Timeouts:     []Timeout{TSec30, TSec60, TSec120},
		Languages:    []Language{LanguageJavaScript, LanguageTypeScript},
		MemoryLimits: []MemoryLimit{MLMegaByte8, MLMegaByte16, MLMegaByte32, MLMegaByte64},
	}
}

func (d *DockerConfig) IsValidTimeout(timeout Timeout) bool {
	return utils.Contains[Timeout](d.Timeouts, timeout)
}

func (d *DockerConfig) IsValidMemoryLimit(memoryLimit MemoryLimit) bool {
	return utils.Contains[MemoryLimit](d.MemoryLimits, memoryLimit)
}

func (d *DockerConfig) IsValidLanguage(language Language) bool {
	return utils.Contains[Language](d.Languages, language)
}
