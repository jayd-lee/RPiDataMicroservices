package device

import (
	"fmt"
	"sync"
	"time"
)

type SensorData struct {
	Temperature float64
	Humidity    float64
	Timestamp   time.Time
}

type Device struct {
	ID       string
	Location string
	LastData SensorData
	mu       sync.RWMutex
}

func NewDevice(id, location string) *Device {
	return &Device{
		ID:       id,
		Location: location,
	}
}

func (d *Device) UpdateData(temp, hum float64) {
	d.mu.Lock()
	defer d.mu.Unlock()

	d.LastData = SensorData{
		Temperature: temp,
		Humidity:    hum,
		Timestamp:   time.Now(),
	}
}

func (d *Device) GetLastData() SensorData {
	d.mu.RLock()
	defer d.mu.RUnlock()
	return d.LastData
}

type DeviceManager struct {
	devices map[string]*Device
	mu      sync.RWMutex
}

func NewDeviceManager() *DeviceManager {
	return &DeviceManager{
		devices: make(map[string]*Device),
	}
}

func (dm *DeviceManager) AddDevice(id, location string) error {
	dm.mu.Lock()
	defer dm.mu.Unlock()

	if _, exists := dm.devices[id]; exists {
		return fmt.Errorf("device with ID %s already exists", id)
	}

	dm.devices[id] = NewDevice(id, location)
	return nil
}

func (dm *DeviceManager) GetDevice(id string) (*Device, error) {
	dm.mu.RLock()
	defer dm.mu.RUnlock()

	device, exists := dm.devices[id]
	if !exists {
		return nil, fmt.Errorf("device with ID %s not found", id)
	}

	return device, nil
}

func (dm *DeviceManager) GetAllDevices() []*Device {
	dm.mu.RLock()
	defer dm.mu.RUnlock()

	devices := make([]*Device, 0, len(dm.devices))
	for _, device := range dm.devices {
		devices = append(devices, device)
	}
	return devices
}
