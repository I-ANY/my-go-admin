package dscp

type UpdatePriorityResultItem struct {
	ServerID *string `json:"serverId"`
	Hostname *string `json:"hostname"`
	Status   *int64  `json:"status"`
	Message  *string `json:"message"`
	Payload  *string `json:"payload"`
}
type IssueDscpConfig struct {
	Default     *int64 `json:"default,omitempty"`
	SameInner   *int64 `json:"same_inner,omitempty"`
	SameOuter   *int64 `json:"same_outer,omitempty"`
	DiffInner   *int64 `json:"diff_inner,omitempty"`
	DiffOuter   *int64 `json:"diff_outer,omitempty"`
	DefaultV6   *int64 `json:"default_v6,omitempty"`
	SameInnerV6 *int64 `json:"same_inner_v6,omitempty"`
	SameOuterV6 *int64 `json:"same_outer_v6,omitempty"`
	DiffInnerV6 *int64 `json:"diff_inner_v6,omitempty"`
	DiffOuterV6 *int64 `json:"diff_outer_v6,omitempty"`
}
