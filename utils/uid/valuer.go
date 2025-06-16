package uid;

import(
	Fmt    "fmt"
	Driver "database/sql/driver"
);



func (uid UID64) Value() (Driver.Value, error) {
	return uid.ToInt(), nil;
}

func (uid *UID64) Scan(src interface{}) error {
	switch t := src.(type) {
	case nil: return nil;
	case uint64:
		val, err := FromInt(t);
		if err != nil { return err; }
		*uid = val;
		break;
	case string:
		val, err := Parse(t);
		if err != nil { return err; }
		*uid = val;
		break;
	default: return Fmt.Errorf("Unable to scan type %T into UID64", t); break;
	}
	return nil;
}



type UID64Slice []UID64;

func (arr UID64Slice) Len()              int  { return len(arr); }
func (arr UID64Slice) Less(x int, y int) bool { return arr[x] < arr[y]; }
func (arr UID64Slice) Swap(x int, y int)      { arr[x], arr[y] = arr[y], arr[x]; }
