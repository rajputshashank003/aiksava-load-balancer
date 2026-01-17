package config

func IsAllowedOrigin(origin string) bool {
	for _, o := range GetAllowedOrigins() {
		if origin == o {
			return true
		}
	}
	return false
}