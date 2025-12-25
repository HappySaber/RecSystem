package anilist

const popularAnimeQuery = `
query ($page: Int) {
  Page(page: $page, perPage: 20){
    media(type: ANIME, sort: POPULARITY_DESC) {
      id
      idMal

      title {
        romaji
        english
        native
      }

      
      
      description(asHtml: false)

      format
      status
      season
      seasonYear

      episodes
      duration

      startDate {
        year
        month
        day
      }

      endDate {
        year
        month
        day
      }

      genres
      tags {
        name
      }

      studios (isMain: true){
        nodes{
          id
          name
        }
      }

      characters (role: MAIN, perPage: 5){
        nodes{
          name{
            full
          }
        }
      }

     #staff( perPage: 5) {
        #nodes {
        #  name {
         #   full
        #  }
      #  }
    #  }

      
      meanScore
      popularity
      favourites

      coverImage {
        large
      }

      trailer {
        site
        id
      }

    }
  }
}
`
