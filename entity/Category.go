package entity

type Category string

const (
	FootballCategory = "football"
	HistoryCategory  = "history"
)

func (c Category) IsValid() bool {
	switch c {
	case FootballCategory:
		return true
	case HistoryCategory:
		return true
	}
	return false
}

func GetCategoryList() []Category {
	return []Category{FootballCategory, HistoryCategory}
}
