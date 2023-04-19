
class OperatorClient:
    def ___init__(self, address):
        self._address = address
    
    def start_trainer(self, grpc_port, proxy_address, trainer_index, spec):
        return f"test id {trainer_index}"
    
    def stop_trainer(self, context):
        return True

    def start_worker_group(self, trainer_address, specs):
        return "test worker group id"
    
    def stop_worker_group(self, context):
        return True
