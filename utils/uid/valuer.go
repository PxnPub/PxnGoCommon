package uid;

import(
	Fmt    "fmt"
	Driver "database/sql/driver"
);



func (uid UID64) Value() (Driver.Value, error) {
	return uid.ToInt(), nil;
}

func (uid *UID64) Scan(src interface{}) error {
	switch src := src.(type) {
		case nil: return nil;
		case uint64:
			val, err := FromInt(src);
			if err != nil { return err; }
			*uid = val;
		case string:
			val, err := Parse(src);
			if err != nil { return err; }
			*uid = val;
		default: return Fmt.Errorf("Unable to scan type %T into UID64", src);
	}
	return nil;
}



type UID64Slice []UID64;

func (arr UID64Slice) Len()              int  { return len(arr); }
func (arr UID64Slice) Less(x int, y int) bool { return arr[x] < arr[y]; }
func (arr UID64Slice) Swap(x int, y int)      { arr[x], arr[y] = arr[y], arr[x]; }
