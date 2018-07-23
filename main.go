package main

import (
	"context"
	"fmt"

	"errors"
	"regexp"
	"os"
	"sync"
	"github.com/fsouza/go-dockerclient"
	"log"
	"bytes"
	"time"
)

var (
	ctx = context.Background()
)


func Connect() *docker.Client {
	endpoint := "unix:///var/run/docker.sock"
	client, err := docker.NewClient(endpoint)
	if err != nil {
		panic(err)
	}
	return client
}

func main() {
	reg := ".*"
	if len(os.Args) > 1 {
		reg = os.Args[1]
	}
	cli := Connect()
	networks, err := cli.ListNetworks()
	if err != nil {
		panic(err)
	}
	netReg, err := regexp.Compile(reg)
	epReg, err := regexp.Compile(".*_testnet")
	if err != nil {
		panic(err)
	}
	var wg sync.WaitGroup
	for _, network := range networks {
		if netReg.MatchString(network.Name) {
			net, err := cli.NetworkInfo(network.ID)
			if err != nil {
				//TODO: Panic is to much
				panic(err)
			}
			fmt.Printf("### %s\n", net.Name)
			for cntId, ep := range net.Containers {
				if epReg.MatchString(cntId) {
					continue
				}
				cnt, err := cli.InspectContainer(cntId)
				if err != nil {
					//TODO: Panic is to much
					panic(err)
				}
				wg.Add(1)
				go startCollector(cli, &wg, ep, cntId, cnt)
			}
		}
	}
	wg.Wait()
	fmt.Println("## WG done")
}

func startCollector(cli *docker.Client, wg *sync.WaitGroup, net docker.Endpoint, cntId string, cnt *docker.Container) {
	_ = cnt
	fmt.Printf("## Start Collector for: %s\n", cntId)
	errChannel := make(chan error, 1)
	statsChannel := make(chan *docker.Stats)
	// Figure out which interface has the IP in net?
	ethFace, err := correlateIface(cli, cntId, net.IPv4Address)
	_ = err
	_ = ethFace
	opts := docker.StatsOptions{
		ID:     cntId,
		Stats:  statsChannel,
		Stream: true,
	}

	go func() {
		errChannel <- cli.Stats(opts)
	}()

	for {
		select {
		case stats, ok := <-statsChannel:
			if !ok {
				err := errors.New(fmt.Sprintf("## Bad response getting stats for container: %s", cntId))
				log.Println(err.Error())
				return
			}

			fmt.Printf("%s\n", assembleRes(cnt, stats))

		}
	}
	fmt.Printf("### -> %s startCollector finished\n", cntId)
}

func correlateIface(cli *docker.Client, cntId, ipv4 string) (iface string, err error) {
	config := docker.CreateExecOptions{
		Container:    cntId,
		AttachStdin:  false,
		AttachStdout: true,
		AttachStderr: true,
		Tty:          true,
		Cmd:          []string{"ip", "-o", "-4", "addr"},
	}
	execObj, err := cli.CreateExec(config)
	if err != nil {
		return
	}
	var stdout, stderr bytes.Buffer
	success := make(chan struct{})
	opts := docker.StartExecOptions{
		OutputStream: &stdout,
		ErrorStream:  &stderr,
		RawTerminal:  true,
		Success:      success,
	}
	go func() {
		if err := cli.StartExec(execObj.ID, opts); err != nil {
			panic(err)
		}
	}()
	<-success
	time.Sleep(5*time.Second)
	fmt.Printf("stdout:%v || stderr:%s\n",stdout.String(), stderr.String())
	return
}
func assembleRes(cnt *docker.Container, stats *docker.Stats) (res string) {
	res = fmt.Sprintf("cntId:%s time:%s", cnt.ID, stats.Read.Format("2006-01-02T15:04:05.999999"))
	for iname, iface := range stats.Networks {
		res = fmt.Sprintf("%s %s.RxBytes:%d %s.TxBytes:%d", res, iname, iface.RxBytes, iname, iface.TxBytes)
	}
	return res
}