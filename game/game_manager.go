package game

import (
	"fmt"
)

// GameManager gère le déroulement d'une partie de Puissance 4.
// Il maintient l'état du jeu, gère les tours des joueurs et de l'IA,
// et compte les victoires/défaites.
type GameManager struct {
	board     Board     // Plateau de jeu
	ai        bool      // true si l'adversaire est une IA
	turn      int       // Numéro du tour actuel
	state     GameState // État actuel de la partie
	winner    string    // Symbole du joueur gagnant ("" si pas de gagnant)
	aiDiff    int       // Niveau de difficulté de l'IA
	lostGames int       // Nombre de parties perdues
	wonGames  int       // Nombre de parties gagnées
}

// GameState représente l'état d'une partie.
type GameState int

const (
	Running GameState = iota // Partie en cours
	Win                      // Victoire du joueur
	Lose                     // Défaite du joueur
	Tie                      // Match nul
)

const (
	PlayerOneColor = "◯" // Symbole du joueur 1
	PlayerTwoColor = "⬤" // Symbole du joueur 2
)

// NewGameManager crée un nouveau gestionnaire de partie.
// Le paramètre ai indique si l'adversaire est contrôlé par l'IA,
// aiDiff définit le niveau de difficulté de l'IA.
func NewGameManager(ai bool, aiDiff int) *GameManager {
	b := *NewBoard()
	return &GameManager{board: b, ai: ai, aiDiff: aiDiff, turn: 0, state: Running, winner: ""}
}

// GetHoleColor renvoie le symbole à la position (i,j) du plateau.
// Renvoie une chaîne vide si la position est hors limites.
func (gm *GameManager) GetHoleColor(i, j int) string {
	if i < 0 || i >= boardHeight || j < 0 || j >= boardWidth {
		return ""
	}
	return gm.board.board[i][j]
}

// currentToken renvoie le symbole du joueur dont c'est le tour.
func (gm *GameManager) currentToken() string {
	if gm.turn%2 == 0 {
		return PlayerOneColor
	}
	return PlayerTwoColor
}

// GetState renvoie l'état actuel de la partie.
func (gm *GameManager) GetState() GameState {
	return gm.state
}

// MakePlayerTurn tente de placer un jeton dans la colonne spécifiée.
// Renvoie (true, nil) si le coup est valide, (false, error) sinon.
func (gm *GameManager) MakePlayerTurn(column int) (bool, error) {
	if column < 0 || column >= boardWidth {
		return false, fmt.Errorf("column %d out of range", column)
	}
	tok := gm.currentToken()
	if gm.board.Drop(column, tok) {
		if gm.board.areFourConnected(tok) {
			gm.state = Win
			gm.winner = tok
			gm.wonGames++
		}
		gm.turn++
		if gm.turn == 42 {
			gm.state = Tie
		}
		return true, nil
	}
	return false, fmt.Errorf("invalid move: column %d is full or invalid", column)
}

// MakeOpponentTurn performs the opponent's move. If the GameManager is configured
// to use an AI (gm.ai == true) the AI will pick a column and providedColumn is ignored.
// For a human opponent (gm.ai == false) the caller must pass the chosen column
// via providedColumn. The method returns the column that was played and an error if
// the move was invalid.
func (gm *GameManager) MakeOpponentTurn(providedColumn int) (int, error) {
	var column int
	if gm.ai {
		column = getAiMove(&gm.board, gm.aiDiff)
	} else {
		if providedColumn < 0 || providedColumn >= boardWidth {
			return -1, fmt.Errorf("no valid column provided for opponent")
		}
		column = providedColumn
	}

	tok := gm.currentToken()
	if !gm.board.Drop(column, tok) {
		return column, fmt.Errorf("invalid move: column %d is full or invalid", column)
	}

	if gm.board.areFourConnected(tok) {
		gm.state = Lose
		gm.lostGames++
		gm.winner = tok
	}
	gm.turn++
	if gm.turn == 42 {
		gm.state = Tie
	}
	return column, nil
}

// WhereConnected renvoie les coordonnées des quatre jetons alignés s'il y a un gagnant.
// Retourne (false, [-1,-1,-1,-1], [-1,-1,-1,-1]) si pas de gagnant.
func (gm *GameManager) WhereConnected() (bool, [4]int, [4]int) {
	if gm.winner == "" {
		return false, [4]int{-1, -1, -1, -1}, [4]int{-1, -1, -1, -1}
	}
	connected, x, y := gm.board.WhereConnected(gm.winner)
	return connected, x, y
}

// ResetGame réinitialise le plateau et l'état de la partie, sans modifier
// le compteur de victoires/défaites.
func (gm *GameManager) ResetGame() {
	gm.board = *NewBoard()
	gm.turn = 0
	gm.state = Running
	gm.winner = ""
}

// GetWonGames renvoie le nombre de parties gagnées.
func (gm *GameManager) GetWonGames() int {
	return gm.wonGames
}

// GetLostGames renvoie le nombre de parties perdues.
func (gm *GameManager) GetLostGames() int {
	return gm.lostGames
}

// IsAI indique si l'adversaire est contrôlé par l'IA.
func (gm *GameManager) IsAI() bool {
	return gm.ai
}
