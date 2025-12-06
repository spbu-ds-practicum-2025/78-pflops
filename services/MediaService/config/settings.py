import os
from dataclasses import dataclass

@dataclass
class Settings:
    GRPC_PORT: int = int(os.getenv("GRPC_PORT", 50051))
    MINIO_ENDPOINT: str = os.getenv("MINIO_ENDPOINT", "minio:9000")
    MINIO_ACCESS_KEY: str = os.getenv("MINIO_ACCESS_KEY", "minioadmin")
    MINIO_SECRET_KEY: str = os.getenv("MINIO_SECRET_KEY", "minioadmin")
    MINIO_BUCKET: str = os.getenv("MINIO_BUCKET", "media-service")

settings = Settings()