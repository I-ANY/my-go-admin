package base

type FlaskApi struct {
	Url      string `mapstructure:"url"`
	ApiToken string `mapstructure:"apiToken"`
}
