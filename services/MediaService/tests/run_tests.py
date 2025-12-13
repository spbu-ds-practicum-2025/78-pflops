#!/usr/bin/env python3
import subprocess
import sys
import os

def run_tests():
    """Запуск тестов через pytest"""
    # Получаем путь к текущему файлу
    current_dir = os.path.dirname(os.path.abspath(__file__))
    
    # Формируем команду для запуска pytest
    cmd = [
        "python", "-m", "pytest",
        "test_models.py",
        "test_helpers.py",
        "test_service.py",
        "test_minio_client.py",
        "test_grpc_server.py",
        "-v",
        "--tb=short",
        "--cov=internal",
        "--cov-report=term-missing",
        "--cov-report=html",
        "--cov-fail-under=80"  # Минимальное покрытие 80%
    ]
    
    print("Запуск тестов...")
    print(f"Команда: {' '.join(cmd)}")
    print()
    
    try:
        # Запускаем pytest
        result = subprocess.run(cmd, cwd=current_dir)
        
        # Возвращаем код завершения
        return result.returncode
        
    except FileNotFoundError:
        print("Ошибка: pytest не найден. Установите его с помощью:")
        print("pip install pytest pytest-cov")
        return 1
    except Exception as e:
        print(f"Ошибка при запуске тестов: {e}")
        return 1

if __name__ == "__main__":
    # Добавляем путь к проекту
    sys.path.insert(0, os.path.dirname(os.path.abspath(__file__)))
    
    # Запускаем тесты
    exit_code = run_tests()
    
    sys.exit(exit_code)