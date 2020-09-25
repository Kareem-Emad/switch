# switch

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Build Status:](https://github.com/Kareem-Emad/switch/workflows/Build/badge.svg)](https://github.com/Kareem-Emad/switch/actions)

Dynamic event managment service based on publisher/subscriber model, it allows your service call sequence to be
as dynamic as a DB change

## setup

To run server

```shell
make run
```

to run worker

```shell
virtualenv --python=python3 venv
venv/bin/pip install -r consumer/requirements.txt
venv/bin/python consumer/main.py
```

## Environment Variables

List of envs needs to be setup before starting the service

- `FAKTORY_URL` url for the faktory server to connect to
- `SWITCH_JWT_SECRET` secret used to verify signature of the publisher/sender
- `SWITCH_DB_NAME` name of the sqlite instance used
- `SWITCH_PRODUCTION_QUEUE` name of the queue used for switch jobs
- `SWITCH_PRODUCTION_QUEUE_NAMESPACE` the job type you want the consumer/producer to use
- `SWITCH_CONSUMER_CONCURRENCY_LEVEL` how many concurrent instances of consumer
- `SWITCH_SERVER_PORT` port used by the switch server

## Features

lama as7a isa, I need to sleep