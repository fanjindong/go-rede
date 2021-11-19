package go_rede

import (
	"context"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
	"time"
)

var (
	ctx  = context.Background()
	rede *Client
)

func TestMain(m *testing.M) {
	rede = NewClient(&Options{
		Namespaces: "rede",
		Addr:       "ml5pub.tsht3.mc.ops:6379",
		Password:   "",
		DB:         0,
	})
	exitCode := m.Run()
	os.Exit(exitCode)
}

func TestClient_Push(t *testing.T) {
	now := float64(time.Now().Unix())
	type input struct {
		member string
		ttl    time.Duration
	}
	tests := []struct {
		input input
		want  float64
	}{
		{input: input{member: "a", ttl: 1 * time.Second}, want: now + 1},
		{input: input{member: "b", ttl: 2 * time.Second}, want: now + 2},
		{input: input{member: "c", ttl: 3 * time.Second}, want: now + 3},
		{input: input{member: "d", ttl: 500 * time.Millisecond}, want: now + 0.5},
		{input: input{member: "e", ttl: 1 * time.Minute}, want: now + 60},
	}

	for _, ts := range tests {
		_, err := rede.Push(ctx, ts.input.member, ts.input.ttl)
		assert.NoError(t, err)
		got, err := rede.ZScore(ctx, rede.Namespaces, ts.input.member).Result()
		assert.NoError(t, err)
		assert.Equal(t, int64(ts.want*1e6), int64(got*1e6))
	}
}

func TestClient_Look(t *testing.T) {
	type input struct {
		member string
		ttl    time.Duration
	}
	tests := []struct {
		input input
		want  float64
	}{
		{input: input{member: "a", ttl: 1 * time.Second}, want: 1},
		{input: input{member: "b", ttl: 2 * time.Second}, want: 2},
		{input: input{member: "c", ttl: 3 * time.Second}, want: 3},
		{input: input{member: "d", ttl: 500 * time.Millisecond}, want: 0.5},
		{input: input{member: "e", ttl: 1 * time.Minute}, want: 60},
	}
	for _, ts := range tests {
		_, _ = rede.Push(ctx, ts.input.member, ts.input.ttl)
		got, err := rede.Look(ctx, ts.input.member)
		assert.NoError(t, err)
		assert.Equal(t, got, ts.want)
	}
}

func TestClient_Ttn(t *testing.T) {
	r, err := rede.Ttn(ctx)
	t.Log(r, err)
}

func TestClient_Poll(t *testing.T) {
	type input struct {
		member string
		ttl    time.Duration
	}
	tests := struct {
		input []input
		sleep time.Duration
		want  []string
	}{
		input: []input{
			{member: "a", ttl: 1 * time.Second},
			{member: "b", ttl: 2 * time.Second},
			{member: "c", ttl: 3 * time.Second},
			{member: "d", ttl: 4 * time.Second},
		},
		sleep: 2 * time.Second,
		want:  []string{"a", "b"},
	}
	rede.Del(ctx, rede.Namespaces)
	for _, ts := range tests.input {
		_, err := rede.Push(ctx, ts.member, ts.ttl)
		assert.NoError(t, err)
	}

	time.Sleep(tests.sleep)

	cur := rede.Poll(ctx)
	i := 0
	for cur.Next() {
		got, err := cur.Get()
		t.Log(got, err)
		assert.NoError(t, err)
		assert.Equal(t, tests.want[i], got)
		i++
	}
}
