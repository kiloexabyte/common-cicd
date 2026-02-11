package cacheactions

// Save stores the given paths under the specified key.
func (Ops) Save(key string, paths []string) error {
	c, err := newClient()
	if err != nil {
		return err
	}
	return c.save(key, paths)
}
