FROM python:3.13.9-slim

WORKDIR /app

# Установка netcat для healthcheck
RUN apt-get update && apt-get install -y netcat-openbsd && rm -rf /var/lib/apt/lists/*

COPY requirements.txt .
RUN pip install --no-cache-dir -r requirements.txt

COPY . .

# Устанавливаем PYTHONPATH чтобы Python видел корень проекта
ENV PYTHONPATH=/app

EXPOSE 50051
CMD ["python", "-u", "tests/run_tests.py"]