package container

import (
	"fmt"
	"math"
	"sync"
	"time"

	api "github.com/fsouza/go-dockerclient"
	"github.com/ulule/deepcopier"
)

type Container struct {
	mu         sync.Mutex `deepcopier:"skip"`
	done       chan bool  `deepcopier:"skip"`
	running    bool       `deepcopier:"skip"`
	conn       *Connector `deepcopier:"skip"`
	lastCpu    float64    `deepcopier:"skip"`
	lastSys    float64    `deepcopier:"skip"`
	cloned     bool       `deepcopier:"skip"`
	ID         string     `json:"id"`
	StartedAt  time.Time  `json:"started_at"`
	Name       string     `json:"name"`
	Status     string     `json:"status"`
	CPUUsage   float32    `json:"cpu_usage"`
	MemLimit   int64      `json:"mem_limit"`
	MemUsage   int64      `json:"mem_usage"`
	MemPercent float32    `json:"mem_percent"`
}

func (c *Container) Start() {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.done == nil {
		c.done = make(chan bool)
	}
	stats := make(chan *api.Stats)

	go func() {
		opts := api.StatsOptions{
			ID:     c.ID,
			Stats:  stats,
			Stream: true,
			Done:   c.done,
		}
		fmt.Printf("Starting watching container %s\n", c.Name)
		c.conn.client.Stats(opts)
		fmt.Printf("Stopping watching container %s\n", c.Name)

		c.mu.Lock()
		defer c.mu.Unlock()
		c.running = false
	}()

	go func() {
		for s := range stats {
			c.readMem(s)
			c.readCPU(s)
		}
	}()
	c.running = true
}

func (c *Container) IsCloned() bool {
	return c.cloned
}

func (c *Container) Clone() *Container {
	c.mu.Lock()
	defer c.mu.Unlock()

	nc := &Container{}
	deepcopier.Copy(c).To(nc)
	nc.cloned = true
	return nc
}

func (c *Container) Stop() *Container {
	c.done <- true
	return c
}

func (c *Container) Running() bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.running
}

func (c *Container) readMem(stats *api.Stats) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if stats != nil {
		c.MemUsage = int64(stats.MemoryStats.Usage - stats.MemoryStats.Stats.Cache)
		c.MemLimit = int64(stats.MemoryStats.Limit)
		c.MemPercent = float32(percent(float64(c.MemUsage), float64(c.MemLimit)))
	}
}

func (c *Container) readCPU(stats *api.Stats) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if stats != nil {
		ncpus := float64(len(stats.CPUStats.CPUUsage.PercpuUsage))
		total := float64(stats.CPUStats.CPUUsage.TotalUsage)
		system := float64(stats.CPUStats.SystemCPUUsage)

		cpudiff := total - c.lastCpu
		syscpudiff := system - c.lastSys

		c.CPUUsage = float32(round((cpudiff / syscpudiff * 100) * ncpus))
		c.lastCpu = total
		c.lastSys = system
	}
}

// return rounded percentage
func percent(val float64, total float64) float64 {
	if total <= 0 {
		return 0
	}
	return round((val / total) * 100)
}

func round(num float64) float64 {
	return math.Trunc(num*100) / 100
}
