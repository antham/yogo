package mail

import (
	"os"
	"testing"

	"github.com/PuerkitoBio/goquery"
	"github.com/antham/yogo/v3/internal/client"
	"github.com/stretchr/testify/assert"
)

func getDoc[M client.MailDoc](t *testing.T, filename string) M {
	dir, err := os.Getwd()
	if err != nil {
		assert.NoError(t, err)
	}

	f, err := os.Open(dir + "/features/" + filename)
	if err != nil {
		assert.NoError(t, err)
	}

	doc, err := goquery.NewDocumentFromReader(f)
	if err != nil {
		assert.NoError(t, err)
	}

	return M(*doc)
}
