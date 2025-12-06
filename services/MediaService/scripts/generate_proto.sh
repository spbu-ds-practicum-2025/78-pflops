#!/bin/bash
python -m grpc_tools.protoc \
    --proto_path=proto \
    --python_out=pkg/pb \
    --grpc_python_out=pkg/pb \
    --mypy_out=pkg/pb \
    proto/media.proto