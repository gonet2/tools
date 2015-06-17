package main

/*
 * This tool is upload configuration to etcd!
 */
import (
	"fmt"
	"github.com/codegangsta/cli"
	"github.com/coreos/go-etcd/etcd"
	"io/ioutil"
	"os"
	"os/signal"
	"path"
	"path/filepath"
	"strings"
	"syscall"
	"time"
)

var (
	configFlags = []cli.Flag{
		cli.StringFlag{
			Name:  "root",
			Value: "/config",
			Usage: "config node key",
		},
		cli.StringFlag{
			Name:  "addr",
			Value: "http://localhost:2379",
			Usage: "etcd host address",
		},

		cli.StringFlag{
			Name:  "file",
			Value: "",
			Usage: "upload single file (file name is key and context is value)",
		},
		cli.StringFlag{
			Name:  "timeout",
			Value: "5",
			Usage: "dial timeout",
		},
	}

	numbersFlags = []cli.Flag{
		cli.StringFlag{
			Name:  "root",
			Value: "/numbers",
			Usage: "root of numbers",
		},
		cli.StringFlag{
			Name:  "pattern",
			Value: "/*.csv",
			Usage: "upload matched files",
		},
		cli.StringFlag{
			Name:  "addr",
			Value: "http://localhost:2379",
			Usage: "etcd host address",
		},
		cli.StringFlag{
			Name:  "dir",
			Value: "",
			Usage: "upload dir to etcd ",
		},
		cli.StringFlag{
			Name:  "timeout",
			Value: "5",
			Usage: "dial timeout",
		},
	}
)

func getEtcdClient(c *cli.Context) (*etcd.Client, string) {
	client := etcd.NewClient(strings.Split(c.String("addr"), ","))
	client.SetDialTimeout(time.Duration(c.Int("timeout")) * time.Second)
	root := c.String("root")
	if !strings.HasPrefix(root, "/") {
		root = ("/" + root)
	}
	return client, root
}

func uploadFiles(c *cli.Context) {
	dir := filepath.Clean(strings.TrimSpace(c.String("dir")))
	if dir == "" {
		fmt.Printf("dir cant not be empty \n")
		cli.ShowAppHelp(c)
		os.Exit(0)
	}

	client, root := getEtcdClient(c)
	fmt.Printf("uploading files : \n")
	err := filepath.Walk(dir, func(filename string, fi os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if fi.IsDir() {
			files, err := filepath.Glob(filename + c.String("pattern"))
			if err != nil {
				fmt.Printf("Read files error  err:%v  dir:%v pattern:%v", err, c.String("dir"), c.String("pattern"))
				os.Exit(0)
			}
			prefix := root + strings.Replace(filename, dir, "", -1)
			for _, v := range files {
				key, value := readFile(v)
				_, err := client.Set(prefix+"/"+key, string(value), 0)
				if err != nil {
					fmt.Printf("upload config failed  file: %v err:%v\n", path.Base(v), err)
					os.Exit(0)
				}
				fmt.Printf("uploaded file %v  \n", filepath.Base(v))
			}
		}
		return nil
	})

	if err != nil {
		fmt.Printf("upload file error %v!", err)
		return
	}
}

func uploadConfig(c *cli.Context) {
	// check flag
	filename := strings.TrimSpace(c.String("file"))
	if filename == "" {
		fmt.Printf("file cant not be empty \n")
		cli.ShowAppHelp(c)
		os.Exit(0)
	}

	// init etcd
	client, root := getEtcdClient(c)
	key, value := readFile(filename)
	fmt.Printf("uploading config file: %v\n", path.Base(filename))
	_, err := client.Set(root+"/"+key, string(value), 0)
	if err != nil {
		fmt.Printf("upload config failed  file: %v err:%v\n", path.Base(filename), err)
		os.Exit(0)
	}
	fmt.Printf("uploading config success file: %v \n", path.Base(filename))
}

func readFile(filename string) (base string, value []byte) {
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		fmt.Printf("no such file or directory: %s \n", filename)
		os.Exit(0)
	}

	file, err := os.OpenFile(filename, os.O_RDONLY, os.ModeType)
	config, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Printf("read file err:%v\n", err)
		os.Exit(0)
	}
	defer file.Close()
	return filepath.Base(filename), config
}

func main() {
	myApp := cli.NewApp()
	myApp.Name = "upload_numbers"
	myApp.Usage = ""
	myApp.Version = "0.0.1"
	myApp.Commands = []cli.Command{
		{Name: "config",
			Usage:       "use it to see a description",
			Description: "upload single file into etcd",
			Action:      uploadConfig,
			Flags:       configFlags,
		},
		{Name: "numbers",
			Usage:       "use it to see a description",
			Description: "upload dir  into etcd",
			Action:      uploadFiles,
			Flags:       numbersFlags,
		},
	}
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	myApp.Run(os.Args)
}
