package painter

import "SuperArch/middleware/taskcontrol"

type ImageScore struct{
	taskcontrol.TaskModuleCommon
	ImageData map[string]string
}
