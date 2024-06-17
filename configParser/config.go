package configParser

type WindowConfig struct {
	Width  int    `json:"width"`
	Height int    `json:"height"`
	Title  string `json:"title"`
}

type TextureInfo struct {
	Name string `json:"name"`
	Path string `json:"path"`
}

type Config struct {
	Window WindowConfig `json:"window"`
}

type IConfigParser interface {
	ParseConfig(data []byte) (*Config, error)
}
