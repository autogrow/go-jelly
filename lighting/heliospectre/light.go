package heliospectre

type Light struct{}

func NewLight(ip string) *Light {
	return new(Light)
}
