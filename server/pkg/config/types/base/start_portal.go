package base

type StarPortal struct {
	Url          string `mapstructure:"url"`
	ApiToken     string `mapstructure:"apiToken"`
	ClientId     string `mapstructure:"clientId"`
	ClientSecret string `mapstructure:"clientSecret"`
}
