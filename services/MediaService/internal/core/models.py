from dataclasses import dataclass
from typing import Optional

@dataclass
class MediaMetadata:
    user_id: str
    mime_type: str
    file_name: str
    upload_date: str