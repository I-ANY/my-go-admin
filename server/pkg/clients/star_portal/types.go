package starPortal

type User struct {
	Name         *string  `json:"name"`
	Title        *string  `json:"title"`
	DeptIDList   []int    `json:"dept_id_list"`
	DeptNameList []string `json:"dept_name_list"`
	Email        *string  `json:"email"`
	Mobile       *string  `json:"mobile"`
	JobNumber    *string  `json:"job_number"`
	WorkPlace    *string  `json:"work_place"`
	LeaderName   *string  `json:"leader_name"`
	LeaderEmail  *string  `json:"leader_email"`
}

//	type User struct {
//		Active           bool     `json:"active"`
//		Admin            bool     `json:"admin"`
//		Avatar           string   `json:"avatar"`
//		Boss             bool     `json:"boss"`
//		DeptIDList       []int    `json:"dept_id_list"`
//		DeptOrder        int64    `json:"dept_order"`
//		Email            string   `json:"email"`
//		ExclusiveAccount bool     `json:"exclusive_account"`
//		Extension        string   `json:"extension"`
//		HideMobile       bool     `json:"hide_mobile"`
//		HiredDate        int64    `json:"hired_date"`
//		JobNumber        string   `json:"job_number"`
//		Leader           bool     `json:"leader"`
//		Mobile           string   `json:"mobile"`
//		Name             string   `json:"name"`
//		Remark           string   `json:"remark"`
//		StateCode        string   `json:"state_code"`
//		Telephone        string   `json:"telephone"`
//		Title            string   `json:"title"`
//		Unionid          string   `json:"unionid"`
//		Userid           string   `json:"userid"`
//		WorkPlace        string   `json:"work_place"`
//		DeptNameList     []string `json:"dept_name_list"`
//	}
type UserTokenInfo struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
	IdToken     string `json:"id_token"`
}

type UserInfo struct {
	Sub               string   `json:"sub"`
	Email             string   `json:"email"`
	EmailVerified     bool     `json:"email_verified"`
	Name              string   `json:"name"`
	GivenName         string   `json:"given_name"`
	PreferredUsername string   `json:"preferred_username"`
	Nickname          string   `json:"nickname"`
	Groups            []string `json:"groups"`
	UserId            string   `yaml:"userid"`
}
