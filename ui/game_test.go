package ui

import (
	"testing"

	"github.com/AbassHammed/c4/game"
	"github.com/hajimehoshi/ebiten/v2"
)

// TestDraw_MenuAndEnterAI vérifie que Draw peut être appelé en menu et en
// saisie de difficulté sans provoquer de panic.
func TestDraw_MenuAndEnterAI(t *testing.T) {
	// preserve global state
	oldState := gameState
	defer func() { gameState = oldState }()

	screen := ebiten.NewImage(640, 640)
	gameState = menu
	g := &Game{}
	g.Draw(screen)

	gameState = enterAIdifficulty
	g.Draw(screen)
}

// TestDraw_GameplayAndWinAndGhost parcourt les chemins de rendu liés au jeu,
// y compris une victoire verticale et le chemin du fantôme (IA).
func TestDraw_GameplayAndWinAndGhost(t *testing.T) {
	// preserve globals
	oldGm := gm
	oldState := gameState
	defer func() { gm = oldGm; gameState = oldState }()

	// Crée une partie locale à deux et joue des coups menant à une victoire
	// verticale du joueur 1 dans la colonne 0.
	gm = game.NewGameManager(false, 0)
	moves := []int{0, 1, 0, 1, 0, 1, 0} // P1:0, P2:1, P1:0, ... -> P1 gets four in column 0
	for i, col := range moves {
		if i%2 == 0 {
			ok, err := gm.MakePlayerTurn(col)
			if err != nil || !ok {
				t.Fatalf("MakePlayerTurn failed at move %d col %d: %v", i, col, err)
			}
		} else {
			_, err := gm.MakeOpponentTurn(col)
			if err != nil {
				t.Fatalf("MakeOpponentTurn failed at move %d col %d: %v", i, col, err)
			}
		}
	}

	if gm.GetState() != game.Win {
		t.Fatalf("expected game state Win but got %v", gm.GetState())
	}

	// Dessine l'écran de fin de partie (exerce drawWinnerDots).
	screen := ebiten.NewImage(640, 640)
	gameState = win
	g := &Game{}
	g.Draw(screen)

	// Teste ensuite le rendu du fantôme pour un adversaire IA.
	gm = game.NewGameManager(true, 1)
	opponentLastCol = 3
	gameState = opponentAnimation
	g.Draw(screen)
}
