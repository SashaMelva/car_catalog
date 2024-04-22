package memory

func (s *Storage) DeleteCarByRegNum(regNum string) error {
	query := `delete cars where reg_num = $1`
	_, err := s.ConnectionDB.Exec(query, regNum)

	if err != nil {
		return err
	}

	return nil
}
