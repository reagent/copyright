package app

import "regexp"

type Route struct {
	Pattern *regexp.Regexp
	Method  string
	Handler Handler
}

type Routes []Route

func (r *Route) MatchURL(target string) (map[string]string, bool) {
	mapping := map[string]string{}
	matches := r.Pattern.FindStringSubmatch(target)

	if len(matches) == 0 {
		return mapping, false
	}

	for idx, name := range r.Pattern.SubexpNames() {
		if len(name) > 0 {
			mapping[name] = matches[idx]
		}
	}

	return mapping, true
}
