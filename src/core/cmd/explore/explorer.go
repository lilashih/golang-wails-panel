package explore

func OpenDir(path string) error {
	cmd := openDirCommand(path)

	if err := cmd.Start(); err != nil {
		return err
	}

	return nil
}
