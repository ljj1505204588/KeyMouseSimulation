package conf

// ConfigOption 配置选项
type ConfigOption func(*configuration)

// WithConfigPath 设置配置文件路径
func WithConfigPath(path string) ConfigOption {
	return func(c *configuration) {
		c.path = path
	}
}
