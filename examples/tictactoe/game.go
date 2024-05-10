package ticktacktoe

import "fmt"

type symbol string

type board [3][3]*symbol

// Player represents a player in a game of tic-tac-toe.
// @rpc:resource
type Player struct {
	Name         string `json:"name"`
	Symbol       symbol `json:"symbol" rpcType:"string"`
	privateValue int
}

// BoardService is a service that allows two players to play a game of tic-tac-toe.
// @rpc:service
type BoardService interface {
	NewGame() error
	Play(player Player, x, y int) error
	Winner() (Player, bool, error)
}

type boardService struct {
	board board
}

func (b *boardService) newGame() error {
	b.board = board{}
	return nil
}

func (b *boardService) play(player Player, x, y int) error {
	if b.board[x][y] != nil {
		return fmt.Errorf("cell already occupied")
	}

	b.board[x][y] = &player.Symbol
	return nil
}

func (b *boardService) winner() (Player, bool, error) {
	for x := range 3 {
		for y := range 3 {
			if b.board[x][y] == nil {
				continue
			}

			s := *b.board[x][y]
			if b.board[x][(y+1)%3] == &s && b.board[x][(y+2)%3] == &s {
				return Player{Symbol: s}, true, nil
			}

			if b.board[(x+1)%3][y] == &s && b.board[(x+2)%3][y] == &s {
				return Player{Symbol: s}, true, nil
			}

			if x == y && b.board[(x+1)%3][(y+1)%3] == &s && b.board[(x+2)%3][(y+2)%3] == &s {
				return Player{Symbol: s}, true, nil
			}

			if x == 2-y && b.board[(x+1)%3][(y-1)%3] == &s && b.board[(x+2)%3][(y-2)%3] == &s {
				return Player{Symbol: s}, true, nil
			}
		}
	}

	return Player{}, false, nil
}
