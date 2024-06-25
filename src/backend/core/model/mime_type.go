package model

import "cognix.ch/api/v2/core/proto"

const (
	MIMEURL = "url"
)

var SupportedMimeTypes = map[string]proto.FileType{
	MIMEURL: proto.FileType_URL,
	"application/vnd.openxmlformats-officedocument.spreadsheetml.sheet": proto.FileType_XLSX,
	"application/vnd.ms-excel": proto.FileType_XLS,
	"application/vnd.openxmlformats-officedocument.wordprocessingml.document": proto.FileType_DOCX,
	"application/msword": proto.FileType_DOC,
	"application/pdf":    proto.FileType_PDF,
	"text/plain":         proto.FileType_TXT,
	"application/vnd.openxmlformats-officedocument.presentationml.presentation": proto.FileType_PPTX,
	"application/vnd.ms-powerpoint":                                             proto.FileType_PPT,
	"application/vnd.ms-xpsdocument":                                            proto.FileType_XPS,
	"application/oxps":                                                          proto.FileType_XPS,
	"application/epub+zip":                                                      proto.FileType_EPUB,
	"application/hwp+zip":                                                       proto.FileType_HWPX,
	"text/markdown":                                                             proto.FileType_MD,
	"application/x-mobipocket-ebook":                                            proto.FileType_MOBI,
	"application/fb2":                                                           proto.FileType_FB2,
	"audio/mpeg":                                                                proto.FileType_MP3,
	"video/mpeg":                                                                proto.FileType_MPEG,
	"video/mp4":                                                                 proto.FileType_MP4,
	"video/mpga":                                                                proto.FileType_MPGA,
	"audio/wav":                                                                 proto.FileType_WAV,
	"video/webm":                                                                proto.FileType_WEBM,
	"video/mov":                                                                 proto.FileType_MOV,
	"video/m4a":                                                                 proto.FileType_M4A,
}

var SupportedExtensions = map[string]string{
	"PDF":  "application/pdf",
	"XLSX": "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet",
	"XLS":  "application/vnd.ms-excel",
	"DOCX": "application/vnd.openxmlformats-officedocument.wordprocessingml.document",
	"DOC":  "application/msword",
	"PPT":  "application/vnd.ms-powerpoint",
	"PPTX": "application/vnd.openxmlformats-officedocument.presentationml.presentation",
	"MD":   "text/markdown",
	"HWPX": "application/hwp+zip",
	"MOBI": "application/x-mobipocket-ebook",
	"FB2":  "application/fb2",
	"MP3":  "audio/mpeg",
	"MPEG": "video/mpeg",
	"MP4":  "video/mp4",
	"MPGA": "video/mpga",
	"WAV":  "audio/wav",
	"WEBM": "video/webm",
	"MOV":  "video/mov",
	"M4A":  "video/m4a",
}

var WhisperFileTypes = map[proto.FileType]bool{
	proto.FileType_MP3:  true,
	proto.FileType_MP4:  true,
	proto.FileType_MPEG: true,
	proto.FileType_MPGA: true,
	proto.FileType_M4A:  true,
	proto.FileType_WAV:  true,
	proto.FileType_WEBM: true,
	proto.FileType_MOV:  true,
}
