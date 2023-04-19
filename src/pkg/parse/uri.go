package parse

import (
	"fmt"
	"log"
	"net/url"
	"strings"
)

func GetIssuerFromUri(uri string) string {
	parsedUrl, err := url.Parse(uri)

	if err != nil {
		log.Fatal(err)
	}

	switch parsedUrl.Scheme {
	case "gke":
		pathList := strings.Split(parsedUrl.Path, "/")
		if len(pathList) != 4 {
			log.Fatalf("too many args in %q", parsedUrl.Path)
		}
		return fmt.Sprintf("https://container.googleapis.com/v1/projects/%s/locations/%s/clusters/%s", pathList[1], pathList[2], pathList[3])
	default:
		log.Fatalf("unimplemented scheme %q", parsedUrl.Scheme)
	}

	return ""
}
