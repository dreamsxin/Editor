package rwundo

import "github.com/jmigpin/editor/util/iout/iorw"

type UndoRedo struct {
	Index int
	D     []byte // deleted bytes of the original op
	I     []byte // inserted bytes of the original op
}

func NewUndoRedoOverwrite(rw iorw.ReadWriter, i, n int, p []byte) (*UndoRedo, error) {
	// copy delete
	b1, err := rw.ReadNAtCopy(i, n)
	if err != nil {
		return nil, err
	}
	// copy insert
	b2 := make([]byte, len(p))
	copy(b2, p)

	if err := rw.Overwrite(i, n, p); err != nil {
		return nil, err
	}
	ur := &UndoRedo{Index: i, D: b1, I: b2}
	return ur, nil
}

//----------

func (ur *UndoRedo) Apply(redo bool, w iorw.Writer) error {
	if redo {
		return w.Overwrite(ur.Index, len(ur.D), ur.I)
	} else {
		return w.Overwrite(ur.Index, len(ur.I), ur.D)
	}
}

func (ur *UndoRedo) IsInsertOnly() bool {
	return len(ur.D) == 0 && len(ur.I) != 0
}
func (ur *UndoRedo) IsDeleteOnly() bool {
	return len(ur.D) != 0 && len(ur.I) == 0
}