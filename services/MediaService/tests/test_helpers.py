# test_helpers.py
import pytest
from unittest.mock import Mock, patch, MagicMock
from datetime import datetime, timedelta
import os
import sys

sys.path.insert(0, os.path.abspath(os.path.join(os.path.dirname(__file__), '..')))

from internal.utils.helpers import (
    generate_uuid,
    calculate_file_hash,
    validate_user_id,
    validate_file_name,
    get_file_extension,
    is_supported_image_type,
    is_supported_document_type,
    format_file_size,
    sanitize_filename,
    parse_date,
    is_expired,
    RateLimiter,
    create_error_response
)


class TestGenerateUUID:
    
    def test_generate_uuid_format(self):
        """Тест формата UUID"""
        # Act
        uuid_str = generate_uuid()
        
        # Assert
        assert isinstance(uuid_str, str)
        assert len(uuid_str) == 36  # Стандартная длина UUID
        # Проверяем формат: xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx
        parts = uuid_str.split('-')
        assert len(parts) == 5
        assert len(parts[0]) == 8
        assert len(parts[1]) == 4
        assert len(parts[2]) == 4
        assert len(parts[3]) == 4
        assert len(parts[4]) == 12
    
    def test_generate_uuid_uniqueness(self):
        """Тест уникальности UUID"""
        # Act
        uuid1 = generate_uuid()
        uuid2 = generate_uuid()
        
        # Assert
        assert uuid1 != uuid2, "UUID должны быть уникальными"


class TestCalculateFileHash:
    
    def test_calculate_file_hash_basic(self):
        """Тест расчета хэша файла"""
        # Arrange
        test_data = b"Hello, World!"
        
        # Act
        result = calculate_file_hash(test_data)
        
        # Assert
        assert isinstance(result, str)
        assert len(result) == 32  # Длина MD5 хэша в hex
        
        # Проверяем известный хэш для "Hello, World!"
        expected_hash = "65a8e27d8879283831b664bd8b7f0ad4"
        assert result == expected_hash
    
    def test_calculate_file_hash_empty(self):
        """Тест расчета хэша пустого файла"""
        # Arrange
        test_data = b""
        
        # Act
        result = calculate_file_hash(test_data)
        
        # Assert
        assert result == "d41d8cd98f00b204e9800998ecf8427e"  # MD5 of empty string
    
    def test_calculate_file_hash_consistency(self):
        """Тест консистентности хэширования"""
        # Arrange
        test_data = b"The same content"
        
        # Act
        hash1 = calculate_file_hash(test_data)
        hash2 = calculate_file_hash(test_data)
        
        # Assert
        assert hash1 == hash2, "Хэши одинаковых данных должны совпадать"


class TestValidateUserId:
    
    def test_validate_user_id_valid(self):
        """Тест валидных user ID"""
        valid_ids = [
            "user123",
            "user_123",
            "user-123",
            "UPPERCASE",
            "lowercase",
            "mixedCase123",
            "a" * 100  # Максимальная длина
        ]
        
        for user_id in valid_ids:
            assert validate_user_id(user_id) is True, f"User ID '{user_id}' должен быть валидным"
    
    def test_validate_user_id_edge_cases(self):
        """Тест граничных случаев"""
        # Граничные символы
        assert validate_user_id("user_123") is True
        assert validate_user_id("user-123") is True
        assert validate_user_id("123456") is True
        assert validate_user_id("_user") is True
        assert validate_user_id("-user") is True


class TestValidateFileName:
    
    def test_validate_file_name_valid(self):
        """Тест валидных имен файлов"""
        valid_names = [
            "file.txt",
            "document.pdf",
            "image.jpg",
            "my_file_name.png",
            "File With Spaces.txt",
            "file-with-dashes.pdf",
            "file.with.dots.jpg",
            "a" * 255  # Максимальная длина
        ]
        
        for file_name in valid_names:
            assert validate_file_name(file_name) is True, f"Имя файла '{file_name}' должно быть валидным"
    
    def test_validate_file_name_invalid(self):
        """Тест невалидных имен файлов"""
        invalid_names = [
            "",  # Пустое имя
            "../etc/passwd",  # Path traversal
            "..\\windows\\system32",  # Windows path traversal
            "file/name.txt",  # Содержит слеш
            "file\\name.txt",  # Содержит обратный слеш
            "file:name.txt",  # Содержит двоеточие
            "file*name.txt",  # Содержит звездочку
            "file?name.txt",  # Содержит вопрос
            'file"name.txt',  # Содержит кавычки
            "file<name.txt",  # Содержит <
            "file>name.txt",  # Содержит >
            "file|name.txt",  # Содержит |
            "a" * 256  # Слишком длинный
        ]
        
        for file_name in invalid_names:
            assert validate_file_name(file_name) is False, f"Имя файла '{file_name}' должно быть невалидным"


class TestGetFileExtension:

    def test_get_file_extension_unknown_type(self):
        """Тест получения расширения для неизвестного MIME типа"""
        # Act
        result = get_file_extension("unknown/type")
        
        # Assert
        assert result is None
    
    def test_get_file_extension_empty_string(self):
        """Тест получения расширения для пустой строки"""
        # Act
        result = get_file_extension("")
        
        # Assert
        assert result is None


class TestIsSupportedImageType:
    
    def test_is_supported_image_type_supported(self):
        """Тест поддержки изображений"""
        supported_types = [
            "image/jpeg",
            "image/jpg",  # Альтернативное написание
            "image/png",
            "image/gif",
            "image/webp",
            "image/svg+xml"
        ]
        
        for mime_type in supported_types:
            assert is_supported_image_type(mime_type) is True, f"MIME тип '{mime_type}' должен поддерживаться"
    
    def test_is_supported_image_type_unsupported(self):
        """Тест неподдерживаемых типов изображений"""
        unsupported_types = [
            "image/bmp",
            "image/tiff",
            "image/ico",
            "application/pdf",
            "text/plain",
            ""
        ]
        
        for mime_type in unsupported_types:
            assert is_supported_image_type(mime_type) is False, f"MIME тип '{mime_type}' не должен поддерживаться"


class TestIsSupportedDocumentType:
    
    def test_is_supported_document_type_supported(self):
        """Тест поддержки документов"""
        supported_types = [
            "application/pdf",
            "text/plain",
            "application/msword",
            "application/vnd.openxmlformats-officedocument.wordprocessingml.document"
        ]
        
        for mime_type in supported_types:
            assert is_supported_document_type(mime_type) is True, f"MIME тип '{mime_type}' должен поддерживаться"
    
    def test_is_supported_document_type_unsupported(self):
        """Тест неподдерживаемых типов документов"""
        unsupported_types = [
            "application/json",
            "application/xml",
            "image/jpeg",
            "text/html",
            ""
        ]
        
        for mime_type in unsupported_types:
            assert is_supported_document_type(mime_type) is False, f"MIME тип '{mime_type}' не должен поддерживаться"


class TestFormatFileSize:
    
    def test_format_file_size_bytes(self):
        """Тест форматирования размера в байтах"""
        test_cases = [
            (0, "0 B"),
            (1, "1.00 B"),
            (500, "500.00 B"),
            (1023, "1023.00 B")
        ]
        
        for size, expected in test_cases:
            result = format_file_size(size)
            assert result == expected, f"Для {size} байт ожидалось '{expected}', получено '{result}'"
    
    def test_format_file_size_kilobytes(self):
        """Тест форматирования размера в килобайтах"""
        test_cases = [
            (1024, "1.00 KB"),
            (1536, "1.50 KB"),  # 1.5 KB
            (1024 * 10, "10.00 KB"),
            (1024 * 1024 - 1, "1024.00 KB")  # Почти 1 MB
        ]
        
        for size, expected in test_cases:
            result = format_file_size(size)
            assert result == expected, f"Для {size} байт ожидалось '{expected}', получено '{result}'"
    
    def test_format_file_size_megabytes(self):
        """Тест форматирования размера в мегабайтах"""
        test_cases = [
            (1024 * 1024, "1.00 MB"),
            (1024 * 1024 * 2.5, "2.50 MB"),  # 2.5 MB
            (1024 * 1024 * 100, "100.00 MB"),
            (1024 * 1024 * 1024 - 1, "1024.00 MB")  # Почти 1 GB
        ]
        
        for size, expected in test_cases:
            result = format_file_size(size)
            assert result == expected, f"Для {size} байт ожидалось '{expected}', получено '{result}'"
    
    def test_format_file_size_gigabytes(self):
        """Тест форматирования размера в гигабайтах"""
        test_cases = [
            (1024 * 1024 * 1024, "1.00 GB"),
            (1024 * 1024 * 1024 * 5, "5.00 GB"),
        ]
        
        for size, expected in test_cases:
            result = format_file_size(size)
            assert result == expected, f"Для {size} байт ожидалось '{expected}', получено '{result}'"
    
    def test_format_file_size_very_large(self):
        """Тест форматирования очень больших размеров"""
        # 2 TB в байтах
        very_large = 2 * 1024 * 1024 * 1024 * 1024
        result = format_file_size(very_large)
        assert "GB" in result  # Ограничиваемся GB в текущей реализации

class TestParseDate:
    
    def test_parse_date_valid_iso(self):
        """Тест парсинга валидной ISO даты"""
        test_cases = [
            ("2023-01-01T12:00:00", datetime(2023, 1, 1, 12, 0, 0)),
            ("2023-01-01T12:00:00.123456", datetime(2023, 1, 1, 12, 0, 0, 123456)),
            ("2023-01-01", datetime(2023, 1, 1)),
        ]
        
        for date_str, expected in test_cases:
            result = parse_date(date_str)
            assert result == expected, f"Для строки '{date_str}' ожидалось {expected}"
    
    def test_parse_date_with_timezone(self):
        """Тест парсинга даты с часовым поясом"""
        # Act
        result = parse_date("2023-01-01T12:00:00Z")
        
        # Assert
        assert result is not None
        # Проверяем, что дата распарсилась (точное сравнение с часовым поясом сложнее)
        assert result.year == 2023
        assert result.month == 1
        assert result.day == 1


class TestIsExpired:
    
    def test_is_expired_invalid_date(self):
        """Тест с невалидной датой"""
        # Act
        result = is_expired("invalid date")
        
        # Assert
        assert result is True  # Невалидная дата считается истекшей


class TestRateLimiter:
    
    def test_rate_limiter_initialization(self):
        """Тест инициализации RateLimiter"""
        # Act
        limiter = RateLimiter(max_requests=10, time_window=60)
        
        # Assert
        assert limiter.max_requests == 10
        assert limiter.time_window == 60
        assert limiter.requests == {}
    
    @patch('internal.utils.helpers.datetime')
    def test_rate_limiter_allowed_under_limit(self, mock_datetime):
        """Тест, когда запросы разрешены (ниже лимита)"""
        # Arrange
        limiter = RateLimiter(max_requests=3, time_window=60)
        user_id = "user123"
        
        # Устанавливаем текущее время
        base_time = datetime(2023, 1, 1, 12, 0, 0)
        mock_datetime.now.return_value = base_time
        
        # Первые два запроса
        assert limiter.is_allowed(user_id) is True
        assert limiter.is_allowed(user_id) is True
        
        # Третий запрос
        assert limiter.is_allowed(user_id) is True
    
    @patch('internal.utils.helpers.datetime')
    def test_rate_limiter_not_allowed_over_limit(self, mock_datetime):
        """Тест, когда запросы запрещены (превышен лимит)"""
        # Arrange
        limiter = RateLimiter(max_requests=2, time_window=60)
        user_id = "user123"
        
        # Устанавливаем текущее время
        base_time = datetime(2023, 1, 1, 12, 0, 0)
        mock_datetime.now.return_value = base_time
        
        # Первые два запроса
        assert limiter.is_allowed(user_id) is True
        assert limiter.is_allowed(user_id) is True
        
        # Третий запрос должен быть отклонен
        assert limiter.is_allowed(user_id) is False
    
    @patch('internal.utils.helpers.datetime')
    def test_rate_limiter_old_requests_removed(self, mock_datetime):
        """Тест удаления старых запросов"""
        # Arrange
        limiter = RateLimiter(max_requests=2, time_window=60)
        user_id = "user123"
        
        # Первый запрос 70 секунд назад (вне временного окна)
        old_time = datetime(2023, 1, 1, 12, 0, 0)
        mock_datetime.now.return_value = old_time
        limiter.is_allowed(user_id)  # Старый запрос
        
        # Второй запрос сейчас
        current_time = datetime(2023, 1, 1, 12, 1, 10)  # 70 секунд позже
        mock_datetime.now.return_value = current_time
        
        # Act & Assert
        # Должен быть разрешен, так как старый запрос удален
        assert limiter.is_allowed(user_id) is True
        
        # Еще один запрос должен быть разрешен (всего 2 в окне)
        assert limiter.is_allowed(user_id) is True
        
        # Третий должен быть отклонен
        assert limiter.is_allowed(user_id) is False
    
    @patch('internal.utils.helpers.datetime')
    def test_rate_limiter_multiple_users(self, mock_datetime):
        """Тест RateLimiter с несколькими пользователями"""
        # Arrange
        limiter = RateLimiter(max_requests=2, time_window=60)
        
        # Устанавливаем время
        mock_datetime.now.return_value = datetime(2023, 1, 1, 12, 0, 0)
        
        # User1 делает 2 запроса
        assert limiter.is_allowed("user1") is True
        assert limiter.is_allowed("user1") is True
        
        # User2 делает 1 запрос
        assert limiter.is_allowed("user2") is True
        
        # User1 превысил лимит
        assert limiter.is_allowed("user1") is False
        
        # User2 еще может
        assert limiter.is_allowed("user2") is True
    
    def test_rate_limiter_edge_cases(self):
        """Тест граничных случаев RateLimiter"""
        # Ноль запросов
        limiter = RateLimiter(max_requests=0, time_window=60)
        assert limiter.is_allowed("user") is False
        
        # Большое временное окно
        limiter = RateLimiter(max_requests=100, time_window=3600)
        assert limiter.is_allowed("user") is True


class TestCreateErrorResponse:
    
    def test_create_error_response_default(self):
        """Тест создания стандартного ответа об ошибке"""
        # Act
        response = create_error_response("Something went wrong")
        
        # Assert
        assert "error" in response
        error = response["error"]
        assert error["code"] == "INTERNAL_ERROR"
        assert error["message"] == "Something went wrong"
        assert "timestamp" in error
        
        # Проверяем формат timestamp
        from datetime import datetime
        try:
            datetime.fromisoformat(error["timestamp"].replace('Z', '+00:00'))
            timestamp_valid = True
        except ValueError:
            timestamp_valid = False
        
        assert timestamp_valid
    
    def test_create_error_response_custom_code(self):
        """Тест создания ответа об ошибке с кастомным кодом"""
        # Act
        response = create_error_response("Not found", code="NOT_FOUND")
        
        # Assert
        assert response["error"]["code"] == "NOT_FOUND"
        assert response["error"]["message"] == "Not found"
    
    def test_create_error_response_empty_message(self):
        """Тест создания ответа с пустым сообщением"""
        # Act
        response = create_error_response("")
        
        # Assert
        assert response["error"]["message"] == ""
    
    def test_create_error_response_special_characters(self):
        """Тест создания ответа со специальными символами"""
        # Act
        message = "Error: Something <went> wrong & нужно исправить"
        response = create_error_response(message)
        
        # Assert
        assert response["error"]["message"] == message


class TestAdditionalHelpers:
    
    def test_mimetypes_initialization(self):
        """Тест инициализации mimetypes с дополнительными расширениями"""
        import mimetypes
        
        # Проверяем, что дополнительные типы добавлены
        assert mimetypes.guess_extension('image/webp') == '.webp'
        assert mimetypes.guess_extension('application/wasm') == '.wasm'


# Комплексные тесты
class TestIntegrationHelpers:
    
    def test_full_file_processing_flow(self):
        """Тест полного потока обработки файла с использованием нескольких функций"""
        # 1. Генерация UUID для имени файла
        file_uuid = generate_uuid()
        
        # 2. Валидация user_id
        user_id = "user_123"
        assert validate_user_id(user_id) is True
        
        # 3. Валидация имени файла
        file_name = f"{file_uuid}.jpg"
        assert validate_file_name(file_name) is True
        
        # 4. Проверка типа файла
        mime_type = "image/jpeg"
        assert is_supported_image_type(mime_type) is True
        
        # 5. Получение расширения
        extension = get_file_extension(mime_type)
        assert extension == ".jpg"
        
        # 6. Создание тестовых данных
        file_data = b"fake image data" * 100
        
        # 7. Расчет хэша
        file_hash = calculate_file_hash(file_data)
        assert len(file_hash) == 32
        
        # 8. Форматирование размера
        size_formatted = format_file_size(len(file_data))
        assert "B" in size_formatted or "KB" in size_formatted
    
    def test_error_handling_flow(self):
        """Тест потока обработки ошибок"""
        # 1. Создание ответа об ошибке
        error_response = create_error_response("File not found", code="FILE_NOT_FOUND")
        
        # 2. Парсинг даты из ошибки
        error_timestamp = error_response["error"]["timestamp"]
        parsed_date = parse_date(error_timestamp)
        
        # 3. Проверка, не истекла ли ошибка (шутка, но тестируем функционал)
        assert is_expired(error_timestamp, expiry_hours=24) is False
        
        # Проверяем структуру ответа
        assert error_response["error"]["code"] == "FILE_NOT_FOUND"
        assert error_response["error"]["message"] == "File not found"
        assert parsed_date is not None


# Тесты производительности (необязательные, но полезные)
class TestPerformance:
    
    def test_generate_uuid_performance(self):
        """Тест производительности генерации UUID"""
        import time
        
        start_time = time.time()
        for _ in range(1000):
            generate_uuid()
        end_time = time.time()
        
        # Генерация 1000 UUID не должна занимать больше 1 секунды
        assert end_time - start_time < 1.0
    
    def test_calculate_hash_performance(self):
        """Тест производительности расчета хэша"""
        import time
        
        test_data = b"x" * 1024 * 1024  # 1 MB данных
        
        start_time = time.time()
        calculate_file_hash(test_data)
        end_time = time.time()
        
        # Хэширование 1 MB не должно занимать больше 0.1 секунды
        assert end_time - start_time < 0.1


if __name__ == "__main__":
    pytest.main(['-v', __file__])