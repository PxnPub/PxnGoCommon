package fs;

import(
	OS       "os"
	Errors   "errors"
	FilePath "path/filepath"
);



func IsFile(file string) bool {
	info, err := OS.Stat(file);
	if err != nil {
		if Errors.Is(err, OS.ErrNotExist) { return false; }
		panic(err);
	}
	return info.Mode().IsRegular();
}

func FindFile(file string, paths...string) string {
	for i := range paths {
		p := FilePath.Join(paths[i], file);
		if IsFile(p) { return p; }
	}
	return "";
}
