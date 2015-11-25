/**
authhttpserver

Copyright (c) 2015 motohoro

This software is released under the MIT License.
http://opensource.org/licenses/mit-license.php
*/

package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
//	"reflect"
	"strings"
	"syscall"
	"unsafe"
    "C"
)

//http://stackoverflow.com/questions/19965795/go-golang-write-log-to-file
func outputLog(s string) {
	f, err := os.OpenFile("logfile.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	log.SetOutput(f)
	log.Println(s)
}

//http://stackoverflow.com/questions/9996767/showing-custom-404-error-page-with-standard-http-package
func errorHandler(w http.ResponseWriter, r *http.Request, status int) {
	w.WriteHeader(status)
	if status == http.StatusNotFound {
		fmt.Fprint(w, "404 NOT FOUND")
	}
}
func handler(w http.ResponseWriter, r *http.Request) {
	//browser
	//  fmt.Println("PATH:",r.URL.Path)
	if r.URL.Path != "/" {
		errorHandler(w, r, http.StatusNotFound)
		return
	}
	if r.URL.Query()["url"][0] == "" {
		errorHandler(w, r, http.StatusNotFound)
		return
	}
	targetURL := r.URL.Query()["url"][0]
	u, err := url.Parse(targetURL)
	if err != nil {
		outputLog(err.Error())
		panic(err)
		errorHandler(w, r, http.StatusNotFound)
		return
	}
	host, port, _ := net.SplitHostPort(u.Host)
	//    fmt.Println(u.User.Username())
	fmt.Println(host)
	fmt.Println(port)
	outputLog("ACCESS:" + targetURL)
	req, _ := http.NewRequest("GET", targetURL, nil)

	/*
	  //get GET parameter http://betterlogic.com/roger/2014/04/golang-go-http-request-how-to-get-get-parameters/
	  fmt.Println("got:", r.URL.Query());
	  fmt.Println("URL:", r.URL.Query()["url"][0]);
	  s2,_ := url.QueryUnescape(r.URL.Query()["url"][0])
	  fmt.Println("URL:", s2);
	  s3 := url.QueryEscape(r.URL.Query()["url"][0])
	  fmt.Println("URL:", s3);
	*/
	if r.URL.Query()["buser"][0] != "" {
		//DLL
		//        dll, err := syscall.LoadDLL(os.Getenv("HOME")+"\\Documents\\Visual Studio 2015\\Projects\\firefoxdecrypt\\Debug\\"+"firefoxdecrypt.dll")
		//        fmt.Println(os.Getenv("HOME"))
		dll, err := syscall.LoadDLL("firefoxdecrypt.dll")
		if err != nil {
			outputLog("err LoadDLL")
			outputLog(err.Error())
			log.Fatal(err)
		}
		defer dll.Release()

		proc, err := dll.FindProc("getAllAuthData")
		if err != nil {
			outputLog("err FindProc")
			outputLog(err.Error())
			log.Fatal(err)
		}

		a, r2, err := proc.Call()
        if err != nil {
		  outputLog("err ProcCall")
		  outputLog(err.Error())
        }
		fmt.Println(r2)
		outputLog("proc Called")

		/*
		   if r2 != 0 && err != nil {
		       fmt.Println("DLL error")
		       log.Fatal(err)//The operation completed successfully.
		   }
		*/

		// https://gist.github.com/mattn/9f0729d2ba2356f38cc6
		//C char* => go string
        /*
		fmt.Println(reflect.TypeOf(a)) //uintptr
		tmp := *(*[8192]byte)(unsafe.Pointer(a))
		fmt.Println(reflect.TypeOf(tmp)) //[8192]uint8
		s := ""
		//        fmt.Println(tmp)
		for n := 0; n < len(tmp); n++ {
			if tmp[n] == 0 {
				s = string(tmp[:n])
				break
			}
		}
        */
        //CGO
        s := C.GoString((*C.char)(unsafe.Pointer(a)))
		//uid,pw,http,realm
		//        fmt.Println(s)
		//        outputLog(s)

		uid := r.URL.Query()["buser"][0]
		pw := ""
		authlines := strings.Split(s, "\n")
		for h := 0; h < len(authlines); h++ {
			authline := strings.Split(authlines[h], ",")
			u2, _ := url.Parse(authline[2])
			host2, _, _ := net.SplitHostPort(u2.Host)
			outputLog(authline[0] + authline[1] + host2)
			if authline[0] == uid && host == host2 {
				pw = authline[1]
				break
			}
			fmt.Println(authline)
		}

		req.SetBasicAuth(uid, pw)
		outputLog(uid + "=" + pw)

	}

	//http://blog.sarabande.jp/post/90736041568
	cookieJar, _ := cookiejar.New(nil)
	client := &http.Client{
		Jar: cookieJar,
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 6.1; rv:38.0) Gecko/20100101 Firefox/38.0")

	res, err := client.Do(req)
	if err != nil {
		outputLog(err.Error())
		log.Fatal(err)
		
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		outputLog(err.Error())
		log.Fatal(err)
		return
	}
	if res.StatusCode>=300 {
		println(res.StatusCode)
		http.Error(w, res.Status, res.StatusCode)
		return
	}
	defer res.Body.Close()
	//    println(string(body))
	//http://qiita.com/futoase/items/ea86b750bbb36d7d859a
	println(res.Header.Get("Content-Type"))

	w.Header().Set("Content-Type", res.Header.Get("Content-Type"))
	fmt.Fprintf(w, string(body))

}

func main2() {
	http.HandleFunc("/", handler) // ハンドラを登録してウェブページを表示させる
	http.ListenAndServe(":8087", nil)

}

func main() {
	//fmt.Print("\x07")
	main2()
}
