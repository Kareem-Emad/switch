
import logging
import sys
from faktory import Worker
from utils import process_job
from constants import QUEUE_NAME, QUEUE_JOB_TYPE, WORKER_CONCURRENCY

logging.getLogger().setLevel(logging.INFO)
logger = logging.getLogger()

stream = logging.StreamHandler(sys.stdout)
stream.setLevel(logging.INFO)
logger.addHandler(stream)

if __name__ == "__main__":

    w = Worker(queues=[QUEUE_NAME], concurrency=WORKER_CONCURRENCY)
    w.register(QUEUE_JOB_TYPE, process_job)
    logger.info('[Switch_Consumer] Sucessfully registered consumer on task queue, waiting for new jobs ....')
    w.run()
