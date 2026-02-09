from keybert import KeyBERT

class KeywordExtractor:
    def __init__(self):
        self.model = KeyBERT("all-MiniLM-L6-v2")

    def extract(self, text: str):
        return self.model.extract_keywords(
            text,
            keyphrase_ngram_range=(1, 2),
            stop_words="english",
            top_n=10
        )