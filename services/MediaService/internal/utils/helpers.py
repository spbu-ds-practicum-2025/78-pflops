"""
Utility functions for MediaService
"""

import uuid
import hashlib
import mimetypes
from typing import Optional
from datetime import datetime, timedelta
import re


def generate_uuid() -> str:
    """Generate a unique UUID string"""
    return str(uuid.uuid4())


def calculate_file_hash(file_bytes: bytes) -> str:
    """Calculate MD5 hash of file bytes for integrity checking"""
    return hashlib.md5(file_bytes).hexdigest()


def validate_user_id(user_id: str) -> bool:
    """Validate user ID format"""
    if not user_id or len(user_id) > 100:
        return False
    # Basic validation - alphanumeric, underscores, hyphens
    return bool(re.match(r'^[a-zA-Z0-9_-]+$', user_id))


def validate_file_name(file_name: str) -> bool:
    """Validate file name for security"""
    if not file_name or len(file_name) > 255:
        return False
    # Prevent path traversal and other unsafe characters
    unsafe_patterns = ['..', '/', '\\', ':', '*', '?', '"', '<', '>', '|']
    return not any(pattern in file_name for pattern in unsafe_patterns)


def get_file_extension(mime_type: str) -> Optional[str]:
    """Get file extension from MIME type"""
    return mimetypes.guess_extension(mime_type)


def is_supported_image_type(mime_type: str) -> bool:
    """Check if MIME type is a supported image format"""
    supported_images = [
        'image/jpeg', 'image/jpg', 'image/png', 
        'image/gif', 'image/webp', 'image/svg+xml'
    ]
    return mime_type in supported_images


def is_supported_document_type(mime_type: str) -> bool:
    """Check if MIME type is a supported document format"""
    supported_docs = [
        'application/pdf',
        'text/plain',
        'application/msword',
        'application/vnd.openxmlformats-officedocument.wordprocessingml.document'
    ]
    return mime_type in supported_docs


def format_file_size(size_in_bytes: int) -> str:
    """Format file size in human-readable format"""
    if size_in_bytes == 0:
        return "0 B"
    
    size_names = ["B", "KB", "MB", "GB"]
    i = 0
    while size_in_bytes >= 1024 and i < len(size_names) - 1:
        size_in_bytes /= 1024.0
        i += 1
    
    return f"{size_in_bytes:.2f} {size_names[i]}"


def sanitize_filename(filename: str) -> str:
    """Sanitize filename by removing unsafe characters"""
    # Remove directory path components
    filename = filename.split('/')[-1].split('\\')[-1]
    
    # Replace unsafe characters with underscore
    filename = re.sub(r'[^\w\-_.]', '_', filename)
    
    # Limit length
    if len(filename) > 200:
        name, ext = os.path.splitext(filename)
        filename = name[:200-len(ext)] + ext
    
    return filename


def parse_date(date_string: str) -> Optional[datetime]:
    """Parse date string to datetime object"""
    try:
        return datetime.fromisoformat(date_string.replace('Z', '+00:00'))
    except (ValueError, TypeError):
        return None


def is_expired(timestamp: str, expiry_hours: int = 24) -> bool:
    """Check if a timestamp has expired"""
    target_date = parse_date(timestamp)
    if not target_date:
        return True
    
    expiry_time = target_date + timedelta(hours=expiry_hours)
    return datetime.now() > expiry_time


class RateLimiter:
    """Simple rate limiter for API calls"""
    
    def __init__(self, max_requests: int, time_window: int):
        self.max_requests = max_requests
        self.time_window = time_window  # in seconds
        self.requests = {}
    
    def is_allowed(self, user_id: str) -> bool:
        """Check if user is allowed to make a request"""
        now = datetime.now()
        user_requests = self.requests.get(user_id, [])
        
        # Remove old requests outside the time window
        user_requests = [req_time for req_time in user_requests 
                        if (now - req_time).seconds < self.time_window]
        
        if len(user_requests) >= self.max_requests:
            return False
        
        user_requests.append(now)
        self.requests[user_id] = user_requests
        return True


def create_error_response(message: str, code: str = "INTERNAL_ERROR") -> dict:
    """Create standardized error response"""
    return {
        "error": {
            "code": code,
            "message": message,
            "timestamp": datetime.now().isoformat()
        }
    }


# Initialize mimetypes with additional extensions
mimetypes.add_type('image/webp', '.webp')
mimetypes.add_type('application/wasm', '.wasm')