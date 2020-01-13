package main

import "time"

func nowUnix() int64 {
	return time.Now().Unix()
}

func getLastMidnight() int64 {
	return time.Now().Truncate(24 * time.Hour).Unix()
}
