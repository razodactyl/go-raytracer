package hittable

import (
	"math"
	"raytracing-tutorial/vec"
)

type Sphere struct {
	center vec.Vector3D
	radius float64
	mat    Material
}

func NewSphere(center vec.Vector3D, radius float64, mat Material) *Sphere {
	return &Sphere{center: center, radius: radius, mat: mat}
}

func (s Sphere) Hit(r vec.Ray3D, tMin float64, tMax float64, rec *HitRecord) bool {
	oc := r.Origin.Subtract(s.center)
	a := r.Direction.LengthSquared()
	halfB := oc.Dot(r.Direction)
	c := oc.LengthSquared() - s.radius*s.radius
	discriminant := halfB*halfB - a*c

	if discriminant > 0 {
		root := float64(math.Sqrt(float64(discriminant)))
		temp := (-halfB - root) / a
		if temp < tMax && temp > tMin {
			rec.t = temp
			rec.P = r.At(rec.t)
			outwardNormal := (rec.P.Subtract(s.center)).DivideScalar(s.radius)
			rec.setFaceNormal(r, outwardNormal)
			rec.Mat = s.mat
			return true
		}
		temp = (-halfB + root) / a
		if temp < tMax && temp > tMin {
			rec.t = temp
			rec.P = r.At(rec.t)
			outwardNormal := (rec.P.Subtract(s.center)).DivideScalar(s.radius)
			rec.setFaceNormal(r, outwardNormal)
			rec.Mat = s.mat
			return true
		}
	}
	return false
}
