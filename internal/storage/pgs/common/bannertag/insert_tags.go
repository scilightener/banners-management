package bannertag

import (
	"strconv"
	"strings"
)

// InsertTagsQuery returns an SQL query for creating tags with provided tagIDs.
func InsertTagsQuery(bannerID int64, tagIDs []int64) string {
	if len(tagIDs) == 0 {
		return ""
	}
	return `INSERT INTO banner_tag (banner_id, tag_id) VALUES ` + tagsQuery(bannerID, tagIDs)
}

func tagsQuery(bannerID int64, tagIDs []int64) string {
	var sb strings.Builder
	for i, tagID := range tagIDs {
		sb.WriteString("(")
		sb.WriteString(strconv.FormatInt(bannerID, 10))
		sb.WriteString(", ")
		sb.WriteString(strconv.FormatInt(tagID, 10))
		sb.WriteString(")")

		if i != len(tagIDs)-1 {
			sb.WriteString(", ")
		}
	}

	return sb.String()
}
