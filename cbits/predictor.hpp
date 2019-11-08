#ifndef __PREDICTOR_HPP__
#define __PREDICTOR_HPP__

#ifdef __cplusplus
extern "C" {
#endif  // __cplusplus

#include <stddef.h>
#include <stdbool.h>

typedef void *PredictorContext;

PredictorContext NewSnpe(char *model_file, int batch, int mode, bool verbose, bool profile);

void SetModeSnpe(int mode);

void InitSnpe();

void PredictSnpe(PredictorContext pred, int* inputData_quantize, float* inputData_float, bool quantize);

float* GetPredictionsSnpe(PredictorContext pred);

void DeleteSnpe(PredictorContext pred);

int GetWidthSnpe(PredictorContext pred);

int GetHeightSnpe(PredictorContext pred);

int GetChannelsSnpe(PredictorContext pred);

int GetPredLenSnpe(PredictorContext pred);

void SetInputSnpe_float(float* out, float* in, int image_height, int image_width, int image_channels, int model_height, int model_width, int model_channels);

void SetInputSnpe_quantize_8_unsigned(uint8_t* out, int* in, int image_height, int image_width, int image_channels, int model_height, int model_width, int model_channels);

void SetInputSnpe_quantize_8_signed(int8_t* out, int* in, int image_height, int image_width, int image_channels, int model_height, int model_width, int model_channels);

#ifdef __cplusplus
}
#endif  // __cplusplus

#endif  // __PREDICTOR_HPP__
