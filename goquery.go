package gathertool

import (
	"github.com/PuerkitoBio/goquery"
	"strings"
)

func NewGoquery(html string) (*goquery.Document, error){
	return goquery.NewDocumentFromReader(strings.NewReader(html))
}