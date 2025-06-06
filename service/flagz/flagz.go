package flagz;

import(
	Fmt  "fmt"
	Flag "flag"
);



func String(name string, token string, defval string) *string {
	Flag.String(name, defval, Fmt.Sprintf("--%s "+token, defval));
}

func Int(name string, token string, defval int) *int {
	Flag.Int(name, defval, Fmt.Sprintf("--%s "+token, defval));
}
