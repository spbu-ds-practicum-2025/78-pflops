from datetime import datetime
from typing import List, Dict, Optional
from internal.storage.minio_client import MinioClient
from internal.core.models import MediaMetadata

class MediaService:
    def __init__(self):
        self.minio_client = MinioClient()
        self.media_metadata: Dict[str, MediaMetadata] = {}
    
    def upload_media(self, user_id: str, file_bytes: bytes, mime_type: str, file_name: str) -> str:
        media_id = self.minio_client.upload_media(
            user_id=user_id,
            file_bytes=file_bytes,
            mime_type=mime_type,
            file_name=file_name
        )
        
        self.media_metadata[media_id] = MediaMetadata(
            user_id=user_id,
            mime_type=mime_type,
            file_name=file_name,
            upload_date=datetime.now().isoformat()
        )
        
        return media_id
    
    def get_media(self, media_id: str) -> Optional[bytes]:
        return self.minio_client.get_media(media_id)
    
    def delete_media(self, media_id: str, user_id: str) -> bool:
        success = self.minio_client.delete_media(media_id, user_id)
        if success:
            self.media_metadata.pop(media_id, None)
        return success
    
    def list_user_media(self, user_id: str) -> List[Dict]:
        return self.minio_client.list_user_media(user_id)
    
    def get_presigned_url(self, media_id: str) -> Optional[str]:
        return self.minio_client.get_presigned_url(media_id)
    
    def get_media_metadata(self, media_id: str) -> Optional[MediaMetadata]:
        return self.media_metadata.get(media_id)