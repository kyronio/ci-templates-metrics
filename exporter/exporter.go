package exporter

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

// RunExporter waits for a .json file located in $METRIC_PATH to be created.
// Then send it to the given $RECEIVER_SERVICE host.
func RunExporter(path, receiverHost string) {
	metric := tryReadFile(path)
	sendMetricToServer(receiverHost, metric)
}

// tryReadFile waits for the metric .json file to be written and then returns it
func tryReadFile(filePath string) []byte {
	fmt.Printf("Waiting for %s to be created.", filePath)
	for {
		metric, err := ioutil.ReadFile(filePath)
		if err == nil {
			fmt.Println("Received metric: " + string(metric))
			return metric
		}
	}
}

// sendMetricToServer receives the metric and then sends it to the received receiver host
func sendMetricToServer(host string, metric []byte) {
	var result map[string]interface{}
	if err := json.Unmarshal(metric, &result); err != nil {
		panic(err)
	}

	req, _ := http.NewRequest("GET", os.Getenv("RECEIVER_SERVICE")+"/steps", nil)
	query := req.URL.Query()
	for k := range result {
		value := replaceValueChar(result[k].(string))
		query.Add(k, value)
	}
	req.URL.RawQuery = query.Encode()
	http.DefaultClient.Do(req)
}

func replaceValueChar(value string) string {
	value = strings.Replace(value, ":", "-", -1)
	value = strings.Replace(value, " ", "_", -1)
	value = strings.Replace(value, "/", "-", -1)
	return value
}
