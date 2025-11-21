FROM python:3.13.9-slim

WORKDIR /mediaservice

COPY requirements.txt .
RUN pip install --no-cache-dir -r requirements.txt

COPY server.py .
COPY minio_client.py .
COPY generated/ ./generated/

EXPOSE 50051
CMD ["python", "server.py"]