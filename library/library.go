package library

import (
	"bytes"
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"encoding/base64"
	"encoding/csv"
	"encoding/hex"
	"fmt"
	"io"
	"math"
	"math/rand"
	"net/http"

	"os"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/avct/uasurfer"
	"github.com/bytedance/sonic"
	"github.com/gocarina/gocsv"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
	"github.com/patrickmn/go-cache"
	"github.com/sirupsen/logrus"
	"github.com/snowzach/rotatefilehook"
	"github.com/xuri/excelize/v2"

	Sendinblue "github.com/sendinblue/APIv3-go-library/v2/lib"
)

type Type int
type Library interface {
	LoadEnv(filenames ...string) error
	Getenv(key string) string
	ReadEnv(cfg interface{}) error
	Getwd() (string, error)
	ReadConfig(path string, cfg interface{}) error
	GetNow() time.Time
	JsonMarshal(v interface{}) ([]byte, error)
	JsonUnmarshal(data []byte, v interface{}) error
	RandRead(b []byte) (n int, err error)
	AESNewCipher(key []byte) (cipher.Block, error)
	GetAES256CBCBlockSize() int
	Base64DecodeString(s string) ([]byte, error)
	GetSlicesByteLen(v []byte) int
	ParseInt(v byte) int
	ParseTime(layout, value string) (time.Time, error)
	HasSuffix(s []byte, suffix []byte) bool
	BytesEqual(a, b []byte) bool
	JWTSignedString(jwtToken *jwt.Token, key interface{}) (string, error)
	CastJWTSigningMethodHMAC(jwtToken *jwt.Token) (*jwt.SigningMethodHMAC, bool)
	CastJWTMapClaims(jwtToken *jwt.Token) (jwt.MapClaims, bool)
	GetNewRotateFileHook(config rotatefilehook.RotateFileConfig) (logrus.Hook, error)
	SendTransacEmail(client *Sendinblue.APIClient, ctx context.Context, sendSmtpEmail Sendinblue.SendSmtpEmail) error
	StringsJoin(strs []string, sep string) string
	StringsReplace(s string, old string, new string, n int) string
	MathCeil(float64) float64
	TrimSpace(s string) string
	ToLower(s string) string
	Sprintf(format string, argrs ...any) string
	GenerateUUID() (string, error)
	ReadFile(path string) ([]byte, error)
	NewDocumentFromReader(r io.Reader) (*goquery.Document, error)
	OuterHtml(s *goquery.Selection) (string, error)
	NewReader(s string) *strings.Reader
	NewBytesReader(b []byte) *bytes.Reader
	Itoa(i int) string
	Atoi(s string) (int, error)
	RandIntn(n int) int
	ToUpper(s string) string
	Contains(s, substr string) bool
	Split(s, sep string) []string
	HashMd5(data []byte) [md5.Size]byte
	Md5ToHex(data []byte) string
	ReadAll(r io.Reader) ([]byte, error)
	NopCloser(r io.Reader) io.ReadCloser
	ExcelizeOpenReader(r io.Reader, opts ...excelize.Options) (*excelize.File, error)
	RandSeed(seed int64)
	Errorf(format string, argrs ...any) error
	ParseUserAgent(userAgent string) *uasurfer.UserAgent
	NewBuffer(buf []byte) *bytes.Buffer
	NewRequest(method string, url string, body io.Reader) (*http.Request, error)
	ReflectNew(typ reflect.Type) reflect.Value
	ReflectTypeOf(i any) reflect.Type
	ReflectStructOf(fields []reflect.StructField) reflect.Type
	ContextWithTimeout(parent context.Context, timeout time.Duration) (context.Context, context.CancelFunc)
	NewRequestWithContext(ctx context.Context, method string, url string, body io.Reader) (*http.Request, error)
	ParseFloat(v string, bitSize int) (float64, error)
	ParseToInt(v string) (int, error)
	Printf(format string, a ...any) (n int, err error)
	Println(a ...any) (n int, err error)
	HttpDo(req *http.Request) (*http.Response, error)
	GetCache(key string) (interface{}, bool)
	SetCache(key, value string, duration int)
	NewBufferString(s string) *bytes.Buffer
	TimeAfter(d time.Duration) <-chan time.Time
	GoCsvUnmarshal(in io.Reader, out interface{}) error
	ReplaceAll(s string, old string, new string) string
	CsvNewWriter(w io.Writer) *csv.Writer
}

type LibraryImpl struct {
	loadEnv               func(filenames ...string) error
	getEnv                func(key string) string
	readEnv               func(cfg interface{}) error
	getWd                 func() (string, error)
	readConfig            func(path string, cfg interface{}) error
	getNow                func() time.Time
	jsonMarshal           func(v interface{}) ([]byte, error)
	jsonUnmarshal         func(data []byte, v interface{}) error
	randRead              func(b []byte) (n int, err error)
	aesNewCipher          func(key []byte) (cipher.Block, error)
	aes256CbcBlockSize    int
	base64DecodeString    func(s string) ([]byte, error)
	byteHasSuffix         func(s []byte, suffix []byte) bool
	byteEqual             func(a, b []byte) bool
	rotateFileHook        func(config rotatefilehook.RotateFileConfig) (logrus.Hook, error)
	stringsJoin           func(strs []string, sep string) string
	stringsReplace        func(s string, old string, new string, n int) string
	mathCeil              func(float64) float64
	trimSpace             func(string) string
	toLower               func(string) string
	sprintF               func(format string, args ...any) string
	newRandomUUID         func() (uuid.UUID, error)
	readFile              func(name string) ([]byte, error)
	newDocumentFromReader func(r io.Reader) (*goquery.Document, error)
	outerHtml             func(s *goquery.Selection) (string, error)
	newReader             func(s string) *strings.Reader
	newBytesReader        func(b []byte) *bytes.Reader
	itoa                  func(i int) string
	atoi                  func(s string) (int, error)
	randIntn              func(n int) int
	toUpper               func(s string) string
	contains              func(s, substr string) bool
	split                 func(s, sep string) []string
	hashMd5               func(data []byte) [md5.Size]byte
	hexMd5                func(data []byte) string
	parseTime             func(layout, value string) (time.Time, error)
	readAll               func(r io.Reader) ([]byte, error)
	nopCloser             func(r io.Reader) io.ReadCloser
	excelizeOpenReader    func(r io.Reader, opts ...excelize.Options) (*excelize.File, error)
	randSeed              func(seed int64)
	errorF                func(format string, args ...any) error
	parseUserAgent        func(ua string) *uasurfer.UserAgent
	newBuffer             func(buf []byte) *bytes.Buffer
	newRequest            func(method string, url string, body io.Reader) (*http.Request, error)
	reflectNew            func(typ reflect.Type) reflect.Value
	reflectTypeOf         func(i any) reflect.Type
	reflectStructOf       func(fields []reflect.StructField) reflect.Type
	contextWithTimeout    func(parent context.Context, timeout time.Duration) (context.Context, context.CancelFunc)
	newRequestWithContext func(ctx context.Context, method string, url string, body io.Reader) (*http.Request, error)
	parseFloat            func(v string, bitSize int) (float64, error)
	parseToInt            func(v string) (int, error)
	printf                func(format string, a ...any) (n int, err error)
	println               func(a ...any) (n int, err error)
	httpDo                func(req *http.Request) (*http.Response, error)
	cache                 *cache.Cache
	newBufferString       func(s string) *bytes.Buffer
	timeAfter             func(d time.Duration) <-chan time.Time
	goCsvUnmarshal        func(in io.Reader, out interface{}) error
	replaceAll            func(s string, old string, new string) string
	csvNewWriter          func(w io.Writer) *csv.Writer
}

func New() Library {
	return &LibraryImpl{
		loadEnv:               godotenv.Load,
		getEnv:                os.Getenv,
		readEnv:               cleanenv.ReadEnv,
		getWd:                 os.Getwd,
		readConfig:            cleanenv.ReadConfig,
		getNow:                time.Now,
		jsonMarshal:           sonic.Marshal,
		jsonUnmarshal:         sonic.Unmarshal,
		randRead:              rand.Read,
		aesNewCipher:          aes.NewCipher,
		aes256CbcBlockSize:    16,
		base64DecodeString:    base64.StdEncoding.DecodeString,
		byteHasSuffix:         bytes.HasSuffix,
		byteEqual:             bytes.Equal,
		rotateFileHook:        rotatefilehook.NewRotateFileHook,
		stringsJoin:           strings.Join,
		stringsReplace:        strings.Replace,
		mathCeil:              math.Ceil,
		trimSpace:             strings.TrimSpace,
		toLower:               strings.ToLower,
		sprintF:               fmt.Sprintf,
		newRandomUUID:         uuid.NewRandom,
		readFile:              os.ReadFile,
		outerHtml:             goquery.OuterHtml,
		newReader:             strings.NewReader,
		newBytesReader:        bytes.NewReader,
		newDocumentFromReader: goquery.NewDocumentFromReader,
		itoa:                  strconv.Itoa,
		atoi:                  strconv.Atoi,
		randIntn:              rand.Intn,
		toUpper:               strings.ToUpper,
		contains:              strings.Contains,
		split:                 strings.Split,
		hashMd5:               md5.Sum,
		hexMd5:                hex.EncodeToString,
		parseTime:             time.Parse,
		readAll:               io.ReadAll,
		nopCloser:             io.NopCloser,
		excelizeOpenReader:    excelize.OpenReader,
		randSeed:              rand.Seed,
		errorF:                fmt.Errorf,
		parseUserAgent:        uasurfer.Parse,
		newBuffer:             bytes.NewBuffer,
		newRequest:            http.NewRequest,
		reflectNew:            reflect.New,
		reflectTypeOf:         reflect.TypeOf,
		reflectStructOf:       reflect.StructOf,
		contextWithTimeout:    context.WithTimeout,
		newRequestWithContext: http.NewRequestWithContext,
		parseFloat:            strconv.ParseFloat,
		parseToInt:            strconv.Atoi,
		printf:                fmt.Printf,
		println:               fmt.Println,
		httpDo:                http.DefaultClient.Do,
		cache:                 cache.New(time.Hour, time.Hour),
		newBufferString:       bytes.NewBufferString,
		timeAfter:             time.After,
		goCsvUnmarshal:        gocsv.Unmarshal,
		replaceAll:            strings.ReplaceAll,
		csvNewWriter:          csv.NewWriter,
	}
}

func (l *LibraryImpl) LoadEnv(filenames ...string) error {
	return l.loadEnv(filenames...)
}

func (l *LibraryImpl) Getenv(key string) string {
	return l.getEnv(key)
}

func (l *LibraryImpl) ReadEnv(cfg interface{}) error {
	return l.readEnv(cfg)
}

func (l *LibraryImpl) Getwd() (string, error) {
	return l.getWd()
}

func (l *LibraryImpl) ReadConfig(path string, cfg interface{}) error {
	return l.readConfig(path, cfg)
}

func (l *LibraryImpl) GetNow() time.Time {
	return l.getNow()
}

func (l *LibraryImpl) JsonMarshal(v interface{}) ([]byte, error) {
	return l.jsonMarshal(v)
}

func (l *LibraryImpl) JsonUnmarshal(data []byte, v interface{}) error {
	return l.jsonUnmarshal(data, v)
}

func (l *LibraryImpl) RandRead(b []byte) (n int, err error) {
	return l.randRead(b)
}

func (l *LibraryImpl) AESNewCipher(key []byte) (cipher.Block, error) {
	return l.aesNewCipher(key)
}

func (l *LibraryImpl) GetAES256CBCBlockSize() int {
	return l.aes256CbcBlockSize
}

func (l *LibraryImpl) Base64DecodeString(s string) ([]byte, error) {
	return l.base64DecodeString(s)
}

func (l *LibraryImpl) GetSlicesByteLen(v []byte) int {
	return len(v)
}

func (l *LibraryImpl) ParseInt(v byte) int {
	return int(v)
}

func (l *LibraryImpl) ParseTime(layout, value string) (time.Time, error) {
	return l.parseTime(layout, value)
}

func (l *LibraryImpl) HasSuffix(s []byte, suffix []byte) bool {
	return l.byteHasSuffix(s, suffix)
}

func (l *LibraryImpl) BytesEqual(a, b []byte) bool {
	return l.byteEqual(a, b)
}

func (l *LibraryImpl) JWTSignedString(jwtToken *jwt.Token, key interface{}) (string, error) {
	return jwtToken.SignedString(key)
}

func (l *LibraryImpl) CastJWTSigningMethodHMAC(jwtToken *jwt.Token) (*jwt.SigningMethodHMAC, bool) {
	method, ok := jwtToken.Method.(*jwt.SigningMethodHMAC)
	if !ok {
		return nil, ok
	}

	return method, ok
}

func (l *LibraryImpl) CastJWTMapClaims(jwtToken *jwt.Token) (jwt.MapClaims, bool) {
	claims, ok := jwtToken.Claims.(jwt.MapClaims)
	if !ok {
		return nil, ok
	}

	return claims, ok
}

func (l *LibraryImpl) GetNewRotateFileHook(config rotatefilehook.RotateFileConfig) (logrus.Hook, error) {
	return l.rotateFileHook(config)
}

func (l *LibraryImpl) SendTransacEmail(client *Sendinblue.APIClient, ctx context.Context, sendSmtpEmail Sendinblue.SendSmtpEmail) error {
	_, _, err := client.TransactionalEmailsApi.SendTransacEmail(ctx, sendSmtpEmail)
	return err
}

func (l *LibraryImpl) StringsJoin(strs []string, sep string) string {
	return l.stringsJoin(strs, sep)
}

func (l *LibraryImpl) StringsReplace(s string, old string, new string, n int) string {
	return l.stringsReplace(s, old, new, n)
}

func (l *LibraryImpl) MathCeil(f float64) float64 {
	return l.mathCeil(f)
}

func (l *LibraryImpl) TrimSpace(s string) string {
	return l.trimSpace(s)
}

func (l *LibraryImpl) ToLower(s string) string {
	return l.toLower(s)
}

func (l *LibraryImpl) Sprintf(format string, argrs ...any) string {
	return l.sprintF(format, argrs...)
}

func (l *LibraryImpl) GenerateUUID() (string, error) {
	// GENERATE RANDOM UUID
	newUUID, err := l.newRandomUUID()
	return newUUID.String(), err
}

func (l *LibraryImpl) ReadFile(path string) ([]byte, error) {
	return l.readFile(path)
}

func (l *LibraryImpl) NewDocumentFromReader(r io.Reader) (*goquery.Document, error) {
	return l.newDocumentFromReader(r)
}

func (l *LibraryImpl) OuterHtml(s *goquery.Selection) (string, error) {
	return l.outerHtml(s)
}

func (l *LibraryImpl) NewReader(s string) *strings.Reader {
	return l.newReader(s)
}

func (l *LibraryImpl) NewBytesReader(b []byte) *bytes.Reader {
	return l.newBytesReader(b)
}

func (l *LibraryImpl) Itoa(i int) string {
	return l.itoa(i)
}

func (l *LibraryImpl) Atoi(s string) (int, error) {
	return l.atoi(s)
}

func (l *LibraryImpl) RandIntn(n int) int {
	return l.randIntn(n)
}

func (l *LibraryImpl) ToUpper(s string) string {
	return l.toUpper(s)
}

func (l *LibraryImpl) Contains(s, substr string) bool {
	return l.contains(s, substr)
}

func (l *LibraryImpl) Split(s, sep string) []string {
	return l.split(s, sep)
}

func (l *LibraryImpl) HashMd5(data []byte) [md5.Size]byte {
	return l.hashMd5(data)
}

func (l *LibraryImpl) Md5ToHex(data []byte) string {
	return l.hexMd5(data)
}

func (l *LibraryImpl) ReadAll(r io.Reader) ([]byte, error) {
	return l.readAll(r)
}

func (l *LibraryImpl) NopCloser(r io.Reader) io.ReadCloser {
	return l.nopCloser(r)
}

func (l *LibraryImpl) ExcelizeOpenReader(r io.Reader, opts ...excelize.Options) (*excelize.File, error) {
	return l.excelizeOpenReader(r, opts...)
}

func (l *LibraryImpl) RandSeed(seed int64) {
	l.randSeed(seed)
}

func (l *LibraryImpl) Errorf(format string, argrs ...any) error {
	return l.errorF(format, argrs...)
}

func (l *LibraryImpl) ParseUserAgent(userAgent string) *uasurfer.UserAgent {
	return l.parseUserAgent(userAgent)
}

func (l *LibraryImpl) NewBuffer(buf []byte) *bytes.Buffer {
	return l.newBuffer(buf)
}

func (l *LibraryImpl) NewRequest(method string, url string, body io.Reader) (*http.Request, error) {
	return l.newRequest(method, url, body)
}

func (l *LibraryImpl) ReflectNew(typ reflect.Type) reflect.Value {
	return l.reflectNew(typ)
}

func (l *LibraryImpl) ReflectTypeOf(i any) reflect.Type {
	return l.reflectTypeOf(i)
}

func (l *LibraryImpl) ReflectStructOf(fields []reflect.StructField) reflect.Type {
	return l.reflectStructOf(fields)
}

func (l *LibraryImpl) ContextWithTimeout(parent context.Context, timeout time.Duration) (context.Context, context.CancelFunc) {
	return l.contextWithTimeout(parent, timeout)
}

func (l *LibraryImpl) NewRequestWithContext(ctx context.Context, method string, url string, body io.Reader) (*http.Request, error) {
	return l.newRequestWithContext(ctx, method, url, body)
}

func (l *LibraryImpl) ParseFloat(v string, bitSize int) (float64, error) {
	return l.parseFloat(v, bitSize)
}

func (l *LibraryImpl) ParseToInt(v string) (int, error) {
	return l.parseToInt(v)
}

func (l *LibraryImpl) Printf(format string, a ...any) (n int, err error) {
	return l.printf(format, a...)
}

func (l *LibraryImpl) Println(a ...any) (n int, err error) {
	return l.println(a...)
}

func (l *LibraryImpl) HttpDo(req *http.Request) (*http.Response, error) {
	return http.DefaultClient.Do(req)
}

func (l *LibraryImpl) GetCache(key string) (interface{}, bool) {
	return l.cache.Get(key)
}

func (l *LibraryImpl) SetCache(key, value string, duration int) {
	l.cache.Set(key, value, time.Duration(duration)*time.Second)
}

func (l *LibraryImpl) NewBufferString(s string) *bytes.Buffer {
	return l.newBufferString(s)
}

func (l *LibraryImpl) TimeAfter(d time.Duration) <-chan time.Time {
	return l.timeAfter(d)
}

func (l *LibraryImpl) GoCsvUnmarshal(in io.Reader, out interface{}) error {
	return l.goCsvUnmarshal(in, out)
}

func (l *LibraryImpl) ReplaceAll(s string, old string, new string) string {
	return l.replaceAll(s, old, new)
}

func (l *LibraryImpl) CsvNewWriter(w io.Writer) *csv.Writer {
	return l.csvNewWriter(w)
}
