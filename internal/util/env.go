package util

import "os"

func GetEnv(k, dv string) string {
	v := os.Getenv(k)
	if v == "" {
		return dv
	}

	return v
}
