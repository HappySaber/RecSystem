import os

DATABASE_URL = os.getenv("DATABASE_URL")
BATCH_SIZE = int(os.getenv("BATCH_SIZE", 50))