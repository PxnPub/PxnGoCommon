package tokenbucket;

import(
	Fmt      "fmt"
	Time     "time"
	Sync     "sync"
	UtilsNet "github.com/PxnPub/PxnGoCommon/utils/net"
);



type TokenBucket struct {
	Tokens    int16
	Hits      int64
	BlockHits int64
}

type RateLimiter struct {
	Mut        Sync.Mutex
	Buckets    map[UtilsNet.TupleIP]*TokenBucket
	Interval   Time.Duration
	TokenStart int16
	HitCost    int16
}



func New(interval Time.Duration, token_start int16, tokens_per_hit int16) *RateLimiter {
	return &RateLimiter{
		Interval:   interval,
		TokenStart: tokens_start,
		HitCost:    tokens_per_hit,
	};
}

func (limiter *RateLimiter) StartTicker() {
	go func() {
		ticker := Time.NewTicker(limiter.Interval);
		defer ticker.Stop();
		for { select {
			case <-ticker.C: limiter.Tick();
		}}
	}();
}



func (limiter *RateLimiter) Tick() {
	if len(limiter.Buckets) == 0 { return; }
	limiter.Mut.Lock();
	defer limiter.Mut.Unlock();
	for key, bucket := range limiter.Buckets {
		// add token to bucket
		bucket.Tokens--;
		// full bucket
		if bucket.Tokens <= 0 {
			delete(limiter.Buckets, key);
			continue;
		}
	}
}



func (limiter *RateLimiter) CheckStr(address string) (bool, error) {
	ip, err := UtilsNet.StringToIntIP(address);
	if err != nil { return true, err; }
	return limiter.CheckInt(ip_h, ip_l), nil;
}

func (limiter *RateLimiter) CheckIntIP(ip *UtilsNet.TupleIP) bool {
	limiter.Mut.Lock();
	defer limiter.Mut.Unlock();
	var bucket *TokenBucket = limiter.Buckets[ip];
	if bucket == nil {
		bucket = &TokenBucket{
			Tokens: 0,
		};
		limiter.Buckets[ip] = bucket;
	}
	if bucket.Tokens >= limiter.MaxTokens {
		bucket.CountBlocked++;
//TODO: remove this
if bucket.CountBlocked > 0 && bucket.CountBlocked % 100 == 0 {
Fmt.Printf("CAP %d\n", bucket.CountBlocked);
}
		return true;
	}
	bucket.Tokens += limiter.TokensPerCall;
	return false;
}
