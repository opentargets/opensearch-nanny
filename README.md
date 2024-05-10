# Opensearch Nanny

This is an Opensearch sidecar intended to be used to keep an eye on the health
of a cluster. It is designed to be run as a sidecar in a Kubernetes pod. The
sidecar will periodically check the health of the cluster and report the status
to a configurable endpoint.

The idea is to simplify the color-based monitoring into a regular 503/200 http
status code response.

## Configuration

The `config.toml` file contains a few configuration options, mostly self-explanatory.

```toml
[opensearch]
health_url = "http://localhost:9200/_cluster/health"
seconds_in_green_for_healthy = 5
ticker_interval = 1

[server]
address = "0.0.0.0"
port = 8080
log_level = "debug"
log_handler = "text"
```

## Building

To build the project, you can use the following command:

```bash
$ go build -o opensearch-nanny cmd/main.go
```

## Running

To run the project, you can use the following command:

```bash
$ ./opensearch-nanny
```

## Local Development

The project includes a `docker-compose.yml` file that can be used to run a local
cluster where you can test the nanny. To start the cluster, you can use the
following command:

```bash
$ docker-compose up
```

## Copyright

Copyright 2014-2024 EMBL - European Bioinformatics Institute, Genentech, GSK, MSD, Pfizer, Sanofi and Wellcome Sanger Institute

This software was developed as part of the Open Targets project. For more
information please see: http://www.opentargets.org

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
