package config

import "os"

const (
    WindowWidth  = 160
    WindowHeight = 40
    ApiURL       = "https://api.anthropic.com/v1/messages"
)

func GetApiKey() string {
    return os.Getenv("ANTHROPIC_API_KEY")
}
