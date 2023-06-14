# Receiver

This µservice receives `HTTP` requests and increases the requested metrics.

## Receiving Metrics

This µservice waits for `HTTP` requests to register metrics:

```cmd
curl -G -d started=2023-05023T13:01:00Z -d status=failed -d project=dummy-project -d name=dotnet_build http://localhost:80/steps
```

This will register the new CI job (if it's the first time its used) and increase it accordingly.
