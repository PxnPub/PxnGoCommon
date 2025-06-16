package uid;

import(
	Fmt    "fmt"
	IO     "io"
	BufIO  "bufio"
	Sync   "sync"
	Time   "time"
	Errors "errors"
	Rand   "crypto/rand"
);



const MaxTimestamp int64 = 0x8FFFFFF_FFFFFFFF;
const MaxCounter   uint8 = 0x3F;
const MaxID        uint8 = 0x3;



type Generator struct {
	mut       Sync.Mutex
	id        uint8
	timestamp int64
	rng       IO.Reader
	count     uint8
}



func New(id int) *Generator {
	return &Generator{
		id:  uint8(id % 4),
		rng: BufIO.NewReaderSize(Rand.Reader, 16),
	};
}



func (gen *Generator) Next() (UID64, error) {
	gen.mut.Lock();
	defer gen.mut.Unlock();
	return gen.NextUnsafe();
}

func (gen *Generator) NextUnsafe() (UID64, error) {
//TODO: could increment timestamp by one
//      would need to do some trickery with the timestamp
	now := Time.Now().UnixMilli();
	if now == gen.timestamp {
		if gen.count == MaxCounter {
			gen.timestamp = now;
			return 0, Errors.New("Gen max ratio exceeded");
		}
		gen.count++;
	} else {
		gen.count = 0;
	}
	gen.timestamp = now;
	b := [1]byte{};
	if _, err := gen.rng.Read(b[:]); err != nil {
		return 0, err;
	}
	rnd := b[0];
	return NewUID64(gen.id, now, rnd, gen.count);
}

func NewUID64(id uint8, time int64, rnd byte, count uint8) (UID64, error) {
	if id > MaxID          { return 0, Fmt.Errorf("ID can't exceed %d", MaxID); }
	if time > MaxTimestamp { return 0, Fmt.Errorf("Timestamp can't exceed 0x%X",MaxTimestamp); }
	if count > MaxCounter  { return 0, Fmt.Errorf("Counter can't exceed 0x%X", MaxCounter); }
	return UID64(
		(uint64(time) << 16) +
		(uint64(rnd)  <<  8) +
		(uint64(id)   <<  6) +
		(uint64(count)),
	), nil;
}
