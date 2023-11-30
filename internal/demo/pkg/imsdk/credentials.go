package imsdk

type Credentials struct {
	Ak string
	Sk string
}

func NewStaticCredentials(ak, sk string) *Credentials {
	return &Credentials{
		Ak: ak,
		Sk: sk,
	}
}

func (c *Credentials) GetAK() string {
	return c.Ak
}

func (c *Credentials) GetSK() string {
	return c.Sk
}
