package main

import "time"

func nowTimestamp() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}

func nowMinusIntervalTimestamp() int64 {
	return (time.Now().UnixNano() / int64(time.Millisecond)) - GameInterval
}
