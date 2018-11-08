package washeet

import (
	"math"
)

func NewLayoutData(originX, originY, maxX, maxY float64) *layoutData {

	if maxX <= originX || maxY <= originY {
		return nil
	}

	return &layoutData{
		startColumn:        int64(0),
		startRow:           int64(0),
		endColumn:          int64(0),
		endRow:             int64(0),
		colStartXCoords:    make([]float64, 0, 1+int(math.Ceil((maxX-originX+1)/DEFAULT_CELL_WIDTH))),
		rowStartYCoords:    make([]float64, 0, 1+int(math.Ceil((maxY-originY+1)/DEFAULT_CELL_HEIGHT))),
		layoutFromStartCol: true,
		layoutFromStartRow: true,
	}
}

func (self *Sheet) computeLayout(layout *layoutData, col, row int64, changeSheetStartCol, changeSheetStartRow bool) {

	if self == nil {
		return
	}

	// Recompute startColumn/startRow/endColumn/endRow colStartXCoords/rowStartYCoords
	layout.layoutFromStartCol = changeSheetStartCol
	layout.layoutFromStartRow = changeSheetStartRow

	if col >= 0 {
		if changeSheetStartCol {
			layout.startColumn = col
		} else {
			layout.endColumn = col
		}
	}

	if row >= 0 {
		if changeSheetStartRow {
			layout.startRow = row
		} else {
			layout.endRow = row
		}
	}

	// Start of actual top left cell to be drawn after accounting
	// for the row/col-headers
	minX, minY := self.origX+DEFAULT_CELL_WIDTH, self.origY+DEFAULT_CELL_HEIGHT

	if layout.layoutFromStartCol {
		layout.endColumn, layout.colStartXCoords = computeCellsCoordsRefStart(
			minX,
			self.maxX,
			layout.startColumn,
			self.dataSource.GetColumnWidth,
			DEFAULT_CELL_WIDTH,
			layout.colStartXCoords,
		)
	} else {
		layout.startColumn, layout.endColumn, layout.colStartXCoords = computeCellsCoordsRefEnd(
			minX,
			self.maxX,
			layout.endColumn,
			self.dataSource.GetColumnWidth,
			DEFAULT_CELL_WIDTH,
			layout.colStartXCoords,
		)
	}

	if layout.layoutFromStartRow {
		layout.endRow, layout.rowStartYCoords = computeCellsCoordsRefStart(
			minY,
			self.maxY,
			layout.startRow,
			self.dataSource.GetRowHeight,
			DEFAULT_CELL_HEIGHT,
			layout.rowStartYCoords,
		)
	} else {
		layout.startRow, layout.endRow, layout.rowStartYCoords = computeCellsCoordsRefEnd(
			minY,
			self.maxY,
			layout.endRow,
			self.dataSource.GetRowHeight,
			DEFAULT_CELL_HEIGHT,
			layout.rowStartYCoords,
		)
	}
}
