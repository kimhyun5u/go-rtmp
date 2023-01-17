package message

type NetConnectionConnection struct {
	CommandObject map[string]interface{}
}

type NetConnectionCreateStream struct {
}

type NetConnectionResult struct {
	Objects []interface{}
}
