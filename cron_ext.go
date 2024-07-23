package cron

import "container/heap"

// RemoveCheckFunc 删除job的检查函数，返回true则删除
type RemoveCheckFunc func(e *Entry) bool

func (c *Cron) removeEntryByJob(cb RemoveCheckFunc) {
	for idx, e := range c.entries {
		if cb(e) {
			heap.Remove(&c.entries, idx)
			c.logger.Info("removed", "entry", e.ID)
		}
	}
}

func (c *Cron) RemoveJob(cb RemoveCheckFunc) {
	c.runningMu.Lock()
	defer c.runningMu.Unlock()
	if c.running {
		c.removeJob <- cb
	} else {
		c.removeEntryByJob(cb)
	}
}
