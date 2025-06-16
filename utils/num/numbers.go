package num;

import(
	Strings "strings"
	StrConv "strconv"
	Errors  "errors"
);



func ToBase36(val uint64) string {
	result := make([]byte, 13);
	copy(result, "0000000000000");
	str := Strings.ToUpper(StrConv.FormatUint(uint64(val), 36));
	size := len(str);
	copy(result[13-size:], str);
	return string(result);
}

func FromBase36(str string) (uint64, error) {
	if len(str) != 13 { return 0, Errors.New("Invalid UID value"); }
	return StrConv.ParseUint(str, 36, 64);
}
