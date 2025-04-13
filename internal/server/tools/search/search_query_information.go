package search

const SearchQueryInformation = `
// Matching
//
// Matching of search keywords is NOT case sensitive.  Keywords are matched in
// any order unless surrounded by double quotes. Searching for playlists will
// return results where the query keyword(s) match any part of the playlist's
// name or description. Only popular public playlists are returned.
//
// Operators
//
// The operator NOT can be used to exclude results.  For example,
// query = "roadhouse NOT blues" returns items that match "roadhouse" but excludes
// those that also contain the keyword "blues".  Similarly, the OR operator can
// be used to broaden the search.  query = "roadhouse OR blues" returns all results
// that include either of the terms.  Only one OR operator can be used in a query.
//
// Operators should be specified in uppercase.
//
// Wildcards
//
// The asterisk (*) character can, with some limitations, be used as a wildcard
// (maximum of 2 per query).  It will match a variable number of non-white-space
// characters.  It cannot be used in a quoted phrase, in a field filter, or as
// the first character of a keyword string.
//
`

const FilterInformation = `
// Field filters
//
// By default, results are returned when a match is found in any field of the
// target object type.  Searches can be made more specific by specifying an album,
// artist, or track field filter.  For example, "album:gold artist:abba type:album"
// will only return results with the text "gold" in the album name and the text
// "abba" in the artist's name.
//
// The field filter "year" can be used with album, artist, and track searches to
// limit the results to a particular year. For example "bob year:2014" or
// "bob year:1980-2020".
//
// The field filter "tag:new" can be used in album searches to retrieve only
// albums released in the last two weeks. The field filter "tag:hipster" can be
// used in album searches to retrieve only albums with the lowest 10% popularity.
//
// Other possible field filters, depending on object types being searched,
// include "genre", "upc", and "isrc".  For example "damian genre:reggae-pop".
//
`
