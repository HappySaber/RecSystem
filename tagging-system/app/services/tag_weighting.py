def normalize_weights(tags):
    if not tags:
        return []
    
    max_weight = max(weight for _, weight in tags)
    return [(t, round(weight / max_weight, 3)) for t, w in tags]