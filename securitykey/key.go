package securitykey

var Key []byte

func init() {
	var err error

	Key, err = getKey()
	if err != nil {
		panic(err)
	}
}

func getKey() ([]byte, error) {
	return []byte("sadfhjsadhflkuasyfowheo2903470298hslhfljsdafohasdlfhlasdflkjsadhflkjsdhalfkjhsaldj"), nil
}
