package utils;

import(
	Time "time"
);



func SleepC() {
	sleep, err := Time.ParseDuration("100ms");
	if err == nil {
		Time.Sleep(sleep);
	}
}
