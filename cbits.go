package snpe

// #include <stdio.h>
// #include <stdlib.h>
// #include "cbits/predictor.hpp"
import "C"
import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"unsafe"

	"github.com/Unknwon/com"
	"github.com/pkg/errors"
	"github.com/rai-project/dlframework"
	"github.com/rai-project/dlframework/framework/feature"
)

// Hardware Modes
const (
	CPU_1_thread = 1
	CPU_2_thread = 2
	CPU_3_thread = 3
	CPU_4_thread = 4
	CPU_5_thread = 5
	CPU_6_thread = 6
	GPU          = 7
	NNAPI        = 8
)

// Predictor Structure definition
type PredictorData struct {
	ctx   C.PredictorContext
	mode  int
	batch int
}

// Make access to mode and batch public
func (pd *PredictorData) Inc() {
	pd.mode++
	pd.batch++
}

// Create new Predictor Structure
func NewPredictorData() *PredictorData {
	return &PredictorData{}
}

// Create new predictor
func New(model string, mode, batch int, verbose bool, profile bool) (*PredictorData, error) {

	modelFile := model
	if !com.IsFile(modelFile) {
		return nil, errors.Errorf("file %s not found", modelFile)
	}

	return &PredictorData{
		ctx: C.NewSnpe(
			C.CString(modelFile),
			C.int(batch),
			C.int(mode),
			C.bool(verbose),
			C.bool(profile),
		),
		mode:  mode,
		batch: batch,
	}, nil
}

// Initialize TFLite
func init() {
	C.InitSnpe()
}

// Run inference
func Predict(p *PredictorData, data []byte, quantize bool) error {

	if len(data) == 0 {
		return fmt.Errorf("image data is empty")
	}

	ptr_quantize := (*C.int)(unsafe.Pointer(&data[0]))
	ptr_float := (*C.float)(unsafe.Pointer(&data[0]))
	if quantize == true {
		C.PredictSnpe(p.ctx, ptr_quantize, ptr_float, true)
	} else {
		C.PredictSnpe(p.ctx, ptr_quantize, ptr_float, false)
	}

	return nil
}

// Return Top-5 predicted label
func ReadPredictionOutput(p *PredictorData, labelFile string) (string, error) {

	batchSize := p.batch
	if batchSize == 0 {
		return "", errors.New("null batch")
	}

	predLen := int(C.GetPredLenSnpe(p.ctx))
	if predLen == 0 {
		return "", errors.New("null predLen")
	}

	length := batchSize * predLen
	if p.ctx == nil {
		return "", errors.New("empty predictor context")
	}

	cPredictions := C.GetPredictionsSnpe(p.ctx)
	if cPredictions == nil {
		return "", errors.New("empty predictions")
	}

	slice := (*[1 << 15]float32)(unsafe.Pointer(cPredictions))[:length:length]

	var labels []string
	f, err := os.Open(labelFile)
	if err != nil {
		panic(err)
	}

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		labels = append(labels, line)
	}

	features := make([]dlframework.Features, batchSize)
	featuresLen := len(slice) / batchSize

	for ii := 0; ii < batchSize; ii++ {
		rprobs := make([]*dlframework.Feature, featuresLen)
		for jj := 0; jj < featuresLen; jj++ {
			rprobs[jj] = feature.New(
				feature.ClassificationIndex(int32(jj)),
				feature.ClassificationLabel(labels[jj]),
				feature.Probability(slice[ii*featuresLen+jj]),
			)
		}
		sort.Sort(dlframework.Features(rprobs))
		features[ii] = rprobs
	}

	top1 := features[0][0]
	top2 := features[0][1]
	top3 := features[0][2]
	top4 := features[0][3]
	top5 := features[0][4]

	top_concatenated := top1.GetClassification().GetLabel() + "|" + top2.GetClassification().GetLabel() + "|" + top3.GetClassification().GetLabel() + "|" + top4.GetClassification().GetLabel() + "|" + top5.GetClassification().GetLabel()

	return top_concatenated, nil

}

// Delete the predictor
func Close(p *PredictorData) {
	C.DeleteSnpe(p.ctx)
}
