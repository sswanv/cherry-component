package redis

import "github.com/redis/go-redis/v9"

type Invoker interface {
	GetRds(id string) (*redis.Client, bool)
}

func (c *Component) GetRds(id string) (*redis.Client, bool) {
	rds, ok := c.rdsMap[id]
	return rds, ok
}
