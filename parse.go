package movabletype

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"
)

// Default
const (
	// If it is not inialized, AllowComments is -1
	DefaultAllowComments = -1

	// If it is not inialized, AllowPings is -1
	DefaultAllowPings = -1
)

// Movable Type Import Format
type Entry struct {
	Author   string
	Title    string
	Basename string
	Status   string

	// 0 or 1. If it is not inialized DefaultAllowComments.
	AllowComments int

	// 0 or 1. If it is not inialized DefaultAllowPings
	AllowPings int

	ConvertBreaks string

	Date time.Time

	PrimaryCategory string

	Category []string

	Body string

	ExtendedBody string
}

// NewMT creates MT.
func NewEntry() *Entry {
	return &Entry{
		AllowComments: DefaultAllowComments,
		AllowPings:    DefaultAllowPings,
	}
}

// Parse creates MT struct from io.Reader
func Parse(r io.Reader) ([]*Entry, error) {
	mts := []*Entry{}

	scanner := bufio.NewScanner(r)

	var err error

	m := NewEntry()

	for scanner.Scan() {
		ss := strings.Split(scanner.Text(), ": ")

		if len(ss) <= 1 {
			value := ss[0]

			if value == "--------" {
				mts = append(mts, m)
				m = NewEntry()
				continue
			}

			if value == "-----" {
				continue
			}

			switch value {
			case "BODY:":
				for scanner.Scan() {
					line := scanner.Text()

					if line == "-----" {
						break
					}

					m.Body += line + "\n"
				}
				break
			case "EXTENDED BODY:":
				for scanner.Scan() {
					line := scanner.Text()

					if line == "-----" {
						break
					}

					m.ExtendedBody += line + "\n"
				}
				break
			}

			continue
		}

		key, value := ss[0], ss[1]

		switch key {
		case "AUTHOR":
			m.Author = value
			break
		case "TITLE":
			m.Title = value
			break
		case "BASENAME":
			m.Basename = value
			break
		case "STATUS":
			if value == "Draft" || value == "Publish" || value == "Future" {
				m.Status = value
			} else {
				return nil, fmt.Errorf("STATUS column is allowed only Draft or Publish or Future. Got %s", value)
			}
			break
		case "ALLOW COMMENTS":
			m.AllowComments, err = strconv.Atoi(value)
			if err != nil {
				return nil, errors.Wrap(err, "ALLOW COMMENTS column is allowed only 0 or 1")
			}
			if m.AllowComments != 0 && m.AllowComments != 1 {
				return nil, fmt.Errorf("ALLOW COMMENTS column is allowed only 0 or 1. Got %d", m.AllowComments)
			}
			break
		case "ALLOW PINGS":
			m.AllowPings, err = strconv.Atoi(value)
			if err != nil {
				return nil, errors.Wrap(err, "ALLOW PINGS column is allowed only 0 or 1")
			}
			if m.AllowComments != 0 && m.AllowComments != 1 {
				return nil, fmt.Errorf("ALLOW PINGS column is allowed only 0 or 1. Got %d", m.AllowPings)
			}
			break
		case "CONVERT BREAKS":
			m.ConvertBreaks = value
			break
		case "DATE":
			if strings.HasSuffix(value, "AM") || strings.HasSuffix(value, "PM") {
				m.Date, err = time.Parse("01/02/2006 03:04:05 PM", value)
			} else {
				m.Date, err = time.Parse("01/02/2006 15:04:05", value)
			}
			if err != nil {
				return nil, errors.Wrap(err, "Parsing error on DATE column")
			}
			break
		case "PRIMARY CATEGORY":
			m.PrimaryCategory = value
			break
		case "CATEGORY":
			m.Category = append(m.Category, value)
			break
		}
	}

	return mts, nil
}
