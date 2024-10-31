package utils

import "math/rand"

func GetRandomUserAgent() string {
	var userAgents = []string{
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_5) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/86.0.4240.80 Safari/537.36",
	}
	return userAgents[rand.Intn(len(userAgents))]
}
