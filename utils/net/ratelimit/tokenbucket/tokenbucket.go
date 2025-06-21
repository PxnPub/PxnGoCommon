package tokenbucket;

import(
	Fmt      "fmt"
	Net      "net"
	Time     "time"
	Sync     "sync"
	TupleIP  "github.com/PxnPub/PxnGoCommon/utils/net/tupleip"
);



type Limiter struct {
	MutTick   Sync.Mutex
	Buckets   map[TupleIP.IP]*Bucket
	Interval  Time.Duration
	HitCost   int32
	TokensMax int32
}

type Bucket struct {
	Tokens    int32
	Hits      int64
	BlockHits int64
}



func New(interval Time.Duration, tokens_per_hit int32,
		tokens_max int32) *Limiter {
	return &Limiter{
		Buckets:   make(map[TupleIP.IP]*Bucket),
		Interval:  interval,
		HitCost:   tokens_per_hit,
		TokensMax: tokens_max,
	};
}

func (limiter *Limiter) StartTicker() {
	go func() {
		ticker := Time.NewTicker(limiter.Interval);
		defer ticker.Stop();
		for { select { case <-ticker.C: limiter.Tick(); }}
	}();
}



func (limiter *Limiter) Tick() {
	if len(limiter.Buckets) == 0 { return; }
	limiter.MutTick.Lock();
	defer limiter.MutTick.Unlock();
	for ip, bucket := range limiter.Buckets {
		// add token to bucket
		bucket.Tokens--;
//Fmt.Printf("  Tok: %s %d\n", ip.String(), bucket.Tokens);
		// full bucket
		if bucket.Tokens <= 0 {
			delete(limiter.Buckets, ip);
			continue;
		}
	}
}



func (limiter *Limiter) CheckNetAddr(addr Net.Addr) (bool, error) {
	host, _, err := Net.SplitHostPort(addr.String());
	if err != nil { return true, err; }
	return limiter.CheckStr(host);
}

func (limiter *Limiter) CheckStr(address string) (bool, error) {
	ip, err := TupleIP.NewFromString(address);
	if err != nil { return true, err; }
	return limiter.CheckTupleIP(ip), nil;
}

func (limiter *Limiter) CheckTupleIP(ip *TupleIP.IP) bool {
	limiter.MutTick.Lock();
	defer limiter.MutTick.Unlock();
	var bucket *Bucket = limiter.Buckets[*ip];
	// new bucket
	if bucket == nil {
		bucket = &Bucket{
			Tokens: 0,
		};
		limiter.Buckets[*ip] = bucket;
	}
	if bucket.Tokens >= limiter.TokensMax {
		bucket.BlockHits++;
		if bucket.BlockHits > 0 && bucket.BlockHits % 100 == 0 {
			Fmt.Printf("Rate Limited %d times!  %s\n",
				bucket.BlockHits, ip.ToStringRaw()); }
		return true;
	}
	bucket.Tokens += limiter.HitCost;
	return false;
}
