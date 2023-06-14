# CI Jobs Metrics

There was a need to export metrics about the usage of our CI templates for KPI. This µservice will help us understand the state of our CI templates and where can it be improved.

One way this can be achieved is by exporting Prometheus metrics and create a KPI dashboard that will show the overall usage of our CI templates.

## Metrics Example

Registering a new CI job will result in the following metrics:

```prometheus

```

## Sending Metrics

This µservice waits for `API` requests to register metrics:

```cmd
curl -G -d started=2023-05023T13:01:00Z -d status=failed -d project=dummy-project -d name=dotnet_build http://localhost:80/steps
```

This will register the new CI job (if it's the first time its used) and increase it accordingly.
