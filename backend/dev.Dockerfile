FROM python:3

RUN mkdir /usr/src/app

WORKDIR /usr/src/app

COPY requirements.txt .

RUN pip install --no-cache-dir -r requirements.txt

COPY . .

ENV SUGGESTIT_POSTGRES_USER suggestit
ENV SUGGESTIT_POSTGRES_PASSWORD password
ENV SUGGESTIT_POSTGRES_HOST db
ENV SUGGESTIT_POSTGRES_PORT 5432
ENV SUGGESTIT_POSTGRES_DB suggestit

EXPOSE 5000

CMD ["sh", "start.sh"]