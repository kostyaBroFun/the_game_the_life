package game

// TODO calculate status
// TODO calculate era

type Pair struct {
	x int
	y int
}

func NewPair(x int, y int) Pair {
	return Pair{
		x: x,
		y: y,
	}
}

func (p Pair) X() int {
	return p.x
}

func (p Pair) Y() int {
	return p.y
}

type ViewFunc func(map[Pair]*Cell)

type Loop struct {
	cells        map[Pair]*Cell
	stop         chan struct{}
	viewFuncList []ViewFunc
}

type Option func(*Loop)

func WithView(viewFunc ViewFunc) Option {
	return func(loop *Loop) {
		loop.viewFuncList = append(loop.viewFuncList, viewFunc)
	}
}

func WithLiveCell(cord Pair) Option {
	return func(loop *Loop) {
		loop.cells[cord] = &Cell{currentStatus: CellLive}
	}
}

func NewLoop(options ...Option) *Loop {
	loop := &Loop{
		cells: make(map[Pair]*Cell),
		stop:  make(chan struct{}),
	}

	for _, option := range options {
		option(loop)
	}

	return loop
}

func (l *Loop) StartTheLife() {
	for {
		select {
		case <-l.stop:
			return
		default:
			l.view()
			l.currentToPast()
			l.recalculateCurrentByPast()
		}
	}
}

func (l *Loop) Stop() {
	l.stop <- struct{}{}
	close(l.stop)
}

func (l *Loop) currentToPast() {
	for _, cell := range l.cells {
		cell.previewStatus = cell.currentStatus
	}
}

func (l *Loop) recalculateCurrentByPast() {
	for point, cell := range l.cells { // TODO cell must have only one curent. loop must think about pust mb
		switch cell.previewStatus {
		case CellLive:
			l.createNeighborsIfNeed(NewPair(point.x-1, point.y-1))
			l.createNeighborsIfNeed(NewPair(point.x-1, point.y))
			l.createNeighborsIfNeed(NewPair(point.x-1, point.y+1))
			l.createNeighborsIfNeed(NewPair(point.x, point.y+1))
			l.createNeighborsIfNeed(NewPair(point.x+1, point.y+1))
			l.createNeighborsIfNeed(NewPair(point.x+1, point.y))
			l.createNeighborsIfNeed(NewPair(point.x+1, point.y-1))
			l.createNeighborsIfNeed(NewPair(point.x, point.y-1))
			countLiveNeighbors := l.calculateNeighbors(point)
			if countLiveNeighbors < 2 || countLiveNeighbors > 3 {
				cell.currentStatus = CellDied
			}
		case CellDied:
			if l.calculateNeighbors(point) == 3 {
				cell.currentStatus = CellLive
			}
		}
	}
}
func (l *Loop) createNeighborsIfNeed(point Pair) {
	if _, ok := l.cells[point]; !ok {
		l.cells[point] = NewCell(CellDied)
	}

	c := l.cells[point]
	if c.previewStatus == CellDied {
		if l.calculateNeighbors(point) == 3 {
			c.currentStatus = CellLive
		}
	}
}

func (l *Loop) calculateNeighbors(point Pair) int {
	countLiveNeighbors := 0

	l.updateNeighborCounter(&countLiveNeighbors, point, -1, -1)
	l.updateNeighborCounter(&countLiveNeighbors, point, -1, 0)
	l.updateNeighborCounter(&countLiveNeighbors, point, -1, 1)
	l.updateNeighborCounter(&countLiveNeighbors, point, 0, 1)
	l.updateNeighborCounter(&countLiveNeighbors, point, 1, 1)
	l.updateNeighborCounter(&countLiveNeighbors, point, 1, 0)
	l.updateNeighborCounter(&countLiveNeighbors, point, 1, -1)
	l.updateNeighborCounter(&countLiveNeighbors, point, 0, -1)

	return countLiveNeighbors
}

func (l *Loop) updateNeighborCounter(countLiveNeighbors *int, point Pair, xShift int, yShift int) {
	if n, ok := l.cells[Pair{x: point.x + xShift, y: point.y + yShift}]; ok {
		if n.previewStatus == CellLive {
			*(countLiveNeighbors)++
		}
	}
}

func (l *Loop) view() {
	for i := 0; i < len(l.viewFuncList); i++ {
		(l.viewFuncList[i])(l.cells)
	}
}
