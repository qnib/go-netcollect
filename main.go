package main

import (
	"context"
	"fmt"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"regexp"
	"os"
	"time"
	"encoding/json"
	"sync"
)

var (
	ctx = context.Background()
)


func main() {
	reg := ".*"
	if len(os.Args) > 1 {
		reg = os.Args[1]
	}
	cli, err := client.NewEnvClient()
	if err != nil {
		panic(err)
	}

	networks, err := cli.NetworkList(ctx, types.NetworkListOptions{})
	if err != nil {
		panic(err)
	}
	netReg, err := regexp.Compile(reg)
	epReg, err := regexp.Compile(".*_testnet")
	if err != nil {
		panic(err)
	}
	var wg sync.WaitGroup
	backChannel := make(chan *types.StatsJSON)
	for _, network := range networks {
		if netReg.MatchString(network.Name) {
			net, err := cli.NetworkInspect(ctx, network.ID, types.NetworkInspectOptions{})
			if err != nil {
				//TODO: Panic is to much
				panic(err)
			}
			fmt.Printf("### %s %v\n", net.Name, net.Containers)
			for cntId, _ := range net.Containers {
				if epReg.MatchString(cntId) {
					continue
				}
				cnt, err := cli.ContainerInspect(ctx, cntId)
				if err != nil {
					//TODO: Panic is to much
					panic(err)
				}
				wg.Add(1)
				go startCollector(cli, &wg, backChannel, cntId, cnt)
			}
		}
	}
	wg.Wait()
	fmt.Println("WG done")
}

func startCollector(cli *client.Client, wg *sync.WaitGroup, data chan *types.StatsJSON, cntId string, cnt types.ContainerJSON) {
	_ = cnt
	fmt.Printf("Start Collector for: %s\n", cntId)
	containerStats, err := cli.ContainerStats(ctx, cntId, true)
	responseBody := containerStats.Body
	//defer responseBody.Close()
	defer wg.Done()
	if err != nil {
		return
	}
	var statsJSON *types.StatsJSON
	dec := json.NewDecoder(responseBody)

	timer := time.NewTicker(1000 * time.Millisecond)
	for {
		select {
		case <-timer.C:
			if err := dec.Decode(&statsJSON); err != nil {
				return
			}
			if statsJSON != nil {
				fmt.Printf("%v\n", statsJSON)
				data <- statsJSON
			}
		case <-ctx.Done():
			return
		}
	}
	fmt.Printf("-> %s startCollector finished\n", cntId)
}
