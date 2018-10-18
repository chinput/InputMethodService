package main

import (
	"archive/zip"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"gopkg.in/mgo.v2/bson"

	"github.com/chinput/InputMethodService/server/config"
	"github.com/chinput/InputMethodService/server/email"
)

type Output struct {
	Status int    `json:"status"`
	Desc   string `json:"desc"`

	Url   string `json:"url"`
	User  string `json:"user"`
	Token string `json:"token"`
}

type Input struct {
	// 请求类型
	Type int `json:"type"`

	// 用于验证登录状态
	Token string `json:"token"`

	// 用于登录
	User string `json:"user"`
	Pass string `json:"pass"`

	// 发送请求的时间，用于混淆加密后的结果
	Time string `json:"time"`

	// 请求下载的用户文件，不同文件间用 | 连接
	ReqFile string `json:"req_file"`

	// 用于验证注册码之类
	Email string `json:"email"`
	Code  string `json:"code"`

	UploadFile     multipart.File
	UploadFileSize int64
}

var (
	outputMap = []string{
		"ok",
		"...",
	}
)

func newOutput(status int, url, user, token string) *Output {
	desc := "error"
	if len(outputMap) > status {
		desc = outputMap[status]
	}

	return &Output{
		Status: status,
		Desc:   desc,
		Url:    url,
		User:   user,
		Token:  token,
	}
}

func writeOutputTo(w io.Writer, status int, url, user, token string) {
	o := newOutput(status, url, user, token)
	data, _ := json.Marshal(o)
	w.Write(data)
}

func writeOutError(w io.Writer, status int) {
	writeOutputTo(w, status, "", "", "")
}

func EmailExistInDB(email_str string) bool {
	u := newUser(nil)
	defer u.End()

	exist_user := u.FindOne(bson.M{
		"email": email_str,
	})

	if exist_user != nil {
		return true
	}

	return false
}

func Register(w http.ResponseWriter, i *Input) {
	if i.Email == "" {
		writeOutError(w, 5)
		return
	}

	if EmailExistInDB(i.Email) {
		writeOutError(w, 15)
		return
	}

	err := SendRegisterCode(i.Email)
	if err != nil {
		writeOutError(w, 6)
		return
	}
	writeOutputTo(w, 0, "", "", "")
}

func CheckTheCode(w http.ResponseWriter, i *Input) {

	regpool.lock.Lock()
	defer regpool.lock.Unlock()

	exist := regpool.group[i.Email]
	if exist == nil {
		writeOutError(w, 10)
		return
	}
	if exist.code != i.Code {
		writeOutError(w, 10)
		return
	}

	if exist.time.Add(time.Minute * 15).Before(time.Now()) {
		writeOutError(w, 11)
		return
	}

	if EmailExistInDB(i.Email) {
		writeOutError(w, 15)
		return
	}

	u := newUser(nil)
	defer u.End()

	path := config.PatternTime(config.DataPath(), time.Now(), "/"+i.Email)

	data := DataUser{
		Username:      i.User,
		Password:      i.Pass,
		DataPath:      path,
		Email:         i.Email,
		Register_Time: time.Now().Unix(),
	}
	id, err := u.Add(data)
	if err != nil {
		writeOutError(w, 9)
		return
	}

	path = config.PatternTime(config.DataPath(), time.Now(), "/"+id)
	u.UpdateById(id, bson.M{
		"data_path": path,
	})

	auth := NewAuth(u)
	token, err := auth.AddAnAuthToken(id)
	if err != nil {
		writeOutError(w, 9)
		return
	}

	writeOutputTo(w, 0, "", i.User, token)
}

func Login(w http.ResponseWriter, i *Input) {
	if i.Email == "" || i.Pass == "" {
		writeOutError(w, 7)
		return
	}

	u := newUser(nil)
	defer u.End()

	cond := bson.M{
		"email":    i.Email,
		"password": i.Pass,
	}
	data := u.FindOne(cond)
	if data == nil {
		writeOutError(w, 8)
		return
	}

	auth := NewAuth(u)
	uid := data["_id"].(bson.ObjectId)
	token, err := auth.AddAnAuthToken(uid.Hex())
	if err != nil {
		writeOutError(w, 9)
		return
	}
	writeOutputTo(w, 0, "", data["username"].(string), token)
}

func bsonToStruct(in bson.M, out interface{}) error {
	buff, err := bson.Marshal(in)
	if err != nil {
		return err
	}
	return bson.Unmarshal(buff, out)
}
func CheckAuthValid(i *Input) (bool, string, *DataUser) {
	auth := NewAuth(nil)
	defer auth.End()
	cond := bson.M{
		"token":       i.Token,
		"logout_time": 0,
	}

	data := auth.FindOne(cond)
	if data == nil {
		return false, "", nil
	}

	data_auth := new(DataAuth)
	bsonToStruct(data, data_auth)

	u := newUser(auth)
	user_data := u.FindById(data_auth.Uid)

	if user_data == nil {
		return false, "", nil
	}

	user_data_struct := new(DataUser)
	bsonToStruct(user_data, user_data_struct)

	user_email := user_data_struct.Email

	if user_email != i.Email {
		return false, "", nil
	}

	return true, data_auth.Uid, user_data_struct
}

type File struct {
	url    string
	path   string
	expire time.Time
	time   time.Time
}

type FilePool struct {
	group map[string]*File
	lock  sync.Mutex
}

func (f *FilePool) Clear() {
	f.lock.Lock()
	defer f.lock.Unlock()
	now := time.Now()
	for key, v := range f.group {
		if v.expire.Before(now) {
			os.Remove(v.path)
			delete(f.group, key)
		}
	}
}

func (f *FilePool) Guard() {
	for {
		<-time.After(time.Second * 10)
		f.Clear()
	}
}

func (f *FilePool) Add(uid string, files []string) (string, error) {
	f.lock.Lock()
	defer f.lock.Unlock()

	exist := f.group[uid]
	if exist != nil {
		if exist.time.Add(time.Second).Before(time.Now()) {
			return exist.url, nil
		}
	}

	path := filepath.Join(config.TmpPath(), uid+".zip")
	url := filepath.Join("/tmp", uid+".zip")

	f.group[uid] = &File{
		path:   path,
		url:    url,
		time:   time.Now(),
		expire: time.Now().Add(time.Second * 5),
	}

	w, err := os.Create(path)
	if err != nil {
		return "", err
	}

	defer w.Close()
	z := zip.NewWriter(w)
	defer z.Close()

	for _, name := range files {
		_, short := filepath.Split(name)
		fp, err := z.Create(short)
		if err != nil {
			return "", err
		}
		r, err := os.Open(name)
		if err != nil {
			return "", err
		}
		_, err = io.Copy(fp, r)
		r.Close()
		if err != nil {
			return "", err
		}
	}

	return url, nil
}

var fpool *FilePool

func init() {
	fpool = new(FilePool)
	fpool.group = make(map[string]*File)
	go fpool.Guard()
}

func DownloadFile(w http.ResponseWriter, i *Input) {
	login, uid, user_info := CheckAuthValid(i)

	if !login || user_info == nil {
		writeOutError(w, 12)
		return
	}
	reqFiles := strings.Split(i.ReqFile, "|")
	reqFilesFull := make([]string, 0, len(reqFiles))

	for _, name := range reqFiles {
		full_path := filepath.Join(user_info.DataPath, name)
		reqFilesFull = append(reqFilesFull, full_path)
	}
	url, err := fpool.Add(uid, reqFilesFull)
	if err != nil {
		writeOutError(w, 13)
		return
	}

	writeOutputTo(w, 0, url, user_info.Username, i.Token)
}

func UploadFile(w http.ResponseWriter, i *Input) {
	if i.UploadFile == nil || i.UploadFileSize == 0 {
		writeOutError(w, 14)
		return
	}

	defer i.UploadFile.Close()

	login, _, user_info := CheckAuthValid(i)

	if !login {
		writeOutError(w, 12)
	}

	r, err := zip.NewReader(i.UploadFile, i.UploadFileSize)
	if err != nil {
		writeOutError(w, 15)
		return
	}

	userPath := user_info.DataPath

	os.MkdirAll(userPath, 0755)

	for _, oneFile := range r.File {
		fullName := filepath.Join(userPath, oneFile.Name)
		wp, err := os.Create(fullName)
		if err != nil {
			writeOutError(w, 13)
			return
		}

		defer wp.Close()

		rd, err := oneFile.Open()
		if err != nil {
			writeOutError(w, 13)
			return
		}
		defer rd.Close()
		io.Copy(wp, rd)

		wp.Close()
		rd.Close()
	}
}

var (
	handlerMap = []func(w http.ResponseWriter, i *Input){
		Register,
		CheckTheCode,
		Login,
		DownloadFile,
		UploadFile,
	}
)

const DEBUG bool = true

func decReqBuffer(buff string) ([]byte, error) {
	if DEBUG {
		return []byte(buff), nil
	}
	return config.DecodeB64(buff)
}

func CommonHandler(w http.ResponseWriter, r *http.Request) {

	buff := r.FormValue("q")

	if buff == "" {
		writeOutError(w, 2)
		return
	}

	reqBuffer, err := decReqBuffer(buff)

	if err != nil {
		writeOutError(w, 2)
		return
	}

	req := new(Input)

	err = json.Unmarshal(reqBuffer, req)

	if err != nil {
		writeOutError(w, 3)
		return
	}

	if req.Email != "" {
		if !email.IsEmailValid(req.Email) {
			writeOutError(w, 4)
		}
	}

	uploadFile, header, err := r.FormFile("upload")
	if err == nil {
		req.UploadFile = uploadFile
		req.UploadFileSize = header.Size
	}

	fn := handlerMap[req.Type]
	if fn != nil {
		fn(w, req)
	}
}

/*
{
	"Type":0,
	"Email":"garfeng_gu@163.com"
}
*/
