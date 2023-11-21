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
		if len(pathList) != 3 {
			log.Fatalf("too many args in %q", parsedUrl.Path)
		}
		return fmt.Sprintf("https://container.googleapis.com/v1/projects/%s/locations/%s/clusters/%s", parsedUrl.Host, pathList[1], pathList[2])
	default:
		log.Fatalf("unimplemented scheme %q", parsedUrl.Scheme)
	}

	return ""
}

func GetSubFromURI(uri string) (string, string, string) {
	parsedUrl, err := url.Parse(uri)

	if err != nil {
		log.Fatal(err)
	}

	return parsedUrl.Query().Get("ns"), parsedUrl.Query().Get("sa"), parsedUrl.Query().Get("po")
}
