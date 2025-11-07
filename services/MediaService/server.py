import grpc
from concurrent import futures
import os
from datetime import datetime

from generated import media_service_pb2, media_service_pb2_grpc
from minio_client import MinioClient

class MediaServiceServicer(media_service_pb2_grpc.MediaServiceServicer):
    def __init__(self):
        self.minio_client = MinioClient()
        # Временное хранилище метаданных (в production заменить на БД)
        self.media_metadata = {}

    def UploadMedia(self, request, context):
        """Реализация UploadMedia RPC"""
        try:
            media_id = self.minio_client.upload_media(
                user_id=request.user_id,
                file_bytes=request.file_bytes,
                mime_type=request.mime_type,
                file_name=request.file_name
            )
            
            # Сохраняем метаданные
            self.media_metadata[media_id] = {
                'user_id': request.user_id,
                'mime_type': request.mime_type,
                'file_name': request.file_name,
                'upload_date': datetime.now().isoformat()
            }
            
            # Генерируем URL
            url = self.minio_client.get_presigned_url(media_id)
            
            return media_service_pb2.UploadMediaResponse(
                media_id=media_id,
                message="Файл успешно загружен",
                url=url or ""
            )
            
        except Exception as e:
            context.set_code(grpc.StatusCode.INTERNAL)
            context.set_details(f"Ошибка загрузки: {str(e)}")
            return media_service_pb2.UploadMediaResponse()

    def GetMedia(self, request, context):
        """Реализация GetMedia RPC"""
        try:
            file_bytes = self.minio_client.get_media(request.media_id)
            metadata = self.media_metadata.get(request.media_id, {})
            
            if file_bytes is None:
                context.set_code(grpc.StatusCode.NOT_FOUND)
                context.set_details("Файл не найден")
                return media_service_pb2.GetMediaResponse()
            
            return media_service_pb2.GetMediaResponse(
                user_id=metadata.get('user_id', ''),
                file_bytes=file_bytes,
                mime_type=metadata.get('mime_type', ''),
                file_name=metadata.get('file_name', '')
            )
            
        except Exception as e:
            context.set_code(grpc.StatusCode.INTERNAL)
            context.set_details(f"Ошибка получения: {str(e)}")
            return media_service_pb2.GetMediaResponse()

    def DeleteMedia(self, request, context):
        """Реализация DeleteMedia RPC"""
        try:
            success = self.minio_client.delete_media(
                request.media_id, 
                request.user_id
            )
            
            if success:
                # Удаляем метаданные
                self.media_metadata.pop(request.media_id, None)
                return media_service_pb2.DeleteMediaResponse(
                    message="Файл успешно удален",
                    success=True
                )
            else:
                context.set_code(grpc.StatusCode.PERMISSION_DENIED)
                return media_service_pb2.DeleteMediaResponse(
                    message="Файл не найден или нет прав доступа",
                    success=False
                )
                
        except Exception as e:
            context.set_code(grpc.StatusCode.INTERNAL)
            context.set_details(f"Ошибка удаления: {str(e)}")
            return media_service_pb2.DeleteMediaResponse()

    def ListMedia(self, request, context):
        """Реализация ListMedia RPC"""
        try:
            media_items = self.minio_client.list_user_media(request.user_id)
            
            response_items = []
            for item in media_items:
                metadata = self.media_metadata.get(item['media_id'], {})
                url = self.minio_client.get_presigned_url(item['media_id'])
                
                response_items.append(
                    media_service_pb2.MediaItem(
                        media_id=item['media_id'],
                        file_name=item['file_name'],
                        mime_type=metadata.get('mime_type', ''),
                        upload_date=metadata.get('upload_date', ''),
                        url=url or ""
                    )
                )
            
            return media_service_pb2.ListMediaResponse(media_items=response_items)
            
        except Exception as e:
            context.set_code(grpc.StatusCode.INTERNAL)
            context.set_details(f"Ошибка получения списка: {str(e)}")
            return media_service_pb2.ListMediaResponse()

    def GetUrl(self, request, context):
        """Реализация GetUrl RPC"""
        try:
            url = self.minio_client.get_presigned_url(request.media_id)
            
            if url:
                return media_service_pb2.GetUrlResponse(
                    url=url,
                    media_id=request.media_id
                )
            else:
                context.set_code(grpc.StatusCode.NOT_FOUND)
                return media_service_pb2.GetUrlResponse()
                
        except Exception as e:
            context.set_code(grpc.StatusCode.INTERNAL)
            context.set_details(f"Ошибка генерации URL: {str(e)}")
            return media_service_pb2.GetUrlResponse()

def serve():
    """Запуск gRPC сервера"""
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))
    media_service_pb2_grpc.add_MediaServiceServicer_to_server(
        MediaServiceServicer(), server
    )
    server.add_insecure_port('[::]:50051')
    
    print("Media Service запущен на порту 50051...")
    server.start()
    server.wait_for_termination()

if __name__ == '__main__':
    serve()