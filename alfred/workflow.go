package alfred

func RunScriptAction(action func(*ScriptActionResponse) error, opts ...ScriptActionResponseOption) {
	sar := NewScriptActionResponse(opts...)
	defer RecoverIfErr(sar)()

	err := action(sar)

	if err != nil {
		sar.SetError(err)
	}

	sar.SendFeedback()
}

func RunScriptFilter(action func(*ScriptFilterResponse) error, opts ...ScriptFilterResponseOption) {
	sfr := NewScriptFilterResponse(opts...)
	defer RecoverIfErr(sfr)()

	err := action(sfr)

	if err != nil {
		sfr.SetError(err)
	}

	sfr.SendFeedback()
}
