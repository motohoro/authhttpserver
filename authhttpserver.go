/**
authhttpserver

Copyright (c) 2015 motohoro

This software is released under the MIT License.
http://opensource.org/licenses/mit-license.php
*/

package main

import (
  "fmt"
  "os"
    "strings"
    "net"
  "net/http"
  "net/url"
    "net/http/cookiejar"
    "io/ioutil"
    "syscall"
    "unsafe"
    "log"
)

//http://stackoverflow.com/questions/9996767/showing-custom-404-error-page-with-standard-http-package
func errorHandler(w http.ResponseWriter, r *http.Request, status int) {
    w.WriteHeader(status)
    if status == http.StatusNotFound {
        fmt.Fprint(w, "404 NOT FOUND")
    }
}
func handler(w http.ResponseWriter, r *http.Request) {
//browser
  fmt.Println("PATH:",r.URL.Path)
  if r.URL.Path !="/"{
    errorHandler(w, r, http.StatusNotFound)
    return;
  }
    if r.URL.Query()["url"][0] == ""{
        errorHandler(w, r, http.StatusNotFound)
        return;
    }
    targetURL := r.URL.Query()["url"][0]
    u, err := url.Parse(targetURL)
    if err != nil {
        panic(err)
        errorHandler(w,r,http.StatusNotFound)
        return;
    }
    host, port, _ := net.SplitHostPort(u.Host)
//    fmt.Println(u.User.Username())
    fmt.Println(host)
    fmt.Println(port)
    req, _:= http.NewRequest("GET", targetURL, nil)
    
  /*
  //get GET parameter http://betterlogic.com/roger/2014/04/golang-go-http-request-how-to-get-get-parameters/
  fmt.Println("got:", r.URL.Query());
  fmt.Println("URL:", r.URL.Query()["url"][0]);
  s2,_ := url.QueryUnescape(r.URL.Query()["url"][0])
  fmt.Println("URL:", s2);
  s3 := url.QueryEscape(r.URL.Query()["url"][0])
  fmt.Println("URL:", s3);
  */
    if r.URL.Query()["buser"][0] != ""{
      //DLL
//        dll, err := syscall.LoadDLL(os.Getenv("HOME")+"\\Documents\\Visual Studio 2015\\Projects\\firefoxdecrypt\\Debug\\"+"firefoxdecrypt.dll")
        fmt.Println(os.Getenv("HOME"))
        dll, err := syscall.LoadDLL("firefoxdecrypt.dll")
        if err != nil {
            log.Fatal(err)
        }
        defer dll.Release()

        proc, err := dll.FindProc("getAllAuthData")
        if err != nil {
            log.Fatal(err)
        }

        a,r2,_:=proc.Call()
        fmt.Println(r2)
        /*
        if r2 != 0 && err != nil {
            fmt.Println("DLL error")
            log.Fatal(err)//The operation completed successfully.
        }
        */
        
        // https://gist.github.com/mattn/9f0729d2ba2356f38cc6
        //C char* => go string
        tmp := *(*[8192]byte)(unsafe.Pointer(a))
        s := ""
        for n := 0; n < len(tmp); n++ {
            if tmp[n] == 0 {
                s = string(tmp[:n])
                break
            }
        }
        //uid,pw,http,realm
        fmt.Println(s)
        
        uid := r.URL.Query()["buser"][0]
        pw := ""
        authlines := strings.Split(s,"\n")
        for h :=0; h<len(s);h++ {
            authline := strings.Split(authlines[h],",")
            u2, _ := url.Parse(authline[2])
            host2,_, _ := net.SplitHostPort(u2.Host)
            if authline[0]==uid && host == host2 {
                pw=authline[1]
                break;
            }
            fmt.Println(authline)
        }
        
        req.SetBasicAuth(uid,pw)
        
    }
    
    //http://blog.sarabande.jp/post/90736041568
    cookieJar, _ := cookiejar.New(nil)
    client := &http.Client {
        Jar: cookieJar,
    }
    req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 6.1; rv:38.0) Gecko/20100101 Firefox/38.0")
    
    
    res, _ := client.Do(req)
    body, _ := ioutil.ReadAll(res.Body)
    defer res.Body.Close()
//    println(string(body))
    //http://qiita.com/futoase/items/ea86b750bbb36d7d859a
    println(res.Header.Get("Content-Type"))

    w.Header().Set("Content-Type",res.Header.Get("Content-Type"))
    fmt.Fprintf(w, string(body))



}

func main() {
  http.HandleFunc("/", handler) // ハンドラを登録してウェブページを表示させる
  http.ListenAndServe(":8087", nil)
}
