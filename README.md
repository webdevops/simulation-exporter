Simulation Exporter
===================

[![license](https://img.shields.io/github/license/webdevops/simulation-exporter.svg)](https://github.com/webdevops/simulation-exporter/blob/master/LICENSE)
[![Docker](https://img.shields.io/badge/docker-webdevops%2Fsimulation--exporter-blue.svg?longCache=true&style=flat&logo=docker)](https://hub.docker.com/r/webdevops/simulation-exporter/)
[![Docker Build Status](https://img.shields.io/docker/build/webdevops/simulation-exporter.svg)](https://hub.docker.com/r/webdevops/simulation-exporter/)

Prometheus exporter for simulated metrics (eg. testing)

Configuration
-------------

Normally no configuration is needed but can be customized using environment variables.

| Environment variable | DefaultValue                    | Description                              |
|----------------------|---------------------------------|------------------------------------------|
| `SCRAPE_TIME`        | `5s`                            | Time (time.Duration) between generations |
| `SERVER_BIND`        | `:8080`                         | IP/Port binding                          |
| `CONFIG`             | `/app/config/node_exporter.yml` | Configuration file (yaml)                |
