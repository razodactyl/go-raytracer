package vec

import (
	"fmt"
	"github.com/razodactyl/go-raytracer/util"
	"math"
)

type Vector3D struct {
	X float64
	Y float64
	Z float64
}

func NewVector3D(x, y, z float64) Vector3D {
	return Vector3D{x, y, z}
}

func FromScalar(v float64) Vector3D {
	return Vector3D{v,v,v}
}

func Zero() Vector3D {
	return Vector3D{0, 0, 0}
}

func Unit() Vector3D {
	return Vector3D{1, 1, 1}
}

func Random() Vector3D {
	return NewVector3D(util.Random(), util.Random(), util.Random())
}

func RandomBetween(min float64, max float64) Vector3D {
	return NewVector3D(util.RandomBetween(min, max), util.RandomBetween(min, max), util.RandomBetween(min, max))
}

func RandomInUnitSphere() Vector3D {
	for {
		p := RandomBetween(-1, 1)
		if p.LengthSquared() >= 1 { continue }
		return p
	}
}

func RandomUnitVector() Vector3D {
	a := util.RandomBetween(0, 2*math.Pi)
	z := util.RandomBetween(-1, 1)
	r := math.Sqrt(float64(1 - z*z))
	return NewVector3D(float64(r*math.Cos(float64(a))), float64(r*math.Sin(float64(a))), z)
}

func RandomInHemisphere(normal Vector3D) Vector3D {
	inUnitSphere := RandomInUnitSphere()
	if inUnitSphere.Dot(normal) > 0.0 { // In the same hemisphere as the normal
		return inUnitSphere
	} else {
		return inUnitSphere.MultiplyScalar(-1)
	}
}

func RandomInUnitDisk() Vector3D {
	for {
		p := NewVector3D(util.RandomBetween(-1,1), util.RandomBetween(-1,1), 0)
		if p.LengthSquared() >= 1 {
			continue
		}
		return p
	}
}

// Bending light.

func Reflect(v Vector3D, n Vector3D) Vector3D {
	return v.Subtract(n.MultiplyScalar(v.Dot(n)*2))
}

func Refract(uv Vector3D, n Vector3D, etaiOverEtat float64) Vector3D {
	cosTheta := uv.MultiplyScalar(-1).Dot(n)
	rOutParallel := uv.Add(n.MultiplyScalar(cosTheta)).MultiplyScalar(etaiOverEtat)
	rOutPerp := n.MultiplyScalar(-math.Sqrt(1.0 - rOutParallel.LengthSquared()))
	return rOutParallel.Add(rOutPerp)
}

// Instance functions

func (v Vector3D) Unit() Vector3D {
	return v.DivideScalar(v.Length())
}

func (v Vector3D) LengthSquared() float64 {
	return v.X * v.X + v.Y * v.Y + v.Z * v.Z
}

func (v Vector3D) Length() float64 {
	return math.Sqrt(v.LengthSquared())
}

func (v Vector3D) Copy() Vector3D {
	return Vector3D{v.X, v.Y, v.Z}
}

func (v Vector3D) Add(v2 Vector3D) Vector3D {
	return Vector3D{v.X + v2.X, v.Y + v2.Y, v.Z + v2.Z}
}

func (v Vector3D) Subtract(v2 Vector3D) Vector3D {
	return Vector3D{v.X - v2.X, v.Y - v2.Y, v.Z - v2.Z}
}

func (v Vector3D) Multiply(v2 Vector3D) Vector3D {
	return Vector3D{v.X * v2.X, v.Y * v2.Y, v.Z * v2.Z}
}

func (v Vector3D) Divide(v2 Vector3D) Vector3D {
	return Vector3D{v.X / v2.X, v.Y / v2.Y, v.Z / v2.Z}
}

func (v Vector3D) Dot(v2 Vector3D) float64 {
	return v.X * v2.X + v.Y * v2.Y + v.Z * v2.Z
}

func (v Vector3D) Cross(v2 Vector3D) Vector3D {
	return Vector3D{
		X: v.Y * v2.Z - v.Z * v2.Y,
		Y: v.Z * v2.X - v.X * v2.Z,
		Z: v.X * v2.Y - v.Y * v2.X,
	}
}

func (v Vector3D) MultiplyScalar(s float64) Vector3D {
	return Vector3D{v.X * s, v.Y * s, v.Z * s}
}

func (v Vector3D) DivideScalar(s float64) Vector3D {
	return Vector3D{v.X / s, v.Y / s, v.Z / s}
}

func (v Vector3D) ColorString(samplesPerPixel int) string {
	// Divide the color total by the number of samples and gamma-correct
	// for a gamma value of 2.0.
	scale := 1.0 / float64(samplesPerPixel)
	r := math.Sqrt(scale * v.X)
	g := math.Sqrt(scale * v.Y)
	b := math.Sqrt(scale * v.Z)

	return fmt.Sprintf("%v %v %v\n", int(256 * util.Clamp(r, 0.0, 0.999)), int(256 * util.Clamp(g, 0.0, 0.999)), int(256 * util.Clamp(b, 0.0, 0.999)))
}
