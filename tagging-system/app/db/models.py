class Movie(Base):
    id = Column(Integer, primary_key=True)
    overview = Column(Text)
    tags_status = Column(String)