#!/usr/bin/env python3
import grpc
from concurrent import futures
import os
import sys
from grpc_reflection.v1alpha import reflection

# Добавляем корень проекта в PYTHONPATH
current_dir = os.path.dirname(os.path.abspath(__file__))
project_root = os.path.dirname(os.path.dirname(os.path.dirname(current_dir)))
sys.path.insert(0, project_root)

from internal.grpc_server.server import MediaServiceServicer
from pkg.pb import media_pb2_grpc, media_pb2
from config.settings import settings

def serve():
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))
    
    # Добавляем наш сервис
    media_pb2_grpc.add_MediaServiceServicer_to_server(
        MediaServiceServicer(), server
    )
    
    # Включаем Reflection API
    SERVICE_NAMES = (
        media_pb2.DESCRIPTOR.services_by_name['MediaService'].full_name,
        reflection.SERVICE_NAME,
    )
    reflection.enable_server_reflection(SERVICE_NAMES, server)
    
    server.add_insecure_port(f'[::]:{settings.GRPC_PORT}')
    
    print(f"Media Service запущен на порту {settings.GRPC_PORT}...")
    print("Reflection API включена")
    server.start()
    server.wait_for_termination()

if __name__ == '__main__':
    serve()