package game

type CellStatus uint8

const (
	CellDied = CellStatus(iota)
	CellLive = CellStatus(iota)
)

type Cell struct {
	currentStatus CellStatus
	previewStatus CellStatus
}

func NewCell(status CellStatus) *Cell {
	return &Cell{
		currentStatus: status,
		previewStatus: status,
	}
}

func (c Cell) CurrentStatus() CellStatus {
	return c.currentStatus
}

func (c Cell) String() string {
	if c.currentStatus == CellDied {
		return "died"
	}

	return "live"
}
