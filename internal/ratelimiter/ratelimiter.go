package ratelimiter

import "time"

// test using autocannon, ex: 4000 connectionRate / num
// npx autocannon -r 4000 --renderStatusCodes localhost:8080/v1/health

type Limiter interface {
	Allow(ip string) (bool, time.Duration)
}

type Config struct {
	RequestsPerTimeFrame int
	TimeFrame            time.Duration
	Enabled              bool
}
