package app

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEncodePriorities(t *testing.T) {
	priorities := []priorityData{
		{
			Id:        "001",
			Color:     1,
			Text:      "Some text",
			IsDeleted: false,
		},
	}
	encodedPriorities := encodePriorities(priorities, 100)
	assert.ObjectsAreEqual(encodedPriorities[0].Text, "U29tZSB0ZXh0")
}

func TestEncodeEmptyPrioritiesRoundtrip(t *testing.T) {
	priorities := []priorityData{}
	decodedPriorities, _ := decodePriorities(encodePriorities(priorities, 100))
	assert.True(t, reflect.DeepEqual(decodedPriorities, priorities))
}

func TestEncodeNonEmptyPrioritiesRoundtrip(t *testing.T) {
	priorities := []priorityData{
		{
			Id:        "001",
			Color:     1,
			Text:      "One",
			IsDeleted: false,
		},
	}
	decodedPriorities, _ := decodePriorities(encodePriorities(priorities, 100))
	assert.True(t, reflect.DeepEqual(decodedPriorities, priorities))
}

func TestDropExtraDeletedPriorities(t *testing.T) {
	priorities := []priorityData{
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
			IsDeleted: true,
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
			IsDeleted: true,
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
			IsDeleted: true,
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
			IsDeleted: true,
		},
		{
			Id:        "009",
			Color:     1,
			Text:      "Nine",
			IsDeleted: false,
		},
	}

	expectedEncoded := []priorityData{
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
			IsDeleted: true,
		},
		{
			Id:        "003",
			Color:     1,
			Text:      "Three",
			IsDeleted: false,
		},
		{
			Id:        "005",
			Color:     1,
			Text:      "Five",
			IsDeleted: false,
		},
		{
			Id:        "007",
			Color:     1,
			Text:      "Seven",
			IsDeleted: false,
		},
		{
			Id:        "009",
			Color:     1,
			Text:      "Nine",
			IsDeleted: false,
		},
	}

	decodedPriorities, _ := decodePriorities(encodePriorities(priorities, 6))
	assert.True(t, reflect.DeepEqual(decodedPriorities, expectedEncoded))
}

func TestDropExtraPriorities(t *testing.T) {
	priorities := []priorityData{
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
			IsDeleted: true,
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
			IsDeleted: true,
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
			IsDeleted: true,
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
			IsDeleted: true,
		},
		{
			Id:        "009",
			Color:     1,
			Text:      "Nine",
			IsDeleted: false,
		},
	}

	expectedEncoded := []priorityData{
		{
			Id:        "001",
			Color:     1,
			Text:      "One",
			IsDeleted: false,
		},
		{
			Id:        "003",
			Color:     1,
			Text:      "Three",
			IsDeleted: false,
		},
		{
			Id:        "005",
			Color:     1,
			Text:      "Five",
			IsDeleted: false,
		},
	}

	decodedPriorities, _ := decodePriorities(encodePriorities(priorities, 3))
	assert.True(t, reflect.DeepEqual(decodedPriorities, expectedEncoded))
}
