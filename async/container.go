package async

import (
	"fmt"
	"sync"
	"time"

	"github.com/mdouchement/wctop/container"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
)

type (
	M      map[string]interface{}
	OsStat struct {
		CPUUsagePercent float32 `json:"cpu_usage_percent"`
		Mem             M       `json:"mem"`
	}

	ContainerStat = container.Container
)

const interval = 1 * time.Second

var (
	mu   sync.Mutex
	done chan bool
	conn = container.NewDocker()
)

func init() {
	done = make(chan bool)
}

func Start() {
	mu.Lock()
	defer mu.Unlock()

	if WsNotifier.Len() > 0 {
		return
	}

	go func() {
		ticker := time.NewTicker(interval)

		fetch()
		for {
			select {
			case <-done:
				return
			case <-ticker.C:
				fetch()
			}
		}
	}()
}

func Stop() {
	time.Sleep(1 * time.Second) // TODO until 1 minute of idle

	mu.Lock()
	defer mu.Unlock()

	if WsNotifier.Len() > 0 {
		// Remaining subscribers
		return
	}

	fmt.Println("Idle: No longer subscribers")

	done <- true
	conn.Reset()
}

func fetch() {
	t, err := top()

	if err != nil {
		WsNotifier.Notify(&Notification{
			Error: err.Error(),
		})
		time.Sleep(10 * time.Second)
		return
	}

	WsNotifier.Notify(&Notification{
		UpdatedAt:       time.Now(),
		OsStats:         t,
		ContainersStats: ctop(),
	})
}

func top() (*OsStat, error) {
	percs, err := cpu.Percent(500*time.Millisecond, false)
	if err != nil {
		return nil, err
	}

	vms, err := mem.VirtualMemory()
	if err != nil {
		return nil, err
	}

	return &OsStat{
		CPUUsagePercent: float32(percs[0]),
		Mem: M{
			"total":     vms.Total,
			"available": vms.Available,
		},
	}, nil
}

func ctop() []*ContainerStat {
	return conn.All()
}
