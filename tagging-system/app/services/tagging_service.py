class TaggingService:
    def tag_movie(self, movie):
        keywords = extractor.extract(movie.overview)
        normalized = normalizer.clean(keywords)
        weighted = weighting.apply(normalized)
        return weighted