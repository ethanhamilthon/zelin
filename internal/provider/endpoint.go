package provider

import "strings"

func (c *Client) Connect(base, key string) {
	c.catalog.Lock()
	defer c.catalog.Unlock()
	c.mu.Lock()
	defer c.mu.Unlock()
	c.URL, c.Key, c.models = strings.TrimRight(base, "/")+"/chat/completions", key, nil
}

func (c *Client) endpoint() (string, string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.URL, c.Key
}

func (c *Client) Base() string {
	url, _ := c.endpoint()
	return strings.TrimSuffix(url, "/chat/completions")
}
