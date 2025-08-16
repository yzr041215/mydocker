package runc

import (
	"engine/internal/config"
	"engine/internal/containers"
	"engine/internal/network"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"
	"strings"
	"syscall"
	"time"

	"github.com/sirupsen/logrus"
	"golang.org/x/sys/unix"
)

func BindVolume(rootfs string, volumes []string) error {
	for _, volume := range volumes {
		if len(volume) == 0 {
			continue
		}
		parts := strings.Split(volume, ":")
		if len(parts) != 3 {
			return fmt.Errorf("invalid volume format: %s", volume)
		}
		src := parts[0]
		dst := parts[1]

		if !strings.HasPrefix(dst, "/") {
			dst = path.Join(rootfs, "/", dst)
		}
		if err := os.MkdirAll(dst, 0755); err != nil {
			return err
		}
		if err := syscall.Mount(src, dst, "", syscall.MS_BIND, ""); err != nil {
			return err
		}
		logrus.Info("bind volume [", volume, "] success!")
	}
	return nil
}
func UnbindVolume(rootfs string, volumes []string) error {
	for _, volume := range volumes {
		if len(volume) == 0 {
			continue
		}
		parts := strings.Split(volume, ":")
		if len(parts) != 3 {
			return fmt.Errorf("invalid volume format: %s", volume)
		}
		//src := parts[0]
		dst := parts[1]

		if !strings.HasPrefix(dst, "/") {
			dst = path.Join(rootfs, "/", dst)
		}
		if err := syscall.Unmount(dst, 0); err != nil {
			return err
		}
		logrus.Info("unbind volume success!")
	}
	return nil
}

func RunContainer(tty bool, image string, cmdArray []string, volume, portMapping []string, networkname string) error {

	rootfs, err := Mount(image)
	if err != nil {
		logrus.Error("Mount rootfs failed:", err)
		return err
	}
	err = BindVolume(rootfs, volume)
	if err != nil {
		logrus.Error("Bind volume failed:", err)
		return err
	}
	defer UnbindVolume(rootfs, volume)
	// 步骤1: 创建容器
	containerId := containers.GenContainerId() // 生成 10 位容器 id
	logrus.Info("=-----------------ContainerId:", containerId)
	var cmd *exec.Cmd
	if cmd, err = Run(tty, rootfs, containerId, cmdArray); err != nil {
		logrus.Error("Run failed:", err)
	} else {
		logrus.Info("Run success!", cmd.Process.Pid)
	}
	var containerIp string
	if networkname != "" {
		containerInfo := &containers.ContainerInfo{
			ContainerID: containerId,
			Pid:         cmd.Process.Pid,
			PortMapping: portMapping,
		}
		ip, err := network.Connect(networkname, containerInfo)
		if err != nil {
			logrus.Error("Connect failed:", err)
			return err
		} else {
			logrus.Info("Connect success! ip:", ip.String())
		}
		containerIp = ip.String()
	}
	info := &containers.ContainerInfo{
		ContainerID: containerId,
		Command:     strings.Join(cmdArray, " "),
		Created:     fmt.Sprintf("%d", time.Now().Unix()),
		Status:      "running",
		Image:       image,
		Pid:         cmd.Process.Pid,
		Volume:      volume,
		NetworkName: networkname,
		PortMapping: portMapping,
		IP:          containerIp,
	}
	containers.SaveContainerInfo(info)
	defer func() {
		if networkname != "" {
			err := network.Disconnect(networkname, info)
			if err != nil {
				logrus.Error("DisConnect failed:", err)
			} else {
				logrus.Info("DisConnect success!")
			}
		}
		//umount rootfs
		if err := syscall.Unmount(rootfs, 0); err != nil {
			logrus.Error("Unmount rootfs failed:", err)
		} else {
			logrus.Info("Unmount rootfs success!")
		}
		if err := syscall.Unmount(rootfs+"/proc", 0); err != nil {
			logrus.Error("Unmount proc failed:", err)
		} else {
			logrus.Info("Unmount proc success!")
		}
		// 保存容器信息
		containers.SaveContainerInfo(&containers.ContainerInfo{
			ContainerID: containerId,
			Command:     strings.Join(cmdArray, " "),
			Created:     fmt.Sprintf("%d", time.Now().Unix()),
			Status:      "exited",
			Image:       image,
			Pid:         cmd.Process.Pid,
			Volume:      volume,
			NetworkName: networkname,
			PortMapping: portMapping,
			IP:          containerIp,
		})
	}()
	if tty {
		logrus.Info("Container run in foreground, pid:", cmd.Process.Pid)
		cmd.Wait()

		logrus.Info("Container exited")
	} else {
		logrus.Info("Container run in background, pid:", cmd.Process.Pid, " ppid :", os.Getpid())
		cmd.Process.Wait()

		logrus.Info("Container Exit!!")
	}
	return nil
}
func Run(tty bool, rootfs string, containerId string, cmdArray []string) (cmd *exec.Cmd, err error) {

	logrus.Info("Main Process Pid :", os.Getpid())

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
		exec chroot %q  /bin/sh  -c 'echo "容器内 hostname: ";whoami; exec %s'
	`, rootfs, rootfs, strings.Join(cmdArray, " ")))

	cmd.Dir = "/" // 工作目录为根目录
	cmd.Env = []string{"PATH=/bin:/usr/bin:/sbin:/usr/sbin", "TERM=xterm"}
	if err := unix.Prctl(unix.PR_CAPBSET_READ, unix.CAP_SYS_CHROOT, 0, 0, 0); err != nil {
		log.Fatal("Prctl failed:", err)
	}
	//CAP_MKNOD
	if err := unix.Prctl(unix.PR_CAPBSET_READ, unix.CAP_MKNOD, 0, 0, 0); err != nil {
		log.Fatal("Prctl failed:", err)
	}
	// 设置新的命名空间
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags:
		//syscall.CLONE_NEWUSER |
		syscall.CLONE_NEWUTS |
			syscall.CLONE_NEWPID |
			syscall.CLONE_NEWNET |
			syscall.CLONE_NEWNS |
			syscall.CLONE_NEWIPC |
			syscall.CLONE_NEWTIME,
		// UidMappings: []syscall.SysProcIDMap{
		// 	{ContainerID: 0, HostID: 1000, Size: 1},
		// },

		// GidMappings: []syscall.SysProcIDMap{
		// 	{ContainerID: 0, HostID: 1000, Size: 1},
		// },
		// Credential: &syscall.Credential{
		// 	Uid: 0, // 容器内 root
		// 	Gid: 0,
		// },

		//GidMappingsEnableSetgroups: false, // 开启 GID 映射
	}

	// 启动命令
	logrus.Info("启动 bash 进程...")
	if tty {
		//前台运行
		logrus.Info("正在前台运行...")
		logrus.Info("containerId:", containerId)
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	} else {
		//重定向输出到文件
		dir := fmt.Sprintf(config.Conf.EnvConf.ImagesDataDir+"/containers/%s", containerId)
		if err := os.MkdirAll(dir, os.ModePerm); err != nil {
			logrus.Info("创建容器目录失败:", err)
			return nil, err
		}
		f, err := os.OpenFile(fmt.Sprintf(config.Conf.EnvConf.ImagesDataDir+"/containers/%s/log.log", containerId), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			logrus.Info("打开日志文件失败:", err)
			return nil, err
		}
		//cmd.Stdin = os.Stdin
		cmd.Stdout = f
		cmd.Stderr = f
	}

	return cmd, cmd.Start()
}
