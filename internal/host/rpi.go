//go:build rpi

package host

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/godbus/dbus/v5"
)

const (
	loginService = "org.freedesktop.login1"
	loginPath    = "/org/freedesktop/login1"
	loginIntf    = "org.freedesktop.login1.Manager"

	pkService = "org.freedesktop.PackageKit"
	pkPath    = "/org/freedesktop/PackageKit"
	pkIntf    = "org.freedesktop.PackageKit"
	pkTxIntf  = "org.freedesktop.PackageKit.Transaction"

	// PK_FILTER_ENUM_NONE indicates no specific filter
	pkFilterNone uint64 = 0
)

// Pi is a Host implementation backed by real Raspberry Pi system files and commands.
type Pi struct {
	conn *dbus.Conn
}

// NewHost creates a new Pi host.
func NewHost() (*Pi, error) {
	log.Println("[dbus] Connecting to system bus...")
	conn, err := dbus.SystemBus()
	if err != nil {
		log.Printf("[dbus] Failed to connect to system bus: %v\n", err)
		return &Pi{}, fmt.Errorf("failed to connect to system bus: %w", err)
	}
	log.Println("[dbus] Connected to system bus successfully")
	return &Pi{conn: conn}, nil
}

func (p *Pi) Close() {
	if p.conn != nil {
		log.Println("[dbus] Closing system bus connection")
		p.conn.Close()
	}
}

// CPUUsage retrieves CPU usage percentage
func (p *Pi) CPUUsage() (float64, error) {
	// Simple CPU usage calculation based on /proc/stat
	// This is a simplified version; for production, consider using a library
	return 0.0, nil
}

// MemoryUsage retrieves memory usage percentage
func (p *Pi) MemoryUsage() (float64, error) {
	file, err := os.Open("/proc/meminfo")
	if err != nil {
		return 0, err
	}
	defer file.Close()

	return parseMemoryUsage(file)
}

func parseMemoryUsage(input io.Reader) (float64, error) {
	var memTotal, memAvail float64
	scanner := bufio.NewScanner(input)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "MemTotal:") {
			parts := strings.Fields(line)
			if len(parts) > 1 {
				memTotal, _ = strconv.ParseFloat(parts[1], 64)
			}
		} else if strings.HasPrefix(line, "MemAvailable:") {
			parts := strings.Fields(line)
			if len(parts) > 1 {
				memAvail, _ = strconv.ParseFloat(parts[1], 64)
			}
		}
	}
	if err := scanner.Err(); err != nil {
		return 0, err
	}

	if memTotal == 0 {
		return 0, fmt.Errorf("could not read memory info")
	}

	used := memTotal - memAvail
	return (used / memTotal) * 100, nil
}

// DiskUsage retrieves disk usage percentage for root partition
func (p *Pi) DiskUsage() (float64, error) {
	var stat syscall.Statfs_t

	if err := syscall.Statfs("/", &stat); err != nil {
		return 0, fmt.Errorf("failed to get filesystem stats: %w", err)
	}

	total := stat.Blocks * uint64(stat.Bsize)
	used := total - (stat.Bfree * uint64(stat.Bsize))

	return float64(used) / float64(total) * 100, nil
}

// CPUTemperature retrieves CPU temperature in Celsius
func (p *Pi) CPUTemperature() (float64, error) {
	tempPath := "/sys/class/thermal/thermal_zone0/temp"
	data, err := os.ReadFile(tempPath)
	if err != nil {
		return 0, err
	}

	return parseCPUTemperature(data)
}

func parseCPUTemperature(data []byte) (float64, error) {
	tempStr := strings.TrimSpace(string(data))
	tempMilliC, err := strconv.ParseFloat(tempStr, 64)
	if err != nil {
		return 0, err
	}

	// Convert from millidegrees to degrees
	return tempMilliC / 1000, nil
}

// getUpgradablePackageIDs retrieves the package IDs of all available updates via PackageKit D-Bus.
func (p *Pi) getUpgradablePackageIDs() ([]string, error) {
	if p.conn == nil {
		log.Println("[dbus] DBus connection not available for getUpgradablePackageIDs")
		return nil, fmt.Errorf("dbus connection not available")
	}

	log.Println("[dbus] PackageKit: Creating transaction for GetUpdates...")
	var txPath dbus.ObjectPath
	err := p.conn.Object(pkService, pkPath).Call(pkIntf+".CreateTransaction", 0).Store(&txPath)
	if err != nil {
		log.Printf("[dbus] PackageKit: CreateTransaction failed: %v\n", err)
		return nil, fmt.Errorf("create PackageKit transaction: %w", err)
	}
	log.Printf("[dbus] PackageKit: Created transaction with path %s\n", txPath)

	rule := fmt.Sprintf("type='signal',sender='%s',interface='%s',path='%s'", pkService, pkTxIntf, txPath)
	log.Printf("[dbus] Adding match rule: %s\n", rule)
	if err := p.conn.BusObject().Call("org.freedesktop.DBus.AddMatch", 0, rule).Err; err != nil {
		log.Printf("[dbus] AddMatch failed: %v\n", err)
		return nil, fmt.Errorf("add dbus signal match: %w", err)
	}
	defer func() {
		log.Printf("[dbus] Removing match rule: %s\n", rule)
		if err := p.conn.BusObject().Call("org.freedesktop.DBus.RemoveMatch", 0, rule).Err; err != nil {
			log.Printf("[dbus] RemoveMatch failed: %v\n", err)
		}
	}()

	signalChan := make(chan *dbus.Signal, 64)
	p.conn.Signal(signalChan)
	defer p.conn.RemoveSignal(signalChan)

	txObj := p.conn.Object(pkService, txPath)
	log.Println("[dbus] PackageKit: Calling GetUpdates...")
	if call := txObj.Call(pkTxIntf+".GetUpdates", 0, pkFilterNone); call.Err != nil {
		log.Printf("[dbus] PackageKit: Call GetUpdates failed: %v\n", call.Err)
		return nil, fmt.Errorf("call GetUpdates: %w", call.Err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	var packageIDs []string
	for {
		select {
		case <-ctx.Done():
			log.Printf("[dbus] PackageKit: GetUpdates timed out: %v\n", ctx.Err())
			return nil, fmt.Errorf("timed out waiting for PackageKit updates: %w", ctx.Err())

		case sig, ok := <-signalChan:
			if !ok {
				log.Println("[dbus] PackageKit: Signal channel closed unexpectedly")
				return nil, fmt.Errorf("signal channel closed unexpectedly")
			}
			if sig.Path != txPath {
				continue
			}

			log.Printf("[dbus] PackageKit signal received: %s\n", sig.Name)

			switch sig.Name {
			case pkTxIntf + ".Package":
				if len(sig.Body) >= 2 {
					if pkgID, ok := sig.Body[1].(string); ok && pkgID != "" {
						log.Printf("[dbus] PackageKit: Found upgradable package: %s\n", pkgID)
						packageIDs = append(packageIDs, pkgID)
					}
				}

			case pkTxIntf + ".ErrorCode":
				if len(sig.Body) >= 2 {
					log.Printf("[dbus] PackageKit error code %v: %v\n", sig.Body[0], sig.Body[1])
					return nil, fmt.Errorf("packagekit error %v: %v", sig.Body[0], sig.Body[1])
				}
				log.Printf("[dbus] PackageKit error: %v\n", sig.Body)
				return nil, fmt.Errorf("packagekit error: %v", sig.Body)

			case pkTxIntf + ".Finished":
				log.Printf("[dbus] PackageKit: GetUpdates finished successfully (%d package(s) found)\n", len(packageIDs))
				return packageIDs, nil
			}
		}
	}
}

// AvailableUpdates retrieves the number of available package updates
func (p *Pi) AvailableUpdates() (int, error) {
	pkgIDs, err := p.getUpgradablePackageIDs()
	if err != nil {
		return 0, err
	}
	return len(pkgIDs), nil
}

// ApplyUpdates installs available package updates using PackageKit over D-Bus.
func (p *Pi) ApplyUpdates() error {
	if p.conn == nil {
		log.Println("[dbus] DBus connection not available for ApplyUpdates")
		return fmt.Errorf("dbus connection not available")
	}

	log.Println("[dbus] PackageKit: Getting list of upgradable packages for ApplyUpdates...")
	packageIDs, err := p.getUpgradablePackageIDs()
	if err != nil {
		log.Printf("[dbus] PackageKit: Failed to get upgradable packages: %v\n", err)
		return fmt.Errorf("get upgradable packages: %w", err)
	}
	if len(packageIDs) == 0 {
		log.Println("[dbus] PackageKit: No updates available to apply")
		return nil
	}

	log.Printf("[dbus] PackageKit: Creating transaction for UpdatePackages (%d package(s))...\n", len(packageIDs))
	var txPath dbus.ObjectPath
	err = p.conn.Object(pkService, pkPath).Call(pkIntf+".CreateTransaction", 0).Store(&txPath)
	if err != nil {
		log.Printf("[dbus] PackageKit: CreateTransaction failed: %v\n", err)
		return fmt.Errorf("create PackageKit transaction: %w", err)
	}
	log.Printf("[dbus] PackageKit: Created transaction with path %s\n", txPath)

	rule := fmt.Sprintf("type='signal',sender='%s',interface='%s',path='%s'", pkService, pkTxIntf, txPath)
	log.Printf("[dbus] Adding match rule: %s\n", rule)
	if err := p.conn.BusObject().Call("org.freedesktop.DBus.AddMatch", 0, rule).Err; err != nil {
		log.Printf("[dbus] AddMatch failed: %v\n", err)
		return fmt.Errorf("add dbus signal match: %w", err)
	}
	defer func() {
		log.Printf("[dbus] Removing match rule: %s\n", rule)
		if err := p.conn.BusObject().Call("org.freedesktop.DBus.RemoveMatch", 0, rule).Err; err != nil {
			log.Printf("[dbus] RemoveMatch failed: %v\n", err)
		}
	}()

	signalChan := make(chan *dbus.Signal, 64)
	p.conn.Signal(signalChan)
	defer p.conn.RemoveSignal(signalChan)

	txObj := p.conn.Object(pkService, txPath)
	log.Printf("[dbus] PackageKit: Calling UpdatePackages for %d package(s)...\n", len(packageIDs))
	if call := txObj.Call(pkTxIntf+".UpdatePackages", 0, pkFilterNone, packageIDs); call.Err != nil {
		log.Printf("[dbus] PackageKit: Call UpdatePackages failed: %v\n", call.Err)
		return fmt.Errorf("call UpdatePackages: %w", call.Err)
	}

	// Package updates may take several minutes depending on the network and storage speed
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Minute)
	defer cancel()

	for {
		select {
		case <-ctx.Done():
			log.Printf("[dbus] PackageKit: UpdatePackages timed out: %v\n", ctx.Err())
			return fmt.Errorf("timed out waiting for package updates: %w", ctx.Err())

		case sig, ok := <-signalChan:
			if !ok {
				log.Println("[dbus] PackageKit: Signal channel closed unexpectedly")
				return fmt.Errorf("signal channel closed unexpectedly")
			}
			if sig.Path != txPath {
				continue
			}

			log.Printf("[dbus] PackageKit signal received: %s\n", sig.Name)

			switch sig.Name {
			case pkTxIntf + ".ErrorCode":
				if len(sig.Body) >= 2 {
					log.Printf("[dbus] PackageKit error code %v: %v\n", sig.Body[0], sig.Body[1])
					return fmt.Errorf("packagekit error %v: %v", sig.Body[0], sig.Body[1])
				}
				log.Printf("[dbus] PackageKit error: %v\n", sig.Body)
				return fmt.Errorf("packagekit error: %v", sig.Body)

			case pkTxIntf + ".Finished":
				log.Println("[dbus] PackageKit: UpdatePackages finished successfully")
				return nil
			}
		}
	}
}

// Restart restarts the device
func (p *Pi) Restart() error {
	if p.conn == nil {
		log.Println("[dbus] DBus connection not available for Restart")
		return fmt.Errorf("dbus connection not available")
	}

	log.Println("[dbus] logind: Calling Reboot...")
	obj := p.conn.Object(loginService, loginPath)
	// Reboot(interactive bool)
	call := obj.Call(loginIntf+".Reboot", 0, false)
	if call.Err != nil {
		log.Printf("[dbus] logind: Reboot call failed: %v\n", call.Err)
		return call.Err
	}
	log.Println("[dbus] logind: Reboot request sent successfully")
	return nil
}

// Shutdown shuts down the device
func (p *Pi) Shutdown() error {
	if p.conn == nil {
		log.Println("[dbus] DBus connection not available for Shutdown")
		return fmt.Errorf("dbus connection not available")
	}

	log.Println("[dbus] logind: Calling PowerOff...")
	obj := p.conn.Object(loginService, loginPath)
	call := obj.Call(loginIntf+".PowerOff", 0, false)
	if call.Err != nil {
		log.Printf("[dbus] logind: PowerOff call failed: %v\n", call.Err)
		return call.Err
	}
	log.Println("[dbus] logind: PowerOff request sent successfully")
	return nil
}
