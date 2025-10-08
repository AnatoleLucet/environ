package environ

// String will ensure the variable is a string
func String(name string) VariableBuilder[string] {
	return VariableBuilder[string]{
		Variable: Variable[string]{Name: name, Type: TypeString},
	}
}

// Int will ensure the variable is a number (int)
func Int(name string) VariableBuilder[int] {
	return VariableBuilder[int]{
		Variable: Variable[int]{Name: name, Type: TypeInt},
	}
}

// Float will ensure the variable is a number (int)
func Float(name string) VariableBuilder[float64] {
	return VariableBuilder[float64]{
		Variable: Variable[float64]{Name: name, Type: TypeFloat},
	}
}

// Boolean will ensure the variable is a boolean ("true", "false", "0", "1", "on", "off")
func Boolean(name string) VariableBuilder[bool] {
	return VariableBuilder[bool]{
		Variable: Variable[bool]{Name: name, Type: TypeBoolean},
	}
}

// Port will ensure the variable is a valid TCP port number (1-65535)
func Port(name string) VariableBuilder[int] {
	return VariableBuilder[int]{
		Variable: Variable[int]{Name: name, Type: TypePort},
	}
}

// Url will ensure the variable is a valid URL with a protocol and a hostname
func Url(name string) VariableBuilder[string] {
	return VariableBuilder[string]{
		Variable: Variable[string]{Name: name, Type: TypeUrl},
	}
}

// Email will ensure the variable is a valid email address
func Email(name string) VariableBuilder[string] {
	return VariableBuilder[string]{
		Variable: Variable[string]{Name: name, Type: TypeEmail},
	}
}
