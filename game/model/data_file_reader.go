package model

import (
	"bufio"
	"os"
	"strconv"
	"strings"
)

const gameFileDir = ".\\game\\levels\\"
const newLine byte = '\n'
const carriageReturn byte = '\r'
const space byte = ' '

type levelReader struct {
	ix int
}

func newGameReaderReader() *levelReader {
	return &levelReader{ix: 1}
}

func newLevelReader(ix int) *levelReader {
	return &levelReader{ix: ix}
}

func LoadGameFiles() *Game {
	game := NewGame()

	for reader := newGameReaderReader(); reader.hasNext(); {
		game.AddLevel(reader.next())
	}

	return game
}

func LoadLevel(ix int) *Level {
	reader := newLevelReader(ix)
	if reader.hasNext() {
		return reader.next()
	}

	panic("level does not exist")
}

func (lr *levelReader) hasNext() bool {
	path := getPath(lr.ix)
	_, err := os.Stat(path)

	return !os.IsNotExist(err)
}

func (lr *levelReader) next() *Level {
	path := getPath(lr.ix)
	reader := openFile(path)

	width := parseInt(readUntil(reader, space))
	height := parseInt(readUntil(reader, carriageReturn))
	readNewLineOrFail(reader)

	player := readPlayer(reader)
	world, boxes := readWorld(reader, width, height)

	lr.ix++
	return NewLevel(width, height, player, world, boxes)
}

func readUntil(reader *bufio.Reader, delimiter byte) string {
	token, err := reader.ReadString(delimiter)
	if err != nil {
		panic(err)
	}

	return strings.Trim(token, string(delimiter))
}

func readPlayer(reader *bufio.Reader) *Player {
	x := parseInt(readUntil(reader, space))
	y := parseInt(readUntil(reader, space))
	direction := parseDirection(readUntil(reader, carriageReturn))
	readNewLineOrFail(reader)

	return NewPlayer(NewPosition(x, y), direction)
}

func readWorld(reader *bufio.Reader, width int, height int) (*[][]Tile, *[]*Position) {
	world := make([][]Tile, height)
	boxes := make([]*Position, 0)

	for y := 0; y < height; y++ {
		world[y] = make([]Tile, width)
		for x := 0; x < width; x++ {
			tile, isBox := readTile(reader)
			world[y][x] = tile
			if isBox {
				boxes = append(boxes, NewPosition(x, y))
			}
		}
		readNewLineOrFail(reader)
	}

	return &world, &boxes
}

func readTile(reader *bufio.Reader) (Tile, bool) {
	char, _, err := reader.ReadRune()
	if err != nil {
		panic(err)
	}

	switch char {
	case '#':
		return Wall, false
	case '.':
		return Floor, false
	case 'D':
		return DropZone, false
	case 'B':
		return Floor, true
	default:
		panic("unexpected character:" + string(char))
	}
}

func readNewLineOrFail(reader *bufio.Reader) {
	char, err := reader.ReadByte()

	if err != nil {
		panic(err)
	}

	if char == carriageReturn {
		char, _ = reader.ReadByte()
	}

	if char != newLine {
		panic("unexpected character:" + string(char))
	}
}

func getPath(ix int) string {
	return gameFileDir + strconv.FormatInt(int64(ix), 10) + ".txt"
}

func parseInt(source string) int {
	parsed, err := strconv.ParseInt(source, 10, 32)
	if err != nil {
		panic(err)
	}

	return int(parsed)
}

func parseDirection(source string) Direction {
	switch source {
	case "N":
		return North
	case "E":
		return East
	case "S":
		return South
	case "W":
		return West
	default:
		panic("unexpected direction: " + source)
	}
}

func openFile(path string) *bufio.Reader {
	file, err := os.Open(path)

	if err != nil {
		panic(err)
	}

	return bufio.NewReader(file)
}
