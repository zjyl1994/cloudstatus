package main

import (
	"strings"
	"time"

	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/disk"
	"github.com/shirou/gopsutil/v4/host"
	"github.com/shirou/gopsutil/v4/load"
	"github.com/shirou/gopsutil/v4/mem"
	"github.com/shirou/gopsutil/v4/net"
	"github.com/zjyl1994/cloudstatus/infra/define"
	"github.com/zjyl1994/cloudstatus/service/sensors"
)

var (
	excludeInterfaceNamePrefix = []string{"lo", "tun", "docker", "veth", "br-", "vmbr", "vnet", "kube"}
	measureTimestamp           int64
	lastDiskRead               uint64
	lastDiskWrite              uint64
	lastNetworkSend            uint64
	lastNetworkRecv            uint64
)

func Measure(interval time.Duration) (*define.StatExchangeFormat, error) {
	var result define.StatExchangeFormat
	// CPU Percent
	if percent, err := cpu.Percent(interval, false); err == nil {
		result.Percent.CPU = percent[0]
	}

	now := time.Now().Unix()
	duration := uint64(now - measureTimestamp)

	// Disk speed
	counters, err := disk.IOCounters()
	if err != nil {
		return nil, err
	}
	var readBytes, writeBytes uint64
	for _, c := range counters {
		readBytes += c.ReadBytes
		writeBytes += c.WriteBytes
	}

	result.Disk.Read = (readBytes - lastDiskRead) / duration
	result.Disk.Write = (writeBytes - lastDiskWrite) / duration

	result.Disk.Read = readBytes
	result.Disk.Write = writeBytes

	lastDiskRead = readBytes
	lastDiskWrite = writeBytes

	// Network speed
	in, out, err := getNetInOut()
	if err != nil {
		return nil, err
	}

	result.Network.Recv = in - lastNetworkRecv
	result.Network.Send = out - lastNetworkSend
	result.Network.Rx = result.Network.Recv / duration
	result.Network.Tx = result.Network.Send / duration
	lastNetworkRecv = in
	lastNetworkSend = out

	measureTimestamp = now

	// Load
	loadAvg, err := load.Avg()
	if err != nil {
		return nil, err
	}
	result.Load.Load1 = loadAvg.Load1
	result.Load.Load5 = loadAvg.Load5
	result.Load.Load15 = loadAvg.Load15
	// Usages
	vm, err := mem.VirtualMemory()
	if err != nil {
		return nil, err
	}
	result.Memory = define.UsageStat{
		Total: vm.Total,
		Used:  vm.Used,
		Free:  vm.Free,
	}
	result.Percent.Mem = vm.UsedPercent

	sm, err := mem.SwapMemory()
	if err != nil {
		return nil, err
	}
	result.Swap = define.UsageStat{
		Total: sm.Total,
		Used:  sm.Used,
		Free:  sm.Free,
	}
	result.Percent.Swap = sm.UsedPercent

	diskUsage, err := disk.Usage("/")
	if err != nil {
		return nil, err
	}
	result.Disk.UsageStat = define.UsageStat{
		Total: diskUsage.Total,
		Used:  diskUsage.Used,
		Free:  diskUsage.Free,
	}
	result.Percent.Disk = diskUsage.UsedPercent
	// host info
	hostinfo, err := host.Info()
	if err != nil {
		return nil, err
	}
	result.Hostname = hostinfo.Hostname
	result.Uptime = hostinfo.Uptime

	// info
	result.Interval = uint64(interval.Seconds())
	result.ReportTime = now

	// temperature
	if temperatures, err := sensors.Get(); err == nil {
		result.Temperature = temperatures
	}
	return &result, nil
}

func getNetInOut() (netIn uint64, netOut uint64, err error) {
	nv, err := net.IOCounters(true)
	if err != nil {
		return 0, 0, err
	}
	for _, v := range nv {
		if matchPrefix(v.Name, excludeInterfaceNamePrefix) {
			continue
		}
		netIn += v.BytesRecv
		netOut += v.BytesSent
	}
	return netIn, netOut, nil
}

func matchPrefix(s string, prefixs []string) bool {
	for _, prefix := range prefixs {
		if strings.HasPrefix(s, prefix) {
			return true
		}
	}
	return false
}
