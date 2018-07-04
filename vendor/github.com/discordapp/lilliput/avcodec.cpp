#include "avcodec.hpp"

#ifdef __cplusplus
extern "C" {
#endif

#include <libavcodec/avcodec.h>
#include <libavformat/avformat.h>
#include <libswscale/swscale.h>
#include <libavutil/imgutils.h>

#ifdef __cplusplus
}
#endif

extern AVInputFormat ff_mov_demuxer;
extern AVInputFormat ff_matroska_demuxer;
extern AVInputFormat ff_mp3_demuxer;
extern AVInputFormat ff_flac_demuxer;
extern AVInputFormat ff_wav_demuxer;
extern AVInputFormat ff_aac_demuxer;
extern AVInputFormat ff_ogg_demuxer;
extern AVCodec ff_h264_decoder;
extern AVCodec ff_mpeg4_decoder;
extern AVCodec ff_vp9_decoder;
extern AVCodec ff_vp8_decoder;
extern AVCodec ff_mp3_decoder;
extern AVCodec ff_flac_decoder;
extern AVCodec ff_aac_decoder;
extern AVCodec ff_vorbis_decoder;

void avcodec_init() {
    av_register_input_format(&ff_mov_demuxer);
    av_register_input_format(&ff_matroska_demuxer);
    av_register_input_format(&ff_mp3_demuxer);
    av_register_input_format(&ff_flac_demuxer);
    av_register_input_format(&ff_wav_demuxer);
    av_register_input_format(&ff_aac_demuxer);
    av_register_input_format(&ff_ogg_demuxer);

    avcodec_register(&ff_h264_decoder);
    avcodec_register(&ff_mpeg4_decoder);
    avcodec_register(&ff_vp9_decoder);
    avcodec_register(&ff_vp8_decoder);
    avcodec_register(&ff_mp3_decoder);
    avcodec_register(&ff_flac_decoder);
    avcodec_register(&ff_aac_decoder);
    avcodec_register(&ff_vorbis_decoder);

    av_log_set_level(AV_LOG_ERROR);
}

struct avcodec_decoder_struct {
    const cv::Mat *mat;
    ptrdiff_t read_index;
    AVFormatContext *container;
    AVCodecContext *codec;
    AVIOContext *avio;
    int video_stream_index;
};

static int avcodec_decoder_read_callback(void *d_void, uint8_t *buf, int buf_size) {
    avcodec_decoder d = static_cast<avcodec_decoder>(d_void);
    size_t buf_len = d->mat->total() - d->read_index;
    size_t read_len = (buf_len > buf_size) ? buf_size : buf_len;
    memmove(buf, d->mat->data + d->read_index, read_len);
    d->read_index += read_len;
    return read_len;
}

static int64_t avcodec_decoder_seek_callback(void *d_void, int64_t offset, int whence) {
    avcodec_decoder d = static_cast<avcodec_decoder>(d_void);
    uint8_t *to;
    switch (whence) {
    case SEEK_SET:
        to = d->mat->data + offset;
        break;
    case SEEK_CUR:
        to = d->mat->data + d->read_index + offset;
        break;
    case SEEK_END:
        to = d->mat->data + d->mat->total() + offset;
        break;
    case AVSEEK_SIZE:
        return d->mat->total();
    default:
        return -1;
    }
    if (to < d->mat->data) {
        return -1;
    }
    if (to >= (d->mat->data + d->mat->total())) {
        return -1;
    }
    d->read_index = (to - d->mat->data);
    return 0;
}

static bool avcodec_decoder_is_audio(const avcodec_decoder d) {
    if (!d->container) {
        return false;
    }
    if (d->container->iformat == &ff_mp3_demuxer) {
        return true;
    }
    if (d->container->iformat == &ff_flac_demuxer) {
        return true;
    }
    if (d->container->iformat == &ff_wav_demuxer) {
        return true;
    }
    if (d->container->iformat == &ff_aac_demuxer) {
        return true;
    }
    if (d->container->iformat == &ff_ogg_demuxer) {
        return true;
    }
    return false;
}

avcodec_decoder avcodec_decoder_create(const opencv_mat buf) {
    avcodec_decoder d = new struct avcodec_decoder_struct();
    memset(d, 0, sizeof(struct avcodec_decoder_struct));
    d->mat = static_cast<const cv::Mat *>(buf);

    d->container = avformat_alloc_context();
    if (!d->container) {
        avcodec_decoder_release(d);
        return NULL;
    }

    d->avio = avio_alloc_context(NULL, 0, 0, d, avcodec_decoder_read_callback,
                                 NULL, avcodec_decoder_seek_callback);
    if (!d->avio) {
        avcodec_decoder_release(d);
        return NULL;
    }
    d->container->pb = d->avio;

    int res = avformat_open_input(&d->container, NULL, NULL, NULL);
    if (res < 0) {
        avformat_free_context(d->container);
        d->container = NULL;
        avcodec_decoder_release(d);
        return NULL;
    }

    res = avformat_find_stream_info(d->container, NULL);
    if (res < 0) {
        avcodec_decoder_release(d);
        return NULL;
    }

    if (avcodec_decoder_is_audio(d)) {
        // in this case, quit out fast since we won't be decoding this anyway
        // (audio is metadata-only)
        return d;
    }

    AVCodecParameters *codec_params = NULL;
    for (int i = 0; i < d->container->nb_streams; i++) {
        if (d->container->streams[i]->codecpar->codec_type == AVMEDIA_TYPE_VIDEO) {
            codec_params = d->container->streams[i]->codecpar;
            d->video_stream_index = i;
            break;
        }
    }
    if (!codec_params) {
        avcodec_decoder_release(d);
        return NULL;
    }

    AVCodec *codec = avcodec_find_decoder(codec_params->codec_id);
    if (!codec) {
        avcodec_decoder_release(d);
        return NULL;
    }

    d->codec = avcodec_alloc_context3(codec);

    res = avcodec_parameters_to_context(d->codec, codec_params);
    if (res < 0) {
        avcodec_decoder_release(d);
        return NULL;
    }

    res = avcodec_open2(d->codec, codec, NULL);
    if (res < 0) {
        avcodec_decoder_release(d);
        return NULL;
    }

    return d;
}

int avcodec_decoder_get_width(const avcodec_decoder d) {
    if (d->codec) {
        return d->codec->width;
    }
    return 0;
}

int avcodec_decoder_get_height(const avcodec_decoder d) {
    if (d->codec) {
        return d->codec->height;
    }
    return 0;
}

int avcodec_decoder_get_orientation(const avcodec_decoder d) {
    if (!d->container) {
        return CV_IMAGE_ORIENTATION_TL;
    }
    if (!d->codec) {
        return CV_IMAGE_ORIENTATION_TL;
    }
    CVImageOrientation orientation = CV_IMAGE_ORIENTATION_TL;
    AVDictionaryEntry *tag = av_dict_get(d->container->streams[d->video_stream_index]->metadata, "rotate", NULL, 0);
    if (!tag) {
        return orientation;
    }
    int rotation = atoi(tag->value);
    switch (rotation) {
    case 90:
        orientation = CV_IMAGE_ORIENTATION_RT;
        break;
    case 180:
        orientation = CV_IMAGE_ORIENTATION_BR;
        break;
    case 270:
        orientation = CV_IMAGE_ORIENTATION_LB;
        break;
    }
    return orientation;
}

float avcodec_decoder_get_duration(const avcodec_decoder d) {
    if (d->container) {
        return d->container->duration / (float)(AV_TIME_BASE);
    }
    return 0;
}

const char *avcodec_decoder_get_description(const avcodec_decoder d) {
    if (d->container) {
        if (d->container->iformat == &ff_mov_demuxer) {
            return "MOV";
        }
        if (d->container->iformat == &ff_matroska_demuxer) {
            return "WEBM";
        }
        if (d->container->iformat == &ff_mp3_demuxer) {
            return "MP3";
        }
        if (d->container->iformat == &ff_flac_demuxer) {
            return "FLAC";
        }
        if (d->container->iformat == &ff_wav_demuxer) {
            return "WAV";
        }
        if (d->container->iformat == &ff_aac_demuxer) {
            return "AAC";
        }
        if (d->container->iformat == &ff_ogg_demuxer) {
            return "OGG";
        }
    }
    return "";
}

static int avcodec_decoder_copy_frame(const avcodec_decoder d, opencv_mat mat, AVFrame *frame) {
    auto cvMat = static_cast<cv::Mat *>(mat);

    int res = avcodec_receive_frame(d->codec, frame);
    if (res >= 0 && frame->width == cvMat->cols && frame->height == cvMat->rows) {
        int stepSize = 4 * frame->width;
        if (frame->width % 32 != 0) {
            int width = frame->width + 32 - (frame->width % 32);
            stepSize = 4 * width;
        }
        if (!opencv_mat_set_row_stride(mat, stepSize)) {
            return -1;
        }

        struct SwsContext *sws = sws_getContext(frame->width, frame->height, (AVPixelFormat)(frame->format),
                                                frame->width, frame->height, AV_PIX_FMT_BGRA,
                                                0, NULL, NULL, NULL);
        int linesizes[] = {stepSize, 0, 0, 0};
        uint8_t *data_ptrs[] = {cvMat->data, NULL, NULL, NULL};
        sws_scale(sws, frame->data, frame->linesize, 0, frame->height, data_ptrs, linesizes);
        sws_freeContext(sws);
    }

    return res;
}

static int avcodec_decoder_decode_packet(const avcodec_decoder d, opencv_mat mat, AVPacket *packet) {
    int res = avcodec_send_packet(d->codec, packet);
    if (res < 0) {
        return res;
    }

    AVFrame *frame = av_frame_alloc();
    if (!frame) {
        return -1;
    }

    res = avcodec_decoder_copy_frame(d, mat, frame);
    av_frame_free(&frame);

    return res;
}

bool avcodec_decoder_decode(const avcodec_decoder d, opencv_mat mat) {
    if (!d) {
        return false;
    }
    if (!d->container) {
        return false;
    }
    if (!d->codec) {
        return false;
    }
    AVPacket packet;
    bool done = false;
    bool success = false;
    while (!done) {
        int res = av_read_frame(d->container, &packet);
        if (res < 0) {
            return false;
        }
        if (packet.stream_index == d->video_stream_index) {
            res = avcodec_decoder_decode_packet(d, mat, &packet);
            if (res >= 0) {
                success = true;
                done = true;
            } else if (res != AVERROR(EAGAIN)) {
                done = true;
            }
        }
        av_packet_unref(&packet);
    }
    return success;
}

void avcodec_decoder_release(avcodec_decoder d) {
    if (d->codec) {
        avcodec_free_context(&d->codec);
    }

    if (d->container) {
        avformat_close_input(&d->container);
    }

    if (d->avio) {
        avio_flush(d->avio);
        av_free(d->avio->buffer);
        av_free(d->avio);
    }

    delete d;
}
