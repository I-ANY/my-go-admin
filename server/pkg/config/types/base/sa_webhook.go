package base

type SaAlertWebhook struct {
	Url         string   `mapstructure:"url"`
	AlertTarget []string `mapstructure:"alertTarget" default:"[\"123\",\"456\"]"`
	DetailUrl   string   `mapstructure:"detailUrl"`
}
