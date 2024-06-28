package pk

import (
	"fmt"

	"github.com/shirou/gopsutil/net"
	"github.com/shirou/gopsutil/process"
)

var SgridProcessUtil = new(sgridProcessUtil)

type sgridProcessUtil struct{}

// findPidByPort finds the process ID (PID) using the specified port.
func (s *sgridProcessUtil) findPidByPort(port string) (int32, error) {
	conns, err := net.Connections("tcp")
	if err != nil {
		return 0, err
	}

	for _, conn := range conns {
		if conn.Laddr.Port == s.portToUint32(port) {
			return conn.Pid, nil
		}
	}
	return 0, fmt.Errorf("no process found on port %s", port)
}

// portToUint32 converts a port string to uint32.
func (s *sgridProcessUtil) portToUint32(port string) uint32 {
	var p uint32
	fmt.Sscan(port, &p)
	return p
}

// killProcess kills the process with the specified PID.
func (s *sgridProcessUtil) killProcess(pid int32) error {
	proc, err := process.NewProcess(pid)
	if err != nil {
		return err
	}
	return proc.Kill()
}

func (s *sgridProcessUtil) QueryProcessPidThenKill(port string) error {
	pid, err := s.findPidByPort(port)
	if err != nil {
		fmt.Println("Error finding PID: ", err)
		return err
	}
	fmt.Println("Found process with PID on port ", pid, port)
	err = s.killProcess(pid)
	if err != nil {
		fmt.Println("Error killing process: ", err)
		return err
	}
	fmt.Println("Successfully killed process with PID ", pid)
	return err
}
