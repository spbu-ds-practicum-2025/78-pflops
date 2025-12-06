from setuptools import setup, find_packages

setup(
    name="media_service",
    version="0.1.0",
    packages=find_packages(),
    install_requires=[
        "grpcio==1.76.0",
        "grpcio-tools==1.76.0",
        "minio==7.2.18",
        "protobuf==6.33.0",
    ],
    python_requires=">=3.9",
)