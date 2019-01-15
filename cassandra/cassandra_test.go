package cassandra

import (
	"errors"
	"fmt"
	"testing"

	"github.com/gocql/gocql"
	"github.com/stretchr/testify/require"
)

type testSession struct {
	name    string
	session *gocql.Session
	err     error
}

func (scenario *testSession) setup(t *testing.T) {
	if scenario.err == nil {
		scenario.session = new(gocql.Session)
	}
}

func (scenario *testSession) tearDown(t *testing.T) {
}

func TestSession(test *testing.T) {
	scenarios := []testSession{
		{
			name: "Creates a new session",
		},
		{
			name: "Returns error because session is blank",
			err:  ErrBlankSession,
		},
	}

	for index, scenario := range scenarios {
		test.Run(
			fmt.Sprintf("[%d]-%s", index, scenario.name),
			func(t *testing.T) {
				scenario.setup(t)
				defer scenario.tearDown(t)

				session, err := newSession(scenario.session)
				require.Equal(t, scenario.err, err, "newsession error")
				if scenario.err == nil {
					require.NotNil(t, session, "session instance")
					require.False(t, session.Closed(), "session is closed")
					session.Close()
				} else {
					require.Nil(t, session, "session invalid instance")
				}
			},
		)
	}
}

type testQuery struct {
	name      string
	session   *gocql.Session
	cql       string
	arguments []interface{}
	err       error
}

func (scenario *testQuery) setup(t *testing.T) {
	scenario.session = new(gocql.Session)
}

func (scenario *testQuery) tearDown(t *testing.T) {
}

func TestQuery(test *testing.T) {
	scenarios := []testQuery{
		{
			name: "Executes query successfully",
			cql:  "select id, text from mock",
		},
		{
			name: "When query returns error",
			cql:  "select id, text from mock",
			err:  errors.New("err_mockquery"),
		},
	}

	for index, scenario := range scenarios {
		test.Run(
			fmt.Sprintf("[%d]-%s", index, scenario.name),
			func(t *testing.T) {
				scenario.setup(t)
				defer scenario.tearDown(t)

				session, err := newSession(scenario.session)
				require.Nil(t, err, "newsession error")
				require.NotNil(t, session, "session instance")
				query := session.Query(scenario.cql, scenario.arguments...)
				require.NotNil(t, query, "query invalid instance")

				query = query.Consistency(gocql.Any)
				query = query.PageSize(100)

				require.Panics(t,
					func() {
						var id, text string
						query.Scan(&id, &text)
					},
				)

				require.Panics(t,
					func() {
						_ = query.Iter()
					},
				)

				require.Panics(t,
					func() {
						_ = query.Exec()

					},
				)

				require.NotZero(t, query.String(), "querystring invalid instance")
				query.Release()
				session.Close()
			},
		)
	}
}
