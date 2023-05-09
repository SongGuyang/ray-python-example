# AutoML Service based on cloud native

## Generate the grpc source files

- pip install grpcio-tools
- python -m grpc_tools.protoc  --python_out=automl/generated --pyi_out=automl/generated --grpc_python_out=automl/generated automl_service.proto -I.

## The start commands of echo components

### proxy
python -m automl.proxy --port=1234 --host-name=127.0.0.1 --operator-address=127.0.0.1:8080

### trainer
python -m automl.trainer --grpc-port=2345 --host-name=127.0.0.1 --proxy-address=127.0.0.1:1234 --trainer-id={trainer_pod_id} --task-id=0  --operator-address=127.0.0.1:8080
### worker
 python -m automl.worker --trainer-address=127.0.0.1:2345 --worker-id={worker_pod_id}
