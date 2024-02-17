package taglib

/*
#cgo pkg-config: taglib
#cgo LDFLAGS: -ltag_c
#include <stdlib.h>
#include <tag_c.h>

// this is needed because cgo doesnt support macros
// its just a wrapper around TAGLIB_COMPLEX_PROPERTY_PICTURE
void taglib_complex_property_set_picture(TagLib_File *file, char *data, unsigned int size, char *desc, char *mime, char *typ)
{
  TAGLIB_COMPLEX_PROPERTY_PICTURE(props, data, size, desc, mime, typ);
  taglib_complex_property_set(file, "PICTURE", props);
}
*/
import "C"
import (
	"errors"
	"time"
	"unsafe"
)

func init() {
	C.taglib_set_string_management_enabled(0)
}

var (
	ErrInvalid   = errors.New("invalid file")
	ErrSave      = errors.New("cannot save file")
	ErrNoPicture = errors.New("no picture")
)

var SupportedExtensions = []string{
	".aac", ".flac", ".mp3", ".mp4", ".m4a", ".m4p", ".ogg", ".oga", ".flac", ".spx", ".opus", ".musepack", ".mpc", ".ape", ".wma", ".tta", ".wav", ".aiff", ".aif",
}

type File struct {
	fp    *C.TagLib_File
	tag   *C.TagLib_Tag
	props *C.TagLib_AudioProperties
}

func (file *File) Save() error {
	var err error
	if C.taglib_file_save(file.fp) == 0 {
		err = ErrSave
	}
	return err
}

func (file *File) Close() {
	C.taglib_file_free(file.fp)
	file.fp = nil
	file.tag = nil
	file.props = nil
}

// Tag properties

func (file *File) Album() string {
	return convertAndFree(C.taglib_tag_album(file.tag))
}
func (file *File) SetAlbum(album string) {
	cs := C.CString(album)
	defer C.free(unsafe.Pointer(cs))
	C.taglib_tag_set_album(file.tag, cs)
}

func (file *File) Artist() string {
	return convertAndFree(C.taglib_tag_artist(file.tag))
}
func (file *File) SetArtist(artist string) {
	cs := C.CString(artist)
	defer C.free(unsafe.Pointer(cs))
	C.taglib_tag_set_artist(file.tag, cs)
}

func (file *File) Comment() string {
	return convertAndFree(C.taglib_tag_comment(file.tag))
}
func (file *File) SetComment(comment string) {
	cs := C.CString(comment)
	defer C.free(unsafe.Pointer(cs))
	C.taglib_tag_set_comment(file.tag, cs)
}

func (file *File) Genre() string {
	return convertAndFree(C.taglib_tag_genre(file.tag))
}
func (file *File) SetGenre(genre string) {
	cs := C.CString(genre)
	defer C.free(unsafe.Pointer(cs))
	C.taglib_tag_set_genre(file.tag, cs)
}

func (file *File) Title() string {
	return convertAndFree(C.taglib_tag_title(file.tag))
}
func (file *File) SetTitle(title string) {
	cs := C.CString(title)
	defer C.free(unsafe.Pointer(cs))
	C.taglib_tag_set_title(file.tag, cs)
}

func (file *File) Track() uint {
	return uint(C.taglib_tag_track(file.tag))
}
func (file *File) SetTrack(track uint) {
	C.taglib_tag_set_track(file.tag, C.uint(track))
}

func (file *File) Year() uint {
	return uint(C.taglib_tag_year(file.tag))
}
func (file *File) SetYear(year uint) {
	C.taglib_tag_set_year(file.tag, C.uint(year))
}

// Support for non-standard tags

// This does the same as file.GetTag("ALBUMARTIST").
// taglib doesn't have a native way of doing this since it only supports the standard tags
func (file *File) AlbumArtist() string {
	return file.GetTag("ALBUMARTIST")
}

// This does the same as file.SetTag("ALBUMARTIST", ...)
// taglib doesn't have a native way of doing this since it only supports the standard tags
func (file *File) SetAlbumArtist(albumArtist string) {
	file.SetTag("ALBUMARTIST", albumArtist)
}

// Used to get non-standard tags
func (file *File) GetTag(name string) string {
	cs := C.CString(name)
	defer C.free(unsafe.Pointer(cs))
	property := C.taglib_property_get(file.fp, cs)
	defer C.free(unsafe.Pointer(property))
	array := toGoStringArray(property)
	if len(array) > 0 {
		return array[0]
	}
	return ""
}

// Used to set non-standard tags
func (file *File) SetTag(name, value string) {
	csName := C.CString(name)
	defer C.free(unsafe.Pointer(csName))
	csValue := C.CString(value)
	defer C.free(unsafe.Pointer(csValue))
	C.taglib_property_set(file.fp, csName, csValue)
}

// Pictures

// Picture is directly taken from TagLibs' Complex_Property_Picture_Data struct
type Picture struct {
	MimeType    string
	PictureType string
	Description string
	Data        []byte
	Size        uint
}

// Gets the first picture from the file
// TODO: support multiple pictures
func (file *File) Picture() (*Picture, error) {
	cs := C.CString("PICTURE")
	defer C.free(unsafe.Pointer(cs))
	property := C.taglib_complex_property_get(file.fp, cs)
	if *property == nil {
		return nil, ErrNoPicture
	}
	defer C.taglib_complex_property_free(property)
	// data doesnt need to be freed since it just contains pointers to the data in the property
	var data C.TagLib_Complex_Property_Picture_Data
	C.taglib_picture_from_complex_property(property, &data)

	return &Picture{
		MimeType:    C.GoString(data.mimeType),
		PictureType: C.GoString(data.pictureType),
		Description: C.GoString(data.description),
		Data:        C.GoBytes(unsafe.Pointer(data.data), C.int(data.size)),
		Size:        uint(data.size),
	}, nil
}

// Sets the first picture in the file
// This removes all other pictures, if any
func (file *File) SetPicture(picture *Picture) error {
	csData := C.CBytes(picture.Data)
	defer C.free(csData)
	csDesc := C.CString(picture.Description)
	defer C.free(unsafe.Pointer(csDesc))
	csMime := C.CString(picture.MimeType)
	defer C.free(unsafe.Pointer(csMime))
	csType := C.CString(picture.PictureType)
	defer C.free(unsafe.Pointer(csType))

	C.taglib_complex_property_set_picture(file.fp, (*C.char)(csData), C.uint(picture.Size), csDesc, csMime, csType)
	return nil
}

// Audio properties

func (file *File) Bitrate() uint {
	return uint(C.taglib_audioproperties_bitrate(file.props))
}

func (file *File) Channels() uint {
	return uint(C.taglib_audioproperties_channels(file.props))
}

func (file *File) Length() time.Duration {
	length := C.taglib_audioproperties_length(file.props)
	return time.Duration(length) * time.Second
}

func (file *File) Samplerate() uint {
	return uint(C.taglib_audioproperties_samplerate(file.props))
}

func Read(filename string) (*File, error) {
	cs := C.CString(filename)
	defer C.free(unsafe.Pointer(cs))
	fp := C.taglib_file_new(cs)
	if fp == nil || C.taglib_file_is_valid(fp) == 0 {
		return nil, ErrInvalid
	}

	return &File{
		fp:    fp,
		tag:   C.taglib_file_tag(fp),
		props: C.taglib_file_audioproperties(fp),
	}, nil
}
