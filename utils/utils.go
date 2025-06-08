package utils;

import(
	Time "time"
);



func SleepX() {
	sleep, err := Time.ParseDuration("10ms");
	if err == nil { Time.Sleep(sleep); }
}

func SleepC() {
	sleep, err := Time.ParseDuration("100ms");
	if err == nil { Time.Sleep(sleep); }
}

func SleepS() {
	sleep, err := Time.ParseDuration("1s");
	if err == nil { Time.Sleep(sleep); }
}
