package game

import (
	"io"
	"os"
	"strings"
	"testing"
)

var _ = func() bool {
	testing.Init()
	return true
}()

func TestBoardWhereConnectedHorizontal(t *testing.T) {
	board := NewBoard()
	board.Drop(5, "*")
	board.Drop(4, "*")
	board.Drop(3, "*")
	board.Drop(2, "*")

	areConnected, row, col := board.WhereConnected("*")
	expectedCol := [4]int{2, 3, 4, 5}
	if col != expectedCol {
		t.Errorf("columns are incorrect, expected %v got %v", expectedCol, col)
	}
	expectedRow := [4]int{5, 5, 5, 5}
	if row != expectedRow {
		t.Errorf("rows are incorrect, expected %v got %v", expectedRow, row)
	}
	if !areConnected {
		t.Errorf("expected four to be connected but return false")
	}
}

func TestBoardWhereConnectedVertical(t *testing.T) {
	board := NewBoard()
	board.Drop(5, "*")
	board.Drop(5, "*")
	board.Drop(5, "*")
	board.Drop(5, "*")

	areConnected, row, col := board.WhereConnected("*")
	expectedCol := [4]int{5, 5, 5, 5}
	if col != expectedCol {
		t.Errorf("columns are incorrect, expected %v got %v", expectedCol, col)
	}
	expectedRow := [4]int{2, 3, 4, 5}
	if row != expectedRow {
		t.Errorf("rows are incorrect, expected %v got %v", expectedRow, row)
	}
	if !areConnected {
		t.Errorf("expected four to be connected but return false")
	}
}

func TestBoardWhereConnectedAscendingDiagonal(t *testing.T) {
	board := NewBoard()
	board.board[5][0] = "*"
	board.board[4][1] = "*"
	board.board[3][2] = "*"
	board.board[2][3] = "*"

	areConnected, row, col := board.WhereConnected("*")
	expectedCol := [4]int{0, 1, 2, 3}
	if col != expectedCol {
		t.Errorf("columns are incorrect, expected %v got %v", expectedCol, col)
	}
	expectedRow := [4]int{5, 4, 3, 2}
	if row != expectedRow {
		t.Errorf("rows are incorrect, expected %v got %v", expectedRow, row)
	}
	if !areConnected {
		t.Errorf("expected four to be connected but return false")
	}
}

func TestBoardWhereConnectedDescendingDiagonal(t *testing.T) {
	board := NewBoard()
	// positions: [5][3], [4][2], [3][1], [2][0]
	board.board[5][3] = "*"
	board.board[4][2] = "*"
	board.board[3][1] = "*"
	board.board[2][0] = "*"

	areConnected, row, col := board.WhereConnected("*")
	expectedCol := [4]int{3, 2, 1, 0}
	if col != expectedCol {
		t.Errorf("columns are incorrect, expected %v got %v", expectedCol, col)
	}
	expectedRow := [4]int{5, 4, 3, 2}
	if row != expectedRow {
		t.Errorf("rows are incorrect, expected %v got %v", expectedRow, row)
	}
	if !areConnected {
		t.Errorf("expected four to be connected but return false")
	}
}

func TestDropBoundaries(t *testing.T) {
	b := NewBoard()
	if b.Drop(-1, "*") {
		t.Errorf("Drop should return false for negative column")
	}
	if b.Drop(len(b.board[0]), "*") {
		t.Errorf("Drop should return false for column >= width")
	}
}

func TestDropColumnFullAndMovesMade(t *testing.T) {
	b := NewBoard()
	col := 0
	// fill the column
	for i := 0; i < len(b.board); i++ {
		if !b.Drop(col, "*") {
			t.Fatalf("expected Drop to succeed at iteration %d", i)
		}
	}
	// now the column should be full
	if b.Drop(col, "*") {
		t.Errorf("Drop should fail when column is full")
	}
	if b.col[col] != len(b.board) {
		t.Errorf("expected col[%d] == %d got %d", col, len(b.board), b.col[col])
	}
}

func TestUndoDrop(t *testing.T) {
	b := NewBoard()
	if !b.Drop(0, "*") {
		t.Fatalf("Drop failed when it should succeed")
	}
	if b.col[0] != 1 || b.movesMade != 1 {
		t.Fatalf("unexpected state after Drop: col[0]=%d movesMade=%d", b.col[0], b.movesMade)
	}
	b.undoDrop(0)
	if b.col[0] != 0 || b.movesMade != 0 {
		t.Fatalf("unexpected state after undoDrop: col[0]=%d movesMade=%d", b.col[0], b.movesMade)
	}
	if b.board[5][0] != emptySpot {
		t.Fatalf("expected bottom cell to be empty after undoDrop, got %q", b.board[5][0])
	}
}

func TestCopyOfBoardDeepCopy(t *testing.T) {
	b := NewBoard()
	if !b.Drop(1, "*") {
		t.Fatalf("initial Drop failed")
	}
	copy := b.copyOfBoard()
	// modify the copy and ensure original is unchanged
	if !copy.Drop(1, "o") {
		t.Fatalf("Drop on copy failed")
	}
	// original should still have only the first drop at row 5
	if b.board[5][1] != "*" {
		t.Fatalf("original board modified after mutating copy: expected '*' at [5][1], got %q", b.board[5][1])
	}
	// the copy should have the new token somewhere above the original one
	if copy.board[4][1] != "o" {
		t.Fatalf("expected copy to have 'o' at [4][1], got %q", copy.board[4][1])
	}
}

func TestNewBoardInitialization(t *testing.T) {
	b := NewBoard()
	if b.movesMade != 0 {
		t.Fatalf("expected movesMade == 0, got %d", b.movesMade)
	}
	for j := 0; j < len(b.col); j++ {
		if b.col[j] != 0 {
			t.Fatalf("expected col[%d] == 0, got %d", j, b.col[j])
		}
	}
	for i := 0; i < len(b.board); i++ {
		for j := 0; j < len(b.board[0]); j++ {
			if b.board[i][j] != emptySpot {
				t.Fatalf("expected board[%d][%d] == emptySpot, got %q", i, j, b.board[i][j])
			}
		}
	}
}

func TestGameOverByMovesMade(t *testing.T) {
	b := NewBoard()
	b.movesMade = 42
	if !b.gameOver() {
		t.Fatalf("expected gameOver to return true when movesMade == 42")
	}
}

func TestWhereConnectedSpecificToPlayer(t *testing.T) {
	b := NewBoard()
	b.board[5][0] = "*"
	b.board[5][1] = "o"
	b.board[5][2] = "*"
	b.board[5][3] = "*"
	b.board[5][4] = "*"
	// there are not four contiguous '*' so this should be false
	areConnected, _, _ := b.WhereConnected("*")
	if areConnected {
		t.Fatalf("expected WhereConnected to be false for '*' since sequence is interrupted by 'o'")
	}
}

func TestAreFourConnectedDirect(t *testing.T) {
	b := NewBoard()
	if !b.Drop(0, "*") || !b.Drop(1, "*") || !b.Drop(2, "*") || !b.Drop(3, "*") {
		t.Fatalf("failed to place four tokens for AreFourConnected test")
	}
	if !b.areFourConnected("*") {
		t.Fatalf("expected areFourConnected to return true for four in a row")
	}
}

func TestPrintBoardOutput(t *testing.T) {
	b := NewBoard()
	// place a token to ensure output contains a non-empty symbol
	if !b.Drop(0, "*") {
		t.Fatalf("initial Drop failed")
	}

	// capture stdout
	oldStdout := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("failed to create pipe: %v", err)
	}
	os.Stdout = w

	b.printBoard()

	_ = w.Close()
	outBytes, err := io.ReadAll(r)
	if err != nil {
		t.Fatalf("failed to read pipe: %v", err)
	}
	os.Stdout = oldStdout
	out := string(outBytes)

	// basic checks on output content
	if !strings.Contains(out, "0 1 2 3 4 5 6") {
		t.Errorf("printBoard output missing column headers: %q", out)
	}
	if !strings.Contains(out, "*") {
		t.Errorf("printBoard output missing placed token '*' : %q", out)
	}
	if !strings.Contains(out, emptySpot) {
		t.Errorf("printBoard output missing empty spot symbol %q: %q", emptySpot, out)
	}
}
