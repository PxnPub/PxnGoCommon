package num;

import(
	Fmt     "fmt"
	Errors  "errors"
	Strings "strings"
	StrConv "strconv"
);



func ParseByteSize(size string) (int64, error) {
	siz := Strings.ToUpper(Strings.TrimSpace(size));
	var factor int64 = 1;
	if Strings.HasSuffix(siz, "B") { siz = Strings.TrimSuffix(siz, "B"); }
	switch {
	case Strings.HasSuffix(siz, "T"): siz = Strings.TrimSuffix(siz, "T"); factor = 1 << 40; break;
	case Strings.HasSuffix(siz, "G"): siz = Strings.TrimSuffix(siz, "G"); factor = 1 << 30; break;
	case Strings.HasSuffix(siz, "M"): siz = Strings.TrimSuffix(siz, "M"); factor = 1 << 20; break;
	case Strings.HasSuffix(siz, "K"): siz = Strings.TrimSuffix(siz, "K"); factor = 1 << 10; break;
	default: break;
	}
	value, err := StrConv.ParseInt(Strings.TrimSpace(siz), 10, 64);
	if err != nil { return 0, Errors.New("Invalid size format"); }
	return (value * factor), nil;
}



func FormatByteSize(size int64) (int64, string) {
	if size > 1<<40 { return (size / 1<<40), "T"; }
	if size > 1<<30 { return (size / 1<<30), "G"; }
	if size > 1<<20 { return (size / 1<<20), "M"; }
	if size > 1<<10 { return (size / 1<<10), "K"; }
	return size, "";
}

func FormatByteSizeString(size int64) string {
	value, unit := FormatByteSize(size);
	return Fmt.Sprintf("%0.1f%s", value, unit);
}
