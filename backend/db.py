from datetime import datetime
from uuid import UUID

import config

from pony.orm import Database, PrimaryKey, Required, Optional, db_session

db = Database()


class Suggestion(db.Entity):
    id = PrimaryKey(UUID, auto=True)
    timestamp = Required(datetime)
    title = Required(str)
    text = Required(str)
    author = Optional(str)


db.bind(
    provider='postgres',
    user=config.POSTGRES_USER,
    password=config.POSTGRES_PASSWORD,
    host=config.POSTGRES_HOST,
    port=config.POSTGRES_PORT,
    database=config.POSTGRES_DB
)

db.generate_mapping(create_tables=True)

@db_session
def suggestion_to_json(suggestion):
    return {
        'id':str(suggestion.id),
        'title':suggestion.title,
        'timestamp':str(suggestion.timestamp),
        'text':suggestion.text,
        'author':suggestion.author
    }

