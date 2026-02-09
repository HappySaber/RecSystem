BLACKLIST = ["film", "movie", "story"]

def normalize(tags):
    result = []
    for tag, weight in tags:
        tag = tag.lower().strip()
        if tag in BLACKLIST:
            continue
        if len(tag) < 3:
            continue
        result.append((tag, weight))
    return result