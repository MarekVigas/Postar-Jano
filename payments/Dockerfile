FROM python:3.11-slim

WORKDIR /app

COPY requirements.txt requirements.txt
RUN pip install -r requirements.txt

COPY cron.yaml cron.yaml
COPY payments.py .

CMD ["yacron", "-c" "/app/cron.yaml"]