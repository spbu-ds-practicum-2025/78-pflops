import grpc
from internal.core.service import MediaService
from pkg.pb import media_pb2, media_pb2_grpc
from config import settings

class MediaServiceServicer(media_pb2_grpc.MediaServiceServicer):
    def __init__(self):
        self.media_service = MediaService()

    def UploadMedia(self, request, context):
        try:
            media_id = self.media_service.upload_media(
                user_id=request.user_id,
                file_bytes=request.file_bytes,
                mime_type=request.mime_type,
                file_name=request.file_name
            )
            
            # Формируем стабильный HTTP URL через nginx, а не presigned-ссылку MinIO
            public_url = f"/media/{settings.MINIO_BUCKET}/{media_id}"

            return media_pb2.UploadMediaResponse(
                media_id=media_id,
                message="Файл успешно загружен",
                url=public_url
            )
            
        except Exception as e:
            context.set_code(grpc.StatusCode.INTERNAL)
            context.set_details(f"Ошибка загрузки: {str(e)}")
            return media_pb2.UploadMediaResponse()

    def GetMedia(self, request, context):
        try:
            file_bytes = self.media_service.get_media(request.media_id)
            metadata = self.media_service.get_media_metadata(request.media_id)
            
            if file_bytes is None:
                context.set_code(grpc.StatusCode.NOT_FOUND)
                context.set_details("Файл не найден")
                return media_pb2.GetMediaResponse()
            
            return media_pb2.GetMediaResponse(
                user_id=metadata.user_id if metadata else '',
                file_bytes=file_bytes,
                mime_type=metadata.mime_type if metadata else '',
                file_name=metadata.file_name if metadata else ''
            )
            
        except Exception as e:
            context.set_code(grpc.StatusCode.INTERNAL)
            context.set_details(f"Ошибка получения: {str(e)}")
            return media_pb2.GetMediaResponse()

    def DeleteMedia(self, request, context):
        try:
            success = self.media_service.delete_media(
                request.media_id, 
                request.user_id
            )
            
            if success:
                return media_pb2.DeleteMediaResponse(
                    message="Файл успешно удален",
                    success=True
                )
            else:
                context.set_code(grpc.StatusCode.PERMISSION_DENIED)
                return media_pb2.DeleteMediaResponse(
                    message="Файл не найден или нет прав доступа",
                    success=False
                )
                
        except Exception as e:
            context.set_code(grpc.StatusCode.INTERNAL)
            context.set_details(f"Ошибка удаления: {str(e)}")
            return media_pb2.DeleteMediaResponse()

    def ListMedia(self, request, context):
        try:
            media_items = self.media_service.list_user_media(request.user_id)
            
            response_items = []
            for item in media_items:
                metadata = self.media_service.get_media_metadata(item['media_id'])
                url = self.media_service.get_presigned_url(item['media_id'])
                
                response_items.append(
                    media_pb2.MediaItem(
                        media_id=item['media_id'],
                        file_name=item['file_name'],
                        mime_type=metadata.mime_type if metadata else '',
                        upload_date=metadata.upload_date if metadata else '',
                        url=url or ""
                    )
                )
            
            return media_pb2.ListMediaResponse(media_items=response_items)
            
        except Exception as e:
            context.set_code(grpc.StatusCode.INTERNAL)
            context.set_details(f"Ошибка получения списка: {str(e)}")
            return media_pb2.ListMediaResponse()

    def GetUrl(self, request, context):
        try:
            url = self.media_service.get_presigned_url(request.media_id)
            
            if url:
                return media_pb2.GetUrlResponse(
                    url=url,
                    media_id=request.media_id
                )
            else:
                context.set_code(grpc.StatusCode.NOT_FOUND)
                return media_pb2.GetUrlResponse()
                
        except Exception as e:
            context.set_code(grpc.StatusCode.INTERNAL)
            context.set_details(f"Ошибка генерации URL: {str(e)}")
            return media_pb2.GetUrlResponse()