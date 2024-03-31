package routes

import (
	"sort"
	"strings"

	"distrifs.dev/indexer/modules/globals"
	"github.com/gin-gonic/gin"
)

type SearchQuery struct {
	Query  string `form:"query" binding:"required"`
	Server string `form:"server" binding:"required"`
}

type CharTarget struct {
	Match    bool
	Original string
	Target   string
}

type FileTarget struct {
	PercentMatched float32
	CharsMatched   int
	CharsMax       int
	Search         string
	Hash           string
	Chars          []CharTarget
}

type FileTargetSearch struct {
	RecursiveResult globals.HashItem
	BreakResult     []FileTarget
}

func searchFileMatch(search string, target string, hash string) FileTarget {
	returnedTarget := FileTarget{}
	// break target up into pieces
	subtarget := strings.Split(target, "")
	subsearch := strings.Split(search, "")
	amountmatched := 0

	// determine greater one of two
	sublen := len(subsearch)
	tarlen := len(subtarget)
	var slicelen int

	if sublen > tarlen {
		slicelen = tarlen
	} else {
		slicelen = sublen
	}

	matched := make([]CharTarget, slicelen)

out:
	for i, v := range subtarget {
		if i > (len(subsearch) - 1) {
			break out
		}

		same := (v == subsearch[i])

		if same {
			amountmatched++
		}

		matched[i] = CharTarget{
			Match:    same,
			Original: subsearch[i],
			Target:   v,
		}
	}

	percentage := float32(amountmatched) / float32(len(subtarget))

	returnedTarget.PercentMatched = percentage
	returnedTarget.Chars = matched
	returnedTarget.CharsMatched = amountmatched
	returnedTarget.CharsMax = len(subtarget)
	returnedTarget.Search = search
	returnedTarget.Hash = hash

	return returnedTarget
}

func searchFileWalkRecursion(v map[string]globals.HashItem, target string) []FileTarget {
	var fileListing []FileTarget
	for _, v := range v {
		if v.FileName != "" {
			fileListing = append(fileListing, searchFileMatch(v.FileName, target, v.Hash))
		}
		if v.SubFiles != nil {
			fnOutput := searchFileWalkRecursion(v.SubFiles, target)
			fileListing = append(fileListing, fnOutput...)
		}
	}
	return fileListing
}

func searchFileWalk(v map[string]globals.HashItem, target string) []FileTarget {
	// fmt.Printf("Visiting %v\n", v)

	fileListing := searchFileWalkRecursion(v, target)

	var sortedListing []FileTarget

	// sort files by percent and remove less than 10% matches
	for _, v := range fileListing {
		// remove if the percent is less than 15%
		if v.PercentMatched < 0.10 {
			continue
		}
		sortedListing = append(sortedListing, v)
	}

	sort.Slice(sortedListing[:], func(i, j int) bool {
		return sortedListing[i].PercentMatched > sortedListing[j].PercentMatched
	})

	return sortedListing
}

func SearchFileRoute(c *gin.Context) {
	// get file by hash

	var hashQuery SearchQuery
	if err := c.ShouldBind(&hashQuery); err != nil {
		c.JSON(400, gin.H{
			"success": false,
			"error":   "bad_uri_request",
		})
		return
	}

	server := globals.IndexerData[hashQuery.Server]

	if !server.Online {
		c.JSON(500, gin.H{
			"success": false,
			"data":    "server_not_online",
		})
		return
	}

	file := searchFileWalk(server.Files.SubFiles, hashQuery.Query)
	// fmt.Println(file)

	// if file.Hash == "" {
	// 	c.JSON(500, gin.H{
	// 		"success": false,
	// 		"data":    "no_file_found",
	// 	})
	// 	return
	// }

	c.JSON(200, gin.H{
		"success": true,
		"data":    file,
	})
}
