package resp

import (
	"fmt"
	"strconv"
	"strings"
)

// Reply ...
type Reply struct {
	Err     error // 只有run runWithLock 会在这里带有值，其他情况不用判断
	Null    bool
	Str     string
	Integer int64
	Replies []*Reply
}

// String just for test
func (p *Reply) String() string {
	if p.Err != nil {
		return fmt.Sprintf("Err: %v", p.Err)
	}
	if p.Null {
		return "<nil>"
	}
	if p.Str != "" {
		return fmt.Sprintf("<String: %v>", p.Str)
	}
	if p.Integer != 0 {
		return fmt.Sprintf("<Int: %v>", p.Integer)
	}
	if len(p.Replies) > 0 {
		var s []string
		for _, v := range p.Replies {
			s = append(s, v.String())
		}
		return fmt.Sprintf("(List: %s)", strings.Join(s, ", "))
	}
	return "<empty>"
}

func (p *Reply) NullInteger() (int64, error) {
	if p.Err != nil {
		return 0, p.Err
	} else if p.Null {
		return 0, ErrKeyNotExist
	}
	return p.Integer, nil
}

func (p *Reply) NullString() (NullString, error) {
	if p.Err != nil {
		return NullString{}, p.Err
	}
	if p.Null {
		return NullString{}, nil
	}

	return NullString{String: p.Str, Valid: true}, nil
}

func (p *Reply) Bool() (bool, error) {
	if p.Err == nil {
		return p.Integer == 1, nil
	}
	return false, p.Err
}

func (p *Reply) Float64() (float64, error) {
	if p.Err != nil {
		return 0, p.Err
	}
	return strconv.ParseFloat(p.Str, 64)
}

func (p *Reply) NullStringSlice() ([]NullString, error) {
	if p.Err != nil {
		return nil, p.Err
	}

	var ns []NullString
	for _, v := range p.Replies {
		if v.Err != nil {
			return nil, v.Err // TODO 这里真的有error吗
		}
		n, _ := v.NullString()
		ns = append(ns, n)
	}
	return ns, nil
}

func (p *Reply) StringSlice() ([]string, error) {
	if p.Err != nil {
		return nil, p.Err
	}

	var s []string
	for _, v := range p.Replies {
		if v.Err != nil {
			return nil, v.Err // TODO 真的有吗
		}
		s = append(s, v.Str)
	}
	return s, nil
}

func (p *Reply) Map() (map[string]string, error) {
	if p.Err != nil {
		return nil, p.Err
	}

	var s = make(map[string]string)
	for i := 0; i+1 < len(p.Replies); i += 2 {
		if p.Replies[i].Err != nil {
			return nil, p.Replies[i].Err // TODO 真的有吗
		}
		if p.Replies[i+1].Err != nil {
			return nil, p.Replies[i+1].Err // TODO 真的有吗
		}
		s[p.Replies[i].Str] = p.Replies[i+1].Str
	}
	return s, nil
}

func (p *Reply) GeoLocationSlice() ([]*GeoLocation, error) {
	if p.Err != nil {
		return nil, p.Err
	}
	var ss []*GeoLocation
	for _, v := range p.Replies {
		if v.Err != nil {
			return nil, v.Err
		}
		if len(v.Replies) < 2 {
			return nil, fmt.Errorf("expect 2 string to parse to geo")
		}
		longitude, err := strconv.ParseFloat(v.Replies[0].Str, 64)
		if err != nil {
			return nil, err
		}
		latitude, err := strconv.ParseFloat(v.Replies[1].Str, 64)
		if err != nil {
			return nil, err
		}
		ss = append(ss, &GeoLocation{Longitude: longitude, Latitude: latitude})
	}
	return ss, nil
}

func (p *Reply) SortedSetSlice() ([]*SortedSet, error) {
	if p.Err != nil {
		return nil, p.Err
	}
	var ss []*SortedSet
	for _, v := range p.Replies {
		if v.Err != nil {
			return nil, v.Err
		}
		ss = append(ss, &SortedSet{Member: v.Str})
	}
	return ss, nil
}

func (p *Reply) SortedSetSliceWithScores() ([]*SortedSet, error) {
	if p.Err != nil {
		return nil, p.Err
	}
	var ss []*SortedSet
	for i := 0; i+1 < len(p.Replies); i += 2 {
		if p.Replies[i].Err != nil {
			return nil, p.Replies[i].Err // TODO 真的有吗
		}
		if p.Replies[i+1].Err != nil {
			return nil, p.Replies[i+1].Err // TODO 真的有吗
		}
		score, err := strconv.ParseFloat(p.Replies[i+1].Str, 64)
		if err != nil {
			return nil, err
		}
		ss = append(ss, &SortedSet{Member: p.Replies[i].Str, Score: score})
	}

	return ss, nil
}

func (p *Reply) Float() (float64, error) {
	if p.Err != nil {
		return 0, p.Err
	}
	return strconv.ParseFloat(p.Str, 64)
}
