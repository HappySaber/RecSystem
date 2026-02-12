class Movie(Base):
    id = Column(Integer, primary_key=True)
    overview = Column(Text)
    tags_status = Column(String)

class Tag(Base):
    __tablename__ = "tags"

    id = Column(Integer, primary_key=True)
    name = Column(String, nullable=False)
    slug = Column(String, nullable=False, unique=True)

class MovieTag(Base):
    __tablename__ = "tags_content"

    movie_id = Column(Integer, ForeignKey("movies.id"), primary_key=True)
    tag_id = Column(Integer, ForeignKey("tags.id"), primary_key=True)
    weight = Column(Float, nullable=False)
