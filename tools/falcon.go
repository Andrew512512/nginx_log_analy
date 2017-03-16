package tools

import (
	"bytes"
	"encoding/json"
	"fmt"

	"time"

	"net/http"

	"github.com/davecgh/go-spew/spew"
	"xsbPro/log"
)

type FalconMessage struct {
	EndPoint    string `json:"endpoint"`    // www.exam
	Metric      string `json:"metric"`
	Timestamp   int    `json:"timestamp"`
	Step        int    `json:"step"`
	Value       int    `json:"value"`
	CounterType string `json:"counterType"` // GAUGE
	Tags        string `json:"tags"`
}

func new_FalconMessage(endpoint, metric string, timestamp, step, value int) *FalconMessage {
	msg := &FalconMessage{
		EndPoint:    endpoint,
		Metric:      metric,
		Timestamp:   timestamp,
		Step:        step,
		Value:       value,
		CounterType: "GAUGE",
	}
	return msg
}

func UploadOnce(endPoint, prefix string) {
	now := time.Now()
	fmt.Println("\n-------------------- Start upload msg... (" + now.Format(time.RFC3339) + "): --------------------")
	timestamp := int(now.Unix())
	messages := []*FalconMessage{}
	for _, info := range (Cache.Items()) {
		log.InfoF("%v", *info)
		messages = append(messages, new_FalconMessage(endPoint, prefix + "_allMsg_" + info.URL, timestamp, 60, info.Success + info.Fail))
		messages = append(messages, new_FalconMessage(endPoint, prefix + "_failMsg_" + info.URL, timestamp, 60, info.Fail))
		messages = append(messages, new_FalconMessage(endPoint, prefix + "_delayMsg_" + info.URL, timestamp, 60, info.Delay))
	}

	//清零缓存
	Cache.Reset()

	json_bs, err := json.Marshal(messages)
	if err != nil {
		fmt.Println("[ERR] marshal err: ", err)
		spew.Dump(messages)
		return
	}

	req, _ := http.NewRequest("POST", "http://127.0.0.1:1988/v1/push", bytes.NewReader(json_bs))
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.SysF("Post data err:%s", err.Error())
		return
	}
	if resp.StatusCode == http.StatusOK {
		log.Trace("[OK] post success")
	} else {
		log.SysF("Post resp err:%s", err.Error())
		spew.Dump(resp)
	}

	fmt.Println("--------------------------------------------------------------------------------")
}
