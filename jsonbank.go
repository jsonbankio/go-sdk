package jsonbank

// Init - initializes the jsonbank instance
func Init(config Config) Instance {
	// Validate config
	// Assign default Host if not provided
	if len(config.Host) <= 0 {
		config.Host = "https://api.jsonbank.io"
	}

	// make instance
	jsb := Instance{}
	// set config
	jsb.config = config
	// set urls
	jsb.SetHost(config.Host)
	// set memory
	jsb.memory = make(map[string]any)

	return jsb
}

// InitWithoutKeys - initializes the jsonbank instance without Keys
func InitWithoutKeys() Instance {
	return Init(Config{})
}
