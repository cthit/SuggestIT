import os

POSTGRES_USER = os.environ.get("SUGGESTIT_POSTGRES_USER", 'postgres')
POSTGRES_PASSWORD = os.environ.get("SUGGESTIT_POSTGRES_PASSWORD", 'password')
POSTGRES_HOST = os.environ.get("SUGGESTIT_POSTGRES_HOST", 'localhost')
POSTGRES_PORT = os.environ.get("SUGGESTIT_POSTGRES_PORT", '5432')
POSTGRES_DB = os.environ.get("SUGGESTIT_POSTGRES_DB", 'postgres')

GOTIFY_URL = os.environ.get("GOTIFY_URL",'http://localhost:1337')
GOTIFY_PRE_SHARED_KEY = os.environ.get("GOTIFY_PRE_SHARED-KEY",'123abc')
