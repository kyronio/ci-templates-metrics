# CI Jobs Metrics

There was a need to export metrics about the usage of our CI templates for KPI. This µservice will help us understand the state of our CI templates and where can it be improved.

One way this can be achieved is by exporting Prometheus metrics and create a KPI dashboard that will show the overall usage of our CI templates.

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

```sh
receiver
```

### Receiving Metrics

This µservice waits for `HTTP` requests to register metrics:

```cmd
curl -G -d started=2023-05023T13:01:00Z -d status=failed -d project=dummy-project -d name=dotnet_build http://localhost:80/steps
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

The `Exporter` µservice will wait for this file and then send it to the [Reciever](/README.md/#receiver) µservice which will export the metrics.
