package flagz;

import(
	Fmt "fmt"
);



func String(ref *string, name string, token string, defval string) {
	Flag.StringVar(ref, name, defval, Fmt.Sprintf("--%s "+token, defval));
}
