package feature

import "github.com/rai-project/dlframework"

func ImageType() Option {
	return Type(dlframework.FeatureType_IMAGE)
}

func Image(e *dlframework.Image) Option {
	return func(o *dlframework.Feature) {
		ImageType()(o)
		o.Feature = &dlframework.Feature_Image{
			Image: e,
		}
	}
}

func ensureImage(o *dlframework.Feature) *dlframework.Image {
	if o.Type != dlframework.FeatureType_IMAGE && !isUnknownType(o) {
		panic("unexpected feature type")
	}
	if o.Feature == nil {
		o.Feature = &dlframework.Feature_Image{}
	}
	img, ok := o.Feature.(*dlframework.Feature_Image)
	if !ok {
		panic("expecting an image feature")
	}
	if img.Image == nil {
		img.Image = &dlframework.Image{}
	}
	ImageType()(o)
	return img.Image
}

func ImageID(id string) Option {
	return func(o *dlframework.Feature) {
		img := ensureImage(o)
		img.ID = id
	}
}

func ImageData(data []byte) Option {
	return func(o *dlframework.Feature) {
		img := ensureImage(o)
		img.Data = data
	}
}
