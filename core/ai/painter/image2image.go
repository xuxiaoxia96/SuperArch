package painter

import "SuperArch/middleware/taskcontrol"

type Image2Image struct{
	taskcontrol.TaskModuleCommon
	Text string
	ImageData string
}
