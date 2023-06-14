# Exporter

This µservice runs a sidecar to the `build` and `helper` containers and shares their volume mount.

At the end of each CI job, there is an `after_script` that echos some variables into a `.json` file:

```sh
after_script:
  - mkdir -p /builds/.metrics
  - echo '{"started":"'"$CI_JOB_STARTED_AT"'","status":"'"CI_JOB_STATUS"'","project":"'"CI_PROJECT_PATH"'","name":"'"$CI_JOB_NAME"'"}' > /builds/.metrics/metrics.json
```

The `Exporter` µservice will wait for this file and then send it to the [Reciever](/receiver/README.md) µservice which will export the metrics.

## Variables

| Name | Description |
| --- | --- |
| `RECEIVER_SERVICE` | The url of the `Receiver` µservice |
| `METRIC_PATH` | The path of the `.json` file to send |
