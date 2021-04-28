package gathertool

import (
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func NewGoquery(html string) (*goquery.Document, error){
	return goquery.NewDocumentFromReader(strings.NewReader(html))
}