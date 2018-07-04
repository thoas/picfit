#ifndef LILLIPUT_AVCODEC_HPP
#define LILLIPUT_AVCODEC_HPP

#include "opencv.hpp"

#ifdef __cplusplus
extern "C" {
#endif

typedef struct avcodec_decoder_struct *avcodec_decoder;

void avcodec_init();

avcodec_decoder avcodec_decoder_create(const opencv_mat buf);
void avcodec_decoder_release(avcodec_decoder d);
int avcodec_decoder_get_width(const avcodec_decoder d);
int avcodec_decoder_get_height(const avcodec_decoder d);
int avcodec_decoder_get_orientation(const avcodec_decoder d);
float avcodec_decoder_get_duration(const avcodec_decoder d);
bool avcodec_decoder_decode(const avcodec_decoder d, opencv_mat mat);
const char *avcodec_decoder_get_description(const avcodec_decoder d);

#ifdef __cplusplus
}
#endif

#endif
