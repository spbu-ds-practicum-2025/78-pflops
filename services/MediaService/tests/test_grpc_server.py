# test_grpc_server.py
import pytest
from unittest.mock import Mock, patch, MagicMock, ANY
import sys
import os

sys.path.insert(0, os.path.abspath(os.path.join(os.path.dirname(__file__), '..')))

# Импортируем сгенерированные proto файлы
try:
    import pkg.pb.media_pb2 as media_pb2
    import pkg.pb.media_pb2_grpc as media_pb2_grpc
except ImportError:
    # Если proto файлы не сгенерированы, создаем моки
    class MockProto:
        class UploadMediaRequest:
            pass
        class UploadMediaResponse:
            pass
        class GetMediaRequest:
            pass
        class GetMediaResponse:
            pass
        class DeleteMediaRequest:
            pass
        class DeleteMediaResponse:
            pass
        class ListMediaRequest:
            pass
        class ListMediaResponse:
            pass
        class MediaItem:
            pass
        class GetUrlRequest:
            pass
        class GetUrlResponse:
            pass
    
    media_pb2 = MockProto()
    media_pb2_grpc = Mock()

from internal.grpc_server.server import MediaServiceServicer
from internal.core.models import MediaMetadata


class TestMediaServiceServicer:
    
    @pytest.fixture
    def mock_media_service(self):
        """Фикстура для мока MediaService"""
        return Mock()
    
    @pytest.fixture
    def servicer(self, mock_media_service):
        """Фикстура для MediaServiceServicer"""
        with patch('internal.grpc_server.server.MediaService') as mock:
            mock.return_value = mock_media_service
            return MediaServiceServicer()
    
    @pytest.fixture
    def mock_context(self):
        """Фикстура для мока gRPC контекста"""
        context = Mock()
        context.set_code = Mock()
        context.set_details = Mock()
        return context
    
    @patch('internal.grpc_server.server.settings')
    def test_upload_media_success(self, mock_settings, servicer, mock_media_service, mock_context):
        """Тест успешной загрузки через gRPC"""
        # Arrange
        request = Mock()
        request.user_id = 'user123'
        request.file_bytes = b'test content'
        request.mime_type = 'image/jpeg'
        request.file_name = 'test.jpg'
        
        mock_media_service.upload_media.return_value = 'user123/uuid-123/test.jpg'
        mock_settings.MINIO_BUCKET = 'media-service'
        
        # Act
        response = servicer.UploadMedia(request, mock_context)
        
        # Assert
        assert response.media_id == 'user123/uuid-123/test.jpg'
        assert response.message == "Файл успешно загружен"
        assert response.url == "/media/media-service/user123/uuid-123/test.jpg"
        mock_media_service.upload_media.assert_called_once_with(
            user_id='user123',
            file_bytes=b'test content',
            mime_type='image/jpeg',
            file_name='test.jpg'
        )
    
    def test_upload_media_failure(self, servicer, mock_media_service, mock_context):
        """Тест неудачной загрузки через gRPC"""
        # Arrange
        request = Mock()
        mock_media_service.upload_media.side_effect = Exception("Storage error")
        
        # Act
        response = servicer.UploadMedia(request, mock_context)
        
        # Assert
        mock_context.set_code.assert_called_once()
        mock_context.set_details.assert_called_once_with("Ошибка загрузки: Storage error")
    
    def test_get_media_success(self, servicer, mock_media_service, mock_context):
        """Тест успешного получения медиа через gRPC"""
        # Arrange
        request = Mock()
        request.media_id = 'user123/uuid-123/test.jpg'
        
        mock_media_service.get_media.return_value = b'file content'
        mock_media_service.get_media_metadata.return_value = MediaMetadata(
            user_id='user123',
            mime_type='image/jpeg',
            file_name='test.jpg',
            upload_date='2023-01-01T12:00:00'
        )
        
        # Act
        response = servicer.GetMedia(request, mock_context)
        
        # Assert
        assert response.file_bytes == b'file content'
        assert response.user_id == 'user123'
        assert response.mime_type == 'image/jpeg'
        assert response.file_name == 'test.jpg'
    
    def test_get_media_not_found(self, servicer, mock_media_service, mock_context):
        """Тест получения несуществующего медиа через gRPC"""
        # Arrange
        request = Mock()
        request.media_id = 'non-existent'
        
        mock_media_service.get_media.return_value = None
        
        # Act
        response = servicer.GetMedia(request, mock_context)
        
        # Assert
        mock_context.set_code.assert_called_once_with(ANY)
        mock_context.set_details.assert_called_once_with("Файл не найден")
    
    def test_get_media_without_metadata(self, servicer, mock_media_service, mock_context):
        """Тест получения медиа без метаданных"""
        # Arrange
        request = Mock()
        request.media_id = 'user123/uuid-123/test.jpg'
        
        mock_media_service.get_media.return_value = b'file content'
        mock_media_service.get_media_metadata.return_value = None
        
        # Act
        response = servicer.GetMedia(request, mock_context)
        
        # Assert
        assert response.file_bytes == b'file content'
        assert response.user_id == ''
        assert response.mime_type == ''
        assert response.file_name == ''
    
    def test_delete_media_success(self, servicer, mock_media_service, mock_context):
        """Тест успешного удаления через gRPC"""
        # Arrange
        request = Mock()
        request.media_id = 'user123/uuid-123/test.jpg'
        request.user_id = 'user123'
        
        mock_media_service.delete_media.return_value = True
        
        # Act
        response = servicer.DeleteMedia(request, mock_context)
        
        # Assert
        assert response.success is True
        assert response.message == "Файл успешно удален"
    
    def test_delete_media_failure(self, servicer, mock_media_service, mock_context):
        """Тест неудачного удаления через gRPC"""
        # Arrange
        request = Mock()
        request.media_id = 'user123/uuid-123/test.jpg'
        request.user_id = 'user123'
        
        mock_media_service.delete_media.return_value = False
        
        # Act
        response = servicer.DeleteMedia(request, mock_context)
        
        # Assert
        assert response.success is False
        assert response.message == "Файл не найден или нет прав доступа"
        mock_context.set_code.assert_called_once_with(ANY)
    
    def test_delete_media_exception(self, servicer, mock_media_service, mock_context):
        """Тест исключения при удалении через gRPC"""
        # Arrange
        request = Mock()
        mock_media_service.delete_media.side_effect = Exception("Delete error")
        
        # Act
        response = servicer.DeleteMedia(request, mock_context)
        
        # Assert
        mock_context.set_code.assert_called_once_with(ANY)
        mock_context.set_details.assert_called_once_with("Ошибка удаления: Delete error")
    
    def test_list_media_success(self, servicer, mock_media_service, mock_context):
        """Тест успешного получения списка медиа через gRPC"""
        # Arrange
        request = Mock()
        request.user_id = 'user123'
        
        mock_media_service.list_user_media.return_value = [
            {'media_id': 'user123/uuid-1/test1.jpg', 'file_name': 'test1.jpg'},
            {'media_id': 'user123/uuid-2/test2.png', 'file_name': 'test2.png'}
        ]
        
        # Настраиваем возврат метаданных
        def get_metadata_side_effect(media_id):
            if media_id == 'user123/uuid-1/test1.jpg':
                return MediaMetadata(
                    user_id='user123',
                    mime_type='image/jpeg',
                    file_name='test1.jpg',
                    upload_date='2023-01-01T10:00:00'
                )
            elif media_id == 'user123/uuid-2/test2.png':
                return MediaMetadata(
                    user_id='user123',
                    mime_type='image/png',
                    file_name='test2.png',
                    upload_date='2023-01-02T12:00:00'
                )
            return None
        
        mock_media_service.get_media_metadata.side_effect = get_metadata_side_effect
        mock_media_service.get_presigned_url.return_value = 'http://presigned.url/file'
        
        # Act
        response = servicer.ListMedia(request, mock_context)
        
        # Assert
        assert len(response.media_items) == 2
        assert response.media_items[0].media_id == 'user123/uuid-1/test1.jpg'
        assert response.media_items[0].file_name == 'test1.jpg'
        assert response.media_items[0].mime_type == 'image/jpeg'
        assert response.media_items[0].upload_date == '2023-01-01T10:00:00'
        assert response.media_items[0].url == 'http://presigned.url/file'
    
    def test_list_media_empty(self, servicer, mock_media_service, mock_context):
        """Тест получения пустого списка через gRPC"""
        # Arrange
        request = Mock()
        request.user_id = 'user123'
        
        mock_media_service.list_user_media.return_value = []
        
        # Act
        response = servicer.ListMedia(request, mock_context)
        
        # Assert
        assert len(response.media_items) == 0
    
    def test_list_media_exception(self, servicer, mock_media_service, mock_context):
        """Тест исключения при получении списка через gRPC"""
        # Arrange
        request = Mock()
        request.user_id = 'user123'
        
        mock_media_service.list_user_media.side_effect = Exception("List error")
        
        # Act
        response = servicer.ListMedia(request, mock_context)
        
        # Assert
        mock_context.set_code.assert_called_once_with(ANY)
        mock_context.set_details.assert_called_once_with("Ошибка получения списка: List error")
    
    def test_get_url_success(self, servicer, mock_media_service, mock_context):
        """Тест успешного получения URL через gRPC"""
        # Arrange
        request = Mock()
        request.media_id = 'user123/uuid-123/test.jpg'
        
        mock_media_service.get_presigned_url.return_value = 'http://presigned.url/test.jpg'
        
        # Act
        response = servicer.GetUrl(request, mock_context)
        
        # Assert
        assert response.url == 'http://presigned.url/test.jpg'
        assert response.media_id == 'user123/uuid-123/test.jpg'
    
    def test_get_url_not_found(self, servicer, mock_media_service, mock_context):
        """Тест получения URL для несуществующего медиа"""
        # Arrange
        request = Mock()
        request.media_id = 'non-existent'
        
        mock_media_service.get_presigned_url.return_value = None
        
        # Act
        response = servicer.GetUrl(request, mock_context)
        
        # Assert
        mock_context.set_code.assert_called_once_with(ANY)
    
    def test_get_url_exception(self, servicer, mock_media_service, mock_context):
        """Тест исключения при получении URL"""
        # Arrange
        request = Mock()
        request.media_id = 'user123/uuid-123/test.jpg'
        
        mock_media_service.get_presigned_url.side_effect = Exception("URL error")
        
        # Act
        response = servicer.GetUrl(request, mock_context)
        
        # Assert
        mock_context.set_code.assert_called_once_with(ANY)
        mock_context.set_details.assert_called_once_with("Ошибка генерации URL: URL error")