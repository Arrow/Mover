package main

import (
	"math/rand"
	"fmt"
	"image/color"
	"github.com/Arrow/display"
	"time"
	"runtime/pprof"
	"flag"
	"os"
	"github.com/golang/glog"
)

const (
	width, height    = 700, 500
	heading          = 25
	border           = 5
	numParticles     = 5
)

func init() {
	t := time.Now()
	rand.Seed(t.Unix())
}

type Mover interface {
	Move()
}

type CVelMover struct {
	vx    float64
	vy    float64
	x     float64
	y     float64
	p     *display.Particle
}

// NewCVelMover generates a new CVelMover with speed vx, vy
func NewCVelMover(d *display.Display) (cvm *CVelMover) {
	cvm    = new(CVelMover)
	cvm.vx = rand.Float64() * 10
	cvm.vy = rand.Float64() * 10
	cvm.x  = rand.Float64() * width
	cvm.y  = rand.Float64() * height
	cvm.p  = d.NewParticle(cvm.x, cvm.y, 2, color.Black)
	return cvm
}

func (cvm *CVelMover) Move() {
	cvm.x += cvm.vx
	if cvm.x < 0 {
		cvm.x = 0
		cvm.vx *= -1.0
	} else if cvm.x > width {
		cvm.x = width
		cvm.vx *= -1.0
	}
	
	cvm.y += cvm.vy 
	if cvm.y < border {
		cvm.y = border
		cvm.vy *= -1.0
	} else if cvm.y > height {
		cvm.y = height
		cvm.vy *= -1.0
	}
	
	cvm.p.Move(cvm.x, cvm.y)
	return
}

var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")

func main() {
	flag.Parse()
	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			glog.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}
	defer glog.Flush()
	
	d, err := display.NewDisplay(width, height, border, heading, "Mover")
	if err != nil {
		glog.Fatal(err)
	}
	c := make([]*CVelMover, numParticles)
	for i := 0; i < numParticles; i++ {
		c[i] = NewCVelMover(d)
	}

	timer := time.Tick(time.Second)
	timerEnd := time.Tick(time.Minute)
	ctr := 0
	fps := 0
	for {
		for _, ci := range c {
			ci.Move()
		}
		d.Frame()
		ctr++
		select {
		case <-timer:
			fps = ctr
			d.SetHeadingText(fmt.Sprint("FPS: ", fps))
			ctr = 0
		case <-timerEnd:
			return
		default:
		}
	}
}
