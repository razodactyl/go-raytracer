package hittable

import (
	"github.com/razodactyl/go-raytracer/util"
	"github.com/razodactyl/go-raytracer/vec"
	"math"
)

type Material interface {
	Scatter(rIn vec.Ray3D, rec HitRecord, attenuation *vec.Vector3D, scattered *vec.Ray3D) bool
}

// Lambertian material

type Lambertian struct {
	Albedo vec.Vector3D
}

func (l Lambertian) Scatter(rIn vec.Ray3D, rec HitRecord, attenuation *vec.Vector3D, scattered *vec.Ray3D) bool {
	scatterDirection := rec.Normal.Add(vec.RandomUnitVector())
	*scattered = vec.NewRay3D(rec.P, scatterDirection)
	*attenuation = l.Albedo
	return true
}

// Metal material

type Metal struct {
	Albedo vec.Vector3D
	Fuzz float64
}

func NewMetal(albedo vec.Vector3D, fuzz float64) *Metal {
	return &Metal{Albedo: albedo, Fuzz: util.Clamp(fuzz, 0, 1)}
}

func (m Metal) Scatter(rIn vec.Ray3D, rec HitRecord, attenuation *vec.Vector3D, scattered *vec.Ray3D) bool {
	reflected := vec.Reflect(rIn.Direction.Unit(), rec.Normal)
	*scattered = vec.NewRay3D(rec.P, reflected.Add(vec.RandomInUnitSphere().MultiplyScalar(m.Fuzz)))
	*attenuation = m.Albedo
	return scattered.Direction.Dot(rec.Normal) > 0
}

// Dielectric material

type Dielectric struct {
	RefractionIndex float64
}

func NewDielectric(refractionIndex float64) *Dielectric {
	return &Dielectric{RefractionIndex: refractionIndex}
}

func Schlick(cosine float64, refractionIndex float64) float64 {
	r0 := (1-refractionIndex) / (1+refractionIndex)
	r0 = r0*r0
	return r0 + (1-r0)*math.Pow((1-cosine), 5)
}

func (d Dielectric) Scatter(rIn vec.Ray3D, rec HitRecord, attenuation *vec.Vector3D, scattered *vec.Ray3D) bool {
	*attenuation = vec.Unit()
	etaiOverEtat := 0.0
	if rec.frontFace {
		etaiOverEtat = 1.0 / d.RefractionIndex
	} else {
		etaiOverEtat = d.RefractionIndex
	}

	unitDirection := rIn.Direction.Unit()

	cosTheta := util.FFMin(unitDirection.MultiplyScalar(-1).Dot(rec.Normal), 1.0)
	sinTheta := math.Sqrt(1.0 - cosTheta*cosTheta)
	if etaiOverEtat * sinTheta > 1.0 {
		reflected := vec.Reflect(unitDirection, rec.Normal)
		*scattered = vec.NewRay3D(rec.P, reflected)
		return true
	}
	reflectProb := Schlick(cosTheta, etaiOverEtat)
	if util.Random() < reflectProb {
		reflected := vec.Reflect(unitDirection, rec.Normal)
		*scattered = vec.NewRay3D(rec.P, reflected)
		return true
	}

	refracted := vec.Refract(unitDirection, rec.Normal, etaiOverEtat)
	*scattered = vec.NewRay3D(rec.P, refracted)
	return true
}

