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
