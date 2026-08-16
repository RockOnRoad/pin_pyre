FROM python:3.12-slim

WORKDIR /app

RUN pip install --no-cache-dir requests

COPY website/ .

CMD ["python", "client.py"]