package washeet

func (self *MouseState) isLeftDown() bool {
	return 0x01 == (self.buttonsState & 0x01)
}

func (self *MouseState) isRightDown() bool {
	return 0x02 == (self.buttonsState & 0x02)
}

func (self *MouseState) setLeftDown() {
	self.buttonsState |= 0x01
}

func (self *MouseState) setLeftUp() {
	self.buttonsState &= 0xFE
}

func (self *MouseState) setRightDown() {
	self.buttonsState |= 0x02
}

func (self *MouseState) setRightUp() {
	self.buttonsState &= 0xFD
}

func (self *MouseState) setLastMouseDownCell(col, row int64) {
	self.lastMouseDownCell.Col, self.lastMouseDownCell.Row = col, row
}

func defaultMouseState() MouseState {
	return MouseState{buttonsState: 0x00, lastMouseDownCell: CellCoords{0, 0}}
}