package securitykey

var Key []byte // read only

func init() {
	var err error

	Key, err = getKey()
	if err != nil {
		panic(err)
	}
}
