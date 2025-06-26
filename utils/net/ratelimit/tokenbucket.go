package ratelimit;

import(
	Log    "log"
	Fmt    "fmt"
	Net    "net"
	Time   "time"
	Sync   "sync"
	PxnNet "github.com/PxnPub/PxnGoCommon/utils/net"
);



const DefaultInterval     = "1s";
const DefaultHitCost      =  3;
const DefaultTokensThresh = 35;
const DefaultTokensCap    = 50;



type TokBuckLim struct {
	MutBuckets   Sync.Mutex
	Buckets      map[PxnNet.TupIP]*Bucket
	Interval     Time.Duration
	HitCost      int32
	TokensThresh int32
	TokensCap    int32
}

type Bucket struct {
	Tokens int32
	Hits   int64
	Blocks int64
}



func NewTokenBucket() *TokBuckLim {
	default_interval, err := Time.ParseDuration(DefaultInterval);
	if err != nil { Log.Panic(err); }
	return &TokBuckLim{
		Buckets:      make(map[PxnNet.TupIP]*Bucket),
		Interval:     default_interval,
		HitCost:      DefaultHitCost,
		TokensThresh: DefaultTokensThresh,
		TokensCap:    DefaultTokensCap,
	};
}

func (tokbuck *TokBuckLim) Start() {
	go func() {
		ticker := Time.NewTicker(tokbuck.Interval);
		defer ticker.Stop();
		for { select { case <-ticker.C: tokbuck.Tick(); }}
	}();
}



func (tokbuck *TokBuckLim) Tick() {
	tokbuck.MutBuckets.Lock();
	defer tokbuck.MutBuckets.Unlock();
	if len(tokbuck.Buckets) == 0 { return; }
	for ip, bucket := range tokbuck.Buckets {
		// add token to bucket
		bucket.Tokens--;
		if bucket.Tokens > tokbuck.TokensCap {
			bucket.Tokens = tokbuck.TokensCap; }
Fmt.Printf("  Tok: %s %d\n", ip.String(), bucket.Tokens);
		// full bucket
		if bucket.Tokens <= 0 {
			delete(tokbuck.Buckets, ip);
			continue;
		}
	}
}



func (tokbuck *TokBuckLim) GetBucket(ip *PxnNet.TupIP) *Bucket {
	bucket, ok := tokbuck.Buckets[*ip];
	if !ok || bucket == nil {
		bucket = tokbuck.NewBucket();
		tokbuck.Buckets[*ip] = bucket;
	}
	return bucket;
}

func (tokbuck *TokBuckLim) NewBucket() *Bucket {
	return &Bucket{
		Tokens: 0,
		Hits:   0,
		Blocks: 0,
	};
}



func (tokbuck *TokBuckLim) CheckNetAddr(addr Net.Addr) (bool, error) {
	host, _, err := Net.SplitHostPort(addr.String());
	if err != nil { return true, err; }
	return tokbuck.CheckAddrStr(host);
}

func (tokbuck *TokBuckLim) CheckAddrStr(addr string) (bool, error) {
	ip_tup := PxnNet.ParseAddrStr(addr);
	if ip_tup == nil { return true, Fmt.Errorf("Invalid IP: %s", addr); }
	return tokbuck.CheckTupleIP(ip_tup), nil;
}

func (tokbuck *TokBuckLim) CheckTupleIP(ip *PxnNet.TupIP) bool {
	tokbuck.MutBuckets.Lock();
	defer tokbuck.MutBuckets.Unlock();
	bucket := tokbuck.GetBucket(ip);
	if bucket == nil { Log.Panic("Failed to get token bucket"); }
	bucket.Tokens += tokbuck.HitCost;
	if bucket.Tokens >= tokbuck.TokensThresh {
		if bucket.Tokens > tokbuck.TokensCap {
			bucket.Tokens = tokbuck.TokensCap; }
		bucket.Blocks++;
		if bucket.Blocks > 0 && bucket.Blocks % 100 == 0 {
			Fmt.Printf("Rate Limited %d times!  %s\n",
				bucket.Blocks, ip.String()); }
		return true;
	}
	return false;
}
