
class OperatorClient:
    def __init__(self, address):
        self._address = address
    
    def start_trainer(self, grpc_port, proxy_address, task_id, spec):
        return f"test_trainer_id"
    
    def stop_trainer(self, trainer_id):
        return True

    def start_worker_group(self, trainer_address, number, specs):
        return "group_id", [f"test_worker_id_{id}" for id in range(number)]
    
    def stop_worker_group(self, group_id):
        return True
