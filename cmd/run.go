package cmd

import (
	"engine/internal/runc"
	"fmt"

	"github.com/urfave/cli"
)

// sudo -E /usr/local/go/bin/go run main.go run -it -image mysql
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
		cli.StringSliceFlag{ // 数据卷
			Name:  "v",
			Usage: "volume,e.g.: -v /ect/conf:/etc/conf -v /ect/logs:/etc/logs",
		},
		cli.StringFlag{ // 数据卷
			Name:  "image",
			Usage: "image name",
		},
		cli.StringSliceFlag{
			Name:  "p",
			Usage: "port mapping,e.g.: -p 8080:80 -p 8081:81",
		},
		cli.StringFlag{
			Name:  "net",
			Usage: "network mode,e.g.: -net networkname",
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

		// if len(context.Args()) < 1 {
		// 	return fmt.Errorf("missing container command")
		// }
		fmt.Println("run command")
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

		fmt.Println("cmdArray:", cmdArray)
		image := context.String("image")
		fmt.Println("image:", image)
		volume := context.StringSlice("v")
		fmt.Println("volume:", volume)
		portMapping := context.StringSlice("p")
		fmt.Println("port:", portMapping)
		networkname := context.String("net")
		fmt.Println("networkMode:", networkname)

		err := runc.RunContainer(tty, image, cmdArray, volume, portMapping, networkname)
		if err != nil {
			return err
		}
		return nil
	},
}
