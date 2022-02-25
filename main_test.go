package ipscraper

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestScraper_Get(t *testing.T) {
	ips, err := New().Get()
	assert.NoError(t, err)
	assert.GreaterOrEqual(t, len(ips), 5)
}
