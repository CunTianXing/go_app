package main

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"
	"time"

	"golang.org/x/sys/unix"
)

func main() {
	pid, _, err := unix.RawSyscall(syscall.SYS_FORK, 0, 0, 0)
	if int(pid) == -1 {
		fmt.Println("Failed to fork:", err)
		os.Exit(1)
	}

	fmt.Println(pid)
	if pid == 0 {
		child()
	} else {
		parent(int(pid))
	}
}

func parent(childPid int) {
	fmt.Println("I'm the parent")
	unix.Wait4(childPid, nil, 0, nil)
	fmt.Println("baby returned")
}

func child() {
	fmt.Println("I'm the baby, gotta love me!")
	time.Sleep(1 * time.Second)
	out, err := exec.Command("date").Output()
	if err != nil {
		fmt.Println("failed to run date:", err)
		os.Exit(1)
	}
	fmt.Printf("The date is %s\n", out)

	ch := make(chan bool)
	go func() {
		<-ch
		fmt.Println("in a goroutine!")
	}()
	ch <- true
}
