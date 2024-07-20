package main

// Skip the app name and return all flags and arguments.
func parseFlags(args []string) (map[string]string, []string) {
	flags := map[string]string{}
	a := []string{}
	for i := 1; i < len(args); i++ {
		switch args[i] {
		case "-h", "--help":
			printHelp()
		case "--version":
			printVersion()
		case "-v", "--verbose":
			flags["debug"] = "true"
		case "-H", "--host":
			flags["api_host"] = args[i+1]
			i++
		case "-p", "--port":
			flags["api_port"] = args[i+1]
			i++
		case "-d", "--db-root":
			flags["db_root"] = args[i+1]
			i++
		case "-c", "--config-file":
			flags["config_file"] = args[i+1]
			i++
		default:
			a = append(a, args[i])
		}
	}

	return flags, a
}
