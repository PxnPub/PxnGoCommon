package san;

import(
	HTML   "html"
	Regexp "regexp"
);



func IsSafeAlphaLower(str string) bool { if str == "" { return true; }; return Regexp.MustCompile("^[a-z]+$"        ).MatchString(str); }
func IsSafeAlphaUpper(str string) bool { if str == "" { return true; }; return Regexp.MustCompile("^[A-Z]+$"        ).MatchString(str); }
func IsSafeAlpha     (str string) bool { if str == "" { return true; }; return Regexp.MustCompile("^[a-zA-Z]+$"     ).MatchString(str); }
func IsSafeAlphaNum  (str string) bool { if str == "" { return true; }; return Regexp.MustCompile("^[a-zA-Z0-9]+$"  ).MatchString(str); }
func IsSafeFilePath  (str string) bool { if str == "" { return true; }; return Regexp.MustCompile("^[a-zA-Z0-9./]+$").MatchString(str); }
func IsSafeDomain    (str string) bool { if str == "" { return true; }; return Regexp.MustCompile("^[a-zA-Z0-9.]+$" ).MatchString(str); }
func IsSafeDomainPort(str string) bool { if str == "" { return true; }; return Regexp.MustCompile("^[a-zA-Z0-9.:]+$").MatchString(str); }
func SafeHTML(str string) string { return HTML.EscapeString(str); }
