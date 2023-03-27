# APM HUB

APM Hub is a powerful tool designed to aggregate logs from various sources and provide a centralized location for querying them.

It is able to integrate with the following sources:

- Elastic search
- Files
- Kubernetes

## Documentation

Read the documentation at [https://docs.flanksource.com/apm-hub/overview/](https://docs.flanksource.com/apm-hub/overview/)

## Samples

Check out the samples directory for example configurations.

## CLI

**Usage**

```bash
apm-hub [command]
```

**Available Commands:**

```
  completion  Generate the autocompletion script for the specified shell
  help        Help about any command
  operator    Start the kubernetes operator
  serve       Start the for querying the logs
  version     Print the version of apm-hub

Flags:
      --db string             Connection string for the postgres database (default "DB_URL")
      --db-log-level string    (default "warn")
      --db-migrations         Run database migrations
      --db-schema string       (default "public")
  -h, --help                  help for apm-hub
      --json-logs             Print logs in json format to stderr
  -v, --loglevel count        Increase logging level

Use "apm-hub [command] --help" for more information about a command.
```
