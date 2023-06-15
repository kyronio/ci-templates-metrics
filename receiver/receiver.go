package receiver

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

var (
	allStepsCounter           = map[string]*prometheus.CounterVec{}
	allStepsDurationHistogram = map[string]*prometheus.HistogramVec{}
	timeZone                  = "Asia/Jerusalem"
)

// handleRequest registers and increases metrics for each CI job request that arrives
func handleRequest(w http.ResponseWriter, r *http.Request) {
	project, status, name, duration := getParameters(r)

	splitJobName := strings.Split(name, "_")

	_, ok := allStepsCounter[fmt.Sprintf("ci_%s_%s_total", splitJobName[0], splitJobName[1])]
	// If key already registered
	if ok {
		increaseMetric(splitJobName[0], splitJobName[1], status, project, duration)
	} else {
		registerNewJob(splitJobName[0], splitJobName[1])
		increaseMetric(splitJobName[0], splitJobName[1], status, project, duration)
	}
}

// increaseMetric increasing both Counter and Histogram metrics for the received CI job
func increaseMetric(framework, step, status, project string, duration float64) {
	allStepsCounter[fmt.Sprintf("ci_%s_%s_total", framework, step)].WithLabelValues(status, project).Add(1)
	allStepsDurationHistogram[fmt.Sprintf("ci_%s_%s_duration", framework, step)].WithLabelValues(status).Observe(duration)
}

// getParameters fetches the parameters from the URL and returns them
func getParameters(r *http.Request) (string, string, string, float64) {
	status := r.URL.Query().Get("status")
	project := r.URL.Query().Get("project")
	name := r.URL.Query().Get("name")
	started := r.URL.Query().Get("started")

	startedParsed := strings.Replace(started, "T", " ", -1)
	startedParsed = strings.Replace(startedParsed, "Z", "", -1)
	startedParsed = strings.Replace(startedParsed, "-", ".", -1)
	startedParsed = strings.Replace(startedParsed, ":", ".", -1)
	startedAt, err := time.ParseInLocation("2006.01.02 15.04.05", startedParsed, time.Local)
	if err != nil {
		fmt.Printf("There was an error parsing the starting time %s\n%s", startedParsed, err)
	}
	duration := time.Now().Sub(startedAt).Seconds()

	fmt.Println(fmt.Printf("%s %s %f %s %s %s", time.Now(), startedAt, duration, status, project, name))
	return project, status, name, duration
}

// registerNewJob registers 2 metrics for each new CI job arrives to this µservice
func registerNewJob(framework, step string) {
	fmt.Println(fmt.Sprintf("Registering metrics for: %s-%s", framework, step))
	allStepsCounter[fmt.Sprintf("ci_%s_%s_total", framework, step)] = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: fmt.Sprintf("ci_%s_%s_total", framework, step),
		Help: fmt.Sprintf("The number of %s %s used", framework, step),
	},
		[]string{"status", "project"},
	)
	allStepsDurationHistogram[fmt.Sprintf("ci_%s_%s_duration", framework, step)] = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:    fmt.Sprintf("ci_%s_%s_duration", framework, step),
		Help:    fmt.Sprintf("The duration of %s %s uses", framework, step),
		Buckets: []float64{60, 120, 300, 600, 900, 1200},
	},
		[]string{"status"},
	)
}

func whileTrue() {
	for {
	}
}

func RunReceiver() {
	os.Setenv("TZ", timeZone)
	loc, _ := time.LoadLocation(timeZone)
	time.Local = loc // -> this is setting the global timezone

	http.HandleFunc("/jobs", handleRequest)
	http.Handle("/", promhttp.Handler())
	http.ListenAndServe(":80", nil)

	fmt.Printf("Exporting metrics on / using port 80\n")
	fmt.Printf("HTTP requests should be send to /jobs\n")

	go whileTrue()
}
