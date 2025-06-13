package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"

	"github.com/spf13/cobra"
)

type CSEResponse struct {
	Kind string `json:"kind"`
	URL  struct {
		Type     string `json:"type"`
		Template string `json:"template"`
	} `json:"url"`
	Queries struct {
		Request []struct {
			Title          string `json:"title"`
			TotalResults   string `json:"totalResults"`
			SearchTerms    string `json:"searchTerms"`
			Count          int    `json:"count"`
			StartIndex     int    `json:"startIndex"`
			InputEncoding  string `json:"inputEncoding"`
			OutputEncoding string `json:"outputEncoding"`
			Safe           string `json:"safe"`
			Cx             string `json:"cx"`
			Filter         string `json:"filter"`
		} `json:"request"`
		NextPage []struct {
			Title          string `json:"title"`
			TotalResults   string `json:"totalResults"`
			SearchTerms    string `json:"searchTerms"`
			Count          int    `json:"count"`
			StartIndex     int    `json:"startIndex"`
			InputEncoding  string `json:"inputEncoding"`
			OutputEncoding string `json:"outputEncoding"`
			Safe           string `json:"safe"`
			Cx             string `json:"cx"`
			Filter         string `json:"filter"`
		} `json:"nextPage"`
	} `json:"queries"`
	Context struct {
		Title string `json:"title"`
	} `json:"context"`
	SearchInformation struct {
		SearchTime            float64 `json:"searchTime"`
		FormattedSearchTime   string  `json:"formattedSearchTime"`
		TotalResults          string  `json:"totalResults"`
		FormattedTotalResults string  `json:"formattedTotalResults"`
	} `json:"searchInformation"`
	Items []struct {
		Kind             string `json:"kind"`
		Title            string `json:"title"`
		HTMLTitle        string `json:"htmlTitle"`
		Link             string `json:"link"`
		DisplayLink      string `json:"displayLink"`
		Snippet          string `json:"snippet"`
		HTMLSnippet      string `json:"htmlSnippet"`
		FormattedURL     string `json:"formattedUrl"`
		HTMLFormattedURL string `json:"htmlFormattedUrl"`
		Pagemap          struct {
			CseThumbnail []struct {
				Src    string `json:"src"`
				Width  string `json:"width"`
				Height string `json:"height"`
			} `json:"cse_thumbnail"`
			Metatags []struct {
				Viewport                   string `json:"viewport"`
				FacebookDomainVerification string `json:"facebook-domain-verification"`
			} `json:"metatags"`
			CseImage []struct {
				Src string `json:"src"`
			} `json:"cse_image"`
		} `json:"pagemap"`
	} `json:"items"`
}

var (
	GOOGLE_API_KEY   string
	CSE_ID           string
	DOMAIN           string
	EXCLUDED_DOMAINS []string
)

func main() {
	var rootCmd = &cobra.Command{
		Use:   "google-index-checker",
		Short: "Check for accidentally indexed subdomains by Google",
		Long:  "A tool to check for accidentally indexed subdomains by Google.",
		Run:   scanGoogleIndex,
	}

	rootCmd.Flags().StringVarP(&GOOGLE_API_KEY, "key", "k", os.Getenv("GOOGLE_API_KEY"), "Google API Key for Custom Search Engine")
	rootCmd.Flags().StringVarP(&CSE_ID, "cse", "c", os.Getenv("GOOGLE_CSE_ID"), "Google Custom Search Engine ID")

	rootCmd.Flags().StringVarP(&DOMAIN, "domain", "d", "", "Base domain to check for indexing (e.g., example.com)")
	err := rootCmd.MarkFlagRequired("domain")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error marking domain flag as required: %v\n", err)
		os.Exit(1)
	}
	rootCmd.Flags().StringArrayVarP(&EXCLUDED_DOMAINS, "not", "n", []string{}, "Subdomains to not check (e.g., www.example.com, wiki.example.com)")

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func CreateQueryParams(domain string, excludedDomains []string) string {
	query := fmt.Sprintf("site:%s", domain)
	for _, subdomain := range excludedDomains {
		query += fmt.Sprintf(" -site:%s", subdomain)
	}
	return query
}

func scanGoogleIndex(cmd *cobra.Command, args []string) {
	query := CreateQueryParams(DOMAIN, EXCLUDED_DOMAINS)
	fmt.Println("Scanning Google Search Index:", query)

	indexedDomains := make(map[string]bool)

	offset := 0
	for {
		resp, err := queryCSE(query, offset)
		if err != nil {
			fmt.Printf("Error querying Google Custom Search Engine: %v\n", err)
			break
		}

		for _, item := range resp.Items {
			u, err := url.Parse(item.Link)
			if err != nil {
				panic(err)
			}
			indexedDomains[u.Host] = true
		}

		if resp.Queries.NextPage == nil {
			break
		}

		offset = resp.Queries.NextPage[0].StartIndex
	}

	fmt.Printf("\n\nScan completed. Found indexed domains:\n")
	for k := range indexedDomains {
		fmt.Println("Indexed Domain: ", k)
	}
}

func queryCSE(query string, offset int) (results CSEResponse, err error) {
	fmt.Printf("Querying Google Custom Search Engine for: %s (Start Index: %d)\n", query, offset)

	req, err := http.NewRequest("GET", "https://customsearch.googleapis.com/customsearch/v1?", nil)
	if err != nil {
		return
	}

	q := req.URL.Query()
	q.Add("key", GOOGLE_API_KEY)
	q.Add("cx", CSE_ID)
	q.Add("q", query)
	q.Add("start", strconv.Itoa(offset))
	req.URL.RawQuery = q.Encode()

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return
	}

	defer resp.Body.Close() //nolint:all
	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("invalid status code: %d", resp.StatusCode)
		return
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}

	err = json.Unmarshal(body, &results)
	if err != nil {
		return
	}

	return
}
