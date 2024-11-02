package utils

import "math/rand"

func GetRandomUserAgent() string {
	var userAgents = []string{
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_0) AppleWebKit/537.30 (KHTML, like Gecko) Chrome/130.0.0.0 Safari/537.30",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_1) AppleWebKit/537.31 (KHTML, like Gecko) Chrome/130.0.0.1 Safari/537.31",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_2) AppleWebKit/537.32 (KHTML, like Gecko) Chrome/130.0.0.2 Safari/537.32",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_3) AppleWebKit/537.33 (KHTML, like Gecko) Chrome/130.0.0.3 Safari/537.33",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_4) AppleWebKit/537.34 (KHTML, like Gecko) Chrome/130.0.0.4 Safari/537.34",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_5) AppleWebKit/537.35 (KHTML, like Gecko) Chrome/130.0.0.5 Safari/537.35",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/130.0.0.6 Safari/537.36",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.37 (KHTML, like Gecko) Chrome/130.0.0.7 Safari/537.37",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_8) AppleWebKit/537.38 (KHTML, like Gecko) Chrome/130.0.0.8 Safari/537.38",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_9) AppleWebKit/537.39 (KHTML, like Gecko) Chrome/130.0.0.9 Safari/537.39",
	}
	return userAgents[rand.Intn(len(userAgents))]
}
