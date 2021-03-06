package washeet

import (
	"math"
)

func newLayoutData(originX, originY, maxX, maxY float64) *layoutData {

	if maxX <= originX || maxY <= originY {
		return nil
	}

	return &layoutData{
		startColumn:        int64(0),
		startRow:           int64(0),
		endColumn:          int64(0),
		endRow:             int64(0),
		colStartXCoords:    make([]float64, 0, 1+int(math.Ceil((maxX-originX+1)/constDefaultCellWidth))),
		rowStartYCoords:    make([]float64, 0, 1+int(math.Ceil((maxY-originY+1)/constDefaultCellHeight))),
		layoutFromStartCol: true,
		layoutFromStartRow: true,
	}
}

func (sheet *Sheet) computeLayout(layout *layoutData, col, row int64, changeSheetStartCol, changeSheetStartRow bool) {

	if sheet == nil {
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
	minX, minY := sheet.origX+constDefaultCellWidth, sheet.origY+constDefaultCellHeight

	if layout.layoutFromStartCol {
		layout.endColumn, layout.colStartXCoords = computeCellsCoordsRefStart(
			minX,
			sheet.maxX,
			layout.startColumn,
			sheet.dataSource.GetColumnWidth,
			constDefaultCellWidth,
			layout.colStartXCoords,
		)
	} else {
		layout.startColumn, layout.endColumn, layout.colStartXCoords = computeCellsCoordsRefEnd(
			minX,
			sheet.maxX,
			layout.endColumn,
			sheet.dataSource.GetColumnWidth,
			constDefaultCellWidth,
			layout.colStartXCoords,
		)
	}

	if layout.layoutFromStartRow {
		layout.endRow, layout.rowStartYCoords = computeCellsCoordsRefStart(
			minY,
			sheet.maxY,
			layout.startRow,
			sheet.dataSource.GetRowHeight,
			constDefaultCellHeight,
			layout.rowStartYCoords,
		)
	} else {
		layout.startRow, layout.endRow, layout.rowStartYCoords = computeCellsCoordsRefEnd(
			minY,
			sheet.maxY,
			layout.endRow,
			sheet.dataSource.GetRowHeight,
			constDefaultCellHeight,
			layout.rowStartYCoords,
		)
	}
}
