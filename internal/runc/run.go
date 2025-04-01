package runc

import (
	"engine/internal/config"
	"engine/internal/containers"
	"fmt"
	"golang.org/x/sys/unix"
	"log"
	"os"
	"os/exec"
	"syscall"
	"time"
)

func RunContainer(tty bool, image string) error {
	//ass
	rootfs, err := Mount(image)
	if err != nil {
		fmt.Println("Mount rootfs failed:", err)
		return err
	}
	containerId := containers.GenContainerId() // 生成 10 位容器 id
	//rootfs := "/home/yzr/Documents/mydocker/data/overlay2/e518df1ad4972f366a0a74b1a2e0859e5262ed03825a700e6a06b4a5d782daa9/merged"
	var cmd *exec.Cmd
	if cmd, err = Run(tty, rootfs, containerId); err != nil {
		fmt.Println("Run failed:", err)
	}

	containers.SaveContainerInfo(&containers.ContainerInfo{
		ContainerID: containerId,
		Command:     "",
		Created:     fmt.Sprintf("%d", time.Now().Unix()),
		Status:      "running",
		Image:       image,
		Pid:         cmd.Process.Pid,
	})

	if tty {
		cmd.Wait()
	} else {
		cmd.Process.Wait()

	}
	fmt.Println("Container exited")
	return nil
}
func Run(tty bool, rootfs string, containerId string) (cmd *exec.Cmd, err error) {

	fmt.Println("Main Process Pid :", os.Getpid())
	//cmd2 := exec.Command("id")
	//
	//cmd2.Stderr = os.Stderr
	//cmd2.Stdout = os.Stdout
	//err = cmd2.Run()
	//if err != nil {
	//	fmt.Println("id 命令执行失败:", err)
	//	return err
	//}
	//fmt.Println("id 命令执行成功")

	cmd = exec.Command("/bin/sh", "-c", fmt.Sprintf(`
		

		# 步骤2: 显示用户映射
		echo '用户映射:'
		cat /proc/self/uid_map
		cat /proc/self/gid_map

		# 步骤3: 设置主机名
		echo '正在设置主机名...'
		#sed -i 's/^root:/root:/' %q/etc/passwd

		# 步骤4: 显示用户信息
		echo '当前用户信息:'
		id
		# 查看当前进程的权能状态
		grep Cap /proc/self/status
		# 步骤5: 进入容器环境
		echo '正在进入容器环境...'
		exec chroot %q  /bin/bash -c 'echo "容器内 hostname: ";whoami; exec /bin/bash'
	`, rootfs, rootfs))

	cmd.Dir = "/" // 工作目录为根目录
	cmd.Env = []string{"PATH=/bin:/usr/bin:/sbin:/usr/sbin", "TERM=xterm"}
	if err := unix.Prctl(unix.PR_CAPBSET_READ, unix.CAP_SYS_CHROOT, 0, 0, 0); err != nil {
		log.Fatal("Prctl failed:", err)
	}
	// 设置新的命名空间
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUSER |
			syscall.CLONE_NEWUTS |
			syscall.CLONE_NEWPID |
			syscall.CLONE_NEWNET |
			syscall.CLONE_NEWNS |
			syscall.CLONE_NEWIPC |
			syscall.CLONE_NEWTIME,
		UidMappings: []syscall.SysProcIDMap{
			{ContainerID: 0, HostID: 1000, Size: 1},
		},

		GidMappings: []syscall.SysProcIDMap{
			{ContainerID: 0, HostID: 1000, Size: 1},
		},
		Credential: &syscall.Credential{
			Uid: 0, // 容器内 root
			Gid: 0,
		},
		//Chroot: rootfs, // 切换根目录as

		GidMappingsEnableSetgroups: false, // 开启 GID 映射
		//Setsid:                     true,  // 新会话
		//Setctty:                    true,
		// 关键：在 fork 后、exec 前设置主机名

	}

	// 启动命令
	fmt.Println("启动 bash 进程...")
	if tty {
		//前台运行
		fmt.Println("正在前台运行...")
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	} else {
		//重定向输出到文件
		f, err := os.OpenFile(fmt.Sprintf(config.Conf.EnvConf.ImagesDataDir+"%s/log.log", containerId), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			fmt.Println("打开日志文件失败:", err)
			return nil, err
		}
		cmd.Stdout = f
		cmd.Stderr = f
	}

	return cmd, cmd.Start()
}
