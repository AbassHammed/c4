package game

import (
    "testing"
)

// Tests unitaires pour GameManager.
//
// Ce fichier contient des tests en français couvrant l'ensemble des chemins
// d'exécution de `game_manager.go` afin d'atteindre 100% de couverture.

// Test des valeurs par défaut après création d'un GameManager.
func TestGameManagerDefaults(t *testing.T) {
    gm := NewGameManager(false, 1)
    if gm == nil {
        t.Fatalf("NewGameManager returned nil")
    }
    if gm.GetState() != Running {
        t.Fatalf("expected state Running, got %v", gm.GetState())
    }
    if gm.turn != 0 {
        t.Fatalf("expected turn 0, got %d", gm.turn)
    }
    if gm.winner != "" {
        t.Fatalf("expected empty winner, got %q", gm.winner)
    }
    if gm.GetWonGames() != 0 || gm.GetLostGames() != 0 {
        t.Fatalf("expected zero stats, got won=%d lost=%d", gm.GetWonGames(), gm.GetLostGames())
    }

    // GetHoleColor hors limites doit retourner chaîne vide
    if c := gm.GetHoleColor(-1, 0); c != "" {
        t.Fatalf("expected empty string for out-of-bounds row, got %q", c)
    }
    if c := gm.GetHoleColor(boardHeight, 0); c != "" {
        t.Fatalf("expected empty string for out-of-bounds row, got %q", c)
    }

    // GetHoleColor dans les limites doit renvoyer la case vide
    if c := gm.GetHoleColor(0, 0); c != emptySpot {
        t.Fatalf("expected emptySpot at (0,0), got %q", c)
    }
}

// Test des mouvements du joueur: colonne invalide, puis victoire verticale.
func TestMakePlayerTurn_InvalidAndWin(t *testing.T) {
    gm := NewGameManager(false, 1)

    // colonne invalide
    if ok, err := gm.MakePlayerTurn(-5); err == nil || ok {
        t.Fatalf("expected error for invalid column, got ok=%v err=%v", ok, err)
    }

    // Préparer une victoire verticale en colonne 0 pour PlayerOneColor
    // poser 3 jetons de PlayerOneColor dans la colonne 0
    for i := 0; i < 3; i++ {
        if !gm.board.Drop(0, PlayerOneColor) {
            t.Fatalf("failed to prefill board at column 0")
        }
    }

    // Le tour courant est toujours le PlayerOne (turn initial 0)
    ok, err := gm.MakePlayerTurn(0)
    if err != nil || !ok {
        t.Fatalf("expected successful move to complete vertical win, got ok=%v err=%v", ok, err)
    }
    if gm.GetState() != Win {
        t.Fatalf("expected state Win after connecting four, got %v", gm.GetState())
    }
    if gm.GetWonGames() != 1 {
        t.Fatalf("expected wonGames == 1, got %d", gm.GetWonGames())
    }
}

// Test des mouvements de l'adversaire pour les modes AI et humain.
func TestMakeOpponentTurn_AIAndHuman(t *testing.T) {
    // AI mode
    gmAI := NewGameManager(true, 1)
    col, err := gmAI.MakeOpponentTurn(-1) // param ignoré quand gm.ai == true
    if err != nil {
        t.Fatalf("expected no error from AI opponent, got %v", err)
    }
    if col < 0 || col >= boardWidth {
        t.Fatalf("AI returned out-of-range column %d", col)
    }
    if gmAI.turn != 1 {
        t.Fatalf("expected gmAI.turn == 1 after opponent move, got %d", gmAI.turn)
    }

    // Human opponent mode: appel sans colonne doit renvoyer une erreur
    gmH := NewGameManager(false, 1)
    if _, err := gmH.MakeOpponentTurn(-1); err == nil {
        t.Fatalf("expected error when no column provided for human opponent")
    }
    // fournissons une colonne valide
    col2, err := gmH.MakeOpponentTurn(2)
    if err != nil {
        t.Fatalf("expected successful human opponent move, got error: %v", err)
    }
    if col2 != 2 {
        t.Fatalf("expected played column 2, got %d", col2)
    }
    if gmH.turn != 1 {
        t.Fatalf("expected gmH.turn == 1 after opponent move, got %d", gmH.turn)
    }
}

// Test de WhereConnected et ResetGame
func TestWhereConnectedAndReset(t *testing.T) {
    gm := NewGameManager(false, 1)

    connected, x, y := gm.WhereConnected()
    if connected {
        t.Fatalf("expected no connection when winner is empty")
    }
    for i := 0; i < 4; i++ {
        if x[i] != -1 || y[i] != -1 {
            t.Fatalf("expected -1 indices when no connection, got x=%v y=%v", x, y)
        }
    }

    // Créer une ligne horizontale gagnante pour PlayerOneColor en colonnes 0..3
    for c := 0; c < 4; c++ {
        if !gm.board.Drop(c, PlayerOneColor) {
            t.Fatalf("failed to drop token for prefill at column %d", c)
        }
    }
    gm.winner = PlayerOneColor
    connected, x, y = gm.WhereConnected()
    if !connected {
        t.Fatalf("expected a connection after placing four tokens horizontally")
    }
    // sur un plateau vierge, les jetons tombent sur la dernière ligne (index 5)
    for i := 0; i < 4; i++ {
        if x[i] != 5 {
            t.Fatalf("expected row index 5 for horizontal connect, got x[%d]=%d", i, x[i])
        }
        if y[i] != i {
            t.Fatalf("expected column index %d for horizontal connect, got y[%d]=%d", i, i, y[i])
        }
    }

    // ResetGame doit réinitialiser l'état et le winner
    gm.ResetGame()
    if gm.GetState() != Running {
        t.Fatalf("expected state Running after ResetGame, got %v", gm.GetState())
    }
    if gm.winner != "" {
        t.Fatalf("expected empty winner after ResetGame, got %q", gm.winner)
    }
    if gm.turn != 0 {
        t.Fatalf("expected turn 0 after ResetGame, got %d", gm.turn)
    }
    // le plateau doit être réinitialisé
    if c := gm.GetHoleColor(0, 0); c != emptySpot {
        t.Fatalf("expected emptySpot at (0,0) after ResetGame, got %q", c)
    }
}
