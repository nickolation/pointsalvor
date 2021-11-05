package pointsalvor

import (
	"errors"
	"fmt"
)

//custom erorrs for validate and performing methods
//Persistant error is simple string with text description of problem
//errTO type error is function with custom description and native golang-error
var (
	errUrl              = errors.New("url request is empty")
	errMethod           = errors.New("method is invalid")
	errJsonMarshal      = errors.New("eror with marshal request body")
	errRequestForm      = errors.New("request is invalid")
	errRequestPerf      = errors.New("error with perform request")
	errModelValid       = errors.New("invalid model")
	errInvalidNameModel = fmt.Errorf("invalid name of model - empty")
	errSwitchType       = fmt.Errorf("error with switch type")

	errDecodeIf = func(e error) error {
		return fmt.Errorf("error with decoding interface - [%v]", e)
	}

	errDecodeMap = func(e error) error {
		return fmt.Errorf("error with decoding map - [%v]", e)
	}

	errConvertTo = func(s string) error {
		return fmt.Errorf("invalid convertation to interface - [%s]", s)
	}

	errKnockTo = func(e error) error {
		return fmt.Errorf("invalid knock to api - [%v]", e)
	}
)
