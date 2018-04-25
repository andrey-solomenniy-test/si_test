package ipcache

import (
	"time"
)

type cacheRecord struct {
	// ip        [6]byte
	code      [2]byte
	whenSaved time.Time
}

type Cache struct {
	// ar      []cacheRecord
	ar      map[string]cacheRecord
	max     int
	curSize int
	exp     int
}

func NewCache(size int, exp int) *Cache {
	c := &Cache{}
	c.max = size
	c.curSize = 0
	c.ar = make(map[string]cacheRecord, size)
	c.exp = exp
	return c
}

func (c *Cache) Find(ip string) (bool, [2]byte) {
	cr, ok := c.ar[ip]
	if !ok {
		return false, [2]byte{}
	}

	t := time.Now()
	diff := t.Sub(cr.whenSaved)
	if diff >= (time.Duration(c.exp) * time.Minute) {
		delete(c.ar, ip)
		c.curSize--
		return false, [2]byte{}
	}

	return true, cr.code
}

func (c *Cache) Save(ip string, code [2]byte) {
	if c.curSize >= c.max {
		return
	}

	c.ar[ip] = cacheRecord{code, time.Now()}
	c.curSize++
}
