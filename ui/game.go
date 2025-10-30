package ui

import (
	"bytes"
	"image"

	"image/color"
	_ "image/jpeg"
	_ "image/png"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/AbassHammed/c4/game"
	"github.com/AbassHammed/c4/images"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text"
	textv2 "github.com/hajimehoshi/ebiten/v2/text/v2"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

var backgroundImage,
	owl,
	redBallImage,
	dot,
	ghost,
	greenBallImage,
	boardImage,
	bats *ebiten.Image

func byteSliceToEbitenImage(arr []byte) *ebiten.Image {
	img, _, err := image.Decode(bytes.NewReader(arr))
	if err != nil {
		log.Fatal(err)
	}
	return ebiten.NewImageFromImage(img)
}

func init() {
	ghost = byteSliceToEbitenImage(images.Ghost_png)
	backgroundImage = byteSliceToEbitenImage(images.Background_png)
	redBallImage = byteSliceToEbitenImage(images.Red_png)
	greenBallImage = byteSliceToEbitenImage(images.Green_png)
	owl = byteSliceToEbitenImage(images.Owl_png)
	dot = byteSliceToEbitenImage(images.Dot_png)
	bats = byteSliceToEbitenImage(images.Bats_png)
	boardImage = byteSliceToEbitenImage(images.Board_png)
	tt, _ := opentype.Parse(images.MPlus1pRegular_ttf)
	mplusNormalFont, _ = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    20,
		DPI:     72,
		Hinting: font.HintingFull,
	})

	// Crée un adaptateur text/v2 Face depuis un golang.org/x/image/font.Face
	// afin d'utiliser text/v2.Draw (qui attend un text.Face).
	// NewGoXFace enveloppe font.Face et fournit la mise en cache des glyphes.
	tvFace = textv2.NewGoXFace(mplusNormalFont)
	initBallYCoords()
}

type Game struct{}

type GameState int

func initBallYCoords() {
	for i := 0; i < 7; i++ {
		for j := 0; j < 6; j++ {
			ballYcoords[i][j] = -tileHeight
		}
	}
}

const (
	yourTurn GameState = iota
	opponentTurn
	win
	lose
	tie
	animation
	opponentAnimation
	menu
	enterAIdifficulty
)

const (
	batsX             = 440
	batsY             = 200
	secondsToMakeTurn = 59
	fps               = 60
	tileHeight        = 65
	tileOffset        = 10
	boardX            = 84
	boardY            = 130
	gravity           = 0.5
)

// colonne choisie par l'adversaire lors du dernier coup
var opponentLastCol int
var frameCount int
var gameState GameState = menu

var ballYcoords [7][6]float64
var ballFallSpeed [7][6]float64

var mplusNormalFont font.Face
var tvFace textv2.Face

// messages affichés pendant une partie
var messages [7]string = [7]string{"Your turn", "Other's turn", "You win!", "You lost.", "Tie.", "...", "..."}

// gm is le gestionnaire de partie (peut être nil si pas de partie en cours)
var gm *game.GameManager

func changeGameStateBasedOnGameManagerState(gmState game.GameState) {
	if gmState != game.Running {
		switch gmState {
		case game.Win:
			gameState = win
		case game.Lose:
			gameState = lose
		case game.Tie:
			gameState = tie
		}
	}
}

func updateBallPos() {
	for i := 0; i < 6; i++ {
		for j := 0; j < 7; j++ {
			if gm == nil {
				continue
			}
			if gm.GetHoleColor(i, j) == game.PlayerTwoColor ||
				gm.GetHoleColor(i, j) == game.PlayerOneColor {
				y, x := i, j
				destY := float64(y) * tileHeight
				fallY := &ballYcoords[x][y]
				fallSpeed := &ballFallSpeed[x][y]

				*fallY += *fallSpeed
				*fallSpeed += gravity
				if *fallY > destY {
					*fallY = destY
					*fallSpeed = 0
				}
			}
		}
	}
}

// logique principale du jeu : transitions d'état et démarrage d'une partie
func (g *Game) Update() error {
	press := inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft)
	// lire les caractères tapés (gère AZERTY et autres dispositions)
	inputRunes := ebiten.AppendInputChars(nil)

	if gameState == yourTurn || gameState == opponentTurn {
		frameCount++
	}

	if gameState == animation || gameState == opponentAnimation {
		frameCount = 0
		updateBallPos()
	}

	if frameCount == fps*secondsToMakeTurn {
		os.Exit(1)
	}

	// mise à jour des positions des billes (désactivée)

	if (gameState == yourTurn || gameState == opponentTurn) && press {
		mouseX, _ := ebiten.CursorPosition()
		if gm != nil {
			prevState := gameState
			ok, _ := gm.MakePlayerTurn(xcoordToColumn(mouseX))
			if ok {
				// show animation for the drop
				gameState = animation
				go func(prev GameState) {
					time.Sleep(1 * time.Second)
					// si la partie est terminée, mettre à jour l'état final
					gmState := gm.GetState()
					if gmState != game.Running {
						changeGameStateBasedOnGameManagerState(gmState)
						return
					}
					// si l'adversaire est une IA, planifier son coup
					if gm.IsAI() {
						// après le coup du joueur, l'IA joue
						gameState = opponentTurn
					} else {
						// jeu local à deux : basculer le tour
						if prev == yourTurn {
							gameState = opponentTurn
						} else {
							gameState = yourTurn
						}
					}
				}(prevState)
			}
		}
	}

	if gameState == opponentTurn {
		// jouer automatiquement uniquement si l'adversaire est IA ; sinon attendre l'entrée utilisateur
		if gm != nil && gm.IsAI() {
			gameState = animation
			go func() {
				col, _ := gm.MakeOpponentTurn(-1)
				opponentLastCol = col
				gameState = opponentAnimation
				time.Sleep(1 * time.Second)
				if gm != nil {
					gmState := gm.GetState()
					if gmState == game.Running {
						gameState = yourTurn
					} else {
						changeGameStateBasedOnGameManagerState(gmState)
					}
				} else {
					gameState = yourTurn
				}
			}()
		}
	}

	if gameState == menu {
		for _, r := range inputRunes {
			switch r {
			case 'a', 'A':
				gameState = enterAIdifficulty
			case 'p', 'P':
				gm = game.NewGameManager(false, 0)
				gameState = yourTurn
			}
		}
	}

	if gameState == enterAIdifficulty {
		runes := ebiten.AppendInputChars(nil)
		if len(runes) == 1 {
			diff := string(runes)
			difficulty, err := strconv.Atoi(diff)
			if err == nil {
				gameState = yourTurn
				gm = game.NewGameManager(true, difficulty+3)
			}
		}
	}
	// Partie locale à deux sur le même clavier : démarrée par 'P' (gérée ci-dessus via inputRunes)

	if isGameOver() && press {
		mouseX, mouseY := ebiten.CursorPosition()
		/*check if mouse is in play again area
		 */
		if mouseX >= 230 && mouseX <= 600 && mouseY >= 500 {

			gmState := gm.GetState()
			gm.ResetGame()
			var s [7][6]float64
			ballFallSpeed = s
			initBallYCoords()
			if gmState == game.Win {
				gameState = opponentTurn
			} else {
				gameState = yourTurn
			}
		}
	}
	return nil
}

// isGameOver returns whether the game is over
func isGameOver() bool {
	return gameState == tie || gameState == win || gameState == lose
}

// dessine l'interface en fonction de l'état de la partie
func (g *Game) Draw(screen *ebiten.Image) {
	screen.DrawImage(backgroundImage, nil)
	op := &ebiten.DrawImageOptions{}

	op.GeoM.Translate(batsX, batsY)
	screen.DrawImage(bats, op)
	op.GeoM.Reset()

	op.GeoM.Translate(boardX, boardY)
	if gameState == menu {
		screen.DrawImage(boardImage, op)
		// Utilise text/v2.Draw avec l'adaptateur GoXFace (tvFace). La position est définie
		// via DrawOptions.DrawImageOptions.GeoM.Translate.
		o1 := &textv2.DrawOptions{}
		o1.DrawImageOptions.GeoM.Translate(float64(boardX), float64(boardY-30))
		textv2.Draw(screen, "[A] - play against AI", tvFace, o1)

		o2 := &textv2.DrawOptions{}
		o2.DrawImageOptions.GeoM.Translate(float64(boardX), float64(570))
		textv2.Draw(screen, "[P] - play local (2 players)", tvFace, o2)
		return
	}

	if gameState == enterAIdifficulty {
		screen.DrawImage(boardImage, op)
		o := &textv2.DrawOptions{}
		o.DrawImageOptions.GeoM.Translate(200, 50)
		textv2.Draw(screen, "Enter difficulty (1-9)", tvFace, o)
		return
	}

	var msg string = messages[gameState]
	text.Draw(screen, "W  "+strconv.Itoa(gm.GetWonGames())+":"+strconv.Itoa(gm.GetLostGames())+"  L", mplusNormalFont, boardX, 50, color.White)
	text.Draw(screen, msg, mplusNormalFont, boardX, 580, color.White)
	text.Draw(screen, "00:"+strconv.Itoa(secondsToMakeTurn-frameCount/fps), mplusNormalFont, 500, 580, color.White)

	drawOwl(screen)
	if gameState == opponentAnimation {
		drawGhost(screen)
	}

	drawBalls(screen)
	screen.DrawImage(boardImage, op)

	if isGameOver() {
		text.Draw(screen, "Click here\nto play again", mplusNormalFont, 250, 580, color.White)
		if gameState != tie {
			drawWinnerDots(screen)
		}
	}
}

// dessine toutes les billes à l'écran
func drawBalls(screen *ebiten.Image) {
	if gm == nil {
		return
	}
	for i := 0; i < 6; i++ {
		for j := 0; j < 7; j++ {
			if gm.GetHoleColor(i, j) == game.PlayerTwoColor {
				drawBall(j, i, game.PlayerTwoColor, screen)
			} else if gm.GetHoleColor(i, j) == game.PlayerOneColor {
				drawBall(j, i, game.PlayerOneColor, screen)
			}
		}
	}
}

// dessine les points indiquant les quatre jetons gagnants
func drawWinnerDots(screen *ebiten.Image) {
	if gm == nil {
		return
	}
	win, dotsY, dotsX := gm.WhereConnected()
	if !win {
		return
	}
	for i := 0; i < 4; i++ {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(boardX+tileOffset, boardY+tileOffset)
		op.GeoM.Translate(float64(dotsX[i])*tileHeight+25, float64(dotsY[i])*tileHeight+25)
		screen.DrawImage(dot, op)
	}
}

// dessine l'image fantôme à l'écran
func drawGhost(screen *ebiten.Image) {
	if gm == nil || !gm.IsAI() {
		return
	}
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(opponentLastCol)*tileHeight+boardX+10, boardY-75)
	screen.DrawImage(ghost, op)
}

// dessine le hibou à l'écran
func drawOwl(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	mouseX, _ := ebiten.CursorPosition()
	if mouseX < boardX {
		mouseX = boardX
	}
	if mouseX > boardX+7*tileHeight {
		mouseX = boardX + 7*tileHeight
	}
	owlX := xcoordToColumn(mouseX)*tileHeight + boardX
	op.GeoM.Translate(float64(owlX), boardY-80)
	screen.DrawImage(owl, op)
}

// dessine une bille à l'écran
func drawBall(x, y int, player string, screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(boardX+tileOffset, boardY+tileOffset)
	op.GeoM.Translate(float64(x)*tileHeight, ballYcoords[x][y])

	if player == game.PlayerTwoColor {
		screen.DrawImage(redBallImage, op)
	} else {
		screen.DrawImage(greenBallImage, op)
	}
}

// updateBallsPos supprimée : la mise à jour des positions est effectuée par updateBallPos

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 640, 640
}

// xcoordToColumn returns the column correspondidng which contains the x coordinate
func xcoordToColumn(x int) int {
	return int(float64(x-tileOffset-boardX) / tileHeight)
}

// StartGuiGame initializes the game and the gui, this is the entry point for the whole game
func StartGuiGame() {
	ebiten.SetWindowSize(640, 640)
	ebiten.SetWindowTitle("Connect four")
	if err := ebiten.RunGame(&Game{}); err != nil {
		log.Fatal(err)
	}
}
