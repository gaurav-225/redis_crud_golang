package serverstatus


type ServerConf struct {
	Port             string `json:"port"`
	Role             string `json:"role"`
	MasterReplid     string `json:"master_replid"`
	MasterReplOffset int    `json:"master_repl_offset"`
	MasterHost       string `json:"master_host"`
	MasterPort       string `json:"master_port"`
	OwnHost			 string `json:"own_host"`
}

func NewServerConf() *ServerConf {
	return &ServerConf{}
}