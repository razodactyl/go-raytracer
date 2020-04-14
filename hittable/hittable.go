package hittable

import (
	"github.com/razodactyl/go-raytracer/vec"
)

type HitRecord struct {
	P         vec.Vector3D
	Normal    vec.Vector3D
	Mat       Material
	t         float64
	frontFace bool
}

func (h *HitRecord) setFaceNormal(r vec.Ray3D, outwardNormal vec.Vector3D) {
	h.frontFace = r.Direction.Dot(outwardNormal) < 0
	if h.frontFace {
		h.Normal = outwardNormal
	} else {
		h.Normal = outwardNormal.MultiplyScalar(-1)
	}
}

type Hittable interface {
	Hit(r vec.Ray3D, tMin float64, tMax float64, rec *HitRecord) bool
}

type HitObjectList struct {
	Objects []Hittable
}

func NewHitObjectList() *HitObjectList {
	return &HitObjectList{}
}

func (l HitObjectList) Clear() {
	l.Objects = []Hittable{}
}

func (l *HitObjectList) Add(h Hittable) {
	l.Objects = append(l.Objects, h)
}

func (l HitObjectList) Hit(r vec.Ray3D, tMin float64, tMax float64, rec *HitRecord) bool {
	tempRec := HitRecord{
		P:         vec.Zero(),
		Normal:    vec.Zero(),
		t:         0,
		frontFace: false,
	}
	hitAnything := false
	closestSoFar := tMax

	for _, o := range l.Objects {
		if o.Hit(r, tMin, closestSoFar, &tempRec) {
			hitAnything = true
			closestSoFar = tempRec.t

			rec.P = tempRec.P
			rec.Normal = tempRec.Normal
			rec.t = tempRec.t
			rec.frontFace = tempRec.frontFace
			rec.Mat = tempRec.Mat

			// had to set values manually (above) - can't set to other object like below?
			//rec = &tempRec
		}
	}

	return hitAnything
}
