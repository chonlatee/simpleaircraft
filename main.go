package main

import (
	"image/color"
	_ "image/png"
	"log"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

const screenHeight = 768
const screenWidth = 1024

type position struct {
	x float64
	y float64
}

type velocity struct {
	x float64
	y float64
}

type Game struct {
	aircraft   *aircraft
	boss       *boss
	bossHP     *bossHP
	aircraftHP *aircraftHP
}

func NewGame() *Game {
	g := &Game{}
	g.init()
	return g
}

func (g *Game) init() {
	g.boss = NewBoss()
	g.bossHP = NewBossHP()
	g.aircraftHP = NewAircraftHP()
	g.aircraft = NewAircraft()
}

func (g *Game) Update() error {
	g.aircraft.Update()
	g.boss.Update()

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{
		R: 230,
		G: 225,
		B: 225,
		A: 230,
	})

	g.aircraft.Draw(screen)
	g.boss.Draw(screen)

	for i, b := range g.aircraft.bullets {
		bullet := g.aircraft.bullets[i]
		if g.boss.IsHit(bullet) {
			b.ResetPosition(g.aircraft)
			g.boss.hit += 1
		}
		b.Draw(screen)
	}

	for i, _ := range g.boss.bullets {
		bullet := g.boss.bullets[i]
		if g.aircraft.isHit(bullet) {
			bullet.ResetPosition(g.boss)
			g.aircraft.hit += 1
		}
		bullet.Draw(screen)

	}

	g.bossHP.Draw(screen, g.boss)
	g.aircraftHP.Draw(screen, g.aircraft)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	rand.Seed(time.Now().UnixNano())
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Aircraft adventure")
	g := NewGame()
	if err := ebiten.RunGame(g); err != nil {
		log.Fatalln("Can't run game: ", err)
	}
}
