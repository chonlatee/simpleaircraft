package main

import (
	"image"
	"image/color"
	"log"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const bossFrameWidth = 100
const bossFrameHeight = 150
const bossFrameOX = 0
const bossFrameOY = 150

type boss struct {
	pos   position
	img   *ebiten.Image
	count int
	vel   velocity
	hp    float64
	hit   int
	bulletNum int
	bullets []*bossBullet
}

func NewBoss() *boss {
	b := &boss{}
	b.init()
	return b
}

func (b *boss) init() {
	fr, err := ebitenutil.OpenFile("boss.png")
	if err != nil {
		log.Fatalln("can't open file: ", err)
	}
	img, _, err := image.Decode(fr)
	imgOrig := ebiten.NewImageFromImage(img)
	b.img = imgOrig
	b.pos.x = screenWidth / 2
	b.pos.y = 20
	b.vel.x = 2
	b.vel.y = 2
	b.bulletNum = 10
	b.bullets = make([]*bossBullet, b.bulletNum)

	for i := 0; i < b.bulletNum; i++ {
		b.bullets[i] = NewBossBullet()
		b.bullets[i].pos.y = b.pos.y + bossFrameHeight + 20
		b.bullets[i].pos.x = b.pos.x + bossFrameWidth/2
		b.bullets[i].vel.x = float64(rand.Intn(4 - 1) + 1)
		b.bullets[i].vel.y = float64(rand.Intn(6 - 3) + 3)
	}
}

func (b *boss) Update() {
	for i, _ := range b.bullets {
		bullet := b.bullets[i]
		if bullet.pos.x >= screenWidth || bullet.pos.x <= 0 {
			bullet.vel.x = bullet.vel.x * -1
		}

		if bullet.pos.y <= 10 {
			bullet.vel.y = bullet.vel.y * -1
		}

		if bullet.pos.y >= screenHeight {
			bullet.pos.y = b.pos.y + bossFrameHeight + 20
			bullet.pos.x = b.pos.x + bossFrameWidth/2
		}

		bullet.pos.x += bullet.vel.x
		bullet.pos.y += bullet.vel.y
	}

	if b.pos.x+100 >= screenWidth || b.pos.x <= 0 {
		b.vel.x = b.vel.x * -1
	}

	if b.pos.y >= 200 || b.pos.y <= 10 {
		b.vel.y = b.vel.y * -1
	}

	b.pos.x += b.vel.x
	b.pos.y += b.vel.y
	b.count++
}

func (b *boss) IsHit(bullet *aircraftBullet) bool {
	return bullet.pos.y <= b.pos.y+float64(bossFrameHeight) &&
		(bullet.pos.x <= b.pos.x+float64(bossFrameWidth) && bullet.pos.x >= b.pos.x)
}

func (b *boss) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	i := (b.count / 20) % 5
	sx, _ := bossFrameOX+i*bossFrameWidth, bossFrameOY
	op.GeoM.Translate(b.pos.x, b.pos.y)
	screen.DrawImage(b.img.SubImage(image.Rect(sx, 0, sx+bossFrameWidth, bossFrameHeight)).(*ebiten.Image), op)
}

type bossHP struct {
	pos position
	img *ebiten.Image
}

func NewBossHP() *bossHP {
	b := &bossHP{}
	b.init()
	return b
}

func (b *bossHP) init() {
	b.pos.x = screenWidth - 450
	b.pos.y = 20
	b.img = ebiten.NewImage(400, 20)
	b.img.Fill(color.RGBA{
		R: 255,
		G: 70,
		B: 70,
		A: 240,
	})
}

func (b *bossHP) Draw(screen *ebiten.Image, bs *boss) {
	scaleHP := 1 - float64(100-(100-(float64(bs.hit)/100)))
	if bs.hit <= 0 {
		scaleHP = 1
	}
	if bs.hit >= 100 {
		scaleHP = 0
	}

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(scaleHP, 1)
	op.GeoM.Translate(b.pos.x, b.pos.y)
	screen.DrawImage(b.img, op)
}

type bossBullet struct {
	pos position
	img *ebiten.Image
	vel velocity
}

func NewBossBullet() *bossBullet {
	b := &bossBullet{}
	b.init()
	return b
}

func (b *bossBullet) init() {
	img := ebiten.NewImage(5, 5)
	b.img = img
	b.img.Fill(color.RGBA{
		R: 255,
		G: 150,
		B: 150,
		A: 255,
	})
	b.vel.x = 3
	b.vel.y = 3
}

func (b *bossBullet) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(b.pos.x, b.pos.y)
	screen.DrawImage(b.img, op)
}

func (b *bossBullet) UpdatePos(x, y int) {

}

func (b *bossBullet) ResetPosition(bs *boss) {
	b.pos.y = bs.pos.y + bossFrameHeight + 20
	b.pos.x = bs.pos.x + float64(bossFrameWidth/2)
}
