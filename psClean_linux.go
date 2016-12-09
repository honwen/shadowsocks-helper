package main

import "syscall"

func psClean(pid int) {
	syscall.Kill(-pid, syscall.SIGKILL)
}
