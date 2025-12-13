# test_models.py
import pytest
from datetime import datetime
from internal.core.models import MediaMetadata


class TestMediaMetadata:
    
    def test_media_metadata_creation(self):
        """Тест создания объекта MediaMetadata"""
        metadata = MediaMetadata(
            user_id="user123",
            mime_type="image/jpeg",
            file_name="test.jpg",
            upload_date="2023-01-01T12:00:00"
        )
        
        assert metadata.user_id == "user123"
        assert metadata.mime_type == "image/jpeg"
        assert metadata.file_name == "test.jpg"
        assert metadata.upload_date == "2023-01-01T12:00:00"
    
    def test_media_metadata_dataclass_features(self):
        """Тест возможностей dataclass"""
        metadata1 = MediaMetadata(
            user_id="user123",
            mime_type="image/jpeg",
            file_name="test.jpg",
            upload_date="2023-01-01T12:00:00"
        )
        
        metadata2 = MediaMetadata(
            user_id="user123",
            mime_type="image/jpeg",
            file_name="test.jpg",
            upload_date="2023-01-01T12:00:00"
        )
        
        # Проверяем равенство
        assert metadata1 == metadata2
        
        # Проверяем repr
        assert "MediaMetadata" in repr(metadata1)