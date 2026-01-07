package parsers

import (
	"time"
)

const DateFormat = "2006-01-02"

type DateParser struct{}

func NewDateParser() *DateParser {
	return &DateParser{}
}

func (p *DateParser) Parse(dateStr string) (*time.Time, error) {
	if dateStr == "" {
		return nil, nil
	}
	t, err := time.Parse(DateFormat, dateStr)
	if err != nil {
		return nil, err
	}
	return &t, nil
}
