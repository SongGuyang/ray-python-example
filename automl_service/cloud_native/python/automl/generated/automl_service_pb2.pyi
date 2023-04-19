from google.protobuf.internal import containers as _containers
from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from typing import ClassVar as _ClassVar, Iterable as _Iterable, Mapping as _Mapping, Optional as _Optional

DESCRIPTOR: _descriptor.FileDescriptor

class DoAutoMLReply(_message.Message):
    __slots__ = ["message", "success", "task_id"]
    MESSAGE_FIELD_NUMBER: _ClassVar[int]
    SUCCESS_FIELD_NUMBER: _ClassVar[int]
    TASK_ID_FIELD_NUMBER: _ClassVar[int]
    message: str
    success: bool
    task_id: int
    def __init__(self, success: bool = ..., task_id: _Optional[int] = ..., message: _Optional[str] = ...) -> None: ...

class DoAutoMLRequest(_message.Message):
    __slots__ = ["data_partition", "data_source", "model_season_lengths", "models"]
    DATA_PARTITION_FIELD_NUMBER: _ClassVar[int]
    DATA_SOURCE_FIELD_NUMBER: _ClassVar[int]
    MODELS_FIELD_NUMBER: _ClassVar[int]
    MODEL_SEASON_LENGTHS_FIELD_NUMBER: _ClassVar[int]
    data_partition: str
    data_source: str
    model_season_lengths: _containers.RepeatedScalarFieldContainer[int]
    models: _containers.RepeatedScalarFieldContainer[str]
    def __init__(self, data_source: _Optional[str] = ..., data_partition: _Optional[str] = ..., model_season_lengths: _Optional[_Iterable[int]] = ..., models: _Optional[_Iterable[str]] = ...) -> None: ...

class GetResultReply(_message.Message):
    __slots__ = ["result", "success"]
    RESULT_FIELD_NUMBER: _ClassVar[int]
    SUCCESS_FIELD_NUMBER: _ClassVar[int]
    result: str
    success: bool
    def __init__(self, success: bool = ..., result: _Optional[str] = ...) -> None: ...

class GetResultRequest(_message.Message):
    __slots__ = ["task_id"]
    TASK_ID_FIELD_NUMBER: _ClassVar[int]
    task_id: int
    def __init__(self, task_id: _Optional[int] = ...) -> None: ...

class SingleTrainReply(_message.Message):
    __slots__ = ["result", "success"]
    class ResultEntry(_message.Message):
        __slots__ = ["key", "value"]
        KEY_FIELD_NUMBER: _ClassVar[int]
        VALUE_FIELD_NUMBER: _ClassVar[int]
        key: str
        value: float
        def __init__(self, key: _Optional[str] = ..., value: _Optional[float] = ...) -> None: ...
    RESULT_FIELD_NUMBER: _ClassVar[int]
    SUCCESS_FIELD_NUMBER: _ClassVar[int]
    result: _containers.ScalarMap[str, float]
    success: bool
    def __init__(self, success: bool = ..., result: _Optional[_Mapping[str, float]] = ...) -> None: ...

class SingleTrainRequest(_message.Message):
    __slots__ = ["df", "freq", "label_column", "metrics", "model", "test_indices", "train_indices"]
    DF_FIELD_NUMBER: _ClassVar[int]
    FREQ_FIELD_NUMBER: _ClassVar[int]
    LABEL_COLUMN_FIELD_NUMBER: _ClassVar[int]
    METRICS_FIELD_NUMBER: _ClassVar[int]
    MODEL_FIELD_NUMBER: _ClassVar[int]
    TEST_INDICES_FIELD_NUMBER: _ClassVar[int]
    TRAIN_INDICES_FIELD_NUMBER: _ClassVar[int]
    df: bytes
    freq: str
    label_column: str
    metrics: bytes
    model: bytes
    test_indices: bytes
    train_indices: bytes
    def __init__(self, model: _Optional[bytes] = ..., df: _Optional[bytes] = ..., train_indices: _Optional[bytes] = ..., test_indices: _Optional[bytes] = ..., label_column: _Optional[str] = ..., metrics: _Optional[bytes] = ..., freq: _Optional[str] = ...) -> None: ...

class TrainerRegisterReply(_message.Message):
    __slots__ = ["data_partition", "data_source", "model_season_lengths", "models", "success"]
    DATA_PARTITION_FIELD_NUMBER: _ClassVar[int]
    DATA_SOURCE_FIELD_NUMBER: _ClassVar[int]
    MODELS_FIELD_NUMBER: _ClassVar[int]
    MODEL_SEASON_LENGTHS_FIELD_NUMBER: _ClassVar[int]
    SUCCESS_FIELD_NUMBER: _ClassVar[int]
    data_partition: str
    data_source: str
    model_season_lengths: int
    models: str
    success: bool
    def __init__(self, success: bool = ..., data_source: _Optional[str] = ..., data_partition: _Optional[str] = ..., model_season_lengths: _Optional[int] = ..., models: _Optional[str] = ...) -> None: ...

class TrainerRegisterRequest(_message.Message):
    __slots__ = ["index"]
    INDEX_FIELD_NUMBER: _ClassVar[int]
    index: int
    def __init__(self, index: _Optional[int] = ...) -> None: ...
