package main

// https://raytracing.github.io/
// https://raytracing.github.io/books/RayTracingInOneWeekend.html

import (
	"fmt"
	"go-raytracer/hittable"
	"go-raytracer/util"
	"go-raytracer/vec"
	"math"
	"os"
)


func check(e error) {
	if e != nil {
		panic(e)
	}
}

func rayColor(r vec.Ray3D, world *hittable.HitObjectList, depth int) vec.Vector3D {
	rec := hittable.HitRecord{}

	// If we've exceeded the ray bounce limit, no more light is gathered.
	if depth <= 0 {
		return vec.Zero()
	}

	if world.Hit(r, 0.001, float64(math.Inf(1)), &rec) {
		scattered := vec.Ray3D{
			Origin:    vec.Vector3D{},
			Direction: vec.Vector3D{},
		}
		attenuation := vec.Zero()
		if rec.Mat.Scatter(r, rec, &attenuation, &scattered) {
			return attenuation.Multiply(rayColor(scattered, world, depth-1))
		}
		return vec.Zero()

		//target := rec.P.Add(vec.RandomInHemisphere(rec.Normal))
		//newRay := vec.NewRay3D(rec.P, target.Subtract(rec.P))
		//return rayColor(newRay, world, depth-1).MultiplyScalar(0.5)

		//return rec.Normal.Add(vec.Unit()).MultiplyScalar(0.5)
	}
	unitDirection := r.Direction.Unit()
	t := 0.5 * (unitDirection.Y + 1)
	return vec.Unit().MultiplyScalar(1.0-t).Add(vec.NewVector3D(0.5, 0.7, 1.0).MultiplyScalar(t))
}

func randomScene() hittable.HitObjectList {
	world := hittable.NewHitObjectList()

	world.Add(hittable.NewSphere(vec.NewVector3D(0,-1000,0), 1000, hittable.Lambertian{vec.NewVector3D(0.1, 0.4, 0.1)}))

	for a := -11.0; a < 11; a++ {
		for b := -11.0; b < 11; b++ {
			chooseMat := util.Random()
			chooseMat = 1
			center := vec.NewVector3D(a + 0.9*util.Random(), 0.2, b + 0.9 * util.Random())
			if (center.Subtract(vec.NewVector3D(4, 0.2, 0))).Length() > 0.9 {
				if chooseMat < 0.8 {
					// diffuse
					albedo := vec.Random().Multiply(vec.Random())
					world.Add(hittable.NewSphere(center, 0.2, hittable.Lambertian{Albedo:albedo}))
				} else if chooseMat < 0.95 {
					// metal
					albedo := vec.RandomBetween(.5, 1)
					fuzz := util.RandomBetween(0, .5)
					world.Add(hittable.NewSphere(center, 0.2, hittable.NewMetal(albedo, fuzz)))
				} else {
					// glass
					world.Add(hittable.NewSphere(center, 0.2, hittable.NewDielectric(util.RandomBetween(1.2, 1.9))))
				}
			}
		}
	}

	world.Add(hittable.NewSphere(vec.NewVector3D(4,1,0), 1.0, hittable.NewDielectric(1.1)))
	world.Add(hittable.NewSphere(vec.NewVector3D(0,1,0), 1.0, hittable.Lambertian{vec.NewVector3D(0.3, 1, 0.3)}))
	world.Add(hittable.NewSphere(vec.NewVector3D(-4,1,0), 1.0, hittable.NewMetal(vec.NewVector3D(0.4, 0.9, 0.4), 0.0)))

	//world.Add(hittable.NewSphere(vec.NewVector3D(0,0,-1), 0.5, hittable.Lambertian{vec.NewVector3D(0.1, 0.2, 0.5)}))
	//world.Add(hittable.NewSphere(vec.NewVector3D(0,-100.5,-1), 100, hittable.Lambertian{vec.NewVector3D(0.8, 0.8, 0.0)}))
	//world.Add(hittable.NewSphere(vec.NewVector3D(1,0,-1), 0.5, hittable.NewMetal(vec.NewVector3D(0.8, 0.6, 0.2), 0.3)))
	//world.Add(hittable.NewSphere(vec.NewVector3D(-1,0,-1), 0.5, hittable.NewDielectric(1.5)))
	//world.Add(hittable.NewSphere(vec.NewVector3D(-1,0,-1), -0.45, hittable.NewDielectric(1.5)))

	return *world
}

func main() {
	const imageWidth = 400
	const imageHeight = 200
	const samplesPerPixel = 100
	const maxDepth = 50

	f, err := os.Create(fmt.Sprintf("./image%vx%v-spp%v-depth%v.ppm", imageWidth, imageHeight, samplesPerPixel, maxDepth))
	check(err)

	defer f.Close()

	f.WriteString(fmt.Sprintf("P3\n%v %v\n255\n", imageWidth, imageHeight))

	world := randomScene()

	aspectRatio := float64(imageWidth / imageHeight)
	lookFrom := vec.NewVector3D(13, 2, 3)
	lookAt := vec.NewVector3D(0,0,0)
	vup := vec.NewVector3D(0, 1, 0)
	//aperture := 2.0
	//distToFocus := (lookFrom.Subtract(lookAt)).Length()
	distToFocus := 10.0
	aperture := 0.1
	camera := vec.NewCamera(lookFrom, lookAt, vup, 30, aspectRatio, aperture, distToFocus)

	for j := imageHeight; j >= 0; j-- {
		fmt.Println("\rScanlines remaining: ", j, " ")
		for i := 0; i < imageWidth; i++ {
			color := vec.Zero()
			for s := 0; s < samplesPerPixel; s++ {
				u := (float64(i) + util.Random()) / imageWidth
				v := (float64(j) + util.Random()) / imageHeight
				r := camera.GetRay(u, v)
				color = color.Add(rayColor(r, &world, maxDepth))
			}

			f.WriteString(color.ColorString(samplesPerPixel))
		}
	}

	fmt.Println("Done.")

	f.Sync()
}
