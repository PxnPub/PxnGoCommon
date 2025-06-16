package utils;

import(
	Time "time"
	Rand "math/rand"
);



func SleepX() {
	sleep, err := Time.ParseDuration("10ms");
	if err == nil { Time.Sleep(sleep); }
}

//TODO: use this
func SleepV() {
	sleep, err := Time.ParseDuration("50ms");
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

func SleepR() {
	Rand.Seed(Time.Now().UnixNano());
	n := Rand.Intn(77) + 1;
	Time.Sleep(Time.Duration(n) * Time.Millisecond);
}
