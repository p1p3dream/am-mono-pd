package lambda

type TaskLauncherMessageBody struct {
	Partner string `json:"partner,omitempty" yaml:"partner,omitempty"`
	Task    string `json:"task,omitempty" yaml:"task,omitempty"`
}
