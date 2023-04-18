# AutoML Service based on cloud native

## Generate the grpc source files

- pip install grpcio-tools
- python -m grpc_tools.protoc  --python_out=automl/generated --pyi_out=automl/generated --grpc_python_out=automl/generated automl_service.proto -I.

## The start commands of echo components

### proxy
python -m automl.proxy --port=80 --grpc-port=1234

### trainer
python -m automl.trainer --grpc-port=2345 --proxy-address="{proxy_ip}:1234" --trainer-index=0 --train-context='{"data_source":"",data_partition: "","model_season_length":[6, 7],"mode":["ZNA", "ZZZ"]}'

### worker
python -m automl.worker --proxy-address="{trainer_ip}:2345" --worker-index=0
