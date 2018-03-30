

### //Callback 是回调应该使用的函数类型。 Event是当回调发生时的当前事件信息。
## type Callback func(*Event)

### // Events 是在NewFSM中定义转换映射的简写。
## type Events []EventDesc

### // Callbacks在NewFSM.a中定义回调的简写
## type Callbacks map[string]Callback

// EventDesc表示初始化FSM时的事件。
//事件可以有一个或多个对执行转换有效的源状态。 如果FSM处于其中一个源状态，它将以指定的目标状态结束，并调用所有已定义的回调。
type EventDesc struct {
	// Name是调用转换时使用的事件名称。
	Name string

	// Src是FSM必须执行状态转换的一部分源状态。
	Src []string

	// 如果转换成功，则Dst是FSM将处于的目标状态。
	Dst string
}


// FSM是保存当前状态的状态机。
//
//必须使用NewFSM创建才能正常运行。
type FSM struct {
	//当前是FSM当前所在的状态。
	current string

	//将事件和源状态映射到目标状态。
	transitions map[eKey]string

	//回调映射事件和targer到回调函数。
	callbacks map[cKey]Callback

	// transition 是直接使用或在异步状态转换中调用转换时使用的内部转换函数。
	transition func()
	//transitionerObj调用FSM的transition（）函数。
	transitionerObj transitioner

	// stateMu guards access 进入目前的状态。
	stateMu sync.RWMutex
	// eventMu guards access to Event() and Transition().
	eventMu sync.Mutex
}
