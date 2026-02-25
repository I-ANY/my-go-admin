package base

type ZapBusiness struct {
	SupplierId          *int    `mapstructure:"supplier_id"`
	SupplierIdDay95     *int    `mapstructure:"supplier_id_day95"`
	SupplierName        string  `mapstructure:"supplier_name"`
	SupplierEname       string  `mapstructure:"supplier_ename"`
	DynamicResourceType string  `mapstructure:"dynamic_resource_type"`
	LiveResourceType    string  `mapstructure:"live_resource_type"`
	InnerIp             string  `mapstructure:"innerIp"`
	SubmitType          int     `mapstructure:"submit_type"`
	Domain              string  `mapstructure:"domain"`
	DomainAccountName   string  `mapstructure:"domain_account_name"`
	DomainPassword      string  `mapstructure:"domain_password"`
	DomainLoginPort     int64   `mapstructure:"domain_login_port"`
	Isps                []isp   `mapstructure:"isp"`
	AccountDay95        account `mapstructure:"account_day95"`
	AccountMonth95      account `mapstructure:"account_month95"`
}

type isp struct {
	Name  string `mapstructure:"name"`
	EName string `mapstructure:"ename"`
}

type account struct {
	Name   string `mapstructure:"name"`
	Secret string `mapstructure:"secret"`
}
