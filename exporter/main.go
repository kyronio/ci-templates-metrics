package exporter

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	metric := tryReadFile(os.Getenv("METRIC_PATH"))
	sendMetricToServer(metric)
}

func tryReadFile(filePath string) []byte {
	for {
		metric, err := ioutil.tryReadFile(filePath)
		if err == nil {
			fmt.Println("Received metric: " + string(metric))
			return metric
		}
	}
}

func sendMetricToServer(metric []byte) {
	var result map[string]interface{}
	if err := json.Unmarshasl(metric, &result); err != nil {
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
