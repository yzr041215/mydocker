package cmd

import (
	"engine/internal/runc"
	"fmt"
	"github.com/urfave/cli"
)

// runCommand 定义了run命令的相关参数和行为
var runCommand = cli.Command{
	Name: "run",
	Usage: `Create a container with namespace and cgroups limit
          mydocker run -it [command]`,
	Flags: []cli.Flag{
		cli.BoolFlag{
			Name:  "it", // 简单起见，这里把 -i 和 -t 参数合并成一个
			Usage: "enable tty",
		},
		cli.BoolFlag{
			Name:  "d",
			Usage: "detach container",
		},
		cli.StringFlag{ // 数据卷
			Name:  "v",
			Usage: "volume,e.g.: -v /ect/conf:/etc/conf",
		},
		cli.StringFlag{ // 数据卷
			Name:  "image",
			Usage: "image name",
		},
		// 省略其他代码
	},
	/*
	   这里是run命令执行的真正函数。
	   1.判断参数是否包含command
	   2.获取用户指定的command
	   3.调用Run function去准备启动容器:
	*/
	Action: func(context *cli.Context) error {
		if len(context.Args()) < 1 {
			return fmt.Errorf("missing container command")
		}

		var cmdArray []string
		for _, arg := range context.Args() {
			cmdArray = append(cmdArray, arg)
		}
		// tty和detach只能同时生效一个
		tty := context.Bool("it")
		detach := context.Bool("d")

		if tty && detach {
			return fmt.Errorf("it and d paramter can not both provided")
		}
		//resConf := &resource.ResourceConfig{
		//	MemoryLimit: context.String("mem"),
		//	CpuCfsQuota: context.Int("cpu"),
		//}
		image := context.String("image")
		volume := context.String("v")
		fmt.Println("volume:", volume)
		fmt.Println("cmdArray:", cmdArray)
		err := runc.RunContainer(tty, image)
		if err != nil {
			return err
		}
		return nil
	},
}
