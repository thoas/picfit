#include "opencv.hpp"
#include <stdbool.h>
#include <opencv2/highgui.hpp>
#include <opencv2/imgproc.hpp>

opencv_mat opencv_mat_create(int width, int height, int type) {
    return new cv::Mat(height, width, type);
}

opencv_mat opencv_mat_create_from_data(int width, int height, int type, void *data, size_t data_len) {
    size_t total_size = width * height * CV_ELEM_SIZE(type);
    if (total_size > data_len) {
        return NULL;
    }
    auto mat = new cv::Mat(height, width, type, data);
    mat->datalimit = (uint8_t*)data + data_len;
    return mat;
}

opencv_mat opencv_mat_create_empty_from_data(int length, void *data) {
    // this is slightly sketchy - what we're going to do is build a 1x0 matrix
    // and then set its data* properties to reflect the capacity (given by length arg here)
    // this tells opencv internally that the Mat can store more but has nothing in it
    // this is directly analogous to Go's len and cap
    auto mat = new cv::Mat(0, 1, CV_8U, data);

    mat->datalimit = mat->data + length;

    return mat;
}

bool opencv_mat_set_row_stride(opencv_mat mat, size_t stride) {
    auto m = static_cast<cv::Mat *>(mat);
    if (m->step == stride) {
        return true;
    }
    size_t width = m->cols;
    size_t height = m->rows;
    auto type = m->type();
    auto width_stride = width * CV_ELEM_SIZE(type);
    if (stride < width_stride) {
        return false;
    }
    if (m->step != width_stride) {
        // refuse to set the stride if it's already set
        // the math for that is confusing and probably unnecessary to figure out
        return false;
    }
    size_t total_size = stride * height;
    if ((m->datastart + total_size) > m->datalimit) {
        // don't exceed end of data array
        return false;
    }
    m->step = stride;
    return true;
}

void opencv_mat_release(opencv_mat mat) {
    auto m = static_cast<cv::Mat *>(mat);
    delete m;
}

int opencv_type_depth(int type) {
    return CV_ELEM_SIZE1(type) * 8;
}

int opencv_type_channels(int type) {
    return CV_MAT_CN(type);
}

int opencv_type_convert_depth(int t, int depth) {
    return CV_MAKETYPE(depth, CV_MAT_CN(t));
}

opencv_decoder opencv_decoder_create(const opencv_mat buf) {
    auto mat = static_cast<const cv::Mat *>(buf);
    cv::ImageDecoder *d = new cv::ImageDecoder(*mat);
    if (d->empty()) {
        delete d;
        d = NULL;
    }
    return d;
}

const char *opencv_decoder_get_description(const opencv_decoder d) {
    auto d_ptr = static_cast<cv::ImageDecoder *>(d);
    return d_ptr->getDescription().c_str();
}

void opencv_decoder_release(opencv_decoder d) {
    auto d_ptr = static_cast<cv::ImageDecoder *>(d);
    delete d_ptr;
}

bool opencv_decoder_read_header(opencv_decoder d) {
    auto d_ptr = static_cast<cv::ImageDecoder *>(d);
    return d_ptr->readHeader();
}

int opencv_decoder_get_width(const opencv_decoder d) {
    auto d_ptr = static_cast<cv::ImageDecoder *>(d);
    return d_ptr->width();
}

int opencv_decoder_get_height(const opencv_decoder d) {
    auto d_ptr = static_cast<cv::ImageDecoder *>(d);
    return d_ptr->height();
}

int opencv_decoder_get_pixel_type(const opencv_decoder d) {
    auto d_ptr = static_cast<cv::ImageDecoder *>(d);
    return d_ptr->type();
}

int opencv_decoder_get_orientation(const opencv_decoder d) {
    auto d_ptr = static_cast<cv::ImageDecoder *>(d);
    return d_ptr->orientation();
}

bool opencv_decoder_read_data(opencv_decoder d, opencv_mat dst) {
    auto d_ptr = static_cast<cv::ImageDecoder *>(d);
    auto *mat = static_cast<cv::Mat *>(dst);
    return d_ptr->readData(*mat);
}

opencv_encoder opencv_encoder_create(const char *ext, opencv_mat dst) {
    auto *mat = static_cast<cv::Mat *>(dst);
    return new cv::ImageEncoder(ext, *mat);
}

void opencv_encoder_release(opencv_encoder e) {
    auto e_ptr = static_cast<cv::ImageEncoder *>(e);
    delete e_ptr;
}

bool opencv_encoder_write(opencv_encoder e, const opencv_mat src, const int *opt, size_t opt_len) {
    auto e_ptr = static_cast<cv::ImageEncoder *>(e);
    auto mat = static_cast<const cv::Mat *>(src);
    std::vector<int> params;
    for (size_t i = 0; i < opt_len; i++) {
        params.push_back(opt[i]);
    }
    return e_ptr->write(*mat, params);
};

void opencv_mat_resize(const opencv_mat src, opencv_mat dst, int width, int height, int interpolation) {
    cv::resize(*static_cast<const cv::Mat *>(src), *static_cast<cv::Mat *>(dst), cv::Size(width, height), 0, 0, interpolation);
}

opencv_mat opencv_mat_crop(const opencv_mat src, int x, int y, int width, int height) {
    auto ret = new cv::Mat;
    *ret = (*static_cast<const cv::Mat *>(src))(cv::Rect(x, y, width, height));
    return ret;
}

void opencv_mat_orientation_transform(CVImageOrientation orientation, opencv_mat mat) {
    auto cvMat = static_cast<cv::Mat *>(mat);
    cv::OrientationTransform(int(orientation), *cvMat);
}

int opencv_mat_get_width(const opencv_mat mat) {
    auto cvMat = static_cast<const cv::Mat *>(mat);
    return cvMat->cols;
}

int opencv_mat_get_height(const opencv_mat mat) {
    auto cvMat = static_cast<const cv::Mat *>(mat);
    return cvMat->rows;
}

void *opencv_mat_get_data(const opencv_mat mat) {
    auto cvMat = static_cast<const cv::Mat *>(mat);
    return cvMat->data;
}
