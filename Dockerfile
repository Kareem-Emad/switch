FROM videra/gopyzmq-alpine:1.14-3.8


WORKDIR /app
COPY . /app


RUN  pip3 install -r /app/consumer/requirements.txt

# ENTRYPOINT ["python3", "/app/consumer/main.py"]
# ENTRYPOINT ["make", "run"]
