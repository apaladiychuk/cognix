from google.protobuf.internal import enum_type_wrapper as _enum_type_wrapper
from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from typing import ClassVar as _ClassVar, Optional as _Optional, Union as _Union

DESCRIPTOR: _descriptor.FileDescriptor

class FileType(int, metaclass=_enum_type_wrapper.EnumTypeWrapper):
    __slots__ = ()
    URL: _ClassVar[FileType]
    PDF: _ClassVar[FileType]
    RTF: _ClassVar[FileType]
    DOC: _ClassVar[FileType]
    XLS: _ClassVar[FileType]
    PPT: _ClassVar[FileType]
    TXT: _ClassVar[FileType]
    MD: _ClassVar[FileType]
URL: FileType
PDF: FileType
RTF: FileType
DOC: FileType
XLS: FileType
PPT: FileType
TXT: FileType
MD: FileType

class ChunkingData(_message.Message):
    __slots__ = ("url", "site_map", "search_for_sitemap", "document_id", "file_type")
    URL_FIELD_NUMBER: _ClassVar[int]
    SITE_MAP_FIELD_NUMBER: _ClassVar[int]
    SEARCH_FOR_SITEMAP_FIELD_NUMBER: _ClassVar[int]
    DOCUMENT_ID_FIELD_NUMBER: _ClassVar[int]
    FILE_TYPE_FIELD_NUMBER: _ClassVar[int]
    url: str
    site_map: str
    search_for_sitemap: bool
    document_id: int
    file_type: FileType
    def __init__(self, url: _Optional[str] = ..., site_map: _Optional[str] = ..., search_for_sitemap: bool = ..., document_id: _Optional[int] = ..., file_type: _Optional[_Union[FileType, str]] = ...) -> None: ...
