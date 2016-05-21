package services

type RulesEngineChecker struct {
	Errors []string
	Status bool
}

func (engine *RulesEngineChecker) Success() (bool, []string) {
	if len(engine.Errors) != 0 {
		return false, engine.Errors
	}

	return true, []string{}
}

func NewRulesEngineChecker()(*RulesEngineChecker){
	return &RulesEngineChecker{Status: true, Errors: []string{} }
}

func (engine *RulesEngineChecker) AddErrors(error string) (bool, []string)  {
	existsInArray := false

	for i:= 0; i < len(engine.Errors); i++ {
		if engine.Errors[i] == error {
			existsInArray = true
			return false, engine.Errors
		}
	}

	if existsInArray == false {
		engine.Errors = append(engine.Errors, error)
		return true, engine.Errors
	}else{
		return false, engine.Errors
	}


}
