package conf

import (
	"go.katupy.org/fixture"
)

var PostgresOptions = &fixture.Config{
	References: map[string]string{},
	TableOptions: map[string]*fixture.TableOptions{
		"data_file_directories": {
			DefaultValues: fixture.Record{},
		},
		"data_file_objects": {
			DefaultValues: fixture.Record{
				"directory_id": "=ref data_file_directories #",
			},
		},
	},
}
