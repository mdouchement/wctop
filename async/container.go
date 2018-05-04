package async

import (
	"fmt"
	"sync"
	"time"

	"github.com/mdouchement/wctop/container"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/mem"
	"github.com/shirou/gopsutil/net"
)

type (
	M      map[string]interface{}
	OsStat struct {
		CPUUsagePercent float32 `json:"cpu_usage_percent"`
		Mem             M       `json:"mem"`
		Net             M       `json:"net"`
		IO              M       `json:"io"`
	}

	ContainerStat = container.Container
)

const interval = 1 * time.Second

var (
	mu     sync.Mutex
	ticker *time.Ticker
	conn   = container.NewDocker()
)

func Start() {
	mu.Lock()
	defer mu.Unlock()

	if WsNotifier.Len() > 0 {
		return
	}
	ticker = time.NewTicker(interval)

	go func() {
		fetch()
		for _ = range ticker.C {
			fetch()
		}
	}()
}

func Stop() {
	time.Sleep(1 * time.Minute) // until 1 minute of idle

	mu.Lock()
	defer mu.Unlock()

	if WsNotifier.Len() > 0 {
		// Remaining subscribers
		return
	}

	fmt.Println("Idle: No longer subscribers")

	ticker.Stop()
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

	nio, err := net.IOCounters(false)
	if err != nil {
		return nil, err
	}

	dio, err := disk.IOCounters()
	if err != nil {
		return nil, err
	}
	var readBytes, writeBytes uint64
	for _, v := range dio {
		readBytes += v.ReadBytes
		writeBytes += v.WriteBytes
	}

	return &OsStat{
		CPUUsagePercent: float32(percs[0]),
		Mem: M{
			"total":     vms.Total,
			"available": vms.Available,
		},
		Net: M{
			"tx": nio[0].BytesSent,
			"rx": nio[0].BytesRecv,
		},
		IO: M{
			"bytes_read":  readBytes,
			"bytes_write": writeBytes,
		},
	}, nil
}

func ctop() []*ContainerStat {
	return conn.All()
}
