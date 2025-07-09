package abodemine

import (
	"github.com/google/uuid"

	"abodemine/projects/datapipe/entities"
)

var PartnerId = uuid.MustParse("7ee2e306-d03f-4f72-abf8-3ac5df4796ab")

const (
	// Use random 9-digit integers (32bit) to ensure new
	// data file types can be added and grouped together
	// with similar types without worrying about
	// collisions and order.
	//
	// Although smaller values would be easier to remember,
	// larger values are used to reduce the chance of
	// collisions with data file types from other partners.
	//
	// New values can be generated with: ugen --digit -l 9.
	DataFileTypeSearchAddress  entities.DataFileType = 299117701
	DataFileTypeSearchAVM      entities.DataFileType = 483488427
	DataFileTypeSearchProperty entities.DataFileType = 658914146
)
