from google.protobuf.internal import containers as _containers
from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from typing import ClassVar as _ClassVar, Iterable as _Iterable, Optional as _Optional

DESCRIPTOR: _descriptor.FileDescriptor

class EmbeddRequest(_message.Message):
    __slots__ = ("content", "model")
    CONTENT_FIELD_NUMBER: _ClassVar[int]
    MODEL_FIELD_NUMBER: _ClassVar[int]
    content: str
    model: str
    def __init__(self, content: _Optional[str] = ..., model: _Optional[str] = ...) -> None: ...

class EmbeddResponse(_message.Message):
    __slots__ = ("vector",)
    VECTOR_FIELD_NUMBER: _ClassVar[int]
    vector: _containers.RepeatedScalarFieldContainer[float]
    def __init__(self, vector: _Optional[_Iterable[float]] = ...) -> None: ...

class EmbeddData(_message.Message):
    __slots__ = ("id", "content", "model", "vector")
    ID_FIELD_NUMBER: _ClassVar[int]
    CONTENT_FIELD_NUMBER: _ClassVar[int]
    MODEL_FIELD_NUMBER: _ClassVar[int]
    VECTOR_FIELD_NUMBER: _ClassVar[int]
    id: int
    content: str
    model: str
    vector: _containers.RepeatedScalarFieldContainer[float]
    def __init__(self, id: _Optional[int] = ..., content: _Optional[str] = ..., model: _Optional[str] = ..., vector: _Optional[_Iterable[float]] = ...) -> None: ...
