package washeet

import (
	//	"fmt"
	"math"
	"syscall/js"
)

func NewSheet(canvasElement, container *js.Value, startX float64, startY float64, maxX float64, maxY float64,
	dSrc SheetDataProvider, dSink SheetModelUpdater) *Sheet {

	// HACK : Adjust for line width of 1.0
	maxX -= 1.0
	maxY -= 1.0

	if canvasElement == nil || startX+DEFAULT_CELL_WIDTH*10 >= maxX ||
		startY+DEFAULT_CELL_HEIGHT*10 >= maxY {
		return nil
	}

	ret := &Sheet{
		document:           js.Global().Get("document"),
		window:             js.Global().Get("window"),
		container:          container,
		canvasElement:      canvasElement,
		canvasContext:      canvasElement.Call("getContext", "2d"),
		origX:              startX,
		origY:              startY,
		maxX:               maxX,
		maxY:               maxY,
		dataSource:         dSrc,
		dataSink:           dSink,
		rafPendingQueue:    make(chan js.Value, SHEET_PAINT_QUEUE_LENGTH),
		startColumn:        int64(0),
		startRow:           int64(0),
		endColumn:          int64(0),
		endRow:             int64(0),
		paintQueue:         make(chan *sheetPaintRequest, SHEET_PAINT_QUEUE_LENGTH),
		colStartXCoords:    make([]float64, 0, 1+int(math.Ceil((maxX-startX+1)/DEFAULT_CELL_WIDTH))),
		rowStartYCoords:    make([]float64, 0, 1+int(math.Ceil((maxY-startY+1)/DEFAULT_CELL_HEIGHT))),
		mark:               MarkData{0, 0, 0, 0},
		stopSignal:         false,
		stopWaitChan:       make(chan bool),
		mouseState:         defaultMouseState(),
		selectionState:     defaultSelectionState(),
		layoutFromStartCol: true,
		layoutFromStartRow: true,
	}

	// TODO : Move these somewhere else
	setFont(&ret.canvasContext, "14px serif")
	setLineWidth(&ret.canvasContext, 1.0)

	ret.setupClipboardTextArea()
	ret.PaintWholeSheet(ret.startColumn, ret.startRow, ret.layoutFromStartCol, ret.layoutFromStartRow)
	ret.setupMouseHandlers()
	ret.setupKeyboardHandlers()

	return ret
}

func (self *Sheet) Start() {

	if self == nil {
		return
	}

	self.stopSignal = false
	go self.processQueue()
}

func (self *Sheet) Stop() {

	if self == nil || self.stopSignal {
		return
	}

	self.teardownKeyboardHandlers()
	self.teardownMouseHandlers()

	self.stopSignal = true
	// clear the widget area.
	// HACK : maxX + 1.0, maxY + 1.0 is the actual limit
	noStrokeFillRectNoAdjust(&self.canvasContext, self.origX, self.origY, self.maxX+1.0, self.maxY+1.0, CELL_DEFAULT_FILL_COLOR)
	// Wait till we get signal from paint-queue when it it has finished
	<-self.stopWaitChan
}

// if col/row = -1 no changes are made before whole-redraw
// changeSheetStartCol/changeSheetStartRow is also used to set self.layoutFromStartCol/self.layoutFromStartRow
func (self *Sheet) PaintWholeSheet(col, row int64, changeSheetStartCol, changeSheetStartRow bool) {
	req := &sheetPaintRequest{
		kind:                sheetPaintWholeSheet,
		col:                 col,
		row:                 row,
		changeSheetStartCol: changeSheetStartCol,
		changeSheetStartRow: changeSheetStartRow,
	}
	self.addPaintRequest(req)
}

func (self *Sheet) PaintCell(col int64, row int64) {

	if self == nil {
		return
	}

	// optimization : don't fill the queue with these
	// if we know they are not going to be painted.
	if col < self.startColumn || col > self.endColumn ||
		row < self.startRow || row > self.endRow {
		return
	}

	self.addPaintRequest(&sheetPaintRequest{
		kind:   sheetPaintCell,
		col:    col,
		row:    row,
		endCol: col,
		endRow: row,
	})
}

func (self *Sheet) PaintCellRange(colStart int64, rowStart int64, colEnd int64, rowEnd int64) {

	if self == nil {
		return
	}

	self.addPaintRequest(&sheetPaintRequest{
		kind:   sheetPaintCellRange,
		col:    colStart,
		row:    rowStart,
		endCol: colEnd,
		endRow: rowEnd,
	})
}

func (self *Sheet) PaintCellSelection(col, row int64) {
	if self == nil {
		return
	}

	self.addPaintRequest(&sheetPaintRequest{
		kind:   sheetPaintSelection,
		col:    col,
		row:    row,
		endCol: col,
		endRow: row,
	})
}

func (self *Sheet) PaintCellRangeSelection(colStart, rowStart, colEnd, rowEnd int64) {
	if self == nil {
		return
	}

	self.addPaintRequest(&sheetPaintRequest{
		kind:   sheetPaintSelection,
		col:    colStart,
		row:    rowStart,
		endCol: colEnd,
		endRow: rowEnd,
	})
}
