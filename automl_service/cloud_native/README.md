# AutoML Service based on cloud native

## Generate the grpc source files

- pip install grpcio-tools
- python -m grpc_tools.protoc  --python_out=automl/generated --pyi_out=automl/generated --grpc_python_out=automl/generated automl_service.proto -I.

## The start commands of echo components

### proxy
python -m automl.proxy --grpc-port=1234 --operator-address=""

### trainer
python -m automl.trainer --grpc-port=2345 --proxy-address="{proxy_ip}:1234" --trainer-id="{service name}" --task-id=0 --operator-address=""
### worker
python -m automl.worker --trainer-address="{trainer_ip}:2345" --worker-id="{pod name}"
