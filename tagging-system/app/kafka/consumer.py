from app.services.tagging import TaggingService

def handle_message(message):
    event = parse_event(message)

    if event.type == "MovieImported":
        TaggingService.process_movie(event)
