package vampires

// createKey creates a levelDB key from the specified path.
func createKey(key string) []byte {
	return []byte("_file://\x00\x01" + key)
}
