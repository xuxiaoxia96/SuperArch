package painter

import "SuperArch/middleware/taskcontrol"

type Txt2Image struct{
	taskcontrol.TaskModuleCommon
	Text string
}
