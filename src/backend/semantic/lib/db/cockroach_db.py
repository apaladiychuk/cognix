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


class DocumentCRUD:
    def __init__(self, connection_string):
        self.engine = create_engine(connection_string)
        Session = sessionmaker(bind=self.engine)
        self.session = Session()
        Base.metadata.create_all(self.engine)

    def insert_document(self, **kwargs):
        new_document = Document(**kwargs)
        self.session.add(new_document)
        self.session.commit()
        return new_document.id

    def select_document(self, document_id):
        return self.session.query(Document).filter_by(id=document_id).first()

    def update_document(self, document_id, **kwargs):
        self.session.query(Document).filter_by(id=document_id).update(kwargs)
        self.session.commit()

    def delete_document(self, document_id):
        self.session.query(Document).filter_by(id=document_id).delete()
        self.session.commit()

