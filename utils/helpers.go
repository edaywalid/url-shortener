package utils

func ToBase62(num uint64) string {
	var base62 = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	var result string
	for num > 0 {
		result = string(base62[num%62]) + result
		num /= 62
	}
	return result
}
