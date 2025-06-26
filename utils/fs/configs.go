package fs;

import(
	IOUtils "io/ioutil"
	YAML    "gopkg.in/yaml.v2"
);



func LoadConfig[T any](file string) (*T, error) {
	data, err := IOUtils.ReadFile(file);
	if err != nil { return nil, err; }
	var config T;
	if err := YAML.Unmarshal(data, &config); err != nil {
		return nil, err; }
	return &config, nil;
}
