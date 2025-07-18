//TODO: finish comment headers
// Copyright 2025 PoiXson
// Use of this source code is governed by
// the AGPL+PXN license, which can be found
// in the LICENSE file.
package utils;

import(
	Time "time"
	Rand "math/rand"
);



// 10ms
func SleepX() {
	sleep, err := Time.ParseDuration("10ms");
	if err == nil { Time.Sleep(sleep); }
}
// 50ms
func SleepV() {
	sleep, err := Time.ParseDuration("50ms");
	if err == nil { Time.Sleep(sleep); }
}
// 100ms
func SleepC() {
	sleep, err := Time.ParseDuration("100ms");
	if err == nil { Time.Sleep(sleep); }
}
// x100ms
func SleepCn(n uint8) {
	sleep, err := Time.ParseDuration("100ms");
	if err == nil {
		for i:=uint8(0); i<n; i++ {
			Time.Sleep(sleep); }}
}
// 1s
func SleepS() {
	sleep, err := Time.ParseDuration("1s");
	if err == nil { Time.Sleep(sleep); }
}
// x1s
func SleepSn(n uint8) {
	sleep, err := Time.ParseDuration("1s");
	if err == nil {
		for i:=uint8(0); i<n; i++ {
			Time.Sleep(sleep); }}
}
// rand ms
func SleepR() {
	Rand.Seed(Time.Now().UnixNano());
	min := Rand.Intn(42)+15;
	max := Rand.Intn(55)+17;
	sleep := Rand.Intn(max+min) + min;
	Time.Sleep(Time.Duration(sleep) * Time.Millisecond);
}



// remove index from array
func RemoveIndex[T any](array []T, index int) []T {
	if index < 0 || index >= len(array) {
		return array; }
	return append(array[:index], array[index+1:]...);
}
