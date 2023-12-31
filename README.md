# CI Jobs Metrics

There was a need to export metrics about the usage of our CI templates for KPI. This µservice will help us understand the state of our CI templates and where can it be improved.

One way this can be achieved is by exporting Prometheus metrics and create a KPI dashboard that will show the overall usage of our CI templates.

## Final Result

Our final ci templates KPI dashboard looks like this:

![KPI dashboard gif](/docs/images/kpi-example.gif)

## Metrics Example

Registering a new CI job will result in the following metrics:

```prometheus
# HELP ci_dotnet_build_total The number of dotnet build used
# TYPE ci_dotnet_build_total counter
ci_dotnet_build_total{project="dummy-project",status="success"} 3
ci_dotnet_build_total{project="dummy-project2",status="success"} 5
ci_dotnet_build_total{project="dummy-project2",status="failed"} 12
ci_dotnet_build_total{project="dummy-project3",status="success"} 8

# HELP ci_helm_restore_duration The duration of the helm restore usage
# TYPE ci_helm_restore_duration histogram
ci_helm_restore_duration_bucket{status="success",le="60"} 60
ci_helm_restore_duration_bucket{status="success",le="120"} 67
ci_helm_restore_duration_bucket{status="success",le="300"} 68
ci_helm_restore_duration_bucket{status="success",le="600"} 72
ci_helm_restore_duration_bucket{status="success",le="900"} 74
ci_helm_restore_duration_bucket{status="success",le="1200"} 75
ci_helm_restore_duration_bucket{status="success",le="+Inf"} 76
ci_helm_restore_duration_sum_{status="success"} 7005.646230537002
```

## How Does It Work

1. User runs a CI job

1. Jobs scripts are completed

1. `after_scripts` saves the metrics into a `.json` file in the mounted volume

1. Sidecar reads the `.json` file and sends it to the `Receiver µservice`

1. The `Receiver µservice` increases the correct `metrics` of the complete job

## Receiver

This µservice receives `HTTP` requests and increases the requested metrics.

```cmd
docker run -p 80:80 <image-name> receiver
```

### Receiving Metrics

This µservice waits for `HTTP` requests to register metrics:

```cmd
curl -G -d started=2023-05-23T13:01:00Z -d status=failed -d project=dummy-project -d name=dotnet_build http://localhost:80/jobs
```

This will register the new CI job (if it's the first time its used) and increase it accordingly.

## Exporter

This µservice runs a sidecar to the `build` and `helper` containers and shares their volume mount.

```sh
exporter <file-path> <receiver-hostname>
```

At the end of each CI job, there is an `after_script` that echos some variables into a `.json` file:

```sh
after_script:
  - mkdir -p /builds/.metrics
  - echo '{"started":"'"$CI_JOB_STARTED_AT"'","status":"'"CI_JOB_STATUS"'","project":"'"CI_PROJECT_PATH"'","name":"'"$CI_JOB_NAME"'"}' > /builds/.metrics/metrics.json
```

The `Exporter` µservice will wait for this file and then send it to the [Receiver](/README.md/#receiver) µservice which will export the metrics.

## See Yourself

### Run the Receiver

First make sure your [Receiver](/README.md/#receiver) is up and running.

### Simulate CI Jobs

Now we mimic a CI jobs usage scenario by running the [curl.bash](/docs/curl.sh) file which will mimic 5 minutes of `dotnet` CI usage.

- Please make sure to override the `starting-time` value! Or all the jobs will arrive with an infinite duration.

### Prometheus

Collect the metrics using a [Prometheus](https://prometheus.io/) instance:

![Prometheus Query](/docs/images/prometheus-dotnet-build-example.PNG)

Or run a new instance using:

```cmd
docker run -p 9090:9090 -v path\to\prometheus.yml:/etc/prometheus/prometheus.yml prom/prometheus
```

### Grafana

Next step is to connect your [Prometheus](https://prometheus.io/) instance to your [Grafana](https://grafana.com/) and create a [dashboard](/dashboards/ci-jobs-status.json):

![Grafana Dashboard](/docs/images/grafana-dotnet-restore-example.PNG)

Or run a new Grafana instance by:

```cmd
docker run -d --name=grafana -p 3000:3000 grafana/grafana
```

#### General Section

General information about the usage generally.

![general row](/docs/images/general-section.PNG)

#### Per Job Section

Information about each CI job.

![per job](/docs/images/per-job.PNG)
