# test_service.py
import pytest
from unittest.mock import Mock, patch, MagicMock, mock_open
from datetime import datetime
import sys
import os

# Добавляем путь к проекту
sys.path.insert(0, os.path.abspath(os.path.join(os.path.dirname(__file__), '..')))

from internal.core.service import MediaService
from internal.core.models import MediaMetadata


class TestMediaService:
    
    @pytest.fixture
    def mock_minio_client(self):
        """Фикстура для мока MinioClient"""
        with patch('internal.core.service.MinioClient') as mock:
            client = Mock()
            mock.return_value = client
            yield client
    
    @pytest.fixture
    def media_service(self, mock_minio_client):
        """Фикстура для MediaService"""
        return MediaService()
    
    @pytest.fixture
    def sample_media_data(self):
        """Пример данных для загрузки медиа"""
        return {
            'user_id': 'user123',
            'file_bytes': b'test file content',
            'mime_type': 'image/jpeg',
            'file_name': 'test.jpg'
        }
    
    @patch('internal.core.service.datetime')
    def test_upload_media_success(self, mock_datetime, media_service, mock_minio_client, sample_media_data):
        """Тест успешной загрузки медиа"""
        # Arrange
        expected_media_id = 'user123/uuid-123/test.jpg'
        fixed_date = '2023-01-01T12:00:00'
        
        # Мокаем datetime.now()
        mock_now = Mock()
        mock_now.isoformat.return_value = fixed_date
        mock_datetime.now.return_value = mock_now
        
        mock_minio_client.upload_media.return_value = expected_media_id
        
        # Act
        result = media_service.upload_media(**sample_media_data)
        
        # Assert
        assert result == expected_media_id
        mock_minio_client.upload_media.assert_called_once_with(
            user_id=sample_media_data['user_id'],
            file_bytes=sample_media_data['file_bytes'],
            mime_type=sample_media_data['mime_type'],
            file_name=sample_media_data['file_name']
        )
        
        # Проверяем метаданные
        assert expected_media_id in media_service.media_metadata
        metadata = media_service.media_metadata[expected_media_id]
        assert isinstance(metadata, MediaMetadata)
        assert metadata.user_id == sample_media_data['user_id']
        assert metadata.mime_type == sample_media_data['mime_type']
        assert metadata.file_name == sample_media_data['file_name']
        assert metadata.upload_date == fixed_date
    
    def test_upload_media_minio_error(self, media_service, mock_minio_client, sample_media_data):
        """Тест ошибки при загрузке в MinIO"""
        # Arrange
        mock_minio_client.upload_media.side_effect = Exception("MinIO connection error")
        
        # Act & Assert
        with pytest.raises(Exception, match="MinIO connection error"):
            media_service.upload_media(**sample_media_data)
    
    def test_get_media_success(self, media_service, mock_minio_client):
        """Тест успешного получения медиа"""
        # Arrange
        media_id = 'user123/uuid-123/test.jpg'
        expected_data = b'file content'
        mock_minio_client.get_media.return_value = expected_data
        
        # Act
        result = media_service.get_media(media_id)
        
        # Assert
        assert result == expected_data
        mock_minio_client.get_media.assert_called_once_with(media_id)
    
    def test_get_media_not_found(self, media_service, mock_minio_client):
        """Тест получения несуществующего медиа"""
        # Arrange
        media_id = 'non-existent-id'
        mock_minio_client.get_media.return_value = None
        
        # Act
        result = media_service.get_media(media_id)
        
        # Assert
        assert result is None
    
    def test_delete_media_success(self, media_service, mock_minio_client):
        """Тест успешного удаления медиа"""
        # Arrange
        media_id = 'user123/uuid-123/test.jpg'
        user_id = 'user123'
        
        # Добавляем метаданные
        media_service.media_metadata[media_id] = MediaMetadata(
            user_id=user_id,
            mime_type='image/jpeg',
            file_name='test.jpg',
            upload_date='2023-01-01T12:00:00'
        )
        
        mock_minio_client.delete_media.return_value = True
        
        # Act
        result = media_service.delete_media(media_id, user_id)
        
        # Assert
        assert result is True
        mock_minio_client.delete_media.assert_called_once_with(media_id, user_id)
        assert media_id not in media_service.media_metadata
    
    def test_delete_media_wrong_user(self, media_service, mock_minio_client):
        """Тест удаления медиа другого пользователя"""
        # Arrange
        media_id = 'user123/uuid-123/test.jpg'
        wrong_user_id = 'user456'
        
        mock_minio_client.delete_media.return_value = False
        
        # Act
        result = media_service.delete_media(media_id, wrong_user_id)
        
        # Assert
        assert result is False
        assert media_id not in media_service.media_metadata
    
    def test_delete_media_not_found_in_minio(self, media_service, mock_minio_client):
        """Тест удаления медиа, которого нет в MinIO"""
        # Arrange
        media_id = 'user123/uuid-123/test.jpg'
        user_id = 'user123'
        
        # Добавляем метаданные
        media_service.media_metadata[media_id] = MediaMetadata(
            user_id=user_id,
            mime_type='image/jpeg',
            file_name='test.jpg',
            upload_date='2023-01-01T12:00:00'
        )
        
        mock_minio_client.delete_media.return_value = False
        
        # Act
        result = media_service.delete_media(media_id, user_id)
        
        # Assert
        assert result is False
        # Метаданные должны остаться, так как удаление не удалось
        assert media_id in media_service.media_metadata
    
    def test_list_user_media_success(self, media_service, mock_minio_client):
        """Тест получения списка медиа пользователя"""
        # Arrange
        user_id = 'user123'
        expected_items = [
            {'media_id': f'{user_id}/uuid-1/test1.jpg', 'file_name': 'test1.jpg', 'size': 1024},
            {'media_id': f'{user_id}/uuid-2/test2.png', 'file_name': 'test2.png', 'size': 2048}
        ]
        
        mock_minio_client.list_user_media.return_value = expected_items
        
        # Act
        result = media_service.list_user_media(user_id)
        
        # Assert
        assert result == expected_items
        mock_minio_client.list_user_media.assert_called_once_with(user_id)
    
    def test_list_user_media_empty(self, media_service, mock_minio_client):
        """Тест получения пустого списка медиа"""
        # Arrange
        user_id = 'user123'
        mock_minio_client.list_user_media.return_value = []
        
        # Act
        result = media_service.list_user_media(user_id)
        
        # Assert
        assert result == []
    
    def test_get_presigned_url_success(self, media_service, mock_minio_client):
        """Тест получения presigned URL"""
        # Arrange
        media_id = 'user123/uuid-123/test.jpg'
        expected_url = 'http://localhost:9000/media-service/user123/uuid-123/test.jpg?token=abc'
        mock_minio_client.get_presigned_url.return_value = expected_url
        
        # Act
        result = media_service.get_presigned_url(media_id)
        
        # Assert
        assert result == expected_url
        mock_minio_client.get_presigned_url.assert_called_once_with(media_id)
    
    def test_get_presigned_url_not_found(self, media_service, mock_minio_client):
        """Тест получения URL для несуществующего медиа"""
        # Arrange
        media_id = 'non-existent'
        mock_minio_client.get_presigned_url.return_value = None
        
        # Act
        result = media_service.get_presigned_url(media_id)
        
        # Assert
        assert result is None
    
    def test_get_media_metadata_success(self, media_service):
        """Тест получения метаданных"""
        # Arrange
        media_id = 'user123/uuid-123/test.jpg'
        expected_metadata = MediaMetadata(
            user_id='user123',
            mime_type='image/jpeg',
            file_name='test.jpg',
            upload_date='2023-01-01T12:00:00'
        )
        
        media_service.media_metadata[media_id] = expected_metadata
        
        # Act
        result = media_service.get_media_metadata(media_id)
        
        # Assert
        assert result == expected_metadata
    
    def test_get_media_metadata_not_found(self, media_service):
        """Тест получения метаданных несуществующего медиа"""
        # Act
        result = media_service.get_media_metadata('non-existent')
        
        # Assert
        assert result is None
    
    def test_concurrent_access_simulation(self, media_service, mock_minio_client):
        """Тест симуляции конкурентного доступа к сервису"""
        # Arrange
        user_ids = ['user1', 'user2', 'user3']
        file_names = ['file1.jpg', 'file2.png', 'file3.pdf']
        
        # Мокаем upload_media для возврата предсказуемых ID
        def mock_upload(user_id, file_bytes, mime_type, file_name):
            return f"{user_id}/mock-id/{file_name}"
        
        mock_minio_client.upload_media.side_effect = mock_upload
        
        # Act - симулируем одновременную загрузку от разных пользователей
        uploaded_ids = []
        for user_id, file_name in zip(user_ids, file_names):
            media_id = media_service.upload_media(
                user_id=user_id,
                file_bytes=b'test content',
                mime_type='image/jpeg',
                file_name=file_name
            )
            uploaded_ids.append(media_id)
        
        # Assert
        assert len(uploaded_ids) == 3
        assert len(media_service.media_metadata) == 3
        
        # Проверяем, что каждый пользователь имеет доступ только к своим файлам
        for i, (user_id, media_id) in enumerate(zip(user_ids, uploaded_ids)):
            # Пользователь может получить свои файлы
            mock_minio_client.get_media.return_value = b'content'
            result = media_service.get_media(media_id)
            assert result == b'content'
            
            # Метаданные доступны
            metadata = media_service.get_media_metadata(media_id)
            assert metadata.user_id == user_id
            assert metadata.file_name == file_names[i]