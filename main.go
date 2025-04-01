package main

import (
	"engine/cmd"
)

func main() {
	cmd.Cmd()
}

// sudo -E /usr/local/go/bin/go run main.go
//func main2() {
//	args := os.Args[1:]
//	if len(args) == 1 {
//		Son()
//		return
//	}
//	rootfs, err := runc.Mount("mysql")
//	if err != nil {
//		fmt.Println(err)
//	}
//
//	cmd := exec.Command("/proc/self/exe", rootfs)
//	cmd.Stdin = os.Stdin
//	cmd.Stdout = os.Stdout
//	cmd.Stderr = os.Stderr
//	fmt.Println("pid ", os.Getpid())
//	err = cmd.Run()
//	if err != nil {
//		return
//	}
//}
//func Son() {
//	rootfs := os.Args[1]
//	if err := runc.Run(rootfs); err != nil {
//		fmt.Println(err)
//	}
//}
