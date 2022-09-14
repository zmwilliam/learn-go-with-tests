package poker

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
)

const PlayerPrompt = "Please enter the number of players: "

type Game interface {
	Start(numberOfPlayers int)
	Finish(winner string)
}

type CLI struct {
	game Game
	in   *bufio.Scanner
	out  io.Writer
}

func (c *CLI) PlayPoker() {
	fmt.Fprintf(c.out, PlayerPrompt)

	numberOfPlayers, _ := strconv.Atoi(c.readLine())

	c.game.Start(numberOfPlayers)

	winner := extractWinner(c.readLine())
	c.game.Finish(winner)
}

func (c *CLI) readLine() string {
	c.in.Scan()
	return c.in.Text()
}

func NewCLI(in io.Reader, out io.Writer, game Game) *CLI {
	return &CLI{
		in:   bufio.NewScanner(in),
		out:  out,
		game: game,
	}
}

func extractWinner(userInput string) string {
	return strings.Replace(userInput, " wins", "", 1)
}
