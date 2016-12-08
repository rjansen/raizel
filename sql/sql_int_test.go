package cassandra

var (
	key1    = "8b06603b-9b0d-4e8c-8aae-10f988639fe6"
	expires = 60
)

func init() {
}

func setup(testConfig *Configuration) error {
	var err error
	if err = Setup(testConfig); err != nil {
		return err
	}
	return nil
}

// func before() error {
// 	var err error
// 	if !setted {
// 		if err = setup(); err != nil {
// 			return err
// 		}
// 	}
// 	persistenceClient, err = pool.Get()
// 	return err
// }
