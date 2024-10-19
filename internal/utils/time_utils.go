package utils

import "time"

func NowUnix() int64 {
	return time.Now().Unix()
}

func GetLastMidnight() int64 {
	return time.Now().Truncate(24 * time.Hour).Unix()
}
