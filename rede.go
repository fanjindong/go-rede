package go_rede

import (
	"context"
	"github.com/go-redis/redis/v8"
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

// Push an Member into the Rede for ttl.Seconds() seconds
func (c *Client) Push(ctx context.Context, member string, ttl time.Duration) (int64, error) {
	z := &redis.Z{
		Score:  float64(time.Now().Unix()) + ttl.Seconds(),
		Member: member,
	}
	result := c.ZAdd(ctx, c.Namespaces, z)
	return result.Result()
}

// Pull the members, remove it from the rede before it expires.
func (c *Client) Pull(ctx context.Context, members ...string) (int64, error) {
	items := make([]interface{}, len(members))
	for i, member := range members {
		items[i] = member
	}
	result := c.ZRem(ctx, c.Namespaces, items...)
	return result.Result()
}

// Look Show the ttl corresponding with element and without removing it from the rede.
func (c *Client) Look(ctx context.Context, member string) (float64, error) {
	result, err := c.ZScore(ctx, c.Namespaces, member).Result()
	if err != nil {
		return 0, err
	}
	return math.Max(result-float64(time.Now().Unix()), 0), nil
}

// Ttn Show the time left (in seconds) until the next element will expire.
func (c *Client) Ttn(ctx context.Context) (float64, error) {
	result, err := c.ZRangeWithScores(ctx, c.Namespaces, 0, 0).Result()
	if len(result) == 0 {
		return -1, nil
	}
	return math.Max(0, result[0].Score-float64(time.Now().Unix())), err
}

//Poll return all the expired members in rede.
// cur := c.Poll(ctx)
// for cur.Next() {
//     member, err := cur.Get()
//     fmt.Println(member, err)
// }
func (c *Client) Poll(ctx context.Context) *pollCursor {
	return newPollCursor(ctx, c)
}

type pollCursor struct {
	ctx   context.Context
	c     *Client
	value string
	err   error
}

func newPollCursor(ctx context.Context, c *Client) *pollCursor {
	return &pollCursor{ctx: ctx, c: c}
}

//Next get the next element in the form of iteration
func (pc *pollCursor) Next() bool {
	if pc.err != nil {
		return false
	}
	zSlice, err := pc.c.ZPopMin(pc.ctx, pc.c.Namespaces, 1).Result()
	if err != nil {
		pc.value, pc.err = "", err
		return true
	}
	if len(zSlice) == 0 {
		return false
	}
	if zSlice[0].Score > float64(time.Now().Unix()) {
		pc.c.ZAdd(pc.ctx, pc.c.Namespaces, &zSlice[0])
		return false
	}
	pc.value = zSlice[0].Member.(string)
	return true
}

func (pc *pollCursor) Get() (string, error) {
	return pc.value, pc.err
}
