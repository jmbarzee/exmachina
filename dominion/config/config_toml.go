package config

// SetupFromTOML produces a default configuration which can be passed to New()
func SetupFromTOML(configFilePath string) error {
	if err := setupDominionConfigFromTOML(configFilePath); err != nil {
		return err
	}
	if err := setupServicesConfigFromTOML(configFilePath); err != nil {
		return err
	}
	return nil
}
