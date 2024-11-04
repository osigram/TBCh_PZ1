package domain

type Blockchain []Block

func (bc *Blockchain) Add(b Block) error {
	if err := b.Validate(*bc); err != nil {
		return err
	}

	*bc = append(*bc, b)

	return nil
}
