package service

type PisigService interface {
	SetSettings(settings map[string]interface{}) (error, bool)
	UpdateSettings(settings map[string]interface{}) (error, bool)

	GroupName() string
	ServiceName() string

	Process(data interface{}) (error, bool)
}
