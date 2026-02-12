class TaggingService:
    def process_movie(self, event):
        with SessionLocal() as db:
            movie = movie_repo.get(db, event.movie_id)

            if not movie or not movie.overview:
                return

            text = self._build_text(movie)
            raw = extractor.extract(text)
            clean = normalize(raw)
            weighted = normalize_weights(clean)

            for tag, weight in weighted:
                slug = make_slug(tag)
                tag_entity = tag_repo.get_or_create(db, tag, slug)
                movie_tag_repo.upsert(db, movie.id, tag_entity.id, weight)

            movie.tags_status = "DONE"
            db.commit()
