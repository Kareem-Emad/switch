# switch

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

## Features

lama as7a isa, I need to sleep