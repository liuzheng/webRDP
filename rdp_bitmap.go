package main

/*
#include <freerdp/graphics.h>
#include "webrdp.h"
static int getsizeof(BYTE* p) {
	return sizeof(p);
}

static BYTE* getBYTEpalette(rdpContext* context) {
	webContext* xfc = (webContext*) context;
	return (BYTE*)xfc->palette;
}
*/
import "C"
import (
	"bufio"
	"bytes"
	"golang.org/x/image/bmp"
	"image/png"
	"log"
	"os"
	"time"
	"unsafe"
)

//export webRdpBitmapNew
func webRdpBitmapNew(context *C.rdpContext, bitmap *C.rdpBitmap) C.BOOL {
	log.Println("webRdpBitmapNew")
	if bitmap.data != nil {
		// l := C.getsizeof(bitmap.data)
		log.Printf("webRdpBitmapNew length::%d", bitmap.length)
		r := bytes.NewReader(C.GoBytes(unsafe.Pointer(bitmap.data), C.int(bitmap.length)))
		image, err := bmp.Decode(r)
		if err != nil {
			log.Println("webRdpBitmapNew Decode error--------------")
			return C.TRUE
		}
		t := time.Now().Format("20060102150405")
		fo, err := os.Create("tmp/" + t)
		defer fo.Close()
		w := bufio.NewWriter(fo)
		e := png.Encode(w, image)
		if e != nil {
			log.Println("webRdpBitmapNew Encode error=============")
		}
	}
	return C.TRUE
}

//export webRdpBitmapFree
func webRdpBitmapFree(context *C.rdpContext, bitmap *C.rdpBitmap) {
	log.Println("webRdpBitmapFree")
}

//export webRdpBitmapPaint
func webRdpBitmapPaint(context *C.rdpContext, bitmap *C.rdpBitmap) C.BOOL {
	log.Println("webRdpBitmapPaint")
	return C.TRUE
}

//export webRdpBitmapDecompress
func webRdpBitmapDecompress(context *C.rdpContext, bitmap *C.rdpBitmap, data *C.BYTE,
	width C.int, height C.int, bpp C.int, length C.int,
	compressed C.BOOL, codecId C.int) C.BOOL {
	log.Println("webRdpBitmapDecompress")
	log.Printf("compressed:%d bpp:%d", compressed, bpp)
	size := width * height * 4
	bitmap.data = (*C.BYTE)(C._aligned_malloc(C.size_t(size), 16))
	if compressed != C.FALSE {
		if bpp < 32 {
			C.freerdp_client_codecs_prepare(context.codecs, C.FREERDP_CODEC_INTERLEAVED)
			C.interleaved_decompress(context.codecs.interleaved, data, C.UINT32(length), bpp,
				&(bitmap.data), C.PIXEL_FORMAT_XRGB32, -1, 0, 0, width, height, C.getBYTEpalette(context))
		} else {
			C.freerdp_client_codecs_prepare(context.codecs, C.FREERDP_CODEC_PLANAR)
			status := C.planar_decompress(context.codecs.planar, data, C.UINT32(length),
				&(bitmap.data), C.PIXEL_FORMAT_XRGB32, -1, 0, 0, width, height, C.TRUE)
			log.Printf("webRdpBitmapDecompress status::::::%d", status)
		}
	} else {
		C.freerdp_image_flip(data, bitmap.data, width, height, bpp)
	}
	bitmap.compressed = C.FALSE
	bitmap.length = C.UINT32(size)
	bitmap.bpp = 32
	return C.TRUE
}

//export webRdpBitmapSetSurface
func webRdpBitmapSetSurface(context *C.rdpContext, bitmap *C.rdpBitmap, primary C.BOOL) C.BOOL {
	log.Println("webRdpBitmapDecompress")
	return C.TRUE
}
