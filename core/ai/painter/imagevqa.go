package painter

import (
	"SuperArch/middleware/taskcontrol"
)

type ImageVQA struct{
	taskcontrol.TaskModuleCommon
	question  string
	ImageData map[string]string
}
