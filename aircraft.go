package main

import (
	"image"
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type aircraft struct {
	pos       position
	img       *ebiten.Image
	bullets   []*aircraftBullet
	bulletNum int
	vel       velocity
	width     int
	height    int
	hit       int
}

func NewAircraft() *aircraft {
	p := &aircraft{}
	p.init()
	return p
}

func (p *aircraft) init() {
	fr, err := ebitenutil.OpenFile("aircraft.png")
	if err != nil {
		log.Fatalln("can't open file: ", err)
	}
	img, _, err := image.Decode(fr)
	imgOrig := ebiten.NewImageFromImage(img)
	p.img = imgOrig
	p.pos.x = screenWidth / 2
	p.pos.y = screenHeight / 2
	p.vel.x = 3
	p.vel.y = 3
	p.width, p.height = p.img.Size()
	p.bulletNum = 1
	p.bullets = make([]*aircraftBullet, p.bulletNum)

	for i := 0; i < p.bulletNum; i++ {
		p.bullets[i] = NewAircraftBullet()
		p.bullets[i].pos.y = p.pos.y - 1
		p.bullets[i].pos.x = p.pos.x + float64(p.width/2)
	}

}

func (p *aircraft) isHit(bullet *bossBullet) bool {
	return bullet.pos.y >= p.pos.y+30 &&
		bullet.pos.y <= p.pos.y+50 &&
		bullet.pos.x >= p.pos.x+(float64(p.width/2)-10) &&
		bullet.pos.x <= p.pos.x+(float64(p.width/2)+10)
}

func (p *aircraft) Update() {
	p.move()
	for i, _ := range p.bullets {
		bullet := p.bullets[i]
		if bullet.pos.y <= 0 {
			bullet.pos.y = p.pos.y - 1
			bullet.pos.x = p.pos.x + float64(p.width/2)
		}
		bullet.pos.y -= bullet.vel.y
	}
}

func (p *aircraft) move() {
	for _, k := range inpututil.PressedKeys() {
		if k == ebiten.KeyUp {
			if p.pos.y >= 210 {
				p.pos.y -= p.vel.y
			}
		} else if k == ebiten.KeyDown {
			if p.pos.y+float64(p.width) <= screenHeight-10 {
				p.pos.y += p.vel.y
			}
		} else if k == ebiten.KeyRight {
			if p.pos.x+float64(p.width) < screenWidth-5 {
				p.pos.x += p.vel.x
			}
		} else if k == ebiten.KeyLeft {
			if p.pos.x > 5 {
				p.pos.x -= p.vel.y
			}
		}
	}
}

func (p *aircraft) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(p.pos.x, p.pos.y)
	screen.DrawImage(p.img, op)
}

type aircraftBullet struct {
	pos position
	img *ebiten.Image
	vel velocity
}

func NewAircraftBullet() *aircraftBullet {
	b := &aircraftBullet{}
	emptyImg := ebiten.NewImage(1, 10)
	emptyImg.Fill(color.Black)
	b.img = emptyImg
	b.vel.y = 4
	return b
}

func (bullet *aircraftBullet) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(bullet.pos.x, bullet.pos.y)
	screen.DrawImage(bullet.img, op)
}

func (bullet *aircraftBullet) ResetPosition(p *aircraft) {
	bullet.pos.y = p.pos.y - 1
	bullet.pos.x = p.pos.x + float64(p.width/2)
}

type aircraftHP struct {
	pos position
	img *ebiten.Image
}

func NewAircraftHP() *aircraftHP {
	p := &aircraftHP{}
	p.init()
	return p
}

func (p *aircraftHP) init() {
	p.pos.x = screenWidth - 400
	p.pos.y = screenHeight - 50
	p.img = ebiten.NewImage(300, 20)
	p.img.Fill(color.RGBA{
		R: 255,
		G: 150,
		B: 150,
		A: 250,
	})
}

func (p *aircraftHP) Draw(screen *ebiten.Image, aircraft *aircraft) {
	scaleAircraftHP := 1 - float64(100-(100-(float64(aircraft.hit)/100)))
	if aircraft.hit <= 0 {
		scaleAircraftHP = 1
	}
	if aircraft.hit >= 100 {
		scaleAircraftHP = 0
	}
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(scaleAircraftHP, 1)
	op.GeoM.Translate(p.pos.x, p.pos.y)
	screen.DrawImage(p.img, op)
}
