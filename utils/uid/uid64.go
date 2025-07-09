package uid;

import(
	PxnNum "github.com/PoiXson/pxnGoCommon/utils/num"
);



// 6 bytes unix timestamp
// 1 byte entropy from /dev/urandom
// 1 byte counter & id
type UID64 uint64;



// to/from string
func Parse(str string) (UID64, error) {
	uid, err := PxnNum.FromBase36(str);
	return UID64(uid), err;
}

func (uid UID64) ToString() string {
	return PxnNum.ToBase36(uint64(uid));
}



// to/from int
func FromInt(val uint64) (UID64, error) {
	return UID64(val), nil;
}

func (uid UID64) ToInt() uint64 {
	return uint64(uid);
}



// get parts
func (uid UID64) GetID() uint8 {
	return ((uint8(uid) & 0b1100_0000) >> 6);
}

func (uid UID64) GetTimestamp() int64 {
	return int64((uint64(uid) & 0xffffffff_ffff0000) >> 16);
}

func (uid UID64) GetRND() uint8 {
	return uint8((uint16(uid) & 0xff00) >> 8);
}

func (uid UID64) GetCounter() uint8 {
	return (uint8(uid) & MaxCounter);
}
