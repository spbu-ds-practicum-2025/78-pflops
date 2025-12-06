import os
from minio import Minio
from minio.error import S3Error
from datetime import timedelta
import uuid
from typing import List, Optional
from config import settings

class MinioClient:
    def __init__(self):
        self.client = Minio(
            os.getenv("MINIO_ENDPOINT", "localhost:9000"),
            access_key=os.getenv("MINIO_ACCESS_KEY", "minioadmin"),
            secret_key=os.getenv("MINIO_SECRET_KEY", "minioadmin"),
            secure=False
        )
        self.bucket_name = os.getenv("MINIO_BUCKET", "media-service")
        self._ensure_bucket_exists()

    def _ensure_bucket_exists(self):
        """Создает бакет если он не существует"""
        try:
            if not self.client.bucket_exists(self.bucket_name):
                self.client.make_bucket(self.bucket_name)
                print(f"Бакет {self.bucket_name} создан")
        except S3Error as e:
            print(f"Ошибка создания бакета: {e}")

    def generate_media_id(self, user_id: str, file_name: str) -> str:
        """Генерация уникального ID для медиа файла"""
        unique_id = str(uuid.uuid4())
        return f"{user_id}/{unique_id}/{file_name}"

    def upload_media(self, user_id: str, file_bytes: bytes, mime_type: str, file_name: str) -> str:
        """Загрузка медиа файла в MinIO"""
        try:
            media_id = self.generate_media_id(user_id, file_name)
            
            # Сохраняем временно в файл для загрузки
            temp_path = f"/tmp/{media_id.replace('/', '_')}"
            os.makedirs(os.path.dirname(temp_path), exist_ok=True)
            
            with open(temp_path, "wb") as f:
                f.write(file_bytes)
            
            self.client.fput_object(
                self.bucket_name,
                media_id,
                temp_path,
                content_type=mime_type
            )
            
            # Удаляем временный файл
            os.remove(temp_path)
            
            return media_id
            
        except S3Error as e:
            print(f"Ошибка загрузки файла: {e}")
            raise

    def get_media(self, media_id: str) -> Optional[bytes]:
        """Получение медиа файла по ID"""
        try:
            response = self.client.get_object(self.bucket_name, media_id)
            file_data = response.read()
            response.close()
            response.release_conn()
            return file_data
        except S3Error as e:
            print(f"Ошибка получения файла: {e}")
            return None

    def delete_media(self, media_id: str, user_id: str) -> bool:
        """Удаление медиа файла"""
        try:
            # Проверяем принадлежность файла пользователю
            if not media_id.startswith(user_id + "/"):
                return False
                
            self.client.remove_object(self.bucket_name, media_id)
            return True
        except S3Error as e:
            print(f"Ошибка удаления файла: {e}")
            return False

    def list_user_media(self, user_id: str) -> List[dict]:
        """Список всех медиа файлов пользователя"""
        try:
            objects = self.client.list_objects(
                self.bucket_name, 
                prefix=user_id + "/",
                recursive=True
            )
            
            media_items = []
            for obj in objects:
                media_id = obj.object_name
                file_name = media_id.split('/')[-1]
                
                media_items.append({
                    'media_id': media_id,
                    'file_name': file_name,
                    'last_modified': obj.last_modified,
                    'size': obj.size
                })
            
            return media_items
        except S3Error as e:
            print(f"Ошибка получения списка файлов: {e}")
            return []

    def get_presigned_url(self, media_id: str, expiry_hours: int = 24) -> Optional[str]:
        """Генерация presigned URL для доступа к файлу"""
        try:
            url = self.client.presigned_get_object(
                self.bucket_name,
                media_id,
                expires=timedelta(hours=expiry_hours)
            )
            return url
        except S3Error as e:
            print(f"Ошибка генерации URL: {e}")
            return None

    def get_media_metadata(self, media_id: str) -> Optional[dict]:
        """Получение метаданных файла"""
        try:
            stat = self.client.stat_object(self.bucket_name, media_id)
            return {
                'media_id': media_id,
                'size': stat.size,
                'content_type': stat.content_type,
                'last_modified': stat.last_modified
            }
        except S3Error:
            return None