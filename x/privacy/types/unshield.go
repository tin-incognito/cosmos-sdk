package types

func (m *MsgUnShield) ValidateByItself() (bool, error) {
	return true, nil
}

func (m *MsgUnShield) ValidateByDb() (bool, error) {
	return true, nil
}

func (m *MsgUnShield) ValidateSanity() (bool, error) {
	return true, nil
}
