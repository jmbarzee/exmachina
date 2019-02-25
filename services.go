package domain

type (
	Service struct {
		Port          int
		ServiceConfig ServiceConfig
	}

	ServiceConfig struct {
		Name     string
		Priority Priority
		Depends  []string
		Traits   []string
	}
)

type Priority string

const (
	Always Priority = "always"
	Needed Priority = "as needed"
)

func (d *Domain) hostServices() {

}
