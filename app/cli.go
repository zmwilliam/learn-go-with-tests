package poker

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
)

const PlayerPrompt = "Please enter the number of players: "

const BadPlayerInputErrMsg = "Bad value received for number of players, please try again with a number"

const BadWinnerInputErrMsg = "invalid winner input, expect format of 'PlayerName wins'"

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
	fmt.Fprint(c.out, PlayerPrompt)

	numberOfPlayers, err := strconv.Atoi(c.readLine())
	if err != nil {
		fmt.Fprint(c.out, BadPlayerInputErrMsg)
		return
	}

	c.game.Start(numberOfPlayers)

	winner, err := extractWinner(c.readLine())
	if err != nil {
		fmt.Fprint(c.out, err)
		return
	}
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

func extractWinner(userInput string) (string, error) {
	if !strings.Contains(userInput, "wins") {
		return "", errors.New(BadWinnerInputErrMsg)
	}

	return strings.Replace(userInput, " wins", "", 1), nil
}
