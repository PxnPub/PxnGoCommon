package flagz;

import(
	Fmt  "fmt"
	Flag "flag"
);



func usage(name string, token string, defval interface{}) string {
	return Fmt.Sprintf("--%s "+token, name, defval);
}



// string
func String(value *string, name string, defval string) {
	Flag.StringVar(value, name, defval, usage(name, "%s", defval));
}



// int
func Int(value *int, name string, defval int) {
	Flag.IntVar(value, name, defval, usage(name, "%d", defval));
}
func UInt(value *uint, name string, defval uint) {
	Flag.UintVar(value, name, defval, usage(name, "%d", defval));
}



// int64
func Int64(value *int64, name string, defval int64) {
	Flag.Int64Var(value, name, defval, usage(name, "%d", defval));
}
func UInt64(value *uint64, name string, defval uint64) {
	Flag.Uint64Var(value, name, defval, usage(name, "%d", defval));
}



// float
func Float(value *float64, name string, defval float64) {
	Flag.Float64Var(value, name, defval, usage(name, "%f", defval));
}



// bool
func Bool(value *bool, name string) {
	Flag.BoolVar(value, name, false, Fmt.Sprintf(name, "--%s", name));
}
