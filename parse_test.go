package movabletype_test

import (
	"bytes"
	"io"
	"reflect"
	"testing"
	"time"

	. "github.com/catatsuy/movabletype"
)

func TestParse(t *testing.T) {
	buf := bytes.NewBufferString(`AUTHOR: catatsuy
TITLE: ポエム
BASENAME: poem
STATUS: Publish
ALLOW COMMENTS: 1
ALLOW PINGS: 1
CONVERT BREAKS: 0
DATE: 04/22/2017 20:41:58
PRIMARY CATEGORY: ブログ
CATEGORY: ポエム
CATEGORY: 技術系
-----
BODY:
<p>body</p>
<p>bodybody</p>
<p>bodybodybody</p>
-----
EXTENDED BODY:
<p>extended body</p>
<p>extended body body</p>
<p>extended body body body</p>
-----
EXCERPT:
ここに概要が表示されます。
-----
--------
AUTHOR: catatsuy
TITLE: 風邪で声を失った話
BASENAME: 2017/04/09/194939
STATUS: Publish
ALLOW COMMENTS: 1
CONVERT BREAKS: 0
DATE: 04/09/2017 07:49:39 PM
CATEGORY: 日常
-----
BODY:
<p>bodybodybody</p>
-----
EXTENDED BODY:
<p>extended body body body</p>
-----
EXCERPT:
ここに概要が表示されます。
-----
--------
`)

	expected := []*Entry{
		&Entry{
			Author:          "catatsuy",
			Title:           "ポエム",
			Basename:        "poem",
			Status:          "Publish",
			AllowComments:   1,
			AllowPings:      1,
			ConvertBreaks:   "0",
			Date:            time.Date(2017, time.April, 22, 20, 41, 58, 0, time.UTC),
			PrimaryCategory: "ブログ",
			Category:        []string{"ポエム", "技術系"},
			Body:            "<p>body</p>\n<p>bodybody</p>\n<p>bodybodybody</p>\n",
			ExtendedBody:    "<p>extended body</p>\n<p>extended body body</p>\n<p>extended body body body</p>\n",
			Excerpt:         "ここに概要が表示されます。\n",
		},
		&Entry{
			Author:        "catatsuy",
			Title:         "風邪で声を失った話",
			Basename:      "2017/04/09/194939",
			Status:        "Publish",
			AllowComments: 1,
			AllowPings:    -1,
			ConvertBreaks: "0",
			Date:          time.Date(2017, time.April, 9, 19, 49, 39, 0, time.UTC),
			Category:      []string{"日常"},
			Body:          "<p>bodybodybody</p>\n",
			ExtendedBody:  "<p>extended body body body</p>\n",
			Excerpt:       "ここに概要が表示されます。\n",
		},
	}

	mts, err := Parse(buf)

	if err != nil {
		t.Fatalf("got error %q", err)
	}

	if !reflect.DeepEqual(mts, expected) {
		t.Errorf("Error parsing, expected %#v; got %#v", expected, mts)
	}
}

func TestParseStatusNotAllowed(t *testing.T) {
	buf := bytes.NewBufferString(`STATUS: Published`)

	_, err := Parse(buf)

	if err == nil || err.Error() != "STATUS column is allowed only Draft or Publish or Future. Got Published" {
		t.Errorf("Error parsing, got %q", err)
	}
}

func TestParseDate(t *testing.T) {
	var featuretests = []struct {
		buf io.Reader
		t   time.Time
	}{
		{
			bytes.NewBufferString("DATE: 04/22/2017 08:41:58 PM\n--------\n"),
			time.Date(2017, time.April, 22, 20, 41, 58, 0, time.UTC),
		},
		{
			bytes.NewBufferString("DATE: 04/22/2017 08:41:58 AM\n--------\n"),
			time.Date(2017, time.April, 22, 8, 41, 58, 0, time.UTC),
		},
		{
			bytes.NewBufferString("DATE: 04/22/2017 20:41:58\n--------\n"),
			time.Date(2017, time.April, 22, 20, 41, 58, 0, time.UTC),
		},
	}

	for _, ft := range featuretests {
		mts, err := Parse(ft.buf)

		if err != nil {
			t.Fatalf("got error %q", err)
		}

		if mts[0].Date != ft.t {
			t.Errorf("m.Date got %v; want %v", mts[0].Date, ft.t)
		}
	}
}

func TestNewMT(t *testing.T) {
	m := NewEntry()

	if m.AllowComments != DefaultAllowComments {
		t.Errorf("By default, AllowComments is %d, got %d", DefaultAllowComments, m.AllowComments)
	}

	if m.AllowPings != DefaultAllowPings {
		t.Errorf("By default, AllowComments is %d, got %d", DefaultAllowPings, m.AllowPings)
	}
}
