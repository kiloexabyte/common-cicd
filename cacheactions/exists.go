package cacheactions

// Exists checks if a cache entry exists for the given key.
func (Ops) Exists(key string, paths []string) (bool, error) {
	c, err := newClient()
	if err != nil {
		return false, err
	}
	return c.exists(key, paths)
}
