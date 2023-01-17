package message

type NetStreamPublish struct {
	CommandObject  interface{}
	PublishingName string
	PublishingType string
}

type NetStreamOnStatus struct {
	CommandObject interface{}
	InfoObject    interface{}
}

type NetStreamOnMetaData struct {
	MetaData map[string]interface{}
}
