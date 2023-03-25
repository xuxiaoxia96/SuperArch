package painter

import "SuperArch/middleware/taskcontrol"

type Image2Txt struct{
	taskcontrol.TaskModuleCommon
	ImageData 		map[string]string
}
