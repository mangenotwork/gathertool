/*
	Description : 对 github.com/PuerkitoBio/goquery的包装
	Author : ManGe
	Version : v0.1
	Date : 2021-04-27
*/

package gathertool

import (
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func NewGoquery(html string) (*goquery.Document, error){
	return goquery.NewDocumentFromReader(strings.NewReader(html))
}