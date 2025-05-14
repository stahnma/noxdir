package structure

// FiltersList aggregates a list of multiple filters.
type FiltersList map[FilterName]EntryFilter

func NewFiltersList() FiltersList {
	return make(FiltersList)
}

// Valid checks the provided *Entry instance against all filters and returns a
// bool value. The "true" value will be returned if the entry passes all defined
// filters.
func (fl *FiltersList) Valid(e *Entry) bool {
	for _, filter := range *fl {
		if !filter.Filter(e) {
			return false
		}
	}

	return true
}

// ToggleFilter adds or removes the provided filter instance. It checks by the
// unique filter name if the FiltersList already contains it. If the filter
// already exists, it will be removed, otherwise, it will be added.
func (fl *FiltersList) ToggleFilter(ef EntryFilter) {
	if _, ok := (*fl)[ef.Name()]; ok {
		delete(*fl, ef.Name())

		return
	}

	(*fl)[ef.Name()] = ef
}

// FilterName defines a custom type for filter name. Each EntryFilter instance
// uses this type for its name.
type FilterName string

const (
	DirsOnlyFilter  FilterName = "DirsOnly"
	FilesOnlyFilter FilterName = "FilesOnly"
)

// EntryFilter defines a contract for filtering a single instance of *Entry.
type EntryFilter interface {
	// Name returns a filter name allowing to uniquely identify the filter.
	Name() FilterName

	// Filter contains the filtering logic and returns a bool value on whether
	// the *Entry instance passed the filtration.
	Filter(e *Entry) bool
}

// DirsFilter filters *Entry by its type and allows directories only.
type DirsFilter struct{}

func (df *DirsFilter) Name() FilterName {
	return DirsOnlyFilter
}

func (df *DirsFilter) Filter(e *Entry) bool {
	return e.IsDir
}

// FilesFilter filters *Entry by its type and allows files only.
type FilesFilter struct{}

func (df *FilesFilter) Name() FilterName {
	return FilesOnlyFilter
}

func (df *FilesFilter) Filter(e *Entry) bool {
	return !e.IsDir
}
