# conftest.py
import pytest
import sys
import os

# Добавляем путь к проекту
sys.path.insert(0, os.path.abspath(os.path.join(os.path.dirname(__file__), '..')))


@pytest.fixture(autouse=True)
def mock_environment_variables(monkeypatch):
    """Мокаем переменные окружения для тестов"""
    monkeypatch.setenv('GRPC_PORT', '50051')
    monkeypatch.setenv('MINIO_ENDPOINT', 'localhost:9000')
    monkeypatch.setenv('MINIO_ACCESS_KEY', 'test_access_key')
    monkeypatch.setenv('MINIO_SECRET_KEY', 'test_secret_key')
    monkeypatch.setenv('MINIO_BUCKET', 'test-bucket')


@pytest.fixture
def sample_media_metadata():
    """Фикстура с примером метаданных"""
    from internal.core.models import MediaMetadata
    
    return MediaMetadata(
        user_id='user123',
        mime_type='image/jpeg',
        file_name='test.jpg',
        upload_date='2023-01-01T12:00:00'
    )


@pytest.fixture
def sample_media_bytes():
    """Фикстура с примером байтов медиа"""
    return b'fake media content'