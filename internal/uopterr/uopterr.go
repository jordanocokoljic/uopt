package uopterr

type DuplicateName string

func (e DuplicateName) Error() string {
	return "name was already registered: " + string(e)
}

type DuplicateFlag string

func (e DuplicateFlag) Error() string {
	return "flag was already registered: " + string(e)
}

type InvalidShortFlag string

func (e InvalidShortFlag) Error() string {
	return "short flag must be a hyphen followed by 1 alphabetic character: " + string(e)
}

type InvalidLongFlag string

func (e InvalidLongFlag) Error() string {
	return "long flag must be two hyphens followed by an alphabetic character, then any number of alphanumeric characters: " + string(e)
}

type NoFlag string

func (e NoFlag) Error() string {
	return "option must have a short or long flag: " + string(e)
}

type UnrecognizedOption string

func (e UnrecognizedOption) Error() string {
	return "unrecognized option: " + string(e)
}

type UnrecognizedArgument string

func (e UnrecognizedArgument) Error() string {
	return "unrecognized argument: " + string(e)
}

type NoCaptureValue string

func (e NoCaptureValue) Error() string {
	return "option value was missing: " + string(e)
}
