package num;

import(
	Fmt     "fmt"
	Errors  "errors"
	Strings "strings"
	StrConv "strconv"
);



func ParseByteSize(size string) (uint64, error) {
	siz := Strings.ToUpper(Strings.TrimSpace(size));
	var factor uint64 = 1;
	switch {
	case Strings.HasSuffix(siz, "TB"): siz = Strings.TrimSuffix(siz, "TB"); factor = 1 << 40; break;
	case Strings.HasSuffix(siz, "GB"): siz = Strings.TrimSuffix(siz, "GB"); factor = 1 << 30; break;
	case Strings.HasSuffix(siz, "MB"): siz = Strings.TrimSuffix(siz, "MB"); factor = 1 << 20; break;
	case Strings.HasSuffix(siz, "KB"): siz = Strings.TrimSuffix(siz, "KB"); factor = 1 << 10; break;
	case Strings.HasSuffix(siz,  "B"): siz = Strings.TrimSuffix(siz,  "B");                   break;
	default:                                                                                  break;
	}
	value, err := StrConv.ParseUint(Strings.TrimSpace(siz), 10, 64);
	if err != nil { return 0, Errors.New("Invalid size format"); }
	return (value * factor), nil;
}



func FormatByteSize(size int64) (uint64, byte) {
	if size > 1<<40 { return Fmt.Sprintf("%dT", size / 1<<40); }
	if size > 1<<30 { return Fmt.Sprintf("%dG", size / 1<<30); }
	if size > 1<<20 { return Fmt.Sprintf("%dM", size / 1<<20); }
	if size > 1<<10 { return Fmt.Sprintf("%dK", size / 1<<10); }
	return size, nil;
}

func FormatByteSizeString(size int64) string {
	value, unit := FormatByteSize(size);
	return Fmt.Sprintf("%0.1f%s", value, unit);
}
