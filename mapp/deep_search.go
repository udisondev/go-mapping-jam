package mapp

func deepFieldSearch(f Field, fieldFullName string) (Field, bool) {
	if f.FullName() == fieldFullName {
		return f, true
	}

	fields := f.Fields()
	if len(fields) == 0 {
		return Field{}, false
	}

	for _, ff := range fields {
		expF, found := deepFieldSearch(ff, fieldFullName)
		if found {
			return expF, true
		}
	}

	return Field{}, false
}
