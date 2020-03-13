package go_rede

import (
	"github.com/go-redis/redis/v7"
	"math"
	"time"
)

// Client is a Redis client representing a pool of zero or more
// underlying connections. It's safe for concurrent use by multiple
// goroutines.
type Client struct {
	Namespaces string
	*redis.Client
}

func NewClient(opt *Options) *Client {
	redisOpt := &redis.Options{
		Network:            opt.Network,
		Addr:               opt.Addr,
		Password:           opt.Password,
		DB:                 opt.DB,
		MaxRetries:         opt.MaxRetries,
		MinRetryBackoff:    opt.MinRetryBackoff,
		MaxRetryBackoff:    opt.MaxRetryBackoff,
		DialTimeout:        opt.DialTimeout,
		ReadTimeout:        opt.ReadTimeout,
		WriteTimeout:       opt.WriteTimeout,
		PoolSize:           opt.PoolSize,
		MinIdleConns:       opt.MinIdleConns,
		MaxConnAge:         opt.MaxConnAge,
		PoolTimeout:        opt.PoolTimeout,
		IdleTimeout:        opt.IdleTimeout,
		IdleCheckFrequency: opt.IdleCheckFrequency,
		TLSConfig:          opt.TLSConfig,
		Limiter:            opt.Limiter,
	}
	return &Client{
		Namespaces: opt.Namespaces,
		Client:     redis.NewClient(redisOpt),
	}
}

//Push an Member into the Rede for ttl.Seconds() seconds
func (c *Client) Push(member interface{}, ttl time.Duration) (int64, error) {
	z := &redis.Z{
		Score:  float64(time.Now().Unix()) + ttl.Seconds(),
		Member: member,
	}
	result := c.ZAdd(c.Namespaces, z)
	return result.Result()
}

//Pull the members, remove it from the rede before it expires.
func (c *Client) Pull(members ...interface{}) (int64, error) {
	result := c.ZRem(c.Namespaces, members...)
	return result.Result()
}

//Show the ttl corresponding with element and without removing it from the rede.
func (c *Client) Look(member string) (float64, error) {
	result, err := c.ZScore(c.Namespaces, member).Result()
	if err != nil {
		return 0, err
	}
	return math.Max(result-float64(time.Now().Unix()), 0), nil
}

////Show the time left (in seconds) until the next element will expire.
//func (c *Client) Ttn() (float64, error) {
//	result, err := c.ZRangeWithScores(c.Namespaces, 0, 0).Result()
//	if len(result) == 0 {
//		return -1, nil
//	}
//	return math.Max(0, result[0].Score-float64(time.Now().Unix())), err
//}

//Pull and return all the expired members in rede.
func (c *Client) Poll() ([]interface{}, error) {
	result := make([]interface{}, 0)
	for {
		zSlice, err := c.ZPopMin(c.Namespaces, 1).Result()
		if err != nil {
			return result, err
		}
		if len(zSlice) == 0 {
			return result, nil
		}

		if zSlice[0].Score > float64(time.Now().Unix()) {
			c.ZAdd(c.Namespaces, &zSlice[0])
			return result, nil
		}
		result = append(result, zSlice[0].Member)
	}
}
