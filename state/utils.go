package state

type DeprecatedOption struct {
	Name        string
	Description string
}

// GetDeprecatedConfigOptions This method is here with an input for reference
// currently not used but helpful for the future
//
//goland:noinspection ALL
func GetDeprecatedConfigOptions(cfg *Config) []DeprecatedOption {

	// below is an example of a deprecated setting. This code can be used in the future if settings are deprecated
	//returnValue := []DeprecatedOption{}
	//if cfg.Telegram.EmojiConfirmation != nil {
	//	returnValue = append(returnValue, DeprecatedOption{
	//		Name:        "[telegram.emoji_confirmation]",
	//		Description: "It has been replaced with [telegram.confirmation_type]",
	//	})
	//
	//	if *cfg.Telegram.EmojiConfirmation {
	//		cfg.Telegram.ConfirmationType = "emoji"
	//	} else {
	//		cfg.Telegram.ConfirmationType = "text"
	//	}
	//	cfg.Telegram.EmojiConfirmation = nil
	//}
	//
	//if len(returnValue) > 0 {
	//	return returnValue
	//} else {
	//	return nil
	//}

	return nil
}
