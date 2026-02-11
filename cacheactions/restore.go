package cacheactions

// Restore retrieves cached paths for the given key.
func (Ops) Restore(key string, paths []string) error {
	c, err := newClient()
	if err != nil {
		return err
	}
	return c.restore(key, paths)
}
