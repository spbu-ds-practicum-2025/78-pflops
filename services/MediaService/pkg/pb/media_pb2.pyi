from google.protobuf.internal import containers as _containers
from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from collections.abc import Iterable as _Iterable, Mapping as _Mapping
from typing import ClassVar as _ClassVar, Optional as _Optional, Union as _Union

DESCRIPTOR: _descriptor.FileDescriptor

class UploadMediaRequest(_message.Message):
    __slots__ = ("user_id", "file_bytes", "mime_type", "file_name")
    USER_ID_FIELD_NUMBER: _ClassVar[int]
    FILE_BYTES_FIELD_NUMBER: _ClassVar[int]
    MIME_TYPE_FIELD_NUMBER: _ClassVar[int]
    FILE_NAME_FIELD_NUMBER: _ClassVar[int]
    user_id: str
    file_bytes: bytes
    mime_type: str
    file_name: str
    def __init__(self, user_id: _Optional[str] = ..., file_bytes: _Optional[bytes] = ..., mime_type: _Optional[str] = ..., file_name: _Optional[str] = ...) -> None: ...

class UploadMediaResponse(_message.Message):
    __slots__ = ("media_id", "message", "url")
    MEDIA_ID_FIELD_NUMBER: _ClassVar[int]
    MESSAGE_FIELD_NUMBER: _ClassVar[int]
    URL_FIELD_NUMBER: _ClassVar[int]
    media_id: str
    message: str
    url: str
    def __init__(self, media_id: _Optional[str] = ..., message: _Optional[str] = ..., url: _Optional[str] = ...) -> None: ...

class GetMediaRequest(_message.Message):
    __slots__ = ("media_id",)
    MEDIA_ID_FIELD_NUMBER: _ClassVar[int]
    media_id: str
    def __init__(self, media_id: _Optional[str] = ...) -> None: ...

class GetMediaResponse(_message.Message):
    __slots__ = ("user_id", "file_bytes", "mime_type", "file_name")
    USER_ID_FIELD_NUMBER: _ClassVar[int]
    FILE_BYTES_FIELD_NUMBER: _ClassVar[int]
    MIME_TYPE_FIELD_NUMBER: _ClassVar[int]
    FILE_NAME_FIELD_NUMBER: _ClassVar[int]
    user_id: str
    file_bytes: bytes
    mime_type: str
    file_name: str
    def __init__(self, user_id: _Optional[str] = ..., file_bytes: _Optional[bytes] = ..., mime_type: _Optional[str] = ..., file_name: _Optional[str] = ...) -> None: ...

class DeleteMediaRequest(_message.Message):
    __slots__ = ("media_id", "user_id")
    MEDIA_ID_FIELD_NUMBER: _ClassVar[int]
    USER_ID_FIELD_NUMBER: _ClassVar[int]
    media_id: str
    user_id: str
    def __init__(self, media_id: _Optional[str] = ..., user_id: _Optional[str] = ...) -> None: ...

class DeleteMediaResponse(_message.Message):
    __slots__ = ("message", "success")
    MESSAGE_FIELD_NUMBER: _ClassVar[int]
    SUCCESS_FIELD_NUMBER: _ClassVar[int]
    message: str
    success: bool
    def __init__(self, message: _Optional[str] = ..., success: bool = ...) -> None: ...

class ListMediaRequest(_message.Message):
    __slots__ = ("user_id",)
    USER_ID_FIELD_NUMBER: _ClassVar[int]
    user_id: str
    def __init__(self, user_id: _Optional[str] = ...) -> None: ...

class ListMediaResponse(_message.Message):
    __slots__ = ("media_items",)
    MEDIA_ITEMS_FIELD_NUMBER: _ClassVar[int]
    media_items: _containers.RepeatedCompositeFieldContainer[MediaItem]
    def __init__(self, media_items: _Optional[_Iterable[_Union[MediaItem, _Mapping]]] = ...) -> None: ...

class MediaItem(_message.Message):
    __slots__ = ("media_id", "file_name", "mime_type", "upload_date", "url")
    MEDIA_ID_FIELD_NUMBER: _ClassVar[int]
    FILE_NAME_FIELD_NUMBER: _ClassVar[int]
    MIME_TYPE_FIELD_NUMBER: _ClassVar[int]
    UPLOAD_DATE_FIELD_NUMBER: _ClassVar[int]
    URL_FIELD_NUMBER: _ClassVar[int]
    media_id: str
    file_name: str
    mime_type: str
    upload_date: str
    url: str
    def __init__(self, media_id: _Optional[str] = ..., file_name: _Optional[str] = ..., mime_type: _Optional[str] = ..., upload_date: _Optional[str] = ..., url: _Optional[str] = ...) -> None: ...

class GetUrlRequest(_message.Message):
    __slots__ = ("media_id",)
    MEDIA_ID_FIELD_NUMBER: _ClassVar[int]
    media_id: str
    def __init__(self, media_id: _Optional[str] = ...) -> None: ...

class GetUrlResponse(_message.Message):
    __slots__ = ("url", "media_id")
    URL_FIELD_NUMBER: _ClassVar[int]
    MEDIA_ID_FIELD_NUMBER: _ClassVar[int]
    url: str
    media_id: str
    def __init__(self, url: _Optional[str] = ..., media_id: _Optional[str] = ...) -> None: ...
