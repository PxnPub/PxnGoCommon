package uid;

import(
	Testing "testing"
	Assert  "github.com/stretchr/testify/assert"
);



func Test_UID64(t *Testing.T) {
	var tests = []struct {
		id    uint8
		time  int64
		rnd   byte
		count uint8
		expect UID64
	}{
		//ID      Time RND Cnt UID
		{ 0,          0, 0, 0, 0x0            },
		{ 0,          1, 0, 0, 0x10000        },
		{ 0, 1747322506, 0, 0, 0x6826068a0000 },
		{ 0, 1747322507, 0, 0, 0x6826068b0000 },
		{ 0, 1747322508, 0, 0, 0x6826068c0000 },
		{ 0, 1747322509, 0, 0, 0x6826068d0000 },
		{ 0, 1747322510, 0, 0, 0x6826068e0000 },
		{ 1, 1747322511, 0, 0, 0x6826068f0040 },
		{ 2, 1747322512, 0, 0, 0x682606900080 },
		{ 3, 1747322514, 0, 0, 0x6826069200c0 },
		{ 0, 1747322515, 1, 0, 0x682606930100 },
		{ 0, 1747322516, 2, 0, 0x682606940200 },
		{ 0, 1747322517, 3, 0, 0x682606950300 },
		{ 0, 1747322518, 4, 0, 0x682606960400 },
		{ 0, 1747322520, 0, 0, 0x682606980000 },
		{ 0, 1747322520, 0, 1, 0x682606980001 },
		{ 0, 1747322520, 0, 2, 0x682606980002 },
		{ 0, 1747322520, 0, 3, 0x682606980003 },
		{ 0, 1747322520, 0, 4, 0x682606980004 },
		{ 0, 1747322520, 0, 5, 0x682606980005 },
	};
	for _, test := range tests {
		result, err := NewUID64(test.id, test.time, test.rnd, test.count);
		Assert.Equal(t, test.expect, result);
		Assert.Equal(t, nil, err);
	}
	// test full stack
	gen := New(0);
	var last UID64 = 0;
	errs := 0;
	for i:=0; i<1000; i++ {
		uid, err := gen.Next();
		if err == nil {
			if i > 0 { Assert.NotEqual(t, last, uid); }
			last = uid;
		} else {
			Assert.EqualError(t, err, "Gen max ratio exceeded");
			errs++;
		}
	}
	Assert.NotEqual(t, 0, errs);
}
