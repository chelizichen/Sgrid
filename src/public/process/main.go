package pk

import (
	"fmt"

	"github.com/shirou/gopsutil/net"
	"github.com/shirou/gopsutil/process"
)

// findPidByPort finds the process ID (PID) using the specified port.
func findPidByPort(port string) (int32, error) {
	conns, err := net.Connections("tcp")
	if err != nil {
		return 0, err
	}

	for _, conn := range conns {
		if conn.Laddr.Port == portToUint32(port) {
			return conn.Pid, nil
		}
	}
	return 0, fmt.Errorf("no process found on port %s", port)
}

// portToUint32 converts a port string to uint32.
func portToUint32(port string) uint32 {
	var p uint32
	fmt.Sscan(port, &p)
	return p
}

// killProcess kills the process with the specified PID.
func killProcess(pid int32) error {
	proc, err := process.NewProcess(pid)
	if err != nil {
		return err
	}
	return proc.Kill()
}

func QueryProcessPidThenKill(port string) error {
	pid, err := findPidByPort(port)
	if err != nil {
		fmt.Println("Error finding PID: ", err)
		return err
	}
	fmt.Println("Found process with PID on port ", pid, port)
	err = killProcess(pid)
	if err != nil {
		fmt.Println("Error killing process: ", err)
		return err
	}
	fmt.Println("Successfully killed process with PID ", pid)
	return err
}
