package poker

import "io"

type Game interface {
	Start(numberOfPlayers int, alertsDestionation io.Writer)
	Finish(winner string)
}
