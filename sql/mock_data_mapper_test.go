package sql

import "time"

type dataMock struct {
	ID         int             `db:"id"`
	Name       string          `db:"name"`
	Age        int             `db:"age"`
	Score      float32         `db:"score"`
	Deleted    bool            `db:"deleted"`
	CreatedAt  time.Time       `db:"created_at"`
	UpdatedAt  time.Time       `db:"updated_at"`
	data       dynamicData     `db:"data"`
	NestedData nestedDataMock  `db:"nested_data"`
	SliceData  []sliceDataMock `db:"slice_data"`
}

type nestedDataMock struct {
	ID        int       `db:"id"`
	Name      string    `db:"name"`
	Deleted   bool      `db:"deleted"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

type sliceDataMock struct {
	ID        int       `db:"id"`
	Name      string    `db:"name"`
	Deleted   bool      `db:"deleted"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}
