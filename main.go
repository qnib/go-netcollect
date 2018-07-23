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
	// FIgure out which interface has the IP in net?
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

func assembleRes(cnt *docker.Container, stats *docker.Stats) (res string) {
	res = fmt.Sprintf("cntId:%s time:%s", cnt.ID, stats.Read.Format("2006-01-02T15:04:05.999999"))
	for iname, iface := range stats.Networks {
		res = fmt.Sprintf("%s %s.RxBytes:%d %s.TxBytes:%d", res, iname, iface.RxBytes, iname, iface.TxBytes)
	}
	return res
}