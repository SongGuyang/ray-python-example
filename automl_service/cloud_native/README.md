# AutoML Service based on cloud native

## Generate the grpc source files

- pip install grpcio-tools
- python -m grpc_tools.protoc  --python_out=automl/generated --pyi_out=automl/generated --grpc_python_out=automl/generated automl_service.proto -I.
