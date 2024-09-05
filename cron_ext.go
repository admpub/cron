package cron

import (
	"container/heap"
)

// RemoveCheckFunc 删除job的检查函数，返回true则删除
type RemoveCheckFunc func(e *Entry) (removeable bool, continueable bool)

func (c *Cron) removeEntryByJob(cb RemoveCheckFunc) {
	var deleted int
	for idx, e := range c.entries {
		removeable, continueable := cb(e)
		if removeable {
			heap.Remove(&c.entries, idx-deleted)
			c.logger.Info("removed", "entry", e.ID)
			if !continueable {
				return
			}
			deleted++
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
