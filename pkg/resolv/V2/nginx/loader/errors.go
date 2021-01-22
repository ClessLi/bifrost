package loader

import "errors"

var (
	ErrInvalidConfig = errors.New("invalid parser.Config")
	ConfigParseError = errors.New("config parse error")
	//ParserTypeError                 = fmt.Errorf("invalid parserType")
	//ParserControlNoParamError       = fmt.Errorf("no valid param has been inputed")
	//ParserControlParamsError        = fmt.Errorf("unkown param has been inputed")
	//ParserControlIndexNotFoundError = fmt.Errorf("index not found")
	//NoBackupRequired                = fmt.Errorf("no backup required")
	//NoReloadRequired                = fmt.Errorf("no reload required")
	//IsInCaches                      = fmt.Errorf("cache already exists")
	//IsNotInCaches                   = fmt.Errorf("cache is not exists")
)
