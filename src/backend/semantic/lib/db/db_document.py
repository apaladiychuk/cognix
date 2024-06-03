from sqlalchemy import create_engine, Column, Integer, BigInteger, Text, Boolean, UUID, TIMESTAMP, ForeignKey, func
from sqlalchemy.ext.declarative import declarative_base
from sqlalchemy.orm import sessionmaker
import uuid

Base = declarative_base()


class Document(Base):
    __tablename__ = 'documents'

    id = Column(Integer, primary_key=True, autoincrement=True)
    parent_id = Column(BigInteger, ForeignKey('documents.id'), nullable=True)
    connector_id = Column(BigInteger, nullable=False)
    source_id = Column(Text, nullable=False)
    url = Column(Text, nullable=True)
    signature = Column(Text, nullable=True)
    chunking_session = Column(UUID(as_uuid=True), nullable=True)
    analyzed = Column(Boolean, nullable=False, default=False)
    creation_date = Column(TIMESTAMP(timezone=False), nullable=False, default=func.now())
    last_update = Column(TIMESTAMP(timezone=False), nullable=True)

    def __repr__(self):
        return (f"<Document(id={self.id}, parent_id={self.parent_id}, connector_id={self.connector_id}, "
                f"source_id={self.source_id}, url={self.url}, signature={self.signature}, "
                f"chunking_session={self.chunking_session}, analyzed={self.analyzed}, "
                f"creation_date={self.creation_date}, last_update={self.last_update})>")


class DocumentCRUD:
    def __init__(self, connection_string):
        self.engine = create_engine(connection_string)
        Session = sessionmaker(bind=self.engine)
        self.session = Session()
        Base.metadata.create_all(self.engine)

    def insert_document(self, **kwargs) -> int:
        new_document = Document(**kwargs)
        self.session.add(new_document)
        self.session.commit()
        return new_document.id

    def select_document(self, document_id) -> Document | None:
        if document_id <= 0:
            raise ValueError("ID value must be positive")
        return self.session.query(Document).filter_by(id=document_id).first()

    def update_document(self, document_id, **kwargs) -> int:
        if document_id <= 0:
            raise ValueError("ID value must be positive")
        updated_docs = self.session.query(Document).filter_by(id=document_id).update(kwargs)
        self.session.commit()
        return updated_docs

    def delete_document(self, document_id) -> int:
        if document_id <= 0:
            raise ValueError("ID value must be positive")
        deleted_docs = self.session.query(Document).filter_by(id=document_id).delete()
        self.session.commit()
        return deleted_docs
