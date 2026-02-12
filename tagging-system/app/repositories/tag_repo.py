class TagRepository:
    def get_by_slug(self, db, slug: str):
        ...

    def create(self, db, name: str, slug: str):
        ...


class MovieTagRepository:
    def upsert(self, db, movie_id, tag_id, weight):
        ...
