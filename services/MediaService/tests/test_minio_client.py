# test_minio_client.py
import pytest
from unittest.mock import Mock, patch, MagicMock, mock_open, ANY
import os
import sys

sys.path.insert(0, os.path.abspath(os.path.join(os.path.dirname(__file__), '..')))

from internal.storage.minio_client import MinioClient


class TestMinioClient:
    
    @pytest.fixture
    def mock_minio_lib(self):
        """Фикстура для мока библиотеки minio"""
        with patch('internal.storage.minio_client.Minio') as mock_minio_class:
            mock_client = Mock()
            mock_minio_class.return_value = mock_client
            yield mock_client
    
    @pytest.fixture
    def minio_client(self, mock_minio_lib):
        """Фикстура для MinioClient с моком библиотеки"""
        with patch.dict(os.environ, {
            'MINIO_ENDPOINT': 'localhost:9000',
            'MINIO_ACCESS_KEY': 'test_key',
            'MINIO_SECRET_KEY': 'test_secret',
            'MINIO_BUCKET': 'test-bucket'
        }):
            client = MinioClient()
            client.client = mock_minio_lib  # Заменяем реальный клиент на мок
            return client
    
    @patch('uuid.uuid4')
    def test_generate_media_id(self, mock_uuid, minio_client):
        """Тест генерации ID медиа"""
        # Arrange
        mock_uuid.return_value = '123e4567-e89b-12d3-a456-426614174000'
        
        # Act
        media_id = minio_client.generate_media_id('user123', 'test.jpg')
        
        # Assert
        assert media_id == 'user123/123e4567-e89b-12d3-a456-426614174000/test.jpg'
    
    @patch('uuid.uuid4')
    @patch('os.makedirs')
    @patch('os.path.dirname')
    @patch('os.path.exists')
    @patch('builtins.open', new_callable=mock_open)
    @patch('os.remove')
    def test_upload_media_success(self, mock_remove, mock_file, mock_exists, 
                                  mock_dirname, mock_makedirs, mock_uuid, 
                                  minio_client, mock_minio_lib):
        """Тест успешной загрузки медиа"""
        # Arrange
        mock_uuid.return_value = '123e4567-e89b-12d3-a456-426614174000'
        mock_exists.return_value = False
        mock_dirname.return_value = '/tmp'
        
        # Act
        media_id = minio_client.upload_media(
            user_id='user123',
            file_bytes=b'test content',
            mime_type='image/jpeg',
            file_name='test.jpg'
        )
        
        # Assert
        expected_media_id = 'user123/123e4567-e89b-12d3-a456-426614174000/test.jpg'
        expected_temp_path = '/tmp/user123_123e4567-e89b-12d3-a456-426614174000_test.jpg'
        
        assert media_id == expected_media_id
        mock_makedirs.assert_called_once_with('/tmp', exist_ok=True)
        mock_file.assert_called_once_with(expected_temp_path, 'wb')
        mock_minio_lib.fput_object.assert_called_once_with(
            'test-bucket',
            expected_media_id,
            expected_temp_path,
            content_type='image/jpeg'
        )
        mock_remove.assert_called_once_with(expected_temp_path)
    
    def test_get_media_success(self, minio_client, mock_minio_lib):
        """Тест успешного получения медиа"""
        # Arrange
        media_id = 'user123/uuid-123/test.jpg'
        mock_response = Mock()
        mock_response.read.return_value = b'file content'
        mock_minio_lib.get_object.return_value = mock_response
        
        # Act
        result = minio_client.get_media(media_id)
        
        # Assert
        assert result == b'file content'
        mock_minio_lib.get_object.assert_called_once_with('test-bucket', media_id)
        mock_response.close.assert_called_once()
        mock_response.release_conn.assert_called_once()
    
    def test_delete_media_success(self, minio_client, mock_minio_lib):
        """Тест успешного удаления медиа"""
        # Arrange
        media_id = 'user123/uuid-123/test.jpg'
        user_id = 'user123'
        
        # Act
        result = minio_client.delete_media(media_id, user_id)
        
        # Assert
        assert result is True
        mock_minio_lib.remove_object.assert_called_once_with('test-bucket', media_id)
    
    def test_delete_media_wrong_user(self, minio_client, mock_minio_lib):
        """Тест удаления медиа другого пользователя"""
        # Arrange
        media_id = 'user123/uuid-123/test.jpg'
        wrong_user_id = 'user456'
        
        # Act
        result = minio_client.delete_media(media_id, wrong_user_id)
        
        # Assert
        assert result is False
        mock_minio_lib.remove_object.assert_not_called()
    
    def test_list_user_media_success(self, minio_client, mock_minio_lib):
        """Тест получения списка медиа пользователя"""
        # Arrange
        user_id = 'user123'
        mock_objects = [
            Mock(object_name=f'{user_id}/uuid-1/file1.jpg', 
                 last_modified='2023-01-01T10:00:00Z', 
                 size=1024),
            Mock(object_name=f'{user_id}/uuid-2/file2.png', 
                 last_modified='2023-01-02T12:00:00Z', 
                 size=2048)
        ]
        mock_minio_lib.list_objects.return_value = mock_objects
        
        # Act
        result = minio_client.list_user_media(user_id)
        
        # Assert
        assert len(result) == 2
        assert result[0] == {
            'media_id': f'{user_id}/uuid-1/file1.jpg',
            'file_name': 'file1.jpg',
            'last_modified': '2023-01-01T10:00:00Z',
            'size': 1024
        }
        mock_minio_lib.list_objects.assert_called_once_with(
            'test-bucket', prefix=f'{user_id}/', recursive=True
        )
    
    def test_list_user_media_empty(self, minio_client, mock_minio_lib):
        """Тест получения пустого списка"""
        # Arrange
        user_id = 'user123'
        mock_minio_lib.list_objects.return_value = []
        
        # Act
        result = minio_client.list_user_media(user_id)
        
        # Assert
        assert result == []
    
    def test_get_presigned_url_success(self, minio_client, mock_minio_lib):
        """Тест получения presigned URL"""
        # Arrange
        media_id = 'user123/uuid-123/test.jpg'
        expected_url = 'http://localhost:9000/test-bucket/user123/uuid-123/test.jpg?token=abc'
        mock_minio_lib.presigned_get_object.return_value = expected_url
        
        # Act
        result = minio_client.get_presigned_url(media_id)
        
        # Assert
        assert result == expected_url
        mock_minio_lib.presigned_get_object.assert_called_once_with(
            'test-bucket', media_id, expires=ANY
        )
    
    def test_get_media_metadata_success(self, minio_client, mock_minio_lib):
        """Тест получения метаданных"""
        # Arrange
        media_id = 'user123/uuid-123/test.jpg'
        mock_stat = Mock(
            size=1024,
            content_type='image/jpeg',
            last_modified='2023-01-01T12:00:00Z'
        )
        mock_minio_lib.stat_object.return_value = mock_stat
        
        # Act
        result = minio_client.get_media_metadata(media_id)
        
        # Assert
        assert result == {
            'media_id': media_id,
            'size': 1024,
            'content_type': 'image/jpeg',
            'last_modified': '2023-01-01T12:00:00Z'
        }