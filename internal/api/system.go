package api

import (
	"bufio"
	"fmt"
	"math"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type SystemHandler struct{}

func NewSystemHandler() *SystemHandler {
	return &SystemHandler{}
}

type SystemStats struct {
	CPUPercent    float64 `json:"cpu_percent"`
	RAMUsedMB     uint64  `json:"ram_used_mb"`
	RAMTotalMB    uint64  `json:"ram_total_mb"`
	RAMPercent    float64 `json:"ram_percent"`
	NetRxBytes    uint64  `json:"net_rx_bytes"`
	NetTxBytes    uint64  `json:"net_tx_bytes"`
	UptimeSeconds float64 `json:"uptime_seconds"`
	LoadAvg1      float64 `json:"load_avg_1"`
	LoadAvg5      float64 `json:"load_avg_5"`
	LoadAvg15     float64 `json:"load_avg_15"`
}

type cpuSample struct {
	total uint64
	idle  uint64
}

func readCPUSample() (cpuSample, error) {
	f, err := os.Open("/proc/stat")
	if err != nil {
		return cpuSample{}, err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		if !strings.HasPrefix(line, "cpu ") {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) < 5 {
			break
		}
		var total, idle uint64
		for i, field := range fields[1:] {
			v, err := strconv.ParseUint(field, 10, 64)
			if err != nil {
				continue
			}
			total += v
			if i == 3 {
				idle = v
			}
		}
		return cpuSample{total: total, idle: idle}, nil
	}
	return cpuSample{}, fmt.Errorf("cpu line not found in /proc/stat")
}

func getCPUPercent() (float64, error) {
	s1, err := readCPUSample()
	if err != nil {
		return 0, err
	}
	time.Sleep(500 * time.Millisecond)
	s2, err := readCPUSample()
	if err != nil {
		return 0, err
	}
	totalDelta := float64(s2.total - s1.total)
	idleDelta := float64(s2.idle - s1.idle)
	if totalDelta == 0 {
		return 0, nil
	}
	return math.Round((1-idleDelta/totalDelta)*1000) / 10, nil
}

func getMemStats() (usedMB, totalMB uint64, percent float64, err error) {
	f, err := os.Open("/proc/meminfo")
	if err != nil {
		return
	}
	defer f.Close()

	var memTotal, memAvailable uint64
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Fields(line)
		if len(fields) < 2 {
			continue
		}
		val, parseErr := strconv.ParseUint(fields[1], 10, 64)
		if parseErr != nil {
			continue
		}
		switch fields[0] {
		case "MemTotal:":
			memTotal = val
		case "MemAvailable:":
			memAvailable = val
		}
	}
	if memTotal == 0 {
		err = fmt.Errorf("MemTotal not found in /proc/meminfo")
		return
	}
	usedKB := memTotal - memAvailable
	totalMB = memTotal / 1024
	usedMB = usedKB / 1024
	percent = math.Round(float64(usedKB)/float64(memTotal)*1000) / 10
	return
}

func getNetStats() (rxBytes, txBytes uint64, err error) {
	f, err := os.Open("/proc/net/dev")
	if err != nil {
		return
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	scanner.Scan() // header line 1
	scanner.Scan() // header line 2
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		colonIdx := strings.Index(line, ":")
		if colonIdx < 0 {
			continue
		}
		iface := strings.TrimSpace(line[:colonIdx])
		if iface == "lo" {
			continue
		}
		fields := strings.Fields(line[colonIdx+1:])
		if len(fields) < 9 {
			continue
		}
		rx, e1 := strconv.ParseUint(fields[0], 10, 64)
		tx, e2 := strconv.ParseUint(fields[8], 10, 64)
		if e1 != nil || e2 != nil {
			continue
		}
		rxBytes += rx
		txBytes += tx
	}
	return
}

func getUptimeAndLoad() (uptime, load1, load5, load15 float64, err error) {
	uptimeData, err := os.ReadFile("/proc/uptime")
	if err != nil {
		return
	}
	fields := strings.Fields(string(uptimeData))
	if len(fields) > 0 {
		uptime, _ = strconv.ParseFloat(fields[0], 64)
	}

	loadData, err := os.ReadFile("/proc/loadavg")
	if err != nil {
		return
	}
	loadFields := strings.Fields(string(loadData))
	if len(loadFields) >= 3 {
		load1, _ = strconv.ParseFloat(loadFields[0], 64)
		load5, _ = strconv.ParseFloat(loadFields[1], 64)
		load15, _ = strconv.ParseFloat(loadFields[2], 64)
	}
	return
}

func (h *SystemHandler) GetStats(c *gin.Context) {
	stats := SystemStats{}

	cpuPct, err := getCPUPercent()
	if err == nil {
		stats.CPUPercent = cpuPct
	}

	usedMB, totalMB, ramPct, err := getMemStats()
	if err == nil {
		stats.RAMUsedMB = usedMB
		stats.RAMTotalMB = totalMB
		stats.RAMPercent = ramPct
	}

	rxBytes, txBytes, err := getNetStats()
	if err == nil {
		stats.NetRxBytes = rxBytes
		stats.NetTxBytes = txBytes
	}

	uptime, load1, load5, load15, err := getUptimeAndLoad()
	if err == nil {
		stats.UptimeSeconds = uptime
		stats.LoadAvg1 = load1
		stats.LoadAvg5 = load5
		stats.LoadAvg15 = load15
	}

	c.JSON(http.StatusOK, stats)
}
