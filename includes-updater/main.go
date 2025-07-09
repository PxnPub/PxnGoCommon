package main;

import(
	OS      "os"
	Fmt     "fmt"
	Flag    "flag"
	Strings "strings"
	HTTP    "net/http"
	IOUtils "io/ioutil"
	JSON    "encoding/json"
);



const DefaultOutFile = "includes.json";
const URL_Resolver = "https://data.jsdelivr.com/v1/packages/gh/%s/resolved?specifier=latest";



func main() {
	var outfile string;
	Flag.StringVar(&outfile, "out-file", "", "--out-file "+DefaultOutFile);
	Flag.Parse();
	if outfile == "" { outfile = DefaultOutFile; }
	print("\n");
	Fmt.Printf("Updating includes..\n");
	versions := make(map[string]string);
	includes := []string{
		"twbs/bootstrap",
		"twbs/icons",
		"floating-ui/floating-ui",
		"DataTables/Dist-DataTables-Bootstrap5",
		"DataTables/Dist-DataTables-Scroller-Bootstrap5",
		"jquery/jquery",
		"apache/echarts",
	};
	for _, key := range includes {
		FindLatestVersion(versions, key); }
	result, err := JSON.MarshalIndent(versions, "", "\t");
	if err != nil { panic(err); }
	Fmt.Printf("\n%s\n", result);
	print("\n");
	out, err := OS.Create(outfile);
	if err != nil { panic(err); }
	defer out.Close();
	if _, err := out.Write(result);     err != nil { panic(err); }
	if _, err := out.WriteString("\n"); err != nil { panic(err); }
	Fmt.Printf("Wrote file: %s\n", outfile);
	print("\n");
}



func FindLatestVersion(versions map[string]string, key string) string {
	url := Fmt.Sprintf(URL_Resolver, key);
	resp, err := HTTP.Get(url);
	if err != nil { panic(err); }
	defer resp.Body.Close();
	if resp.StatusCode != HTTP.StatusOK { panic(Fmt.Errorf(
		"Failed to fetch latest version for %s, got %d", key, resp.StatusCode)); }
	body, err := IOUtils.ReadAll(resp.Body);
	if err != nil { panic(err); }
	var json map[string]interface{};
	if err := JSON.Unmarshal(body, &json); err != nil { panic(err); }
	vers, ok := json["version"];
	if !ok { panic(Fmt.Errorf(
		"Failed to fetch latest version for %s, field not found", key)); }
	version := vers.(string);
	padsize := 24 - len(key);
	padding := "";
	if padsize > 0 { padding = Strings.Repeat(" ", padsize); }
	Fmt.Printf("  %s%s %s\n", key, padding, version);
	versions[key] = version;
	return version;
}
