package config

type RedisConfig struct {
	host string
	port int
}

func redisConfig() *RedisConfig {
	return &RedisConfig{
		host: getString(redisHost, true),
		port: getInt(redisPort, true),
	}
}

func (redis *RedisConfig) Host() string {
	return redis.host
}

func (redis *RedisConfig) Port() int {
	return redis.port
}
