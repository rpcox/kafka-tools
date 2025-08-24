package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
	"time"
)

func main() {
	fh, err := os.Open("data/consumer.stats")
	if err != nil {
		log.Fatal(err)
	}

	data, err := io.ReadAll(fh)
	if err != nil {
		log.Fatal(err)
	}

	var stats LibRdKafkaStats
	err = json.Unmarshal(data, &stats)
	if err != nil {
		log.Fatal(err)
	}

	sformat := "%20s: %s\n"
	dformat := "%20s: %d\n"
	fmt.Printf(sformat, "Name", stats.Name)
	fmt.Printf(sformat, "ClientID", stats.ClientId)
	fmt.Printf(sformat, "Type", stats.Type)
	ts := time.Unix(stats.Ts/1000000, (stats.Ts%1000000)/1000000000) // Convert microseconds to seconds and remaining to nanoseconds
	fmt.Printf("%20s: %d (%s)\n", "Ts", stats.Ts, ts.UTC().Format("2006-01-02 15:04:05 MST"))
	ts = time.Unix(stats.Time, 0)
	fmt.Printf("%20s: %d (%s)\n", "Time", stats.Time, ts.UTC().Format("2006-01-02 15:04:05 MST"))
	fmt.Printf("%20s: %d (%s)\n", "Age", stats.Age, time.Duration((stats.Age/100000)*time.Hour.Microseconds()))
	fmt.Printf(dformat, "ReplyQ", stats.ReplyQ)
	fmt.Printf(dformat, "MsgCnt", stats.MsgCnt)
	fmt.Printf(dformat, "MsgSize", stats.MsgSize)
	fmt.Printf(dformat, "MsgMax", stats.MsgMax)
	fmt.Printf(dformat, "MsgSizeMax", stats.MsgSizeMax)
	fmt.Printf(dformat, "Tx", stats.Tx)
	fmt.Printf(dformat, "TxBytes", stats.TxBytes)
	fmt.Printf(dformat, "Rx", stats.Rx)
	fmt.Printf(dformat, "RxBytes", stats.RxBytes)
	fmt.Printf(dformat, "TxMsgs", stats.TxMsgs)
	fmt.Printf(dformat, "TxMsgsBytes", stats.TxMsgBytes)
	fmt.Printf(dformat, "RxMsgs", stats.RxMsgs)
	fmt.Printf(dformat, "RxMsgsBytes", stats.RxMsgBytes)
	fmt.Printf(dformat, "SimpleCnt", stats.SimpleCnt)
	fmt.Printf(dformat, "MetaDataCacheCnt", stats.MetaDataCacheCnt)
	fmt.Printf("\n%20s\n", "BROKERS:")
	for k, v := range stats.Brokers {
		re := regexp.MustCompile(`^GroupCoordinator`)
		if v.NodeId > -1 && !re.Match([]byte(k)) { // skip bootstrap nodes and GroupCoordinator for now
			fmt.Printf("++ %20s\n", k)
			fmt.Printf(dformat, "NodeId", v.NodeId)
			fmt.Printf(sformat, "State", v.State)
			// Just check we can hit the maps
			fmt.Printf(dformat, "Req[Fetch]", v.Req["Fetch"])
			fmt.Printf(dformat, "IntLatency[min]", v.IntLatency["min"])
			fmt.Printf(dformat, "OutBufLatency[min]", v.OutBufLatency["min"])
			fmt.Printf("%20s: %d (%s)\n", "RTT[avg]", v.RTT["avg"], time.Duration((v.RTT["avg"]/100000)*time.Hour.Microseconds()))
			fmt.Printf("%20s: %d (%s)\n", "Throttle[avg]", v.Throttle["avg"], time.Duration((v.Throttle["avg"]/100000)*time.Hour.Microseconds()))
			fmt.Printf("\n%20s\n", " BROKER PARTITIONS:")
			for _, v := range v.TopPars {
				fmt.Printf("%20s: %d\n", v.Topic, v.Partition)
			}
			fmt.Println()
		}
	}
	fmt.Printf("\n%20s\n", "TOPICS:")
	for k, v := range stats.Topics {
		fmt.Printf("++ %20s\n", k)
		for p, v := range v.Partitions {
			if p != "-1" {
				fmt.Printf(dformat, "P-"+p+" Consumer Lag: ", v.ConsumerLag)
			}
		}
	}

}
