package app

import "testing"

func TestPriorityListLengthValidationEmptyList(t *testing.T) {
	priorities := priorityListData{
		Items: []priorityData{},
	}

	isValid := isPriorityListLengthValid(priorities)
	if isValid != true {
		t.Errorf("Did not get expected result. Expected: %v, actual: %v", true, isValid)
	}
}

func TestPriorityListLengthValidationExactly9ItemsAllActive(t *testing.T) {
	priorities := priorityListData{
		Items: []priorityData{
			{
				Id:        "001",
				Color:     1,
				Text:      "One",
				IsDeleted: false,
			},
			{
				Id:        "002",
				Color:     1,
				Text:      "Two",
				IsDeleted: false,
			},
			{
				Id:        "003",
				Color:     1,
				Text:      "Three",
				IsDeleted: false,
			},
			{
				Id:        "004",
				Color:     1,
				Text:      "Four",
				IsDeleted: false,
			},
			{
				Id:        "005",
				Color:     1,
				Text:      "Five",
				IsDeleted: false,
			},
			{
				Id:        "006",
				Color:     1,
				Text:      "Six",
				IsDeleted: false,
			},
			{
				Id:        "007",
				Color:     1,
				Text:      "Seven",
				IsDeleted: false,
			},
			{
				Id:        "008",
				Color:     1,
				Text:      "Eight",
				IsDeleted: false,
			},
			{
				Id:        "009",
				Color:     1,
				Text:      "Nine",
				IsDeleted: false,
			},
		},
	}

	isValid := isPriorityListLengthValid(priorities)
	if isValid != true {
		t.Errorf("Did not get expected result. Expected: %v, actual: %v", true, isValid)
	}
}

func TestPriorityListLengthValidation9ActiveItems1InactiveItem(t *testing.T) {
	priorities := priorityListData{
		Items: []priorityData{
			{
				Id:        "001",
				Color:     1,
				Text:      "One",
				IsDeleted: false,
			},
			{
				Id:        "002",
				Color:     1,
				Text:      "Two",
				IsDeleted: false,
			},
			{
				Id:        "003",
				Color:     1,
				Text:      "Three",
				IsDeleted: false,
			},
			{
				Id:        "004",
				Color:     1,
				Text:      "Four",
				IsDeleted: false,
			},
			{
				Id:        "005",
				Color:     1,
				Text:      "Five",
				IsDeleted: false,
			},
			{
				Id:        "006",
				Color:     1,
				Text:      "Six",
				IsDeleted: false,
			},
			{
				Id:        "007",
				Color:     1,
				Text:      "Seven",
				IsDeleted: false,
			},
			{
				Id:        "008",
				Color:     1,
				Text:      "Eight",
				IsDeleted: false,
			},
			{
				Id:        "009",
				Color:     1,
				Text:      "Nine",
				IsDeleted: false,
			},
			{
				Id:        "010",
				Color:     1,
				Text:      "Ten",
				IsDeleted: true,
			},
		},
	}

	isValid := isPriorityListLengthValid(priorities)
	if isValid != true {
		t.Errorf("Did not get expected result. Expected: %v, actual: %v", true, isValid)
	}
}

func TestPriorityListLengthValidationExactlyOver9ActiveItems(t *testing.T) {
	priorities := priorityListData{
		Items: []priorityData{
			{
				Id:        "001",
				Color:     1,
				Text:      "One",
				IsDeleted: false,
			},
			{
				Id:        "002",
				Color:     1,
				Text:      "Two",
				IsDeleted: false,
			},
			{
				Id:        "003",
				Color:     1,
				Text:      "Three",
				IsDeleted: false,
			},
			{
				Id:        "004",
				Color:     1,
				Text:      "Four",
				IsDeleted: false,
			},
			{
				Id:        "005",
				Color:     1,
				Text:      "Five",
				IsDeleted: false,
			},
			{
				Id:        "006",
				Color:     1,
				Text:      "Six",
				IsDeleted: false,
			},
			{
				Id:        "007",
				Color:     1,
				Text:      "Seven",
				IsDeleted: false,
			},
			{
				Id:        "008",
				Color:     1,
				Text:      "Eight",
				IsDeleted: false,
			},
			{
				Id:        "009",
				Color:     1,
				Text:      "Nine",
				IsDeleted: false,
			},
			{
				Id:        "010",
				Color:     1,
				Text:      "Ten",
				IsDeleted: false,
			},
		},
	}

	isValid := isPriorityListLengthValid(priorities)
	if isValid != false {
		t.Errorf("Did not get expected result. Expected: %v, actual: %v", true, isValid)
	}
}
