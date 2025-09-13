package config

type Environment string

const (
	EnvLocal Environment = "local"
	EnvDev   Environment = "dev"
	EnvStage Environment = "stage"
	EnvProd  Environment = "prod"
)

func (e Environment) String() string {
	return string(e)
}

func (e Environment) IsProduction() bool {
	return e == EnvProd
}

func (e Environment) IsDevelopment() bool {
	return e == EnvLocal || e == EnvDev
}

func (e Environment) IsValid() bool {
	return e == EnvLocal || e == EnvDev || e == EnvStage || e == EnvProd
}
