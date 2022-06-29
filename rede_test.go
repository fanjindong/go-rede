package go_rede

import (
	"context"
	"os"
	"reflect"
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
		Addr:       os.Getenv("REDIS_ADDR"),
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
		input   input
		want    float64
		wantErr bool
	}{
		{input: input{member: "a", ttl: 1 * time.Second}, want: now + 1},
		{input: input{member: "b", ttl: 2 * time.Second}, want: now + 2},
		{input: input{member: "c", ttl: 3 * time.Second}, want: now + 3},
		{input: input{member: "d", ttl: 500 * time.Millisecond}, want: now + 0.5},
		{input: input{member: "e", ttl: 1 * time.Minute}, want: now + 60},
	}

	for _, ts := range tests {
		_, err := rede.Push(ctx, ts.input.member, ts.input.ttl)
		if (err != nil) != ts.wantErr {
			t.Errorf("Push() error = %v, wantErr %v", err, ts.wantErr)
			return
		}
		got, err := rede.ZScore(ctx, rede.Namespaces, ts.input.member).Result()
		if (err != nil) != ts.wantErr {
			t.Errorf("ZScore() error = %v, wantErr %v", err, ts.wantErr)
			return
		}
		if !reflect.DeepEqual(int64(ts.want*1e6), int64(got*1e6)) {
			t.Errorf("ZScore() got = %v, want %v", int64(got*1e6), int64(ts.want*1e6))
		}
	}
}

func TestClient_Look(t *testing.T) {
	type input struct {
		member string
		ttl    time.Duration
	}
	tests := []struct {
		input   input
		want    float64
		wantErr bool
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
		if (err != nil) != ts.wantErr {
			t.Errorf("Look() error = %v, wantErr %v", err, ts.wantErr)
			return
		}
		if !reflect.DeepEqual(got, ts.want) {
			t.Errorf("ZScore() got = %v, want %v", got, ts.want)
		}
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
		input   []input
		sleep   time.Duration
		want    []string
		wantErr bool
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
		_, _ = rede.Push(ctx, ts.member, ts.ttl)
	}

	time.Sleep(tests.sleep)

	cur := rede.Poll(ctx)
	i := 0
	for cur.Next() {
		got, err := cur.Get()
		t.Log(got, err)
		if (err != nil) != tests.wantErr {
			t.Errorf("Poll() error = %v, wantErr %v", err, tests.wantErr)
			return
		}
		if !reflect.DeepEqual(tests.want[i], got) {
			t.Errorf("Next() got = %v, want %v", got, tests.want[i])
		}
		i++
	}
}
