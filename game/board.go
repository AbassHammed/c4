package game

import (
	"fmt"
	"strings"
)

// Board représente l'état du plateau de Puissance 4.
//
// Champs :
// - board : matrice (hauteur x largeur) contenant les symboles des cases.
// - col : slice indiquant combien de jetons sont déjà placés par colonne.
// - movesMade : nombre total de coups joués sur le plateau.
type Board struct {
	board     [][]string
	col       []int
	movesMade int
}

// Constantes de configuration du plateau : largeur, hauteur et symbole
// utilisé pour représenter une case vide.
const (
	boardWidth  = 7
	boardHeight = 6
	emptySpot   = "∟"
)

// gameOver retourne true si la partie est terminée :
// - soit le nombre maximal de coups a été atteint (toutes les cases remplies),
// - soit un joueur a quatre jetons connectés.
func (b *Board) gameOver() bool {
	return b.movesMade == 42 || b.areFourConnected("<>") || b.areFourConnected("<>")
}

// copyOfBoard renvoie une copie profonde du plateau courant. La copie
// est indépendante de l'original (modifications ultérieures n'affectent pas
// l'instance source).
func (b *Board) copyOfBoard() *Board {
	boardCopy := NewBoard()
	for i := 0; i < boardHeight; i++ {
		for j := 0; j < boardWidth; j++ {
			boardCopy.board[i][j] = b.board[i][j]
		}
	}
	boardCopy.col = make([]int, boardWidth)
	copy(boardCopy.col, b.col)
	return boardCopy
}

// NewBoard crée et retourne un nouveau plateau pour une partie de
// Puissance 4. Toutes les cases sont initialisées avec le symbole
// représentant une case vide et les compteurs de colonnes sont remis à zéro.
func NewBoard() *Board {
	var b *Board
	b = new(Board)
	b.movesMade = 0
	b.col = make([]int, boardWidth)
	// initialize the connect 4 board
	for i := 0; i < boardHeight; i++ {
		row := make([]string, boardWidth)

		for i := 0; i < len(row); i++ {
			row[i] = emptySpot
		}
		b.board = append(b.board, row)
	}
	return b
}

// printBoard affiche le plateau sur la sortie standard de façon lisible.
func (b *Board) printBoard() {
	space := strings.Repeat(" ", 20)
	fmt.Print(space)
	for i := 0; i < len(b.board[0]); i++ {
		fmt.Printf("%d ", i)
	}
	fmt.Println()
	for i := 0; i < len(b.board); i++ {
		fmt.Print(space)
		for j := 0; j < len(b.board[0]); j++ {
			fmt.Print(b.board[i][j] + " ")
		}
		fmt.Println()
	}
}

// undoDrop annule le dernier dépôt dans la colonne spécifiée en
// décrémentant le compteur de la colonne et réinitialisant la case
// correspondante à la valeur vide.
func (b *Board) undoDrop(column int) {
	b.col[column]--
	b.board[5-b.col[column]][column] = emptySpot
	b.movesMade--
}

// Drop tente de placer le jeton du joueur dans la colonne indiquée.
// Retourne true si le jeton a été placé avec succès, false sinon
// (colonne invalide ou pleine).
func (b *Board) Drop(column int, player string) bool {
	if column < len(b.board[0]) && (column >= 0) && b.col[column] < len(b.board) {
		b.board[5-b.col[column]][column] = player
		b.col[column]++
		b.movesMade++
		return true
	}
	return false
}

// WhereConnected recherche s'il existe quatre jetons consécutifs du
// joueur fourni. Si trouvé, retourne true ainsi que deux tableaux de
// 4 entiers représentant les indices de ligne et de colonne des
// quatre positions; sinon retourne false et des tableaux remplis de -1.
func (b *Board) WhereConnected(player string) (bool, [4]int, [4]int) {
	for j := 0; j < len(b.board[0])-3; j++ {
		for i := 0; i < len(b.board); i++ {
			if b.board[i][j] == player &&
				b.board[i][j+1] == player &&
				b.board[i][j+2] == player &&
				b.board[i][j+3] == player {
				return true, [4]int{i, i, i, i}, [4]int{j, j + 1, j + 2, j + 3}
			}
		}
	}
	// verticalCheck
	for i := 0; i < len(b.board)-3; i++ {
		for j := 0; j < len(b.board[0]); j++ {
			if b.board[i][j] == player &&
				b.board[i+1][j] == player &&
				b.board[i+2][j] == player &&
				b.board[i+3][j] == player {
				return true, [4]int{i, i + 1, i + 2, i + 3}, [4]int{j, j, j, j}
			}
		}
	}
	// ascendingDiagonalCheck
	for i := 3; i < len(b.board); i++ {
		for j := 0; j < len(b.board[0])-3; j++ {
			if b.board[i][j] == player &&
				b.board[i-1][j+1] == player &&
				b.board[i-2][j+2] == player &&
				b.board[i-3][j+3] == player {
				return true, [4]int{i, i - 1, i - 2, i - 3}, [4]int{j, j + 1, j + 2, j + 3}
			}
		}
	}
	// descendingDiagonalCheck
	for i := 3; i < len(b.board); i++ {
		for j := 3; j < len(b.board[0]); j++ {
			if b.board[i][j] == player &&
				b.board[i-1][j-1] == player &&
				b.board[i-2][j-2] == player &&
				b.board[i-3][j-3] == player {
				return true, [4]int{i, i - 1, i - 2, i - 3}, [4]int{j, j - 1, j - 2, j - 3}
			}
		}
	}
	return false, [4]int{-1, -1, -1, -1}, [4]int{-1, -1, -1, -1}
}

// areFourConnected retourne true si le joueur a quatre jetons
// connectés sur le plateau (utilise WhereConnected pour la détection).
func (b *Board) areFourConnected(player string) bool {
	connected, _, _ := b.WhereConnected(player)
	return connected
}
