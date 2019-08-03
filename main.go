package main

import (
    "html/template"
    "path/filepath"
    "net/http"
    "net/url"
    "net"
    "crypto/tls"
    "log"
    "github.com/gorilla/mux"
    "fmt"
    "io/ioutil"
    "encoding/json"
    "os"
    "google.golang.org/api/googleapi/transport"
    "google.golang.org/api/youtube/v3"
    bolt "github.com/johnnadratowski/golang-neo4j-bolt-driver"
    "github.com/DDRBoxman/go-amazon-product-api"
    xj "github.com/basgys/goxml2json"
    "strings"
    "gopkg.in/mgo.v2"
    "gopkg.in/mgo.v2/bson"
    "time"
    "bytes"
    //"golang.org/x/crypto/bcrypt"
    "crypto/aes"
    "crypto/cipher"
    "crypto/rand"
    "errors"
    "io"
    R "math/rand"
    "strconv"
    "encoding/hex"
    "github.com/logpacker/PayPal-Go-SDK"
    //"github.com/leebenson/paypal"
    //"github.com/kyokomi/cloudinary"
    //"golang.org/x/net/context"
    //"github.com/google/uuid"
    //"github.com/gotsunami/go-cloudinary"
)

func randInt31(low, hi int64) int64 {
    return low + R.Int63n(hi-low)
}


func encrypt(plaintext []byte, key []byte) ([]byte, error) {
    c, err := aes.NewCipher(key)
    if err != nil {
        return nil, err
    }

    gcm, err := cipher.NewGCM(c)
    if err != nil {
        return nil, err
    }

    nonce := make([]byte, gcm.NonceSize())
    if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
        return nil, err
    }

    return gcm.Seal(nonce, nonce, plaintext, nil), nil
}


/*func decrypt(ciphertext []byte, key []byte) ([]byte, error) {
    c, err := aes.NewCipher(key)
    if err != nil {
        return nil, err
    }

    gcm, err := cipher.NewGCM(c)
    if err != nil {
        return nil, err
    }

    nonceSize := gcm.NonceSize()
    if len(ciphertext) < nonceSize {
        return nil, errors.New("ciphertext too short")
    }

    nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
    return gcm.Open(nil, nonce, ciphertext, nil)
}*/


func decrypt(ciphertext1 string, key1 string) ([]byte, error) {
    //ciphertext:=[]byte(ciphertext1)
    ciphertext,_:=hex.DecodeString(ciphertext1)
    key:=[]byte(key1)
    c, err := aes.NewCipher(key)
    if err != nil {
        return nil, err
    }

    gcm, err := cipher.NewGCM(c)
    if err != nil {
        return nil, err
    }

    nonceSize := gcm.NonceSize()
    if len(ciphertext) < nonceSize {
        return nil, errors.New("ciphertext too short")
    }

    nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
    return gcm.Open(nil, nonce, ciphertext, nil)
}


type Person struct {
  FirstName string
  LastName string
}
type Me struct {
  MyDetails []interface{}
}

type AuthedPlaylist struct {
  PlaylistOf []interface{}
}

type MyLinkedFaceAllTogether struct {
  MyAllLinkedFace []interface{}
}
type GlobalAllMine struct {
    Cdn *[]string
    Me []string
    Linked [][]string
    Playlist []string
}
type GlobalAll struct {
    Cdn *[]string
    Me *Me
    Linked *MyLinkedFaceAllTogether
}
func (g *GlobalAll) setGlobalsAll(cd []string,me *Me,linked *MyLinkedFaceAllTogether) {
    for _,v:=range cd{
        *g.Cdn=append(*g.Cdn,v)
    }
    *g.Me=*me
    *g.Linked=*linked
}
func (g *GlobalAll) getGlobalsAll() ([]string,Me,MyLinkedFaceAllTogether) {
    return *g.Cdn,*g.Me,*g.Linked
}
type Global struct {
    cdn *[]string
}
type Auth struct {
    Email string
    Password string
}
func (g *Global) setGlobals(m ... string) {
    for _,v:=range m{
        *g.cdn=append(*g.cdn,v)
    }
}
func (g *Global) getGlobals() []string {
    return *g.cdn
}
type Music struct {
  Content []interface{}
}

type MyLinkedFaces struct {
  All []interface{}
}

type Fauzi struct {
  Linkd *[][]string
}

func serveTemplate(w http.ResponseWriter, r *http.Request) {
  //http.Redirect(w, r, "/notatalldevelopped", 301)
    fp := filepath.Join("templat", "/login/loginMain.rita")
    cd := filepath.Join("templat", "/login/cdn.rita")
    st := filepath.Join("templat", "/login/loginStyle.rita")
    lh := filepath.Join("templat", "/login/loginHead.rita")
    lb := filepath.Join("templat", "/login/loginBody.rita")
    sn:= filepath.Join("templat", "/login/sanitize.rita")
    pl:= filepath.Join("templat", "/login/placeholderLog.rita")
    t, err:= template.ParseFiles(fp,cd,st,lh,lb,sn,pl)
    if err != nil {
        panic(err)
    }
    l:=make([]string,0)
    bn:=Global{cdn:&l}
    bn.setGlobals("https://cdnjs.cloudflare.com/ajax/libs/react/15.0.2/react.js","https://cdnjs.cloudflare.com/ajax/libs/react/15.0.2/react-dom.js","https://cdnjs.cloudflare.com/ajax/libs/babel-core/5.8.23/browser.min.js")
    mnb:=bn.getGlobals();
    w.Header().Set("Content-Type", "text/html")
    w.Header().Set("charset", "utf-8")
    w.Header().Set("Access-Control-Allow-Origin", "*")
    t.ExecuteTemplate(w, "login", mnb)
}
func serveSignupTemplate(w http.ResponseWriter,r *http.Request)  {
  fp := filepath.Join("templat", "/signup/signupMain.rita")
  cd := filepath.Join("templat", "/signup/cdn.rita")
  st := filepath.Join("templat", "/signup/signupStyle.rita")
  sh := filepath.Join("templat", "/signup/signupHead.rita")
  sb := filepath.Join("templat", "/signup/signupBody.rita")
  sn:= filepath.Join("templat", "/signup/verify.rita")
  pl:= filepath.Join("templat", "/signup/placeholder.rita")
  t, err:= template.ParseFiles(fp,cd,st,sh,sb,sn,pl)
  if err != nil {
      panic(err)
  }
  l:=make([]string,0)
  bn:=Global{cdn:&l}
  bn.setGlobals("https://cdnjs.cloudflare.com/ajax/libs/react/15.0.2/react.js","https://cdnjs.cloudflare.com/ajax/libs/react/15.0.2/react-dom.js","https://cdnjs.cloudflare.com/ajax/libs/babel-core/5.8.23/browser.min.js")
  mnb:=bn.getGlobals();
  w.Header().Set("Content-Type", "text/html")
  w.Header().Set("charset", "utf-8")
  w.Header().Set("Access-Control-Allow-Origin", "*")
  t.ExecuteTemplate(w, "signup", mnb)
}



func serveMainTemplate(w http.ResponseWriter,r *http.Request)  {
   http.Redirect(w, r, "/login", 301)
  /*fp := filepath.Join("templat", "/home/homeMain.rita")
  cd := filepath.Join("templat", "/home/cdn.rita")
  st := filepath.Join("templat", "/home/homeStyle.rita")
  sh := filepath.Join("templat", "/home/homeHead.rita")
  sb := filepath.Join("templat", "/home/homeBody.rita")
  sbOne := filepath.Join("templat", "/home/homeBodyPartOne.rita")
  sbThree := filepath.Join("templat", "/home/homeBodyPartThree.rita")
  pl:= filepath.Join("templat", "/home/placeholderHome.rita")
  chst := filepath.Join("templat", "/home/chatContainerStyle.rita")
  t, err:= template.ParseFiles(fp,cd,st,sh,pl,sb,sbOne,sbThree,chst)
  if err != nil {
      panic(err)
  }
  l:=make([]string,0)
  bn:=Global{cdn:&l}
  bn.setGlobals("https://cdnjs.cloudflare.com/ajax/libs/react/15.0.2/react.js","https://cdnjs.cloudflare.com/ajax/libs/react/15.0.2/react-dom.js","https://cdnjs.cloudflare.com/ajax/libs/babel-core/5.8.23/browser.min.js","https://cdnjs.cloudflare.com/ajax/libs/rxjs/5.4.0/Rx.min.js")
  mnb:=bn.getGlobals()
  w.Header().Set("Content-Type", "text/html")
  w.Header().Set("charset", "utf-8")
  w.Header().Set("Access-Control-Allow-Origin", "*")
  t.ExecuteTemplate(w, "home", mnb)*/



  //new

  /*c1 := make(chan []interface{})
  c2 := make(chan []interface{})

  go findMe(&c1,"9831296420")
  go findMyLinkedFaces(&c2,"9831296420")

  var ajay *Me
  var vijay *MyLinkedFaceAllTogether

  for i := 0; i < 2; i++ {
          select {
            case msg1 := <-c1:
                  ajay=&Me{MyDetails:msg1}
                  fmt.Println(ajay)
            case msg2 := <-c2:
                  vijay=&MyLinkedFaceAllTogether{MyAllLinkedFace:msg2}
                  fmt.Println(vijay)
            }
    }
    MeForMyself:=make([]string,0)
    AllInternal:=make([][]string,0)
    //BllInternal:=make([][]string,0)
    for _,r:=range vijay.MyAllLinkedFace{
      Internal:=make([]string,0)
      //TempHolder:=make([]string,0)
      for _,rl:=range r.([]interface{}){
        //mbcv:=strings.Join(r.([]interface),",")
        Internal=append(Internal,rl.(string))
      }
      //mbcv:=strings.Join(Internal,",")
      //TempHolder=append(TempHolder,mbcv)
      //mndir:="["+mbcv+"]"
      AllInternal=append(AllInternal,Internal)
      //fmt.Println(mbcv)
      //fmt.Println(AllInternal)
    }
    for _,rmn:=range ajay.MyDetails{
      for _,rj:=range rmn.([]interface{}){
        MeForMyself=append(MeForMyself,rj.(string))
      }
    }
    fmt.Println(AllInternal)
    //gst:=strings.Join(AllInternal,",")
    //gstr:=`[`+gst+`]`
    //fmt.Println(gst)
    //lk:=&Fauzi{Linkd:&AllInternal}
    //b, errm := json.Marshal(lk)
    b,errm:=JSONMarshal(lk)
	   if errm != nil {
		     panic(errm)
	      }
        funcMap := template.FuncMap{
            "decreased": func(i int) int {
                return i - 1
            },
        }
    l:=make([]string,0)
    l=append(l,"https://cdnjs.cloudflare.com/ajax/libs/react/15.0.2/react.js")
    l=append(l,"https://cdnjs.cloudflare.com/ajax/libs/react/15.0.2/react-dom.js")
    l=append(l,"https://cdnjs.cloudflare.com/ajax/libs/babel-core/5.8.23/browser.min.js")
    l=append(l,"https://cdnjs.cloudflare.com/ajax/libs/rxjs/5.4.0/Rx.min.js")
    bn:=GlobalAllMine{Cdn:&l,Me:MeForMyself,Linked:AllInternal}
    fp := filepath.Join("templat", "/home/homeMain.rita")
    cd := filepath.Join("templat", "/home/cdn.rita")
    st := filepath.Join("templat", "/home/homeStyle.rita")
    sh := filepath.Join("templat", "/home/homeHead.rita")
    sb := filepath.Join("templat", "/home/homeBody.rita")
    sbOne := filepath.Join("templat", "/home/homeBodyPartOne.rita")
    sbThree := filepath.Join("templat", "/home/homeBodyPartThree.rita")
    pl:= filepath.Join("templat", "/home/placeholderHome.rita")
    chst := filepath.Join("templat", "/home/chatContainerStyle.rita")
    mona := filepath.Join("templat", "/home/popupListener.rita")
    t:= template.Must(template.New("home").Funcs(funcMap).ParseFiles(fp,cd,st,sh,pl,sb,sbOne,sbThree,chst,mona))
    if err != nil {
        panic(err)
    }
    w.Header().Set("Content-Type", "text/html")
    w.Header().Set("charset", "utf-8")
    w.Header().Set("Access-Control-Allow-Origin", "*")
    t.ExecuteTemplate(w, "home", bn)*/
}

func AuthinticateMeYaar(a *chan []interface{},id string){
  s1:="Match (ee:Rita) where ee.email ='"
  b:=id
  s2:="' with ee optional Match (ee)-[:HAS_MANDATORY_DP]->(gg:ProfilePic) return collect([ee.password]) as all"
  s12 := fmt.Sprint(s1,b,s2)
  fmt.Println(s12)
  driver := bolt.NewDriver()
    conn, err := driver.OpenNeo("bolt://rita:b.PuhuqVThYfCn.fvurl1e25g7fzyCI@hobby-panhpmpgjildgbkepcdcklol.dbs.graphenedb.com:24786?tls=true")
    defer conn.Close()
    dur,_:=time.ParseDuration("600s")
    conn.SetTimeout(dur)
    if err != nil {
      panic(err)
    }
    data, _, _, _ := conn.QueryNeoAll(s12, nil)
    //fmt.Println("hooo")
    //fmt.Printf("COLUMNS: %#v\n", rowsMetadata["fields"].([]interface{}))  // COLUMNS: n.foo,n.bar
    //fmt.Printf("FIELDS: %s\n", data[0][0].([]interface{})) // FIELDS: 1 2.2
        *a <- data[0][0].([]interface{})
}

type PasswordMy struct {
  Password []interface{}
}




func serveAuth(w http.ResponseWriter, r *http.Request) {

  r.ParseForm()
  e:=r.Form.Get("_em_o_phn")
  p:=r.Form.Get("_pwd")
  fmt.Println(e)
  fmt.Println(p)
  c6 := make(chan []interface{})
  go AuthinticateMeYaar(&c6,e)
  msg6 := <-c6
    c:=&PasswordMy{Password:msg6}

    for _,ght:=range c.Password{
      for _,dnt:=range ght.([]interface{}){
        if p==dnt.(string) {
          dilTodKe:=generateHashAndReplicateToDbPratibhaPlease(&e)
          c1 := make(chan []interface{})
          c2 := make(chan []interface{})
          c3 := make(chan []interface{})
          go findMe(&c1,e)
          go findMyLinkedFaces(&c2,e)
          go findMyPlaylist(&c3,e)
          var ajay *Me
          var vijay *MyLinkedFaceAllTogether
          var sanjay *AuthedPlaylist
          for i := 0; i < 3; i++ {
                  select {
                    case msg1 := <-c1:
                          ajay=&Me{MyDetails:msg1}
                          fmt.Println(ajay)
                    case msg2 := <-c2:
                          vijay=&MyLinkedFaceAllTogether{MyAllLinkedFace:msg2}
                          fmt.Println(vijay)
                    case msg3 := <-c3:
                          sanjay=&AuthedPlaylist{PlaylistOf:msg3}
                          fmt.Println(sanjay)
                          fmt.Println("hi sanjay")
                          //fmt.Println(sanjay.PlaylistOf)
                    }
            }
            MeForMyself:=make([]string,0)
            MyGettingPlaylist:=make([]string,0)
            AllInternal:=make([][]string,0)
            for _,r:=range vijay.MyAllLinkedFace{
              Internal:=make([]string,0)
              for _,rl:=range r.([]interface{}){
                Internal=append(Internal,rl.(string))
              }
              AllInternal=append(AllInternal,Internal)
            }
            for _,rmn:=range ajay.MyDetails{
              for _,rj:=range rmn.([]interface{}){
                MeForMyself=append(MeForMyself,rj.(string))
              }
            }
            fmt.Println("hello")
            fmt.Println(sanjay)
            for _,gnm:=range sanjay.PlaylistOf{
                MyGettingPlaylist=append(MyGettingPlaylist,gnm.(string))
            }

            fmt.Println(AllInternal)
                funcMap := template.FuncMap{
                    "decreased": func(i int) int {
                        return i - 1
                    },
                }


            MeForMyself[1]=dilTodKe[0]
            MeForMyself=append(MeForMyself,dilTodKe[1])


            l:=make([]string,0)
            l=append(l,"https://cdnjs.cloudflare.com/ajax/libs/react/15.0.2/react.js")
            l=append(l,"https://cdnjs.cloudflare.com/ajax/libs/react/15.0.2/react-dom.js")
            l=append(l,"https://cdnjs.cloudflare.com/ajax/libs/babel-core/5.8.23/browser.min.js")
            l=append(l,"https://cdnjs.cloudflare.com/ajax/libs/rxjs/5.4.0/Rx.min.js")
            l=append(l,"https://code.jquery.com/jquery-3.2.1.min.js")
            l=append(l,"https://cdnjs.cloudflare.com/ajax/libs/materialize/0.100.1/js/materialize.min.js")
            //l=append(l,"https://cdnjs.cloudflare.com/ajax/libs/cloudinary-core/2.3.0/cloudinary-core.js")
            bn:=GlobalAllMine{Cdn:&l,Me:MeForMyself,Linked:AllInternal,Playlist:MyGettingPlaylist}
            fp := filepath.Join("templat", "/home/homeMain.rita")
            cd := filepath.Join("templat", "/home/cdn.rita")
            st := filepath.Join("templat", "/home/homeStyle.rita")
            sh := filepath.Join("templat", "/home/homeHead.rita")
            sb := filepath.Join("templat", "/home/homeBody.rita")
            sbOne := filepath.Join("templat", "/home/homeBodyPartOne.rita")
            sbThree := filepath.Join("templat", "/home/homeBodyPartThree.rita")
            pl:= filepath.Join("templat", "/home/placeholderHome.rita")
            chst := filepath.Join("templat", "/home/chatContainerStyle.rita")
            mona := filepath.Join("templat", "/home/popupListener.rita")
            t:= template.Must(template.New("home").Funcs(funcMap).ParseFiles(fp,cd,st,sh,pl,sb,sbOne,sbThree,chst,mona))
            w.Header().Set("Content-Type", "text/html")
            w.Header().Set("charset", "utf-8")
            w.Header().Set("Access-Control-Allow-Origin", "*")
            t.ExecuteTemplate(w, "home", bn)

        }else {
          fmt.Println("err")
        }
      }
    }

}

func serveAnyTemplate(w http.ResponseWriter, r *http.Request) {

  vars := mux.Vars(r)
  mkiloj:=vars["query"]
  hashedData:=vars["hash"]
  pKey:=vars["Pkey"]
  fmt.Println(mkiloj)

  insaan:=hashedData
  mainInsaan:=pKey
  fmt.Println("aaj se tera")
  fmt.Println(mainInsaan)
  mainInsaanSaasur, _ := strconv.ParseInt(mainInsaan, 10, 64)
  rest:=authorizeThisFace(&insaan,&mainInsaanSaasur)

  if rest!="ERROR" {
    yt := make(chan []byte)
      //seb:="crime patrol"
      i := 20
  var maxResults int64
              query:= mkiloj
              maxResults=int64(i)

    go youTubeVideo(&yt,&query,&maxResults);

    msg:=<-yt
    fmt.Println("byte")
    fmt.Println(string(msg))
    w.Write(msg)
  }else{
    m:=HttpStatusForbiddenUnAuthorize{
    HttpStatus:403,HttpText:"Authentication Failed"}
    resp, _ := json.Marshal(m)
    w.WriteHeader(http.StatusForbidden)
    w.Header().Set("Content-Type", "application/json")
    w.Header().Set("charset", "utf-8")
    w.Header().Set("Access-Control-Allow-Origin", "*")
    w.Write([]byte(resp))
  }

}

func serveMyLinedFaceTemplate(w http.ResponseWriter, r *http.Request) {
    s1:="Match (ee:Rita)-[:LINKED]-(ff:Rita) where ee.email ='"
    b:="9831296420"
    s2:="' with ff optional Match (ff)-[:HAS_MANDATORY_DP]->(gg:ProfilePic) return collect([ff.name,ff.email,gg.title]) as all"
    s12 := fmt.Sprint(s1,b,s2)
    fmt.Println(s12)
    driver := bolt.NewDriver()
    	conn, err := driver.OpenNeo("bolt://rita:b.PuhuqVThYfCn.fvurl1e25g7fzyCI@hobby-panhpmpgjildgbkepcdcklol.dbs.graphenedb.com:24786?tls=true")
    	defer conn.Close()
      dur,_:=time.ParseDuration("600s")
      conn.SetTimeout(dur)
    	if err != nil {
    		panic(err)
    	}
    	data, rowsMetadata, _, _ := conn.QueryNeoAll(s12, nil)
    	fmt.Println("hooo")
    	fmt.Printf("COLUMNS: %#v\n", rowsMetadata["fields"].([]interface{}))  // COLUMNS: n.foo,n.bar
    	fmt.Printf("FIELDS: %s\n", data[0][0].([]interface{})) // FIELDS: 1 2.2
      m:=MyLinkedFaces{
        All:data[0][0].([]interface{})}
        resp, _ := json.Marshal(m)
        fmt.Println(string(resp))
      w.Header().Set("Content-Type", "application/json")
      w.Header().Set("charset", "utf-8")
      w.Header().Set("Access-Control-Allow-Origin", "*")
      w.Write([]byte(resp))
}
func findMe(a *chan []interface{},id string){
  s1:="Match (ee:Rita) where ee.email ='"
  b:=id
  s2:="' with ee optional Match (ee)-[:HAS_MANDATORY_DP]->(gg:ProfilePic) return collect([ee.name,ee.email,gg.title,ee.gender]) as all"
  s12 := fmt.Sprint(s1,b,s2)
  fmt.Println(s12)
  driver := bolt.NewDriver()
    conn, err := driver.OpenNeo("bolt://rita:b.PuhuqVThYfCn.fvurl1e25g7fzyCI@hobby-panhpmpgjildgbkepcdcklol.dbs.graphenedb.com:24786?tls=true")
    defer conn.Close()
    dur,_:=time.ParseDuration("600s")
    conn.SetTimeout(dur)
    if err != nil {
      panic(err)
    }
    data, _, _, _ := conn.QueryNeoAll(s12, nil)
    //fmt.Println("hooo")
    //fmt.Printf("COLUMNS: %#v\n", rowsMetadata["fields"].([]interface{}))  // COLUMNS: n.foo,n.bar
    //fmt.Printf("FIELDS: %s\n", data[0][0].([]interface{})) // FIELDS: 1 2.2
        *a <- data[0][0].([]interface{})
}
func findMyLinkedFaces(a *chan []interface{},id string){

  s1:="Match (ee:Rita)-[r:LINKED]->(ff:Rita) where ee.email ='"
  b:=id
  s2:="' with ee,ff,r match (ff)-[:HAS_MANDATORY_DP]->(gg:ProfilePic) return collect([ff.name,ff.email,gg.title,r.tlfSocial,r.tlfProfessional,r.tlfConsumer,r.tlfReader,r.tlfDater,r.tlfDonater,r.tlfSinner]) as all union Match (ee:Rita)-[r:LINKED]->(ff:Rita) where ff.email ='"
  s34567:=id
  s554433:="' with ee,ff,r  match (ee)-[:HAS_MANDATORY_DP]->(gg:ProfilePic) return collect([ee.name,ee.email,gg.title,r.tlfSocial,r.tlfProfessional,r.tlfConsumer,r.tlfReader,r.tlfDater,r.tlfDonater,r.tlfSinner]) as all"
  s12 := fmt.Sprint(s1,b,s2,s34567,s554433)
  fmt.Println(s12)
  driver := bolt.NewDriver()
    conn, err := driver.OpenNeo("bolt://rita:b.PuhuqVThYfCn.fvurl1e25g7fzyCI@hobby-panhpmpgjildgbkepcdcklol.dbs.graphenedb.com:24786?tls=true")
    defer conn.Close()
    dur,_:=time.ParseDuration("600s")
    conn.SetTimeout(dur)
    if err != nil {
      panic(err)
      fmt.Println("lollllllllll")
    }
    data, _, _, _ := conn.QueryNeoAll(s12, nil)
    //fmt.Println("hooo")
    //fmt.Printf("COLUMNS: %#v\n", rowsMetadata["fields"].([]interface{}))  // COLUMNS: n.foo,n.bar
    //fmt.Printf("FIELDS: %s\n", data[0][0].([]interface{})) // FIELDS: 1 2.2

    var sealApp []interface{}
    if len(data)==0 || len(data)==1 {
      sealApp = data[0][0].([]interface{})
    }else{
      sealApp = append(data[0][0].([]interface{}), data[1][0].([]interface{})...)
    }
    *a <- sealApp

        //*a <- data[0][0].([]interface{})
}

func findMyPlaylist(a *chan []interface{},id string)  {
  s1:="MATCH (ee:Rita)-[:HAS_YOUTUBE_PLAYLIST]->(ll) where ee.email='"
  b:=id
  s2:="' return collect(ll.title) as all"
  s12 := fmt.Sprint(s1,b,s2)
  fmt.Println(s12)
  driver := bolt.NewDriver()
    conn, err := driver.OpenNeo("bolt://rita:b.PuhuqVThYfCn.fvurl1e25g7fzyCI@hobby-panhpmpgjildgbkepcdcklol.dbs.graphenedb.com:24786?tls=true")
    defer conn.Close()
    dur,_:=time.ParseDuration("600s")
    conn.SetTimeout(dur)
    if err != nil {
      panic(err)
    }
    data, _, _, _ := conn.QueryNeoAll(s12, nil)
    //fmt.Println("hooo")
    //fmt.Printf("COLUMNS: %#v\n", rowsMetadata["fields"].([]interface{}))  // COLUMNS: n.foo,n.bar
    //fmt.Printf("FIELDS: %s\n", data[0][0].([]interface{})) // FIELDS: 1 2.2
        *a <- data[0][0].([]interface{})
}

func serveDataTemplate(w http.ResponseWriter, r *http.Request) {
  c1 := make(chan []interface{})
  c2 := make(chan []interface{})

  go findMe(&c1,"9831296420")
  go findMyLinkedFaces(&c2,"9831296420")

  var ajay *Me
  var vijay *MyLinkedFaceAllTogether

  for i := 0; i < 2; i++ {
        	select {
        		case msg1 := <-c1:
                  ajay=&Me{MyDetails:msg1}
                  fmt.Println(ajay)
        		case msg2 := <-c2:
            			vijay=&MyLinkedFaceAllTogether{MyAllLinkedFace:msg2}
                  fmt.Println(vijay)
        		}
    }

    l:=make([]string,0)
    l=append(l,"https://cdnjs.cloudflare.com/ajax/libs/react/15.0.2/react.js")
    l=append(l,"https://cdnjs.cloudflare.com/ajax/libs/react/15.0.2/react-dom.js")
    l=append(l,"https://cdnjs.cloudflare.com/ajax/libs/babel-core/5.8.23/browser.min.js")
    l=append(l,"https://cdnjs.cloudflare.com/ajax/libs/rxjs/5.4.0/Rx.min.js")
    bn:=GlobalAll{Cdn:&l,Me:ajay,Linked:vijay}
    fmt.Println("hello writing my")
    fmt.Println(bn.Cdn)
    fmt.Println("hello writing my")
    //bn.setGlobals("https://cdnjs.cloudflare.com/ajax/libs/react/15.0.2/react.js","https://cdnjs.cloudflare.com/ajax/libs/react/15.0.2/react-dom.js","https://cdnjs.cloudflare.com/ajax/libs/babel-core/5.8.23/browser.min.js","https://cdnjs.cloudflare.com/ajax/libs/rxjs/5.4.0/Rx.min.js")
    a,b,c:=bn.getGlobalsAll()
    fmt.Println(a)
    fmt.Println(b)
    fmt.Println(c)
}


//youtube more videos Handler

type MoreYouTubeVideoHealper struct {
  Id string
  Title rune
}


type MoreYouTubeVideo struct {
  SkipCount int32
  Current string
  Me string
  Pkey string
}


func serveMoreVideosYoutube(w http.ResponseWriter, r *http.Request){
  var jsx MoreYouTubeVideo
  fmt.Println("kkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkk")
  token := r.FormValue("data")
//u, _ := url.Parse(token)
  fmt.Println(token)
  json.Unmarshal([]byte(token),&jsx)


  insaan:=jsx.Me
  mainInsaan:=jsx.Pkey
  fmt.Println("aaj se tera")
  fmt.Println(mainInsaan)
  mainInsaanSaasur, _ := strconv.ParseInt(mainInsaan, 10, 64)
  rest:=authorizeThisFace(&insaan,&mainInsaanSaasur)


  if rest!="ERROR" {
    yt := make(chan []byte)
      //seb:="crime patrol"

      i := 50
        mkiloj:=jsx.Current
      var maxResults int64
                query:= mkiloj
                maxResults=int64(i)
    fmt.Println(maxResults)
      go youTubeVideo(&yt,&query,&maxResults);

      msg:=<-yt
      fmt.Println("byte")
      fmt.Println(string(msg))
      w.Write(msg)



    fmt.Println(jsx)
  }else{
    m:=HttpStatusForbiddenUnAuthorize{
    HttpStatus:403,HttpText:"Authentication Failed"}
    resp, _ := json.Marshal(m)
    w.WriteHeader(http.StatusForbidden)
    w.Header().Set("Content-Type", "application/json")
    w.Header().Set("charset", "utf-8")
    w.Header().Set("Access-Control-Allow-Origin", "*")
    w.Write([]byte(resp))
  }
}









// end youtube more videos handler




//youtube search










const developerKey = "AIzaSyCDsAnh9YH2s0qKjIhsyC5K0Xou3lfavpM"


type YouTube struct {
  Id string `json:"id"`
  Title string `json:"title"`
}
type AllYoutube struct {
  All []YouTube `json:"all"`
}

func youTubeVideo(a *chan []byte,c *string,e *int64){


          client := &http.Client{
                  Transport: &transport.APIKey{Key: developerKey},
          }

          service, err := youtube.New(client)
          if err != nil {
                  panic(err)
          }

          // Make the API call to YouTube.
          call := service.Search.List("id,snippet").
                  Q(*c).
                  MaxResults(*e)
          response, err := call.Do()
          if err != nil {
                  panic(err)
          }

          // Group video, channel, and playlist results in separate lists.
          videos := make(map[string]string)
          channels := make(map[string]string)
          playlists := make(map[string]string)
          sArr:=make([]YouTube,0)
          // Iterate through each item and add it to the correct list.
          for _, item := range response.Items {
                  switch item.Id.Kind {
                  case "youtube#video":
                            c:=YouTube{Id:item.Id.VideoId,Title:item.Snippet.Title}
                            sArr=append(sArr,c)
                            //jsonString, _ := json.Marshal(datas)
                            //fmt.Println(jsonString)
                          //videos[item.Id.VideoId] = item.Snippet.Title
                  case "youtube#channel":
                          channels[item.Id.ChannelId] = item.Snippet.Title
                  case "youtube#playlist":
                          playlists[item.Id.PlaylistId] = item.Snippet.Title
                  }
          }
          //*a<-videos
          datas:=AllYoutube{All:sArr}
          jsonString, _ := json.Marshal(datas)


          *a <- jsonString
          //fmt.Println(videos)
          printIDs("Videos", videos)

}

//end youtube search

//youtube load more videos


func youTubeMoreVideo(a *chan []byte,c *string,e *int64){


          client := &http.Client{
                  Transport: &transport.APIKey{Key: developerKey},
          }

          service, err := youtube.New(client)
          if err != nil {
                  panic(err)
          }

          // Make the API call to YouTube.
          call := service.Search.List("id,snippet").
                  Q(*c).
                  MaxResults(*e)
          response, err := call.Do()
          if err != nil {
                  panic(err)
          }

          // Group video, channel, and playlist results in separate lists.
          videos := make(map[string]string)
          channels := make(map[string]string)
          playlists := make(map[string]string)
          sArr:=make([]YouTube,0)
          // Iterate through each item and add it to the correct list.
          for _, item := range response.Items {
                  switch item.Id.Kind {
                  case "youtube#video":
                            c:=YouTube{Id:item.Id.VideoId,Title:item.Snippet.Title}
                            sArr=append(sArr,c)
                            //jsonString, _ := json.Marshal(datas)
                            //fmt.Println(jsonString)
                          //videos[item.Id.VideoId] = item.Snippet.Title
                  case "youtube#channel":
                          channels[item.Id.ChannelId] = item.Snippet.Title
                  case "youtube#playlist":
                          playlists[item.Id.PlaylistId] = item.Snippet.Title
                  }
          }
          //*a<-videos
          datas:=AllYoutube{All:sArr}
          jsonString, _ := json.Marshal(datas)


          *a <- jsonString
          //fmt.Println(videos)
          printIDs("Videos", videos)

}








//end youtube load more videos

//times of india

func timesOfIndia(w http.ResponseWriter, r *http.Request){
  vars := mux.Vars(r)
  mkiloj:=vars["newsType"]
  fmt.Println(mkiloj)
  resp, err := http.Get("https://newsapi.org/v1/articles?source="+mkiloj+"&sortBy=top&apiKey=9b9b05565f3e424fa70f46c06b6d10c8")
if err != nil {
	panic(err)
}
defer resp.Body.Close()
body, err := ioutil.ReadAll(resp.Body)
if err != nil {
	panic(err)
}
//fmt.Println(string(body))
w.Write(body)
}



//end times of india


func redirectHandler(w http.ResponseWriter, r *http.Request){
  newUrl:="http://localhost:5000/signup"
  http.Redirect(w, r, newUrl, 301)
  fmt.Println("biccccccccccccccccccccccccccccccccccccccccccccccccccccccccccc")
}

type PlaylistMarsheler struct {
  Me string
  Playlist string
  Pkey string
}

type PlaylistResponder struct {
  All []interface{} `json:all`
}

func newPlaylistCreationHandler(w http.ResponseWriter, r *http.Request)  {
  var playlistMe *PlaylistMarsheler
  token := r.FormValue("data")
  json.Unmarshal([]byte(token), &playlistMe)
  insaan:=playlistMe.Me
  mainInsaan:=playlistMe.Pkey
  fmt.Println("aaj se tera")
  fmt.Println(mainInsaan)
  mainInsaanSaasur, _ := strconv.ParseInt(mainInsaan, 10, 64)
  rest:=authorizeThisFace(&insaan,&mainInsaanSaasur)
  if rest!="ERROR" {
    s1:="Match (ee:Rita) where ee.email ='"
    b:=rest
    s2:="' with ee MERGE (gg:YoutubePlaylist{title:'"
    fh:=playlistMe.Playlist
    s3:="'}) with ee,gg MERGE (ee)-[:HAS_YOUTUBE_PLAYLIST]->(gg) with ee MATCH (ee)-[:HAS_YOUTUBE_PLAYLIST]->(ll) return collect(ll.title) as all"
    s12 := fmt.Sprint(s1,b,s2,fh,s3)
    fmt.Println(s12)
    driver := bolt.NewDriver()
      conn, err := driver.OpenNeo("bolt://rita:b.PuhuqVThYfCn.fvurl1e25g7fzyCI@hobby-panhpmpgjildgbkepcdcklol.dbs.graphenedb.com:24786?tls=true")
      defer conn.Close()
      if err != nil {
        panic(err)
      }
      data, _, _, _ := conn.QueryNeoAll(s12, nil)
      //fmt.Println("hooo")
      //fmt.Printf("COLUMNS: %#v\n", rowsMetadata["fields"].([]interface{}))  // COLUMNS: n.foo,n.bar
      //fmt.Printf("FIELDS: %s\n", data[0][0].([]interface{})) // FIELDS: 1 2.2
      m:=PlaylistResponder{
        All:data[0][0].([]interface{})}
        resp, _ := json.Marshal(m)
        fmt.Println(string(resp))
      w.Header().Set("Content-Type", "application/json")
      w.Header().Set("charset", "utf-8")
      w.Header().Set("Access-Control-Allow-Origin", "*")
      w.Write([]byte(resp))
  }else{
    m:=HttpStatusForbiddenUnAuthorize{
    HttpStatus:403,HttpText:"Authentication Failed"}
    resp, _ := json.Marshal(m)
    w.WriteHeader(http.StatusForbidden)
    w.Header().Set("Content-Type", "application/json")
    w.Header().Set("charset", "utf-8")
    w.Header().Set("Access-Control-Allow-Origin", "*")
    w.Write([]byte(resp))
  }
}

func CreateSpaceForMe(a *chan []interface{},name string,email string,password string,gender string)  {
  s1:="CREATE (ee:Rita{name:'"
  cvIn:=name
  cv1:="',email:'"
  cvIn2:=email
  cv2:="',password:'"
  cvIn3:=password
  cv3:="',gender:'"
  cvIn4:=gender
  cv4:="',interest:'',dob:''}),(ff:ProfilePic{title:''}),(ee)-[:HAS_MANDATORY_DP]->(ff) return collect(ee.email)"
  s12 := fmt.Sprint(s1,cvIn,cv1,cvIn2,cv2,cvIn3,cv3,cvIn4,cv4)
  fmt.Println(s12)
  driver := bolt.NewDriver()
    conn, err := driver.OpenNeo("bolt://rita:b.PuhuqVThYfCn.fvurl1e25g7fzyCI@hobby-panhpmpgjildgbkepcdcklol.dbs.graphenedb.com:24786?tls=true")
    defer conn.Close()
    dur,_:=time.ParseDuration("600s")
    conn.SetTimeout(dur)
    if err != nil {
      panic(err)
    }
    data, rowsMetadata, _, errx := conn.QueryNeoAll(s12, nil)
    if errx !=nil {
      var interfaceSlice []interface{} = make([]interface{}, 0)
      *a <- interfaceSlice
    }else {
      *a <- data[0][0].([]interface{})
    }
    fmt.Println("hooo")
    //fmt.Println(f)
    //fmt.Println(errx)
    fmt.Println(data)
    fmt.Println(rowsMetadata)
    //fmt.Printf("COLUMNS: %#v\n", rowsMetadata["fields"].([]interface{}))  // COLUMNS: n.foo,n.bar
    //fmt.Printf("FIELDS: %s\n", data) // FIELDS: 1 2.2
        //*a <- data[0][0].([]interface{})
}

func serveAuthAndSignUp(w http.ResponseWriter, r *http.Request)  {
    r.ParseForm()
    name:=r.Form.Get("name")
    email:=r.Form.Get("_em_o_phn")
    password:=r.Form.Get("_pwd")
    repassword:=r.Form.Get("_pwd_re")
    gender:=r.Form.Get("gender")
    fmt.Println(repassword)
    chnlSign := make(chan []interface{})
    go CreateSpaceForMe(&chnlSign,name,email,password,gender)
    msgSign := <-chnlSign
    if len(msgSign)==0{
      fmt.Println("User already exist")
    }else{
      fmt.Println("New User")
      fmt.Println(msgSign[0])
      e:=msgSign[0].(string)

      // main page after validation begining












      c1 := make(chan []interface{})
      c2 := make(chan []interface{})
      c3 := make(chan []interface{})
      go findMe(&c1,e)
      go findMyLinkedFaces(&c2,e)
      go findMyPlaylist(&c3,e)
      var ajay *Me
      var vijay *MyLinkedFaceAllTogether
      var sanjay *AuthedPlaylist
      for i := 0; i < 3; i++ {
              select {
                case msg1 := <-c1:
                      ajay=&Me{MyDetails:msg1}
                      fmt.Println(ajay)
                case msg2 := <-c2:
                      vijay=&MyLinkedFaceAllTogether{MyAllLinkedFace:msg2}
                      fmt.Println(vijay)
                case msg3 := <-c3:
                      sanjay=&AuthedPlaylist{PlaylistOf:msg3}
                      fmt.Println(sanjay)
                      fmt.Println("hi sanjay")
                      //fmt.Println(sanjay.PlaylistOf)
                }
        }
        MeForMyself:=make([]string,0)
        MyGettingPlaylist:=make([]string,0)
        AllInternal:=make([][]string,0)
        for _,r:=range vijay.MyAllLinkedFace{
          Internal:=make([]string,0)
          for _,rl:=range r.([]interface{}){
            Internal=append(Internal,rl.(string))
          }
          AllInternal=append(AllInternal,Internal)
        }
        for _,rmn:=range ajay.MyDetails{
          for _,rj:=range rmn.([]interface{}){
            MeForMyself=append(MeForMyself,rj.(string))
          }
        }
        fmt.Println("hello")
        fmt.Println(sanjay)
        for _,gnm:=range sanjay.PlaylistOf{
            //fmt.Println("oooooooooooooooooooooooooooooooooooo")
            //fmt.Println(gnm)
            MyGettingPlaylist=append(MyGettingPlaylist,gnm.(string))
        }

        fmt.Println(AllInternal)
            funcMap := template.FuncMap{
                "decreased": func(i int) int {
                    return i - 1
                },
            }
        l:=make([]string,0)
        l=append(l,"https://cdnjs.cloudflare.com/ajax/libs/react/15.0.2/react.js")
        l=append(l,"https://cdnjs.cloudflare.com/ajax/libs/react/15.0.2/react-dom.js")
        l=append(l,"https://cdnjs.cloudflare.com/ajax/libs/babel-core/5.8.23/browser.min.js")
        l=append(l,"https://cdnjs.cloudflare.com/ajax/libs/rxjs/5.4.0/Rx.min.js")
        bn:=GlobalAllMine{Cdn:&l,Me:MeForMyself,Linked:AllInternal,Playlist:MyGettingPlaylist}
        fp := filepath.Join("templat", "/home/homeMain.rita")
        cd := filepath.Join("templat", "/home/cdn.rita")
        st := filepath.Join("templat", "/home/homeStyle.rita")
        sh := filepath.Join("templat", "/home/homeHead.rita")
        sb := filepath.Join("templat", "/home/homeBody.rita")
        sbOne := filepath.Join("templat", "/home/homeBodyPartOne.rita")
        sbThree := filepath.Join("templat", "/home/homeBodyPartThree.rita")
        pl:= filepath.Join("templat", "/home/placeholderHome.rita")
        chst := filepath.Join("templat", "/home/chatContainerStyle.rita")
        mona := filepath.Join("templat", "/home/popupListener.rita")
        t:= template.Must(template.New("home").Funcs(funcMap).ParseFiles(fp,cd,st,sh,pl,sb,sbOne,sbThree,chst,mona))
        w.Header().Set("Content-Type", "text/html")
        w.Header().Set("charset", "utf-8")
        w.Header().Set("Access-Control-Allow-Origin", "*")
        t.ExecuteTemplate(w, "home", bn)















      // end of main page


    }
}

func uploadAndProcessMyNewDp(w http.ResponseWriter, r *http.Request)  {
  //var results []string
  if r.Method == "POST" {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Error reading request body",
				http.StatusInternalServerError)
		}
    fmt.Println(body)
    //ctx := context.Background()
    //ctx = NewContext(ctx, "cloudinary://864654217542164:fdQqxrCeKl_OJdwR84Bw9LhuUhM@hnruvsvqz")

	    //cloudinary.UploadStaticImage(ctx, "ajay", bytes.NewBuffer(body))
		//results = append(results, string(body))
    //fmt.Println(results)
		//fmt.Fprint(w, "POST done")
	} else {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}
}

type Poster struct {
  Content string
  Who string
  Attachment string
  When string
  Type string
  Pkey string
}


type PostedAllData struct {
  All []interface{} `json:all`
}


func saveMyPostWithAttachmentHandler(w http.ResponseWriter, r *http.Request)  {
  //token := r.FormValue("data")
  //fmt.Println(token)
  var posterMe *Poster
  body, _ := ioutil.ReadAll(r.Body)
  myFirstEncodedString,_ := url.QueryUnescape(string(body))
  fmt.Println(myFirstEncodedString)
  json.Unmarshal([]byte(myFirstEncodedString), &posterMe)

  insaan:=posterMe.Who
  mainInsaan:=posterMe.Pkey
  fmt.Println("aaj se tera")
  fmt.Println(mainInsaan)
  mainInsaanSaasur, _ := strconv.ParseInt(mainInsaan, 10, 64)
  rest:=authorizeThisFace(&insaan,&mainInsaanSaasur)
  //fmt.Println(posterMe.Content)
  //fmt.Println(string(body))
  //r.ParseForm()
  //fmt.Println(r.Form)
  //e:=r.Form.Get("content")
  //p:=r.Form.Get("_pwd")
  //fmt.Println(e)
  //var posterMe *Poster
  //json.Unmarshal([]byte(token), &posterMe)
  if rest!="ERROR" {
    s1:="CREATE (ee:PostedData{who:'"
    b1:=rest
    s11:="',attachment:'"
    b2:=posterMe.Attachment
    s12:="',when:'"
    b3:=posterMe.When
    s13:="',content:'"
    b4:=posterMe.Content
    slx1:="',type:'"
    slx2:=posterMe.Type
    s14:="',rate:'0.0',ignore:'0',action:'0'}) WITH ee MATCH (ff:Rita) where ff.email ='"
    b5:=rest
    s15:="' with ee,ff MERGE (ff)-[:POSTED_THIS]-(ee) return collect([ee.who,ee.attachment,ee.when,ee.content,ee.type])"
    s17 := fmt.Sprint(s1,b1,s11,b2,s12,b3,s13,b4,slx1,slx2,s14,b5,s15)
    fmt.Println(s17)
    driver := bolt.NewDriver()
      conn, err := driver.OpenNeo("bolt://rita:b.PuhuqVThYfCn.fvurl1e25g7fzyCI@hobby-panhpmpgjildgbkepcdcklol.dbs.graphenedb.com:24786?tls=true")
      defer conn.Close()
      dur,_:=time.ParseDuration("600s")
      conn.SetTimeout(dur)
      if err != nil {
        panic(err)
      }
      data, _, _, _ := conn.QueryNeoAll(s17, nil)
      //fmt.Println("hooo")
      //fmt.Printf("COLUMNS: %#v\n", rowsMetadata["fields"].([]interface{}))  // COLUMNS: n.foo,n.bar
      fmt.Printf("FIELDS: %s\n", data) // FIELDS: 1 2.2
      m:=PostedAllData{
        All:data[0][0].([]interface{})}
        resp, _ := json.Marshal(m)
        fmt.Println(string(resp))
      w.Header().Set("Content-Type", "application/json")
      w.Header().Set("charset", "utf-8")
      w.Header().Set("Access-Control-Allow-Origin", "*")
      w.Write([]byte(resp))
  }else{
    m:=HttpStatusForbiddenUnAuthorize{
    HttpStatus:403,HttpText:"Authentication Failed"}
    resp, _ := json.Marshal(m)
    w.WriteHeader(http.StatusForbidden)
    w.Header().Set("Content-Type", "application/json")
    w.Header().Set("charset", "utf-8")
    w.Header().Set("Access-Control-Allow-Origin", "*")
    w.Write([]byte(resp))
  }
}


type GetMyAllPostPlease struct {
  Me string
  PostType string
  Pkey string
}

type AllMyNotifiedPost struct {
  All []interface{}
}

func getMyAllRelatedPostBillaHandler(w http.ResponseWriter, r *http.Request)  {
  token := r.FormValue("data")
  var allMyPosts *GetMyAllPostPlease
  json.Unmarshal([]byte(token), &allMyPosts)

  insaan:=allMyPosts.Me
  mainInsaan:=allMyPosts.Pkey
  fmt.Println("aaj se tera")
  fmt.Println(mainInsaan)
  mainInsaanSaasur, _ := strconv.ParseInt(mainInsaan, 10, 64)
  rest:=authorizeThisFace(&insaan,&mainInsaanSaasur)

  if rest!="ERROR" {
    s1:="Match (ee:Rita)-[r:LINKED]-(ff:Rita) where ee.email ='"
    b1:=rest
    s2:="' with ee,ff,r match (ee)-[:POSTED_THIS]->(gg:PostedData) where gg.type='"
    sMubarak:=allMyPosts.PostType
    dardMubarak:="' return collect(distinct [ee.name,gg.who,gg.attachment,gg.when,gg.content,id(gg),gg.action,gg.rate,gg.ignore]) as all union Match (ee:Rita)-[r:LINKED]-(ff:Rita) where ff.email ='"
    b2:=rest
    s3:="' with ee,ff,r  match (ee)-[:POSTED_THIS]->(gg:PostedData) where gg.type='"
    sHangama:=allMyPosts.PostType
    sQbarpa:="' return collect(distinct [ee.name,gg.who,gg.attachment,gg.when,gg.content,id(gg),gg.action,gg.rate,gg.ignore]) as all"
    s17 := fmt.Sprint(s1,b1,s2,sMubarak,dardMubarak,b2,s3,sHangama,sQbarpa)
    fmt.Println(s17)
    driver := bolt.NewDriver()
      conn, err := driver.OpenNeo("bolt://rita:b.PuhuqVThYfCn.fvurl1e25g7fzyCI@hobby-panhpmpgjildgbkepcdcklol.dbs.graphenedb.com:24786?tls=true")
      defer conn.Close()
      dur,_:=time.ParseDuration("600s")
      conn.SetTimeout(dur)
      if err != nil {
        panic(err)
      }
      data, _, _, _ := conn.QueryNeoAll(s17, nil)
      //fmt.Println("hooo")
      //fmt.Printf("COLUMNS: %#v\n", rowsMetadata["fields"].([]interface{}))  // COLUMNS: n.foo,n.bar
      fmt.Printf("FIELDS: %s\n", data) // FIELDS: 1 2.2
      var sealApp []interface{}
      if len(data)==0 || len(data)==1 {
        sealApp = data[0][0].([]interface{})
      }else{
        sealApp = append(data[0][0].([]interface{}), data[1][0].([]interface{})...)
      }
      m:=AllMyNotifiedPost{
        All:sealApp}
        resp, _ := json.Marshal(m)
        fmt.Println(string(resp))
      w.Header().Set("Content-Type", "application/json")
      w.Header().Set("charset", "utf-8")
      w.Header().Set("Access-Control-Allow-Origin", "*")
      w.Write([]byte(resp))
  }else{
    m:=HttpStatusForbiddenUnAuthorize{
    HttpStatus:403,HttpText:"Authentication Failed"}
    resp, _ := json.Marshal(m)
    w.WriteHeader(http.StatusForbidden)
    w.Header().Set("Content-Type", "application/json")
    w.Header().Set("charset", "utf-8")
    w.Header().Set("Access-Control-Allow-Origin", "*")
    w.Write([]byte(resp))
  }
}


func getMyAllRelatedPostBillaHandlerEventSource(w http.ResponseWriter, r *http.Request)  {
  token := r.FormValue("data")
  var allMyPosts *GetMyAllPostPlease
  json.Unmarshal([]byte(token), &allMyPosts)

  s1:="Match (ee:Rita)-[r:LINKED]-(ff:Rita) where ee.email ='"
  b1:=allMyPosts.Me
  s2:="' with ee,ff,r match (ee)-[:POSTED_THIS]->(gg:PostedData) where gg.type='"
  sMubarak:=allMyPosts.PostType
  dardMubarak:="' return collect(distinct [ee.name,gg.who,gg.attachment,gg.when,gg.content,id(gg),gg.action,gg.rate,gg.ignore]) as all union Match (ee:Rita)-[r:LINKED]-(ff:Rita) where ff.email ='"
  b2:=allMyPosts.Me
  s3:="' with ee,ff,r  match (ee)-[:POSTED_THIS]->(gg:PostedData) where gg.type='"
  sHangama:=allMyPosts.PostType
  sQbarpa:="' return collect(distinct [ee.name,gg.who,gg.attachment,gg.when,gg.content,id(gg),gg.action,gg.rate,gg.ignore]) as all"
  s17 := fmt.Sprint(s1,b1,s2,sMubarak,dardMubarak,b2,s3,sHangama,sQbarpa)
  fmt.Println(s17)
  driver := bolt.NewDriver()
    conn, err := driver.OpenNeo("bolt://rita:b.PuhuqVThYfCn.fvurl1e25g7fzyCI@hobby-panhpmpgjildgbkepcdcklol.dbs.graphenedb.com:24786?tls=true")
    defer conn.Close()
    if err != nil {
      panic(err)
    }
    data, _, _, _ := conn.QueryNeoAll(s17, nil)
    //fmt.Println("hooo")
    //fmt.Printf("COLUMNS: %#v\n", rowsMetadata["fields"].([]interface{}))  // COLUMNS: n.foo,n.bar
    fmt.Printf("FIELDS: %s\n", data) // FIELDS: 1 2.2
    var sealApp []interface{}
    if len(data)==0 || len(data)==1 {
      sealApp = data[0][0].([]interface{})
    }else{
      sealApp = append(data[0][0].([]interface{}), data[1][0].([]interface{})...)
    }
    m:=AllMyNotifiedPost{
      All:sealApp}
      resp, _ := json.Marshal(m)
      fmt.Println(string(resp))
    w.Header().Set("Content-Type", "application/json")
    w.Header().Set("charset", "utf-8")
    w.Header().Set("Access-Control-Allow-Origin", "*")
    w.Write([]byte(resp))
}



type GetMyAllAvailableFacesToLinkPlease struct {
  Me string
  Pkey string
}
type AllMyNotifiedFacesList struct {
  All []interface{}
}

func getMyAllRelatedFacesBillaHandler(w http.ResponseWriter, r *http.Request)  {
  token := r.FormValue("data")
  var allMyFaces *GetMyAllAvailableFacesToLinkPlease
  json.Unmarshal([]byte(token), &allMyFaces)


  insaan:=allMyFaces.Me
  mainInsaan:=allMyFaces.Pkey
  fmt.Println("aaj se tera")
  fmt.Println(mainInsaan)
  mainInsaanSaasur, _ := strconv.ParseInt(mainInsaan, 10, 64)
  rest:=authorizeThisFace(&insaan,&mainInsaanSaasur)


  if rest!="ERROR" {
    s1:="match (mm:Rita) where mm.email='"
    b1:=rest
    s2:="' with (mm) optional match (nn:Rita) where not (nn)<-[:LINKED]->(mm) and nn.email <>  '"
    s21:=rest
    s23:="' with nn optional match (nn)-[:HAS_MANDATORY_DP]->(ll) return collect(distinct [nn.name,nn.email,nn.gender,ll.title])[0..5]"
    s17 := fmt.Sprint(s1,b1,s2,s21,s23)
    fmt.Println(s17)
    driver := bolt.NewDriver()
      conn, err := driver.OpenNeo("bolt://rita:b.PuhuqVThYfCn.fvurl1e25g7fzyCI@hobby-panhpmpgjildgbkepcdcklol.dbs.graphenedb.com:24786?tls=true")
      defer conn.Close()
      dur,_:=time.ParseDuration("600s")
      conn.SetTimeout(dur)
      if err != nil {
        panic(err)
      }
      data, _, _, _ := conn.QueryNeoAll(s17, nil)
      //fmt.Println("hooo")
      //fmt.Printf("COLUMNS: %#v\n", rowsMetadata["fields"].([]interface{}))  // COLUMNS: n.foo,n.bar
      fmt.Printf("FIELDS: %s\n", data) // FIELDS: 1 2.2
      m:=AllMyNotifiedFacesList{
        All:data[0][0].([]interface{})}
        resp, _ := json.Marshal(m)
        fmt.Println(string(resp))
      w.Header().Set("Content-Type", "application/json")
      w.Header().Set("charset", "utf-8")
      w.Header().Set("Access-Control-Allow-Origin", "*")
      w.Write([]byte(resp))
  }else{
    m:=HttpStatusForbiddenUnAuthorize{
    HttpStatus:403,HttpText:"Authentication Failed"}
    resp, _ := json.Marshal(m)
    w.WriteHeader(http.StatusForbidden)
    w.Header().Set("Content-Type", "application/json")
    w.Header().Set("charset", "utf-8")
    w.Header().Set("Access-Control-Allow-Origin", "*")
    w.Write([]byte(resp))
  }
}
type BookShareBookIsbn struct {
  Author []string
  AvailableToDownload int
  BriefSynopsis string
  DownloadFormat []string
  DtbookSize int
  FreelyAvailable int
  Id int
  Images int
  Isbn13 string
  Publisher string
  Title string
}

type BookResult struct {
  Limit int
  NumPages int
  Page int
  Result []BookShareBookIsbn
  TotalResults int
}

type BookAsBook struct {
  List BookResult
}

type BookWrapper struct {
  Book BookAsBook
  Messages []string
  Version string
}

type WrappedBook struct {
  Bookshare BookWrapper
}

func testBooks(){

  /*var wrappedBook WrappedBook
  resp, err := http.Get("https://api.bookshare.org/book/popular/page/1/limit/100/format/json?api_key=xf3vuftt48kmqq8r3nctwxp7")
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
  fmt.Println("get:\n",string(body))
  if err1 := json.Unmarshal(body, &wrappedBook); err1 != nil {
        panic(err1)
    }
    isb:=wrappedBook.Bookshare.Book.List.Result[4].Isbn13
    fmt.Println(isb)
    var buffer bytes.Buffer
    buffer.WriteString("https://openlibrary.org/api/books?bibkeys=ISBN:")
    buffer.WriteString(isb)
    buffer.WriteString("&jscmd=data")
    fmt.Println(buffer.String())
    respB, errB := http.Get(buffer.String())
    if errB != nil {
  		panic(errB)
  	}
    defer respB.Body.Close()
  	bodyB, errC := ioutil.ReadAll(respB.Body)
    if errC != nil {
  		panic(errC)
  	}
    fmt.Println(string(bodyB))*/
  //w.Write([]byte(body))
}

func serveBooks(w http.ResponseWriter, r *http.Request){
  //var wrappedBook WrappedBook
  /*resp, err := http.Get("https://api.bookshare.org/book/popular/page/1/limit/100/format/json?api_key=xf3vuftt48kmqq8r3nctwxp7")
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
  fmt.Println("get:\n",string(body))*/
  var buffer bytes.Buffer
  buffer.WriteString("https://openlibrary.org/subjects/accessible_book.json")
  buffer.WriteString("?details=true&ebooks=true&limit=100&offset=0")
  resp, err := http.Get(buffer.String())
  if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, err1 := ioutil.ReadAll(resp.Body)
  if err1 != nil {
		panic(err1)
	}
  w.Write(body)
}

type AllItem struct {
  Woman string
  Man string
  Kid string
}



func serveWomanAmazon(womanApi *amazonproduct.AmazonProductAPI,womanChannel *chan []byte){
  result,err := womanApi.ItemSearchByKeyword("women's+fashion", 1)
  if (err != nil) {
    panic(err)
  }
  if err == nil {
    //aws := new(ItemLookupResponse)
    //xml.Unmarshal([]byte(result), aws)
    //TODO: Use "aws" freely :-)
  }
  xml := strings.NewReader(result)
  json, err1 := xj.Convert(xml)
  if err1 != nil {
    panic("That's embarrassing...")
  }
  *womanChannel<-json.Bytes()
}

func serveManAmazon(manApi *amazonproduct.AmazonProductAPI,manChannel *chan []byte){
  result,err := manApi.ItemSearchByKeyword("men's+fashion", 1)
  if (err != nil) {
    panic(err)
  }
  if err == nil {
    //aws := new(ItemLookupResponse)
    //xml.Unmarshal([]byte(result), aws)
    //TODO: Use "aws" freely :-)
  }
  xml := strings.NewReader(result)
  json, err1 := xj.Convert(xml)
  if err1 != nil {
    panic("That's embarrassing...")
  }
  *manChannel<-json.Bytes()
}

func serveKidsAmazon(kidApi *amazonproduct.AmazonProductAPI,kidChannel *chan []byte){
  result,err := kidApi.ItemSearchByKeyword("kid's+fashion", 1)
  if (err != nil) {
    panic(err)
  }
  if err == nil {
    //aws := new(ItemLookupResponse)
    //xml.Unmarshal([]byte(result), aws)
    //TODO: Use "aws" freely :-)
  }
  xml := strings.NewReader(result)
  json, err1 := xj.Convert(xml)
  if err1 != nil {
    panic("That's embarrassing...")
  }
  *kidChannel<-json.Bytes()
}

func amazonTestApiResponseGroup(kidApi *amazonproduct.AmazonProductAPI,kidChannel *chan []byte){
  m:= make(map[string]string)
  m["Keywords"]="kid's+fashion"
  result,err := kidApi.ItemSearch("Apparel",m)
  if (err != nil) {
    panic(err)
  }
  if err == nil {
    //aws := new(ItemLookupResponse)
    //xml.Unmarshal([]byte(result), aws)
    //TODO: Use "aws" freely :-)
  }
  xml := strings.NewReader(result)
  json, err1 := xj.Convert(xml)
  if err1 != nil {
    panic("That's embarrassing...")
  }
  *kidChannel<-json.Bytes()
}

func queryMongodbForAmazon(mongoChannel *chan string){
  dialInfo := &mgo.DialInfo{
    Addrs:    []string{"rita-shard-00-00-qk9t0.mongodb.net:27017","rita-shard-00-01-qk9t0.mongodb.net:27017","rita-shard-00-02-qk9t0.mongodb.net:27017"},
    Database: "test",
    Username: "badcodercpp",
    Password: "Ajayjha93",
    ReplicaSetName: "Rita-shard-0",
    Source: "admin",
    DialServer: func(addr *mgo.ServerAddr) (net.Conn, error) {
        return tls.Dial("tcp", addr.String(), &tls.Config{})
    },
    Timeout: time.Second * 10000,
  }
  session, err := mgo.DialWithInfo(dialInfo)

  if(err != nil){
    panic(err)
  }
  c := session.DB("test").C("amazonDbase")
    // Query One  4
    result := AmazonData{}
    err = c.Find(nil).One(&result)
    if err != nil {
      panic(err)
    }
    alllllss:=result.Data
    *mongoChannel<-alllllss
}

func serveAmazonQuery(w http.ResponseWriter, r *http.Request){
    chl1 := make(chan string)
    go queryMongodbForAmazon(&chl1)
    var tempAlling string
    for i := 0; i < 1; i++ {
      select {
        case msg1 := <-chl1:
            tempAlling = msg1
        }
    }
    w.Header().Set("Content-Type", "json")
    w.Write([]byte(tempAlling))
}

func serveAmazonQueryChanged(w http.ResponseWriter, r *http.Request){
  var api amazonproduct.AmazonProductAPI
  api.AccessKey = "AKIAJXGZLO3JAU726QXA"
  api.SecretKey = "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"
  api.Host = "webservices.amazon.in"
  api.AssociateTag = "ajayjha-21"
  api.Client = &http.Client{} // optional
  ch1 := make(chan []byte)
  ch2 := make(chan []byte)
  ch3 := make(chan []byte)
  ch4 := make(chan []byte)
  go serveWomanAmazon(&api,&ch1)
  go serveManAmazon(&api,&ch2)
  go serveKidsAmazon(&api,&ch3)
  go amazonTestApiResponseGroup(&api,&ch4)
  var tempWoman []byte
  var tempMan []byte
  var tempKid []byte

  //ch2 := make(chan []interface{})
  //ch3 := make(chan []interface{})
  //var aws ItemLookupResponse
  /*result,err := api.ItemSearchByKeyword("women's+fashion", 1)
  if (err != nil) {
    panic(err)
  }
  if err == nil {
    //aws := new(ItemLookupResponse)
    //xml.Unmarshal([]byte(result), aws)
    //TODO: Use "aws" freely :-)
  }
  xml := strings.NewReader(result)
  json, err1 := xj.Convert(xml)
  if err1 != nil {
    panic("That's embarrassing...")
  }*/
  for i := 0; i < 4; i++ {
    select {
      case msg1 := <-ch1:
            tempWoman = msg1
      case msg2 := <-ch2:
            tempMan = msg2
      case msg3 := <-ch3:
            tempKid = msg3
      case msg4 := <-ch4:
            fmt.Println(string(msg4))
      }
  }
  lk:=&AllItem{Woman:string(tempWoman),Man:string(tempMan),Kid:string(tempKid)}
  all, _ := json.Marshal(lk)
  w.Header().Set("Content-Type", "json")
  //fmt.Println(json.String())
  w.Write(all)
}


func serveWomanAmazonWebService(womanApi *amazonproduct.AmazonProductAPI,womanChannel *chan []byte){
  m:= make(map[string]string)
  m["Keywords"]="women's+fashion"
  m["ResponseGroup"]="Large,Variations"
  m["ItemPage"]="1"
  result,err := womanApi.ItemSearch("Apparel",m)
  if (err != nil) {
    panic(err)
  }
  xml := strings.NewReader(result)
  json, err1 := xj.Convert(xml)
  if err1 != nil {
    panic("That's embarrassing...")
  }
  *womanChannel<-json.Bytes()
}

func serveManAmazonWebService(manApi *amazonproduct.AmazonProductAPI,manChannel *chan []byte){
  m:= make(map[string]string)
  m["Keywords"]="men's+fashion"
  m["ResponseGroup"]="Large,Variations"
  m["ItemPage"]="1"
  result,err := manApi.ItemSearch("Apparel",m)
  if (err != nil) {
    panic(err)
  }
  xml := strings.NewReader(result)
  json, err1 := xj.Convert(xml)
  if err1 != nil {
    panic("That's embarrassing...")
  }
  *manChannel<-json.Bytes()
}

func serveKidsAmazonWebService(kidApi *amazonproduct.AmazonProductAPI,kidChannel *chan []byte){
  m:= make(map[string]string)
  m["Keywords"]="kid's+fashion"
  m["ResponseGroup"]="Large,Variations"
  m["ItemPage"]="1"
  result,err := kidApi.ItemSearch("Apparel",m)
  if (err != nil) {
    panic(err)
  }
  xml := strings.NewReader(result)
  json, err1 := xj.Convert(xml)
  if err1 != nil {
    panic("That's embarrassing...")
  }
  *kidChannel<-json.Bytes()
}


func syncNeo4jAmazonWomen(saveData *[]byte){
  nn:=string(*saveData)
  nss:= strings.TrimSpace(nn)
  s1:="merge (mm:Amazon{type:'women',amazonData:'''"
  b1:=nss
  s42:="'''}) return mm"
  s17 := fmt.Sprint(s1,b1,s42)
  //fmt.Println(s17)
  driver := bolt.NewDriver()
    conn, err := driver.OpenNeo("bolt://rita:b.PuhuqVThYfCn.fvurl1e25g7fzyCI@hobby-panhpmpgjildgbkepcdcklol.dbs.graphenedb.com:24786?tls=true")
    defer conn.Close()
    if err != nil {
      panic(err)
    }
    data, _, _, err1 := conn.QueryNeoAll(s17, nil)
    if err1 != nil {
      panic(err1)
    }
    //fmt.Println("hooo")
    //fmt.Printf("COLUMNS: %#v\n", rowsMetadata["fields"].([]interface{}))  // COLUMNS: n.foo,n.bar
    fmt.Printf("FIELDS: %s\n", data) // FIELDS: 1 2.2
}

func syncNeo4jAmazonMen(saveData *[]byte){
  nn:=string(*saveData)
  nss:= strings.TrimSpace(nn)
  s1:="merge (mm:Amazon{type:'men',amazonData:'''"
  b1:=nss
  s42:="'''}) return mm"
  s17 := fmt.Sprint(s1,b1,s42)
  //fmt.Println(s17)
  driver := bolt.NewDriver()
    conn, err := driver.OpenNeo("bolt://rita:b.PuhuqVThYfCn.fvurl1e25g7fzyCI@hobby-panhpmpgjildgbkepcdcklol.dbs.graphenedb.com:24786?tls=true")
    defer conn.Close()
    if err != nil {
      panic(err)
    }
    data, _, _, err1 := conn.QueryNeoAll(s17, nil)
    if err1 != nil {
      panic(err1)
    }
    //fmt.Println("hooo")
    //fmt.Printf("COLUMNS: %#v\n", rowsMetadata["fields"].([]interface{}))  // COLUMNS: n.foo,n.bar
    fmt.Printf("FIELDS: %s\n", data) // FIELDS: 1 2.2
}

func syncNeo4jAmazonKid(saveData *[]byte){
  nn:=string(*saveData)
  nss:= strings.TrimSpace(nn)
  s1:="merge (mm:Amazon{type:'kid',amazonData:'''"
  b1:=nss
  s42:="'''}) return mm"
  s17 := fmt.Sprint(s1,b1,s42)
  //fmt.Println(s17)
  driver := bolt.NewDriver()
    conn, err := driver.OpenNeo("bolt://rita:b.PuhuqVThYfCn.fvurl1e25g7fzyCI@hobby-panhpmpgjildgbkepcdcklol.dbs.graphenedb.com:24786?tls=true")
    defer conn.Close()
    if err != nil {
      panic(err)
    }
    data, _, _, err1 := conn.QueryNeoAll(s17, nil)
    if err1 != nil {
      panic(err1)
    }
    //fmt.Println("hooo")
    //fmt.Printf("COLUMNS: %#v\n", rowsMetadata["fields"].([]interface{}))  // COLUMNS: n.foo,n.bar
    fmt.Printf("FIELDS: %s\n", data) // FIELDS: 1 2.2
}


type AmazonData struct {
  Data string
}

func serveAmazonQueryWebService(amazonChannel *chan string){
  var api amazonproduct.AmazonProductAPI
  api.AccessKey = "AKIAJXGZLO3JAU726QXA"
  api.SecretKey = "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"
  api.Host = "webservices.amazon.in"
  api.AssociateTag = "ajayjha-21"
  api.Client = &http.Client{} // optional
  chs1 := make(chan []byte)
  chs2 := make(chan []byte)
  chs3 := make(chan []byte)
  chn1 := make(chan string)
  go serveWomanAmazonWebService(&api,&chs1)
  go serveManAmazonWebService(&api,&chs2)
  go serveKidsAmazonWebService(&api,&chs3)
  var tempWoman []byte
  var tempMan []byte
  var tempKid []byte
  var tempAlls string
  //ch2 := make(chan []interface{})
  //ch3 := make(chan []interface{})
  //var aws ItemLookupResponse
  /*result,err := api.ItemSearchByKeyword("women's+fashion", 1)
  if (err != nil) {
    panic(err)
  }
  if err == nil {
    //aws := new(ItemLookupResponse)
    //xml.Unmarshal([]byte(result), aws)
    //TODO: Use "aws" freely :-)
  }
  xml := strings.NewReader(result)
  json, err1 := xj.Convert(xml)
  if err1 != nil {
    panic("That's embarrassing...")
  }*/
  for i := 0; i < 3; i++ {
    select {
      case msg1 := <-chs1:
            tempWoman = msg1
            //go syncNeo4jAmazonWomen(&tempWoman)
            fmt.Println("firsta arrived")
      case msg2 := <-chs2:
            tempMan = msg2
            fmt.Println("firsta arrived")
            //go syncNeo4jAmazonMen(&tempMan)
      case msg3 := <-chs3:
            tempKid = msg3
            fmt.Println("firsta arrived")
            //go syncNeo4jAmazonKid(&tempKid)
      }
  }
  lk:=&AllItem{Woman:string(tempWoman),Man:string(tempMan),Kid:string(tempKid)}
  all, err := json.Marshal(lk)
  if(err != nil){
    panic(err)
  }
  fmt.Println("amazon printed")
  go connectMong(string(all),&chn1)
  for i := 0; i < 1; i++ {
    select {
      case msgn1 := <-chn1:
            tempAlls = msgn1
            fmt.Println("firsta arrived")
      }
  }
  *amazonChannel<-tempAlls
}



func connectMong(str string,mongoChannel *chan string){
  fmt.Println("executing mango")

  dialInfo := &mgo.DialInfo{
    Addrs:    []string{"rita-shard-00-00-qk9t0.mongodb.net:27017","rita-shard-00-01-qk9t0.mongodb.net:27017","rita-shard-00-02-qk9t0.mongodb.net:27017"},
    Database: "test",
    Username: "badcodercpp",
    Password: "Ajayjha93",
    ReplicaSetName: "Rita-shard-0",
    Source: "admin",
    DialServer: func(addr *mgo.ServerAddr) (net.Conn, error) {
        return tls.Dial("tcp", addr.String(), &tls.Config{})
    },
    Timeout: time.Second * 10000,
  }
  session, err := mgo.DialWithInfo(dialInfo)


  fmt.Println(session)
  if(err != nil){
    panic(err)
  }
  c := session.DB("test").C("amazonDbase")
  _ , errl:=c.RemoveAll(nil)
  if errl != nil {
    panic(errl)
  }
  err = c.Insert(&AmazonData{Data: str})
  fmt.Println(c)
    if err != nil {
      panic(err)
    }

    // Query One  4
    result := AmazonData{}
    err = c.Find(nil).One(&result)
    if err != nil {
      panic(err)
    }
    fmt.Println("Phone", result)
    alllllss:=result.Data
    *mongoChannel<-alllllss
}

func syncAmazonMongoWebService(w http.ResponseWriter, r *http.Request){
  chn := make(chan string)
  var tempAll string
  go serveAmazonQueryWebService(&chn)
  for i := 0; i < 1; i++ {
    select {
      case msg1 := <-chn:
            tempAll = msg1
    }
  }
  w.Header().Set("Content-Type", "json")
  w.Write([]byte(tempAll))
}

type BeautifulPratibhaFace struct {
  Face []interface{}
}
type BeautifulPratibhaPlaylist struct {
  Playlist []interface{}
}
type BeautifulPratibhaPost struct {
  Post []interface{}
}
type BeautifulPratibhaDetail struct {
  Detail []interface{}
}

type BeautifulPratibha struct {
  Face BeautifulPratibhaFace
  Playlist BeautifulPratibhaPlaylist
  Post BeautifulPratibhaPost
  Detail BeautifulPratibhaDetail
}

func serveWowPratibhaYouLooksLikeAnAngelDetail(aadat *chan []interface{},id *string)  {
  //sealApp := append(data[0][0].([]interface{}), data[1][0].([]interface{})...)
  s1:="Match (ee:Rita) where ee.email ='"
  b:=*id
  s2:="' with ee optional Match (ee)-[:HAS_MANDATORY_DP]->(gg:ProfilePic) return collect([ee.name,ee.email,gg.title,ee.gender]) as all"
  s12 := fmt.Sprint(s1,b,s2)
  fmt.Println(s12)
  driver := bolt.NewDriver()
    conn, err := driver.OpenNeo("bolt://rita:b.PuhuqVThYfCn.fvurl1e25g7fzyCI@hobby-panhpmpgjildgbkepcdcklol.dbs.graphenedb.com:24786?tls=true")
    defer conn.Close()
    dur,_:=time.ParseDuration("600s")
    conn.SetTimeout(dur)
    if err != nil {
      panic(err)
      fmt.Println("lollllllllll")
    }
    data, _, _, _ := conn.QueryNeoAll(s12, nil)
    //fmt.Printf("FIELDS: %s\n", data) // FIELDS: 1 2.2


    //fmt.Println("hooo")
    //fmt.Printf("COLUMNS: %#v\n", rowsMetadata["fields"].([]interface{}))  // COLUMNS: n.foo,n.bar
    //fmt.Printf("FIELDS: %s\n", data) // FIELDS: 1 2.2
    *aadat <- data[0][0].([]interface{})
}

func serveWowPratibhaYouLooksLikeAnAngelFace(aadat *chan []interface{},id *string)  {
  //sealApp := append(data[0][0].([]interface{}), data[1][0].([]interface{})...)
  s1:="Match (ee:Rita)-[r:LINKED]->(ff:Rita) where ee.email ='"
  b:=*id
  s2:="' with ee,ff,r match (ff)-[:HAS_MANDATORY_DP]->(gg:ProfilePic) return collect([ff.name,ff.email,gg.title,r.tlfSocial,r.tlfProfessional,r.tlfConsumer,r.tlfReader,r.tlfDater,r.tlfDonater,r.tlfSinner]) as all union Match (ee:Rita)-[r:LINKED]->(ff:Rita) where ff.email ='"
  s34567:=*id
  s554433:="' with ee,ff,r  match (ee)-[:HAS_MANDATORY_DP]->(gg:ProfilePic) return collect([ee.name,ee.email,gg.title,r.tlfSocial,r.tlfProfessional,r.tlfConsumer,r.tlfReader,r.tlfDater,r.tlfDonater,r.tlfSinner]) as all"
  s12 := fmt.Sprint(s1,b,s2,s34567,s554433)
  fmt.Println(s12)
  driver := bolt.NewDriver()
    conn, err := driver.OpenNeo("bolt://rita:b.PuhuqVThYfCn.fvurl1e25g7fzyCI@hobby-panhpmpgjildgbkepcdcklol.dbs.graphenedb.com:24786?tls=true")
    defer conn.Close()
    dur,_:=time.ParseDuration("600s")
    conn.SetTimeout(dur)
    if err != nil {
      panic(err)
      fmt.Println("lollllllllll")
    }
    data, _, _, _ := conn.QueryNeoAll(s12, nil)
    fmt.Printf("FIELDS: %s\n", data) // FIELDS: 1 2.2
    var sealApp []interface{}
    if len(data)==0 || len(data)==1 {
      sealApp = data[0][0].([]interface{})
    }else{
      sealApp = append(data[0][0].([]interface{}), data[1][0].([]interface{})...)
    }

    //fmt.Println("hooo")
    //fmt.Printf("COLUMNS: %#v\n", rowsMetadata["fields"].([]interface{}))  // COLUMNS: n.foo,n.bar
    //fmt.Printf("FIELDS: %s\n", data) // FIELDS: 1 2.2
    *aadat <- sealApp
}

func serveWowPratibhaYouLooksLikeAnAngelPlaylist(aadat *chan []interface{},id *string)  {
  saa1:="MATCH (ee:Rita)-[:HAS_YOUTUBE_PLAYLIST]->(ll) where ee.email='"
  baa:=*id
  saa2:="' return collect(ll.title) as all"
  saa12 := fmt.Sprint(saa1,baa,saa2)
  fmt.Println(saa12)
  driveraa := bolt.NewDriver()
    connaa, erraa := driveraa.OpenNeo("bolt://rita:b.PuhuqVThYfCn.fvurl1e25g7fzyCI@hobby-panhpmpgjildgbkepcdcklol.dbs.graphenedb.com:24786?tls=true")
    defer connaa.Close()
    dur,_:=time.ParseDuration("600s")
    connaa.SetTimeout(dur)
    if erraa != nil {
      panic(erraa)
    }
    dataaa, _, _, _ := connaa.QueryNeoAll(saa12, nil)
    //fmt.Println("hooo")
    //fmt.Printf("COLUMNS: %#v\n", rowsMetadata["fields"].([]interface{}))  // COLUMNS: n.foo,n.bar
    //fmt.Printf("FIELDS: %s\n", data[0][0].([]interface{})) // FIELDS: 1 2.2
        *aadat <- dataaa[0][0].([]interface{})
}

func serveWowPratibhaYouLooksLikeAnAngelPost(aadat *chan []interface{},id *string)  {
  s1:="Match (ff:Rita)-[:POSTED_THIS]->(gg:PostedData) where ff.email ='"
  b1:=*id
  s2:="' return collect(distinct [ff.name,gg.who,gg.attachment,gg.when,gg.content,id(gg),gg.action,gg.rate,gg.ignore]) as all"

  s17 := fmt.Sprint(s1,b1,s2)
  fmt.Println(s17)
  driver := bolt.NewDriver()
    conn, err := driver.OpenNeo("bolt://rita:b.PuhuqVThYfCn.fvurl1e25g7fzyCI@hobby-panhpmpgjildgbkepcdcklol.dbs.graphenedb.com:24786?tls=true")
    defer conn.Close()
    dur,_:=time.ParseDuration("600s")
    conn.SetTimeout(dur)
    if err != nil {
      panic(err)
    }
    data, _, _, _ := conn.QueryNeoAll(s17, nil)
    //fmt.Println("hooo")
    //fmt.Printf("COLUMNS: %#v\n", rowsMetadata["fields"].([]interface{}))  // COLUMNS: n.foo,n.bar
    //fmt.Printf("FIELDS: %s\n", data) // FIELDS: 1 2.2
    var sealApp []interface{}
    if len(data)==0 || len(data)==1 {
      sealApp = data[0][0].([]interface{})
    }else{
      sealApp = append(data[0][0].([]interface{}), data[1][0].([]interface{})...)
    }
        *aadat <- sealApp
}


type GehraHaiApnaPyar struct {
  Me string
  Pkey string
  Id string
}


func serveWowPratibhaYouLooksLikeAnAngel(w http.ResponseWriter, r *http.Request)  {
    //id:= r.FormValue("id")

    r.ParseForm()
    data:=r.Form.Get("data")
    var LetMeSeeYouPratibhaIlu GehraHaiApnaPyar
    json.Unmarshal([]byte(data),&LetMeSeeYouPratibhaIlu)
    id:=LetMeSeeYouPratibhaIlu.Id
    insaan:=LetMeSeeYouPratibhaIlu.Me
    mainInsaan:=LetMeSeeYouPratibhaIlu.Pkey
    fmt.Println("aaj se tera")
    fmt.Println(mainInsaan)
    mainInsaanSaasur, _ := strconv.ParseInt(mainInsaan, 10, 64)
    rest:=authorizeThisFace(&insaan,&mainInsaanSaasur)


    if rest!="ERROR" {
      chnl1 := make(chan []interface{})
      chnl2 := make(chan []interface{})
      chnl3 := make(chan []interface{})
      chnl4 := make(chan []interface{})
      go serveWowPratibhaYouLooksLikeAnAngelFace(&chnl1,&id)
      go serveWowPratibhaYouLooksLikeAnAngelPlaylist(&chnl2,&id)
      go serveWowPratibhaYouLooksLikeAnAngelPost(&chnl3,&id)
      go serveWowPratibhaYouLooksLikeAnAngelDetail(&chnl4,&id)
      var hamsafarFace BeautifulPratibhaFace
      var hamsafarPlaylist BeautifulPratibhaPlaylist
      var hamsafarPost BeautifulPratibhaPost
      var hamsafarDetail BeautifulPratibhaDetail

      for i := 0; i < 4; i++ {
        select {
          case msg1 := <-chnl1:
                hamsafarFace = BeautifulPratibhaFace{
                Face:msg1}
          case msg2 := <-chnl2:
                hamsafarPlaylist = BeautifulPratibhaPlaylist{
                Playlist:msg2}
          case msg3 := <-chnl3:
                hamsafarPost = BeautifulPratibhaPost{
                Post:msg3}
          case msg4 := <-chnl4:
                hamsafarDetail = BeautifulPratibhaDetail{
                Detail:msg4}
          }
      }
      m:=BeautifulPratibha{
      Face:hamsafarFace,Playlist:hamsafarPlaylist,Post:hamsafarPost,Detail:hamsafarDetail}
      resp, _ := json.Marshal(m)

      fmt.Println(string(resp))
      w.Header().Set("Content-Type", "application/json")
      w.Header().Set("charset", "utf-8")
      w.Header().Set("Access-Control-Allow-Origin", "*")
      w.Write([]byte(resp))
    }else{
      m:=HttpStatusForbiddenUnAuthorize{
      HttpStatus:403,HttpText:"Authentication Failed"}
      resp, _ := json.Marshal(m)
      w.WriteHeader(http.StatusForbidden)
      w.Header().Set("Content-Type", "application/json")
      w.Header().Set("charset", "utf-8")
      w.Header().Set("Access-Control-Allow-Origin", "*")
      w.Write([]byte(resp))
    }


          //*a <- data[0][0].([]interface{})
}

//face of pratibha


type IwanaSeeYourFaceDailyPratibha struct {
  Me string
  Pkey string
}

func serveWowPratibhaYouLooksLikeAnAngelPratibha(w http.ResponseWriter, r *http.Request)  {
    //id:= r.FormValue("id")
    r.ParseForm()
    data:=r.Form.Get("data")
    var LetMeSeeYouPratibha IwanaSeeYourFaceDailyPratibha
    json.Unmarshal([]byte(data),&LetMeSeeYouPratibha)

    insaan:=LetMeSeeYouPratibha.Me
    mainInsaan:=LetMeSeeYouPratibha.Pkey
    fmt.Println("aaj se tera")
    fmt.Println(mainInsaan)
    mainInsaanSaasur, _ := strconv.ParseInt(mainInsaan, 10, 64)
    rest:=authorizeThisFace(&insaan,&mainInsaanSaasur)


    if rest!="ERROR" {
      chnl1 := make(chan []interface{})
      chnl2 := make(chan []interface{})
      chnl3 := make(chan []interface{})
      chnl4 := make(chan []interface{})
      go serveWowPratibhaYouLooksLikeAnAngelFace(&chnl1,&rest)
      go serveWowPratibhaYouLooksLikeAnAngelPlaylist(&chnl2,&rest)
      go serveWowPratibhaYouLooksLikeAnAngelPost(&chnl3,&rest)
      go serveWowPratibhaYouLooksLikeAnAngelDetail(&chnl4,&rest)
      var hamsafarFace BeautifulPratibhaFace
      var hamsafarPlaylist BeautifulPratibhaPlaylist
      var hamsafarPost BeautifulPratibhaPost
      var hamsafarDetail BeautifulPratibhaDetail

      for i := 0; i < 4; i++ {
        select {
          case msg1 := <-chnl1:
                hamsafarFace = BeautifulPratibhaFace{
                Face:msg1}
          case msg2 := <-chnl2:
                hamsafarPlaylist = BeautifulPratibhaPlaylist{
                Playlist:msg2}
          case msg3 := <-chnl3:
                hamsafarPost = BeautifulPratibhaPost{
                Post:msg3}
          case msg4 := <-chnl4:
                hamsafarDetail = BeautifulPratibhaDetail{
                Detail:msg4}
          }
      }
      m:=BeautifulPratibha{
      Face:hamsafarFace,Playlist:hamsafarPlaylist,Post:hamsafarPost,Detail:hamsafarDetail}
      resp, _ := json.Marshal(m)

      fmt.Println(string(resp))
      w.Header().Set("Content-Type", "application/json")
      w.Header().Set("charset", "utf-8")
      w.Header().Set("Access-Control-Allow-Origin", "*")
      w.Write([]byte(resp))
    }else{
      m:=HttpStatusForbiddenUnAuthorize{
      HttpStatus:403,HttpText:"Authentication Failed"}
      resp, _ := json.Marshal(m)
      w.WriteHeader(http.StatusForbidden)
      w.Header().Set("Content-Type", "application/json")
      w.Header().Set("charset", "utf-8")
      w.Header().Set("Access-Control-Allow-Origin", "*")
      w.Write([]byte(resp))
    }


          //*a <- data[0][0].([]interface{})
}


//end face of pratibha

type CreateMyShopSteadyStateDataBeforeFlush struct {
  ShopName string
  ShopType []string
  ShopStartDate string
  ShopMobile string
  ShopPin []string
  ShopAddFirst string
  ShopAddSecond string
  ShopOpenTime string
  ShopCloseTime string
  DeliveryFlag string
  DeliveryTimeFrom string
  DeliveryTimeUpto string
  SallDeliveryMyPinCode []string
  AvgDeliverCharge string
  MaxDeliveryTime string
  FeasiblePaymentOption []string
  CollactibleFlag string
  CollactiblePreserveTime string
  ShopIcon string
  Owner string
  Pkey string
}



func saveShopDataToNeoDatabase(aadat *chan []interface{},meriRani *CreateMyShopSteadyStateDataBeforeFlush,meOwwwwn *string)  {
  s1:="Match (ee:Rita) where ee.email ='"
  b1:=*meOwwwwn
  s2:="' with ee MERGE (ee)-[:HAS_OWN_SHOP]-(ff:ShopContainer{ShopName:'"
  bx:=meriRani.ShopName
  s4411:="',ShopType:["
  bx2:=""
  for ii1,val:=range meriRani.ShopType {
    if ii1==len(meriRani.ShopType)-1 {
      bx2=bx2+"'"+val+"'"
    }else{
      bx2=bx2+"'"+val+"',"
    }

  }
  s4412:="],ShopStartDate:'"
  bx3:=meriRani.ShopStartDate
  s4413:="',ShopMobile:'"
  bx4:=meriRani.ShopStartDate
  s4414:="',ShopPin:["
  bx5:=""
  for i1,val1:=range meriRani.ShopPin {
    if i1==len(meriRani.ShopPin)-1 {
      bx5=bx5+"'"+val1+"'"
    }else{
      bx5=bx5+"'"+val1+"',"
    }

  }
  s4415:="],ShopAddFirst:'"
  bx6:=meriRani.ShopAddFirst
  s4416:="',ShopAddSecond:'"
  bx7:=meriRani.ShopAddSecond
  s4417:="',ShopOpenTime:'"
  bx8:=meriRani.ShopOpenTime
  s4418:="',ShopCloseTime:'"
  bx9:=meriRani.ShopCloseTime
  s4419:="',DeliveryFlag:'"
  bx10:=meriRani.DeliveryFlag
  s4420:="',DeliveryTimeFrom:'"
  bx11:=meriRani.DeliveryTimeFrom
  s4421:="',DeliveryTimeUpto:'"
  bx12:=meriRani.DeliveryTimeUpto
  s4422:="',SallDeliveryMyPinCode:["
  bx13:=""
  for i2,val2:=range meriRani.SallDeliveryMyPinCode {
    if i2==len(meriRani.SallDeliveryMyPinCode)-1 {
      bx13=bx13+"'"+val2+"'"
    }else{
      bx13=bx13+"'"+val2+"',"
    }

  }
  s4423:="],AvgDeliverCharge:'"
  bx14:=meriRani.AvgDeliverCharge
  s4424:="',MaxDeliveryTime:'"
  bx15:=meriRani.MaxDeliveryTime
  s4425:="',FeasiblePaymentOption:["
  bx16:=""
  for i3,val3:=range meriRani.FeasiblePaymentOption {
    if i3==len(meriRani.FeasiblePaymentOption)-1 {
      bx16=bx16+"'"+val3+"'"
    }else{
      bx16=bx16+"'"+val3+"',"
    }

  }
  s4426:="],CollactibleFlag:'"
  bx17:=meriRani.CollactibleFlag
  s4427:="',CollactiblePreserveTime:'"
  bx18:=meriRani.CollactiblePreserveTime
  s4428:="',ShopIcon:'"
  bx19:=meriRani.ShopIcon
  s4429:="',Owner:'"
  bx20:=*meOwwwwn
  s4430:="'}) with ee optional match (ee)-[:HAS_OWN_SHOP]-(gg:ShopContainer) return collect([gg.ShopName,gg.ShopType,gg.ShopStartDate,gg.ShopMobile,gg.ShopPin,gg.ShopAddFirst,gg.ShopAddSecond,gg.ShopOpenTime,gg.ShopCloseTime,gg.DeliveryFlag,gg.DeliveryTimeFrom,gg.DeliveryTimeUpto,gg.SallDeliveryMyPinCode,gg.AvgDeliverCharge,gg.MaxDeliveryTime,gg.FeasiblePaymentOption,gg.CollactibleFlag,gg.CollactiblePreserveTime,gg.ShopIcon,gg.Owner,id(gg)]) as all"
  s17 := fmt.Sprint(s1,b1,s2,bx,s4411,bx2,s4412,bx3,s4413,bx4,s4414,bx5,s4415,bx6,s4416,bx7,s4417,bx8,s4418,bx9,s4419,bx10,s4420,bx11,s4421,bx12,s4422,bx13,s4423,bx14,s4424,bx15,s4425,bx16,s4426,bx17,s4427,bx18,s4428,bx19,s4429,bx20,s4430)
  fmt.Println(s17)
  driver := bolt.NewDriver()
    conn, err := driver.OpenNeo("bolt://rita:b.PuhuqVThYfCn.fvurl1e25g7fzyCI@hobby-panhpmpgjildgbkepcdcklol.dbs.graphenedb.com:24786?tls=true")
    defer conn.Close()
    dur,_:=time.ParseDuration("600s")
    conn.SetTimeout(dur)
    if err != nil {
      panic(err)
    }
    data, _, _, _ := conn.QueryNeoAll(s17, nil)
    fmt.Println(data)
    *aadat<-data[0][0].([]interface{})
}


type PratibhaShopDetailsAll struct {
  All []interface{}
}


func iLoveYouPratibhaSharmaAndIWillGetYouShopCreateHandler(w http.ResponseWriter, r *http.Request) {
  r.ParseForm()
  data:=r.Form.Get("data")
  var MyShopContainer CreateMyShopSteadyStateDataBeforeFlush
  json.Unmarshal([]byte(data),&MyShopContainer)
  chnlShop := make(chan []interface{})
  //var writerPratibha []interface{}


  insaan:=MyShopContainer.Owner
  mainInsaan:=MyShopContainer.Pkey
  fmt.Println("aaj se tera")
  fmt.Println(mainInsaan)
  mainInsaanSaasur, _ := strconv.ParseInt(mainInsaan, 10, 64)
  rest:=authorizeThisFace(&insaan,&mainInsaanSaasur)

  if rest!="ERROR" {
    go saveShopDataToNeoDatabase(&chnlShop,&MyShopContainer,&rest)
    for i := 0; i < 1; i++ {
      select {
        case msg1 := <-chnlShop:
          m:=PratibhaShopDetailsAll{
            All:msg1}
            resp, _ := json.Marshal(m)
            w.Header().Set("Content-Type", "application/json")
            w.Header().Set("charset", "utf-8")
            w.Header().Set("Access-Control-Allow-Origin", "*")
            w.Write([]byte(resp))
        }
    }
  }else{
    m:=HttpStatusForbiddenUnAuthorize{
    HttpStatus:403,HttpText:"Authentication Failed"}
    resp, _ := json.Marshal(m)
    w.WriteHeader(http.StatusForbidden)
    w.Header().Set("Content-Type", "application/json")
    w.Header().Set("charset", "utf-8")
    w.Header().Set("Access-Control-Allow-Origin", "*")
    w.Write([]byte(resp))
  }


  //fmt.Println(string(data))
  //fmt.Println(MyShopContainer.ShopName)

}

type HoldMeTemperoryPratibhaPease struct {
  Me string
  Pkey string
}

type IndiShopDetailAjaxPratibha struct {
  All []interface{}
}

func getMyOwnShopDetailsPratibhaPleaseLoveYouHandler(w http.ResponseWriter, r *http.Request)  {
  r.ParseForm()
  data:=r.Form.Get("data")
  var MyTempHolder HoldMeTemperoryPratibhaPease
  json.Unmarshal([]byte(data),&MyTempHolder)


  insaan:=MyTempHolder.Me
  mainInsaan:=MyTempHolder.Pkey
  fmt.Println("aaj se tera")
  fmt.Println(mainInsaan)
  mainInsaanSaasur, _ := strconv.ParseInt(mainInsaan, 10, 64)
  rest:=authorizeThisFace(&insaan,&mainInsaanSaasur)


  if rest!="ERROR" {
    s1:="Match (ee:Rita)-[:HAS_OWN_SHOP]->(gg:ShopContainer) where ee.email ='"
    b1:=rest
    s2:="' return collect([gg.ShopName,gg.ShopType,gg.ShopStartDate,gg.ShopMobile,gg.ShopPin,gg.ShopAddFirst,gg.ShopAddSecond,gg.ShopOpenTime,gg.ShopCloseTime,gg.DeliveryFlag,gg.DeliveryTimeFrom,gg.DeliveryTimeUpto,gg.SallDeliveryMyPinCode,gg.AvgDeliverCharge,gg.MaxDeliveryTime,gg.FeasiblePaymentOption,gg.CollactibleFlag,gg.CollactiblePreserveTime,gg.ShopIcon,gg.Owner,id(gg)]) as all"
    s17 := fmt.Sprint(s1,b1,s2)
    fmt.Println(s17)
    driver := bolt.NewDriver()
      conn, err := driver.OpenNeo("bolt://rita:b.PuhuqVThYfCn.fvurl1e25g7fzyCI@hobby-panhpmpgjildgbkepcdcklol.dbs.graphenedb.com:24786?tls=true")
      defer conn.Close()
      dur,_:=time.ParseDuration("600s")
      conn.SetTimeout(dur)
      if err != nil {
        panic(err)
      }
      dataOf, _, _, _ := conn.QueryNeoAll(s17, nil)
      m:=IndiShopDetailAjaxPratibha{
        All:dataOf[0][0].([]interface{})}
        resp, _ := json.Marshal(m)
        w.Header().Set("Content-Type", "application/json")
        w.Header().Set("charset", "utf-8")
        w.Header().Set("Access-Control-Allow-Origin", "*")
        w.Write([]byte(resp))
  }else{
    m:=HttpStatusForbiddenUnAuthorize{
    HttpStatus:403,HttpText:"Authentication Failed"}
    resp, _ := json.Marshal(m)
    w.WriteHeader(http.StatusForbidden)
    w.Header().Set("Content-Type", "application/json")
    w.Header().Set("charset", "utf-8")
    w.Header().Set("Access-Control-Allow-Origin", "*")
    w.Write([]byte(resp))
  }
}


type MyTempItemWillBeAddedToShopPratibha struct {
  ItemName string
  ItemType []string
  ItemCategory []string
  ItemQuantity string
  ItemPrice string
  MfgDate string
  ExpDate string
  MaxAvailability string
  IconArray []string
  Me int
  Pkey string
  OwwnMe string
}


type IndiItemGetResponseDbPratibha struct {
  All []interface{}
}


func addItemToMyShopPratibhaPleaseLUHandler(w http.ResponseWriter, r *http.Request)  {
  r.ParseForm()
  data:=r.Form.Get("data")
  fmt.Println(data)
  var MyTempItemHolder MyTempItemWillBeAddedToShopPratibha
  json.Unmarshal([]byte(data),&MyTempItemHolder)
  fmt.Println(MyTempItemHolder.Me)


  insaan:=MyTempItemHolder.OwwnMe
  mainInsaan:=MyTempItemHolder.Pkey
  fmt.Println("aaj se tera")
  fmt.Println(mainInsaan)
  mainInsaanSaasur, _ := strconv.ParseInt(mainInsaan, 10, 64)
  rest:=authorizeThisFace(&insaan,&mainInsaanSaasur)


  if rest!="ERROR" {
    s1:="Match (gg:ShopContainer) where id(gg)="
    b1:=MyTempItemHolder.Me
    s2:=" MERGE (gg)-[:HAS_UNIQUE_ITEM]-(ee:ItemContainer{ItemName:'"
    b2:=MyTempItemHolder.ItemName
    s3:="',ItemType:["
    b3:=""
    for i3,val3:=range MyTempItemHolder.ItemType {
      if i3==len(MyTempItemHolder.ItemType)-1 {
        b3=b3+"'"+val3+"'"
      }else{
        b3=b3+"'"+val3+"',"
      }
    }
    s4:="],ItemCategory:["
    b4:=""
    for i4,val4:=range MyTempItemHolder.ItemCategory {
      if i4==len(MyTempItemHolder.ItemCategory)-1 {
        b4=b4+"'"+val4+"'"
      }else{
        b4=b4+"'"+val4+"',"
      }
    }
    s5:="],ItemQuantity:'"
    b5:=MyTempItemHolder.ItemQuantity
    s6:="',ItemPrice:'"
    b6:=MyTempItemHolder.ItemPrice
    s7:="',MfgDate:'"
    b7:=MyTempItemHolder.MfgDate
    s8:="',ExpDate:'"
    b8:=MyTempItemHolder.ExpDate
    s9:="',MaxAvailability:'"
    b9:=MyTempItemHolder.MaxAvailability
    s10:="',IconArray:["
    b10:=""
    for i5,val5:=range MyTempItemHolder.IconArray {
      if i5==len(MyTempItemHolder.IconArray)-1 {
        b10=b10+"'"+val5+"'"
      }else{
        b10=b10+"'"+val5+"',"
      }
    }
    s11:="],Me:'"
    b11:=MyTempItemHolder.Me
    s12:="'}) return collect([ee.ItemName,ee.ItemType,ee.ItemCategory,ee.ItemQuantity,ee.ItemPrice,ee.MfgDate,ee.ExpDate,ee.MaxAvailability,ee.IconArray,ee.Me,id(ee)]) as all"
    s17 := fmt.Sprint(s1,b1,s2,b2,s3,b3,s4,b4,s5,b5,s6,b6,s7,b7,s8,b8,s9,b9,s10,b10,s11,b11,s12)
    fmt.Println(s17)
    driver := bolt.NewDriver()
      conn, err := driver.OpenNeo("bolt://rita:b.PuhuqVThYfCn.fvurl1e25g7fzyCI@hobby-panhpmpgjildgbkepcdcklol.dbs.graphenedb.com:24786?tls=true")
      defer conn.Close()
      dur,_:=time.ParseDuration("600s")
      conn.SetTimeout(dur)
      if err != nil {
        panic(err)
      }
      dataOf, _, _, _ := conn.QueryNeoAll(s17, nil)
      m:=IndiItemGetResponseDbPratibha{
        All:dataOf[0][0].([]interface{})}
        resp, _ := json.Marshal(m)
        w.Header().Set("Content-Type", "application/json")
        w.Header().Set("charset", "utf-8")
        w.Header().Set("Access-Control-Allow-Origin", "*")
        w.Write([]byte(resp))
  }else{
    m:=HttpStatusForbiddenUnAuthorize{
    HttpStatus:403,HttpText:"Authentication Failed"}
    resp, _ := json.Marshal(m)
    w.WriteHeader(http.StatusForbidden)
    w.Header().Set("Content-Type", "application/json")
    w.Header().Set("charset", "utf-8")
    w.Header().Set("Access-Control-Allow-Origin", "*")
    w.Write([]byte(resp))
  }


}


type ReturnMyAllRelatedItemsContainer struct {
  All []interface{}
}


type BearWithPratibhaDrink struct {
  Me string
  Pkey string
}


func getMyRelatedItemsToBuyPratibhaPleaseHandler(w http.ResponseWriter, r *http.Request)  {
  r.ParseForm()
  data:=r.Form.Get("data")
  fmt.Println(data)
  var MyTempBearHolder BearWithPratibhaDrink
  json.Unmarshal([]byte(data),&MyTempBearHolder)
  fmt.Println(MyTempBearHolder.Me)


  insaan:=MyTempBearHolder.Me
  mainInsaan:=MyTempBearHolder.Pkey
  fmt.Println("aaj se tera")
  fmt.Println(mainInsaan)
  mainInsaanSaasur, _ := strconv.ParseInt(mainInsaan, 10, 64)
  rest:=authorizeThisFace(&insaan,&mainInsaanSaasur)


  if rest!="ERROR" {
    s1:="Match (ff:Rita)-[:HAS_OWN_SHOP]->(gg:ShopContainer) where ff.email <>'"
    b1:=rest
    s2:="' or (ff.email)='"
    sWd:=rest
    sWe:="' with gg OPTIONAL MATCH (gg)-[:HAS_UNIQUE_ITEM]-(ee:ItemContainer) return collect([ee.ItemName,ee.ItemType,ee.ItemCategory,ee.ItemQuantity,ee.ItemPrice,ee.MfgDate,ee.ExpDate,ee.MaxAvailability,ee.IconArray,ee.Me,id(ee),gg.ShopName,gg.ShopType,gg.ShopStartDate,gg.ShopMobile,gg.ShopPin,gg.ShopAddFirst,gg.ShopAddSecond,gg.ShopOpenTime,gg.ShopCloseTime,gg.DeliveryFlag,gg.DeliveryTimeFrom,gg.DeliveryTimeUpto,gg.SallDeliveryMyPinCode,gg.AvgDeliverCharge,gg.MaxDeliveryTime,gg.FeasiblePaymentOption,gg.CollactibleFlag,gg.CollactiblePreserveTime,gg.ShopIcon,gg.Owner,id(gg),id(ee)]) as all"
    s17 := fmt.Sprint(s1,b1,s2,sWd,sWe)
    fmt.Println(s17)
    driver := bolt.NewDriver()
      conn, err := driver.OpenNeo("bolt://rita:b.PuhuqVThYfCn.fvurl1e25g7fzyCI@hobby-panhpmpgjildgbkepcdcklol.dbs.graphenedb.com:24786?tls=true")
      defer conn.Close()
      dur,_:=time.ParseDuration("600s")
      conn.SetTimeout(dur)
      if err != nil {
        panic(err)
      }
      dataOf, _, _, _ := conn.QueryNeoAll(s17, nil)
      m:=ReturnMyAllRelatedItemsContainer{
        All:dataOf[0][0].([]interface{})}
        resp, _ := json.Marshal(m)
        w.Header().Set("Content-Type", "application/json")
        w.Header().Set("charset", "utf-8")
        w.Header().Set("Access-Control-Allow-Origin", "*")
        w.Write([]byte(resp))
  }else{
    m:=HttpStatusForbiddenUnAuthorize{
    HttpStatus:403,HttpText:"Authentication Failed"}
    resp, _ := json.Marshal(m)
    w.WriteHeader(http.StatusForbidden)
    w.Header().Set("Content-Type", "application/json")
    w.Header().Set("charset", "utf-8")
    w.Header().Set("Access-Control-Allow-Origin", "*")
    w.Write([]byte(resp))
  }


}

func notFound(w http.ResponseWriter, r *http.Request)  {
  http.Redirect(w, r, "/login", 301)
}

type LinkedFaceRequestToPratibha struct {
  MeOw string
  CurrentG string
  SocialOpt bool
  ProfessionalOpt bool
  ConsumerOpt bool
  ReaderOpt bool
  DaterOpt bool
  DonatorOpt bool
  SinnerOpt bool
  Pkey string
}

func serveLinkMyFaceWithThemPratibhaPlease(w http.ResponseWriter, r *http.Request)  {
  r.ParseForm()
  data:=r.Form.Get("data")
  fmt.Println(data)
  var FaceContainerOpt LinkedFaceRequestToPratibha
  json.Unmarshal([]byte(data),&FaceContainerOpt)
  fmt.Println(FaceContainerOpt.MeOw)

  insaan:=FaceContainerOpt.MeOw
  mainInsaan:=FaceContainerOpt.Pkey
  fmt.Println("aaj se tera")
  fmt.Println(mainInsaan)
  mainInsaanSaasur, _ := strconv.ParseInt(mainInsaan, 10, 64)
  rest:=authorizeThisFace(&insaan,&mainInsaanSaasur)

  if rest!="ERROR" {
    w.Header().Set("Content-Type", "application/json")
    w.Header().Set("charset", "utf-8")
    w.Header().Set("Access-Control-Allow-Origin", "*")
    w.Write([]byte(data))
  }else{
    m:=HttpStatusForbiddenUnAuthorize{
    HttpStatus:403,HttpText:"Authentication Failed"}
    resp, _ := json.Marshal(m)
    w.WriteHeader(http.StatusForbidden)
    w.Header().Set("Content-Type", "application/json")
    w.Header().Set("charset", "utf-8")
    w.Header().Set("Access-Control-Allow-Origin", "*")
    w.Write([]byte(resp))
  }

}

type ForgiveMePratibhaPlease struct {
  All []interface{}
}

func getMyRelatedActionFacePratibhaSorry(aadat *chan []interface{},id *string,referableId *string)  {
  s1:="Match (ee:Rita)-[:LINKED]->(ff:Rita) where ff.name =~'.*(?i)"
  b:=*id
  s2:=".*' and ee.email ='"
  bPyar:=*referableId
  sPyar:="' with ff optional Match (ff)-[:HAS_MANDATORY_DP]->(gg:ProfilePic) return collect(distinct [ff.name,ff.email,gg.title])[0..10] as all"

  s17 := fmt.Sprint(s1,b,s2,bPyar,sPyar)
  fmt.Println(s17)
  driver := bolt.NewDriver()
    conn, err := driver.OpenNeo("bolt://rita:b.PuhuqVThYfCn.fvurl1e25g7fzyCI@hobby-panhpmpgjildgbkepcdcklol.dbs.graphenedb.com:24786?tls=true")
    defer conn.Close()
    dur,_:=time.ParseDuration("600s")
    conn.SetTimeout(dur)
    if err != nil {
      panic(err)
    }
    dataOf, _, _, _ := conn.QueryNeoAll(s17, nil)
    *aadat<-dataOf[0][0].([]interface{})
}

func getMyCommonActionFacePratibhaHandler(w http.ResponseWriter, r *http.Request)  {
  vars := mux.Vars(r)
  mkiloj:=vars["query"]
  referableId:=vars["reference_id"]
  refPub:=vars["Pkey"]
  fmt.Println(referableId)
  fmt.Println(mkiloj)

  insaan:=referableId
  mainInsaan:=refPub
  fmt.Println("aaj se tera")
  fmt.Println(mainInsaan)
  mainInsaanSaasur, _ := strconv.ParseInt(mainInsaan, 10, 64)
  rest:=authorizeThisFace(&insaan,&mainInsaanSaasur)

  if rest!="ERROR" {
    chnlface := make(chan []interface{})
    go getMyRelatedActionFacePratibhaSorry(&chnlface,&mkiloj,&rest)
    //var HamsafarFacePratibha ForgiveMePratibhaPlease
    var IloveMyJindagi []interface{}
    for i := 0; i < 1; i++ {
      select {
        case msg := <-chnlface:
              IloveMyJindagi=msg
        }
    }
    m:=ForgiveMePratibhaPlease{
    All:IloveMyJindagi}
    resp, _ := json.Marshal(m)
    fmt.Println(string(resp))
    w.Header().Set("Content-Type", "application/json")
    w.Header().Set("charset", "utf-8")
    w.Header().Set("Access-Control-Allow-Origin", "*")
    w.Write([]byte(resp))
  }else{
    m:=HttpStatusForbiddenUnAuthorize{
    HttpStatus:403,HttpText:"Authentication Failed"}
    resp, _ := json.Marshal(m)
    w.WriteHeader(http.StatusForbidden)
    w.Header().Set("Content-Type", "application/json")
    w.Header().Set("charset", "utf-8")
    w.Header().Set("Access-Control-Allow-Origin", "*")
    w.Write([]byte(resp))
  }
}

type MyFacesSearchAnyOneHolder struct {
  All []interface{}
}

func serveFaceTemplate(w http.ResponseWriter, r *http.Request) {

vars := mux.Vars(r)
mkiloj:=vars["query"]
fmt.Println(mkiloj)
s1:="Match (ee:Rita) where ee.name =~'.*."
b:=mkiloj
s2:=".*' with ee optional match (ee)-[:HAS_MANDATORY_DP]->(ff:ProfilePic) return collect(distinct [ee.name,ee.email,ff.title]) as all"
s17 := fmt.Sprint(s1,b,s2)
fmt.Println(s17)
driver := bolt.NewDriver()
  conn, err := driver.OpenNeo("bolt://rita:b.PuhuqVThYfCn.fvurl1e25g7fzyCI@hobby-panhpmpgjildgbkepcdcklol.dbs.graphenedb.com:24786?tls=true")
  defer conn.Close()
  dur,_:=time.ParseDuration("600s")
  conn.SetTimeout(dur)
  if err != nil {
    panic(err)
  }
  dataOf, _, _, _ := conn.QueryNeoAll(s17, nil)
  //*aadat<-dataOf[0][0].([]interface{})
  m:=MyFacesSearchAnyOneHolder{
  All:dataOf[0][0].([]interface{})}
  resp, _ := json.Marshal(m)
  fmt.Println(string(resp))
  w.Header().Set("Content-Type", "application/json")
  w.Header().Set("charset", "utf-8")
  w.Header().Set("Access-Control-Allow-Origin", "*")
  w.Write([]byte(resp))
}

type SaveMyYouTubeUpdatePratibhaPlease struct {
  Me string
  Playlist []string
  SongName string
  SongId string
  Pkey string
}

type PlaylistStatusCodeResponse struct {
  Status string
}

func satayeMenuKyonOmyYaraIloveYouHandler(w http.ResponseWriter, r *http.Request)  {
  r.ParseForm()
  data:=r.Form.Get("data")
  fmt.Println(data)
  var SavableResult SaveMyYouTubeUpdatePratibhaPlease
  json.Unmarshal([]byte(data),&SavableResult)
  fmt.Println(SavableResult.Playlist)


  insaan:=SavableResult.Me
  mainInsaan:=SavableResult.Pkey
  fmt.Println("aaj se tera")
  fmt.Println(mainInsaan)
  mainInsaanSaasur, _ := strconv.ParseInt(mainInsaan, 10, 64)
  rest:=authorizeThisFace(&insaan,&mainInsaanSaasur)


  if rest!="ERROR" {
    s1:="match (ee:Rita)-[:HAS_YOUTUBE_PLAYLIST]->(ff:YoutubePlaylist) where ee.email='"
    cvIn:=rest
    cv4:="' and ff.title in ["
    b10:=""
    for i5,val5:=range SavableResult.Playlist {
      if i5==len(SavableResult.Playlist)-1 {
        b10=b10+"'"+val5+"'"
      }else{
        b10=b10+"'"+val5+"',"
      }
    }
    cvSabTera:="] merge (ff)-[:ACTUAL_PLAYLIST_SONG]->(mm:PlaylistSong{name:'"
    cvMera:=SavableResult.SongName
    cvAller:="',source:'"
    cvCaller:=SavableResult.SongId
    cvFaller:="'}) return ee,ff,mm"
    s12 := fmt.Sprint(s1,cvIn,cv4,b10,cvSabTera,cvMera,cvAller,cvCaller,cvFaller)
    fmt.Println(s12)
    driver := bolt.NewDriver()
      conn, err := driver.OpenNeo("bolt://rita:b.PuhuqVThYfCn.fvurl1e25g7fzyCI@hobby-panhpmpgjildgbkepcdcklol.dbs.graphenedb.com:24786?tls=true")
      defer conn.Close()
      dur,_:=time.ParseDuration("600s")
      conn.SetTimeout(dur)
      if err != nil {
        //panic(err)
      }
      _, _, _, errx := conn.QueryNeoAll(s12, nil)
      if errx != nil {
        //panic(err)
        m:=PlaylistStatusCodeResponse{
        Status:"NO"}
        resp, _ := json.Marshal(m)
        w.Header().Set("Content-Type", "application/json")
        w.Header().Set("charset", "utf-8")
        w.Header().Set("Access-Control-Allow-Origin", "*")
        w.Write([]byte(resp))
      }else {
        m:=PlaylistStatusCodeResponse{
        Status:"OK"}
        resp, _ := json.Marshal(m)
        w.Header().Set("Content-Type", "application/json")
        w.Header().Set("charset", "utf-8")
        w.Header().Set("Access-Control-Allow-Origin", "*")
        w.Write([]byte(resp))
      }
  }else{
    m:=HttpStatusForbiddenUnAuthorize{
    HttpStatus:403,HttpText:"Authentication Failed"}
    resp, _ := json.Marshal(m)
    w.WriteHeader(http.StatusForbidden)
    w.Header().Set("Content-Type", "application/json")
    w.Header().Set("charset", "utf-8")
    w.Header().Set("Access-Control-Allow-Origin", "*")
    w.Write([]byte(resp))
  }

}

type LetMeGetAllSongFromPlaylistPratibha struct {
  Me string
  Playlist string
  Pkey string
}

type ReturnableResultPlaylistToClient struct {
  All []interface{}
}

func letMeSeeMyPlaylistSongPratibhaPleaseHandler(w http.ResponseWriter, r *http.Request)  {
  r.ParseForm()
  data:=r.Form.Get("data")
  fmt.Println(data)
  var GetableableResult LetMeGetAllSongFromPlaylistPratibha
  json.Unmarshal([]byte(data),&GetableableResult)
  fmt.Println(GetableableResult.Playlist)
  insaan:=GetableableResult.Me
  mainInsaan:=GetableableResult.Pkey
  fmt.Println("aaj se tera")
  fmt.Println(mainInsaan)
  mainInsaanSaasur, _ := strconv.ParseInt(mainInsaan, 10, 64)
  rest:=authorizeThisFace(&insaan,&mainInsaanSaasur)
  if rest!="ERROR" {
    s1:="match (ee:Rita)-[:HAS_YOUTUBE_PLAYLIST]->(ff:YoutubePlaylist) where ee.email='"
    cvIn:=rest
    cv4:="' and ff.title='"
    cvAller:=GetableableResult.Playlist
    cvTaller:="' with ff optional match (ff)-[:ACTUAL_PLAYLIST_SONG]->(gg:PlaylistSong) return collect(distinct [gg.name,gg.source]) as All"
    s12 := fmt.Sprint(s1,cvIn,cv4,cvAller,cvTaller)
    fmt.Println(s12)
    driver := bolt.NewDriver()
      conn, err := driver.OpenNeo("bolt://rita:b.PuhuqVThYfCn.fvurl1e25g7fzyCI@hobby-panhpmpgjildgbkepcdcklol.dbs.graphenedb.com:24786?tls=true")
      defer conn.Close()
      dur,_:=time.ParseDuration("600s")
      conn.SetTimeout(dur)
      if err != nil {
        //panic(err)
      }
      dataF, _, _, _ := conn.QueryNeoAll(s12, nil)
      m:=ReturnableResultPlaylistToClient{
      All:dataF[0][0].([]interface{})}
      resp, _ := json.Marshal(m)
      w.Header().Set("Content-Type", "application/json")
      w.Header().Set("charset", "utf-8")
      w.Header().Set("Access-Control-Allow-Origin", "*")
      w.Write([]byte(resp))
  }else{
    m:=HttpStatusForbiddenUnAuthorize{
    HttpStatus:403,HttpText:"Authentication Failed"}
    resp, _ := json.Marshal(m)
    w.WriteHeader(http.StatusForbidden)
    w.Header().Set("Content-Type", "application/json")
    w.Header().Set("charset", "utf-8")
    w.Header().Set("Access-Control-Allow-Origin", "*")
    w.Write([]byte(resp))
  }
}


type DateRequestorPratibha struct {
  Me string
  Target string
  EmotionalMessage string
  Pkey string
}


type ReturnableDateReqTempPratibha struct {
  All []interface{}
}

func sendMyDateRequestPratibhaPleaseHandler(w http.ResponseWriter, r *http.Request)  {
  r.ParseForm()
  data:=r.Form.Get("data")
  fmt.Println(data)
  var GetableableDateReqData DateRequestorPratibha
  json.Unmarshal([]byte(data),&GetableableDateReqData)

  insaan:=GetableableDateReqData.Me
  mainInsaan:=GetableableDateReqData.Pkey
  fmt.Println("aaj se tera")
  fmt.Println(mainInsaan)
  mainInsaanSaasur, _ := strconv.ParseInt(mainInsaan, 10, 64)
  rest:=authorizeThisFace(&insaan,&mainInsaanSaasur)

  if rest!="ERROR" {
    fmt.Println("ghgjhgjhghjgjhghjgjgjggjhgjhgrfrerterertetertertetetertertetetetretretretettertert")
    s1:="match (ee:Rita) where ee.email='"
    sTanha:=rest
    bTanha:="' with ee optional match (ff:Rita) where ff.email='"
    bGgullak:=GetableableDateReqData.Target
    sGgullak:="' with ee,ff merge (ee)-[r:WANT_DATING{emotionalMessage:'"
    bGuru:=GetableableDateReqData.EmotionalMessage
    sGuru:="',target:'"
    bWaheG:=GetableableDateReqData.Target
    sWaheG:="'}]->(ff) return collect(distinct [r.target]) as all"

    s12 := fmt.Sprint(s1,sTanha,bTanha,bGgullak,sGgullak,bGuru,sGuru,bWaheG,sWaheG)
    fmt.Println(s12)
    driver := bolt.NewDriver()
      conn, err := driver.OpenNeo("bolt://rita:b.PuhuqVThYfCn.fvurl1e25g7fzyCI@hobby-panhpmpgjildgbkepcdcklol.dbs.graphenedb.com:24786?tls=true")
      defer conn.Close()
      dur,_:=time.ParseDuration("600s")
      conn.SetTimeout(dur)
      if err != nil {
        //panic(err)
      }
      dataF, _, _, errTThh := conn.QueryNeoAll(s12, nil)
      if errTThh != nil {
        m:=HttpStatusForbiddenUnAuthorize{
        HttpStatus:400,HttpText:"Bad request type"}
        resp, _ := json.Marshal(m)
        w.WriteHeader(http.StatusForbidden)
        w.Header().Set("Content-Type", "application/json")
        w.Header().Set("charset", "utf-8")
        w.Header().Set("Access-Control-Allow-Origin", "*")
        w.Write([]byte(resp))
      }else {
        m:=ReturnableDateReqTempPratibha{
        All:dataF[0][0].([]interface{})}
        resp, _ := json.Marshal(m)
        w.Header().Set("Content-Type", "application/json")
        w.Header().Set("charset", "utf-8")
        w.Header().Set("Access-Control-Allow-Origin", "*")
        w.Write([]byte(resp))
      }

  }else{
    m:=HttpStatusForbiddenUnAuthorize{
    HttpStatus:403,HttpText:"Authentication Failed"}
    resp, _ := json.Marshal(m)
    w.WriteHeader(http.StatusForbidden)
    w.Header().Set("Content-Type", "application/json")
    w.Header().Set("charset", "utf-8")
    w.Header().Set("Access-Control-Allow-Origin", "*")
    w.Write([]byte(resp))
  }
}


func getMyPostForProfilePratibhaPleaseHandler(w http.ResponseWriter, r *http.Request)  {
  r.ParseForm()
  data:=r.Form.Get("data")
  fmt.Println(data)
}


type HoledMeTempToTheServe struct {
  Me string
  MyAllLinkedFaces []interface{}
}




func replicateMyMessageToDbMongoHandler(w http.ResponseWriter, r *http.Request)  {
  r.ParseForm()
  data:=r.Form.Get("data")
  fmt.Println(data)
  /*s1:="Match (ee:Rita)-[:LINKED]-(ff:Rita) where ff.name =~'.*(?i)"
  b:=*id
  s2:=".*' and ee.email <> '"
  bPyar:=*referableId
  sPyar:="' with ff optional Match (ff)-[:HAS_MANDATORY_DP]->(gg:ProfilePic) return collect(distinct [ff.name,ff.email,gg.title])[0..10] as all union "
  s17 := fmt.Sprint(s1,b,s2,bPyar,sPyar)
  fmt.Println(s17)
  driver := bolt.NewDriver()
    conn, err := driver.OpenNeo("bolt://rita:b.PuhuqVThYfCn.fvurl1e25g7fzyCI@hobby-panhpmpgjildgbkepcdcklol.dbs.graphenedb.com:24786?tls=true")
    defer conn.Close()
    if err != nil {
      panic(err)
    }
    dataOf, _, _, _ := conn.QueryNeoAll(s17, nil)
    *aadat<-dataOf[0][0].([]interface{})*/
}

type MessageDataAllerHolder struct {
  To string
  From string
  Content string
  Domain string
  When string
  Pkey string
}

type MessageDataAllerHolderCopy struct {
  To string
  From string
  Content string
  Domain string
  When string
  Pkey string
}


type ReturnableMessageAjaxable struct {
  All []interface{}
}

func saveMyMessagePratibhaPleaseHelpHandler(w http.ResponseWriter, r *http.Request)  {
  r.ParseForm()
  data:=r.Form.Get("data")
  fmt.Println(data)
  var MessageDataStorableGraph MessageDataAllerHolder
  json.Unmarshal([]byte(data),&MessageDataStorableGraph)
  fmt.Println(MessageDataStorableGraph.From)

  insaan:=MessageDataStorableGraph.From
  mainInsaan:=MessageDataStorableGraph.Pkey
  fmt.Println("aaj se tera")
  fmt.Println(mainInsaan)
  mainInsaanSaasur, _ := strconv.ParseInt(mainInsaan, 10, 64)
  rest:=authorizeThisFace(&insaan,&mainInsaanSaasur)


  if rest!="ERROR" {
    s1:="Match (ee:Rita)-[:HAS_MANDATORY_DP]->(ff:ProfilePic) where ee.email ='"
    b:=rest
    s2:="' with ee,ff merge (ee)-[:SENT_MESSAGE]-(gg:MessageData{to:'"
    bTenu:=MessageDataStorableGraph.To
    sTenu:="',content:'"
    bMenu:=MessageDataStorableGraph.Content
    sMenu:="',domain:'"
    bKismat:=MessageDataStorableGraph.Domain
    sKismat:="',when:'"
    bHoho:=MessageDataStorableGraph.When
    bBadla:="',from:'"
    bBewafa:=rest
    sBewafa:="'}) return collect([gg.to,ee.email,gg.content,gg.domain,ff.title,gg.when]) as all "
    s17 := fmt.Sprint(s1,b,s2,bTenu,sTenu,bMenu,sMenu,bKismat,sKismat,bHoho,bBadla,bBewafa,sBewafa)
    fmt.Println(s17)
    driver := bolt.NewDriver()
      conn, err := driver.OpenNeo("bolt://rita:b.PuhuqVThYfCn.fvurl1e25g7fzyCI@hobby-panhpmpgjildgbkepcdcklol.dbs.graphenedb.com:24786?tls=true")
      defer conn.Close()
      dur,_:=time.ParseDuration("600s")
      conn.SetTimeout(dur)
      if err != nil {
        panic(err)
      }
      dataOf, _, _, _ := conn.QueryNeoAll(s17, nil)
      //*aadat<-dataOf[0][0].([]interface{})
      m:=ReturnableMessageAjaxable{
      All:dataOf[0][0].([]interface{})}
      resp, _ := json.Marshal(m)
      w.Header().Set("Content-Type", "application/json")
      w.Header().Set("charset", "utf-8")
      w.Header().Set("Access-Control-Allow-Origin", "*")
      w.Write([]byte(resp))
  }else{
    m:=HttpStatusForbiddenUnAuthorize{
    HttpStatus:403,HttpText:"Authentication Failed"}
    resp, _ := json.Marshal(m)
    w.WriteHeader(http.StatusForbidden)
    w.Header().Set("Content-Type", "application/json")
    w.Header().Set("charset", "utf-8")
    w.Header().Set("Access-Control-Allow-Origin", "*")
    w.Write([]byte(resp))
  }


}

type MessageSourceForMe struct {
  Me string
  Domain string
  PKey string
}

type MessageSourceDataAjaxablePratibhaType struct {
  All []interface{}
}

type HttpStatusForbiddenUnAuthorize struct{
  HttpStatus int
  HttpText string
}

func retriveMyEventSourceForMessagePratibhaHandler(w http.ResponseWriter, r *http.Request)  {
  r.ParseForm()
  data:=r.Form.Get("data")
  fmt.Println(data)
  var TuNahiHaiTeriYadeMujheSatatiHai MessageSourceForMe
  json.Unmarshal([]byte(data),&TuNahiHaiTeriYadeMujheSatatiHai)
  fmt.Println(TuNahiHaiTeriYadeMujheSatatiHai.Me)
  insaan:=TuNahiHaiTeriYadeMujheSatatiHai.Me
  mainInsaan:=TuNahiHaiTeriYadeMujheSatatiHai.PKey
  fmt.Println("aaj se tera")
  fmt.Println(mainInsaan)
  mainInsaanSaasur, _ := strconv.ParseInt(mainInsaan, 10, 64)
  rest:=authorizeThisFace(&insaan,&mainInsaanSaasur)
  if rest!="ERROR" {
    s1:="Match (ee:Rita)-[:SENT_MESSAGE]->(ff:MessageData) where ee.email ='"
    b:=rest
    s2:="' and ff.domain='"
    bDillagi:=TuNahiHaiTeriYadeMujheSatatiHai.Domain
    sDillagi:="' with ee,ff optional match (ee)-[:HAS_MANDATORY_DP]->(gg:ProfilePic) return collect(distinct [ee.email,gg.title,ff.to,ff.from,ff.content,ff.domain,ff.when]) as all union match (ee:Rita)-[:SENT_MESSAGE]->(ff:MessageData) where ff.to='"
    bBewafa:=rest
    sBewafa:="' and ff.domain='"
    bJuda:=TuNahiHaiTeriYadeMujheSatatiHai.Domain
    sJuda:="' with ee,ff optional match (ee)-[:HAS_MANDATORY_DP]->(gg:ProfilePic) return collect(distinct [ee.email,gg.title,ff.to,ff.from,ff.content,ff.domain,ff.when]) as all"
    s17 := fmt.Sprint(s1,b,s2,bDillagi,sDillagi,bBewafa,sBewafa,bJuda,sJuda)
    fmt.Println(s17)
    driver := bolt.NewDriver()
      conn, err := driver.OpenNeo("bolt://rita:b.PuhuqVThYfCn.fvurl1e25g7fzyCI@hobby-panhpmpgjildgbkepcdcklol.dbs.graphenedb.com:24786?tls=true")
      defer conn.Close()
      dur,_:=time.ParseDuration("600s")
      conn.SetTimeout(dur)
      if err != nil {
        panic(err)
      }
      dataOf, _, _, _ := conn.QueryNeoAll(s17, nil)
      //*aadat<-dataOf[0][0].([]interface{})
      var sealApp []interface{}
      if len(dataOf)==0 || len(dataOf)==1 {
        sealApp = dataOf[0][0].([]interface{})
      }else{
        sealApp = append(dataOf[0][0].([]interface{}), dataOf[1][0].([]interface{})...)
      }
      m:=MessageSourceDataAjaxablePratibhaType{
      All:sealApp}
      resp, _ := json.Marshal(m)
      fmt.Println(string(resp))
      w.Header().Set("Content-Type", "application/json")
      w.Header().Set("charset", "utf-8")
      w.Header().Set("Access-Control-Allow-Origin", "*")
      w.Write([]byte(resp))
  }else{
    m:=HttpStatusForbiddenUnAuthorize{
    HttpStatus:403,HttpText:"Authentication Failed"}
    resp, _ := json.Marshal(m)
    w.WriteHeader(http.StatusForbidden)
    w.Header().Set("Content-Type", "application/json")
    w.Header().Set("charset", "utf-8")
    w.Header().Set("Access-Control-Allow-Origin", "*")
    w.Write([]byte(resp))
  }


}


type KhudaKeLiyeChhorDoAbYeParda struct {
  FirstUser string
  SecondUser string
  Domain string
  Me string
  Pkey string
}

func myFaceDetailsForExpandedChatPratibhaPlease(channelT *chan []interface{},idF *string,idS *string)  {
  s1:="Match (ee:Rita)-[:LINKED]->(ff:Rita) where ee.email ='"
  b:=*idF
  s2:="' and ff.email='"
  bDillagi:=*idS
  sDillagi:="' with ee,ff optional match (ee)-[:HAS_MANDATORY_DP]->(gg:ProfilePic) with ee,ff,gg optional match (ff)-[:HAS_MANDATORY_DP]->(hh:ProfilePic) return collect(distinct [[ee.email,ee.name,ee.gender,gg.title],[ff.email,ff.name,ff.gender,hh.title]]) as all union Match (ee:Rita)-[:LINKED]->(ff:Rita) where ee.email ='"
  bDhua:=*idS
  s2Dhua:="' and ff.email='"
  bDillagiDhua:=*idF
  sDillagiDhua:="' with ee,ff optional match (ee)-[:HAS_MANDATORY_DP]->(gg:ProfilePic) with ee,ff,gg optional match (ff)-[:HAS_MANDATORY_DP]->(hh:ProfilePic) return collect(distinct [[ee.email,ee.name,ee.gender,gg.title],[ff.email,ff.name,ff.gender,hh.title]]) as all"
  s17 := fmt.Sprint(s1,b,s2,bDillagi,sDillagi,bDhua,s2Dhua,bDillagiDhua,sDillagiDhua)
  fmt.Println(s17)
  driver := bolt.NewDriver()
    conn, err := driver.OpenNeo("bolt://rita:b.PuhuqVThYfCn.fvurl1e25g7fzyCI@hobby-panhpmpgjildgbkepcdcklol.dbs.graphenedb.com:24786?tls=true")
    defer conn.Close()
    dur,_:=time.ParseDuration("600s")
    conn.SetTimeout(dur)
    if err != nil {
      panic(err)
    }
    dataOf, _, _, _ := conn.QueryNeoAll(s17, nil)
    var sealApp []interface{}
    if len(dataOf)==0 || len(dataOf)==1 {
      sealApp = dataOf[0][0].([]interface{})
    }else{
      sealApp = append(dataOf[0][0].([]interface{}), dataOf[1][0].([]interface{})...)
    }
    *channelT<-sealApp
}

func tuneAisiPilayiMazaAaGayaChatDataPratibha(channelT *chan []interface{},firstUser *string,secondUser *string,domain *string)  {
  s1:="Match (ee:Rita)-[:SENT_MESSAGE]->(ff:MessageData) where ee.email ='"
  b:=*firstUser
  s2:="' and ff.domain='"
  bDillagi:=*domain
  sDillagi:="' and ff.to='"
  aaEmp:=*secondUser
  bbEmp:="' with ee,ff optional match (ee)-[:HAS_MANDATORY_DP]->(gg:ProfilePic) return collect(distinct [ee.email,gg.title,ff.to,ff.from,ff.content,ff.domain,ff.when]) as all union match (ee:Rita)-[:SENT_MESSAGE]->(ff:MessageData) where ee.email='"
  bBewafa:=*secondUser
  sBewafa:="' and ff.domain='"
  bJuda:=*domain
  sJuda:="' and ff.to='"
  woMuskuraDe:=*firstUser
  agarHamKaheTo:="' with ee,ff optional match (ee)-[:HAS_MANDATORY_DP]->(gg:ProfilePic) return collect(distinct [ee.email,gg.title,ff.to,ff.from,ff.content,ff.domain,ff.when]) as all"
  s17 := fmt.Sprint(s1,b,s2,bDillagi,sDillagi,aaEmp,bbEmp,bBewafa,sBewafa,bJuda,sJuda,woMuskuraDe,agarHamKaheTo)
  fmt.Println(s17)
  driver := bolt.NewDriver()
    conn, err := driver.OpenNeo("bolt://rita:b.PuhuqVThYfCn.fvurl1e25g7fzyCI@hobby-panhpmpgjildgbkepcdcklol.dbs.graphenedb.com:24786?tls=true")
    defer conn.Close()
    dur,_:=time.ParseDuration("600s")
    conn.SetTimeout(dur)
    if err != nil {
      panic(err)
    }
    dataOf, _, _, _ := conn.QueryNeoAll(s17, nil)
    //*aadat<-dataOf[0][0].([]interface{})
    var sealApp []interface{}
    if len(dataOf)==0 || len(dataOf)==1 {
      sealApp = dataOf[0][0].([]interface{})
    }else{
      sealApp = append(dataOf[0][0].([]interface{}), dataOf[1][0].([]interface{})...)
    }
    *channelT<-sealApp
}

type ResponseWritableToChatAndDetailsP struct {
  Detail []interface{}
  ChatableData []interface{}
}

func tumheDillaggibhulJaniParegiYaarHandler(w http.ResponseWriter, r *http.Request)  {
  r.ParseForm()
  data:=r.Form.Get("data")
  fmt.Println(data)
  var hameshaHueDekhkarMuskuraoTo KhudaKeLiyeChhorDoAbYeParda
  json.Unmarshal([]byte(data),&hameshaHueDekhkarMuskuraoTo)
  fmt.Println(hameshaHueDekhkarMuskuraoTo.FirstUser)

  insaan:=hameshaHueDekhkarMuskuraoTo.Me
  mainInsaan:=hameshaHueDekhkarMuskuraoTo.Pkey
  fmt.Println("aaj se tera")
  fmt.Println(mainInsaan)
  mainInsaanSaasur, _ := strconv.ParseInt(mainInsaan, 10, 64)
  rest:=authorizeThisFace(&insaan,&mainInsaanSaasur)

  if rest!="ERROR" {
    chnlface := make(chan []interface{})
    chnlchat :=make(chan []interface{})
    chatableF:=hameshaHueDekhkarMuskuraoTo.FirstUser
    chatableS:=hameshaHueDekhkarMuskuraoTo.SecondUser
    chatableD:=hameshaHueDekhkarMuskuraoTo.Domain
    go myFaceDetailsForExpandedChatPratibhaPlease(&chnlface,&chatableF,&chatableS)
    go tuneAisiPilayiMazaAaGayaChatDataPratibha(&chnlchat,&chatableF,&chatableS,&chatableD)
    //var HamsafarFacePratibha ForgiveMePratibhaPlease
    var IloveMyJindagi []interface{}
    var ButPratibhaIsMyJindagi []interface{}
    for i := 0; i < 2; i++ {
      select {
        case msg := <-chnlface:
              IloveMyJindagi=msg
        case msg1 := <-chnlchat:
              ButPratibhaIsMyJindagi=msg1
        }
    }
    m:=ResponseWritableToChatAndDetailsP{
    Detail:IloveMyJindagi,ChatableData:ButPratibhaIsMyJindagi}
    resp, _ := json.Marshal(m)
    //fmt.Println(string(resp))
    w.Header().Set("Content-Type", "application/json")
    w.Header().Set("charset", "utf-8")
    w.Header().Set("Access-Control-Allow-Origin", "*")
    w.Write([]byte(resp))
  }else{
    m:=HttpStatusForbiddenUnAuthorize{
    HttpStatus:403,HttpText:"Authentication Failed"}
    resp, _ := json.Marshal(m)
    w.WriteHeader(http.StatusForbidden)
    w.Header().Set("Content-Type", "application/json")
    w.Header().Set("charset", "utf-8")
    w.Header().Set("Access-Control-Allow-Origin", "*")
    w.Write([]byte(resp))
  }

}

func generatePublicKey(channelPost *chan int64,low int64, hi int64)  {
  ll:=low + R.Int63n(hi-low)
  *channelPost<-ll
}

func generatePrivateKey(channelPost *chan int64,low int64, hi int64)  {
  ll:=low + R.Int63n(hi-low)
  *channelPost<-ll
}

func generateHashAndReplicateToDbPratibhaPlease(emailData *string) []string {
  	v:=time.Now().UnixNano()
  	R.Seed(int64(v))
    tempEm:=*emailData
    var publicKey,privateKey int64
    chnlPublic := make(chan int64)
    chnlPrivate :=make(chan int64)
    go generatePublicKey(&chnlPublic,1000000000000000,9999999999999999)
  	go generatePrivateKey(&chnlPrivate,1000000000000000,9999999999999999)
    for i := 0; i < 2; i++ {
      select {
        case msg := <-chnlPublic:
              publicKey=msg
        case msg1 := <-chnlPrivate:
              privateKey=msg1
        }
    }
    retArrStr:=replicateRSAToDb(&publicKey,&privateKey,&tempEm)
    fmt.Println("returning hash")
    fmt.Println(retArrStr)
    return retArrStr
}

type TestAuth struct {
  Me string
  PublicKey string
}


func replicateRSAToDb(publicKey *int64,privateKey *int64,emailData *string) []string {
  randColl:=fmt.Sprintf("%d%d",*publicKey,*privateKey)
  //vvb:=*publicKey
  fmt.Println(randColl)
  text := []byte(*emailData)
  key := []byte(randColl)
  ciphertext, err := encrypt(text, key)
  if err != nil {
      log.Fatal(err)
  }
  fmt.Printf("%s => %x\n", text, ciphertext)
  mEnc:=fmt.Sprintf("%x",ciphertext)
  fmt.Println(ciphertext)
  s1:="CREATE (ee:Auth{hash:'"
  b:=mEnc
  sDillagi:="',privateKey:'"
  sIknow:=*privateKey
  okLook:="'}) return collect([ee.hash]) as all"

  s17 := fmt.Sprint(s1,b,sDillagi,sIknow,okLook)
  fmt.Println(s17)
  driver := bolt.NewDriver()
    conn, err := driver.OpenNeo("bolt://rita:b.PuhuqVThYfCn.fvurl1e25g7fzyCI@hobby-panhpmpgjildgbkepcdcklol.dbs.graphenedb.com:24786?tls=true")
    defer conn.Close()
    dur,_:=time.ParseDuration("600s")
    conn.SetTimeout(dur)
    if err != nil {
      panic(err)
    }
    dataOf, _, _, _ := conn.QueryNeoAll(s17, nil)
    nkk:=dataOf[0][0].([]interface{})
    fmt.Println(nkk)
    fmt.Println("authorizing face")
    secArr:=make([]string, 0)
    for _,ght:=range nkk{
      for _,dnt:=range ght.([]interface{}){
        secArr=append(secArr,dnt.(string))
      }
    }
    strPub:=fmt.Sprintf("%d",*publicKey)
    secArr=append(secArr,strPub)

    //fmt.Println(secArr)
    //fmt.Println("hello")
    //dex,_:=decrypt(ciphertext,key)
    //fmt.Println(string(dex))
    //fmt.Println()

    //authorizedHash:=authorizeThisFace(&ciphertext,&vvb)
    //fmt.Println("dekhte dekhte")
    //fmt.Println(authorizedHash)
    fmt.Println("returning same hash")
    fmt.Println(secArr)
    return secArr
}

func authorizeThisFace(hash *string,publicKey *int64) string {
  //ll:=fmt.Sprintf("%s",string(*hash))
  //fmt.Println(ll)
  s1:="Match (ee:Auth) where ee.hash='"
  b:=*hash
  okLook:="' return collect([ee.hash,ee.privateKey]) as all"

  s17 := fmt.Sprint(s1,b,okLook)
  fmt.Println(s17)
  eRRTemp:=""
  var fDecrypt string
  driver := bolt.NewDriver()
    conn, err := driver.OpenNeo("bolt://rita:b.PuhuqVThYfCn.fvurl1e25g7fzyCI@hobby-panhpmpgjildgbkepcdcklol.dbs.graphenedb.com:24786?tls=true")
    defer conn.Close()
    dur,_:=time.ParseDuration("600s")
    conn.SetTimeout(dur)
    if err != nil {
      eRRTemp="ERROR"
    }else{
      dataOf, _, _, errxA := conn.QueryNeoAll(s17, nil)
      if errxA != nil {
        eRRTemp="ERROR"
      }else{
        if len(dataOf)>0 {
          nkk:=dataOf[0][0].([]interface{})
          fmt.Println("authorizable face")
          fmt.Println(nkk)
          //mkk:=nkk[0]
          //fmt.Println(mkk)
          secArr:=make([]string, 0)
          for _,ght:=range nkk{
            for _,dnt:=range ght.([]interface{}){
              secArr=append(secArr,dnt.(string))
            }
          }
          fmt.Println(secArr)
          if len(secArr)>0 {
            kli, errL := strconv.ParseInt(secArr[1], 10, 64)
            if errL!=nil {
              eRRTemp="ERROR"
            }else{
              fmt.Println(*publicKey)
              fmt.Println(kli)
              newKeyForHash:=fmt.Sprintf("%d%d",*publicKey,kli)
              fmt.Println(newKeyForHash)
              fmt.Println(*hash)
              sttr:=*hash
              fmt.Println([]byte(sttr))
              fmt.Println("your value")
              decryptedData,errM:=decrypt(sttr,newKeyForHash)
              if errM!=nil {
                eRRTemp="ERROR"
              }else{
                fmt.Println("kusur hai")
                fmt.Println(eRRTemp)
                fDecrypt=string(decryptedData)
              }
            }
          }else{
            eRRTemp="ERROR"
          }
        }else{
          eRRTemp="ERROR"
        }
      }
    }
    if fDecrypt!="" {
      if eRRTemp=="ERROR" {
        return eRRTemp
      }else{
        return fDecrypt
      }
    }else{
      return eRRTemp
    }
}


func myTokenAuthTest(w http.ResponseWriter, r *http.Request)  {
  r.ParseForm()
  data:=r.Form.Get("data")
  fmt.Println(data)
}

func makePaypalPayment() (int,error) {
  client, err := paypalsdk.NewClient("AR6ERqJofXz37bYK4pbuUmSsH5J8B6is1LMCubv8KKaTn_wpFk7TPVkIU7R-01wcMFJXTYUhtzNXosqo", "ENiKdPNAW-DxK5cAot4mxTnFM4X0WXa_cg3gXhXC1d1EPkMBLKWypqOLohFxg5Hoo1YbusnaTOY4EHdQ", paypalsdk.APIBaseSandBox)
client.SetLog(os.Stdout) // Set log to terminal stdout

accessToken, err := client.GetAccessToken()
fmt.Println(accessToken)
if err!=nil {
  return 0,err
}

/*
*/
//fmt.Println(accessToken)
  p := paypalsdk.Payment{
  Intent: "sale",
  Payer: &paypalsdk.Payer{
      PaymentMethod: "credit_card",
      FundingInstruments: []paypalsdk.FundingInstrument{paypalsdk.FundingInstrument{
          CreditCard: &paypalsdk.CreditCard{
            Number:      "4375514342139002",
            Type:        "visa",
            ExpireMonth: "02",
            ExpireYear:  "2020",
            CVV2:        "123",
            FirstName:   "AJAY",
            LastName:    "KUMAR",
          },
      }},
  },
  Transactions: []paypalsdk.Transaction{paypalsdk.Transaction{
      Amount: &paypalsdk.Amount{
          Currency: "USD",
          Total:    "0.01",
      },
      Description: "My Payment",
  }},
  RedirectURLs: &paypalsdk.RedirectURLs{
      ReturnURL: "http://localhost:5000/paypal_success",
      CancelURL: "http://localhost:5000/paypal_cancel",
  },
}
  //paymentResponse, err := client.CreatePayment(p)
  _, err1 := client.CreatePayment(p)
  if err!=nil {
    return 0,err1
  }
  //fmt.Println(paymentResponse)
  return 200,nil
}

type HoldMyTempDataRelatedMyJobPratibha struct {
  Name string
  Type []string
  StartDate string
  Mobile string
  PinCode []string
  AddressOne string
  AddressTwo string
  OpenTime string
  CloseTime string
  Website string
  RegistrationNumber string
  Payscale string
  ExtraBenefit string
  Onsite string
  Transport string
  Insurance string
  RegularPerk string
  OfficeTrip string
  OnsiteCountry string
  Me string
  Pkey string
  JobIcon string
}


func createMyJobPratibhaPleaseForMeSorryHandler(w http.ResponseWriter, r *http.Request)  {
  r.ParseForm()
  data:=r.Form.Get("data")
  fmt.Println(data)
  var hameshaHueDekhkarMuskuraoTo HoldMyTempDataRelatedMyJobPratibha
  json.Unmarshal([]byte(data),&hameshaHueDekhkarMuskuraoTo)
  fmt.Println(hameshaHueDekhkarMuskuraoTo.Name)

  insaan:=hameshaHueDekhkarMuskuraoTo.Me
  mainInsaan:=hameshaHueDekhkarMuskuraoTo.Pkey
  fmt.Println("aaj se tera")
  fmt.Println(mainInsaan)
  mainInsaanSaasur, _ := strconv.ParseInt(mainInsaan, 10, 64)
  rest:=authorizeThisFace(&insaan,&mainInsaanSaasur)

  if rest!="ERROR" {


    s1:="Match (ee:Rita) where ee.email='"
    b1:=rest
    s2:="' MERGE (ee)-[:HAS_JOB]-(ee:JobData{jobName:'"
    b2:=hameshaHueDekhkarMuskuraoTo.Name
    s3:="',jobType:["
    b3:=""
    for i3,val3:=range hameshaHueDekhkarMuskuraoTo.Type {
      if i3==len(hameshaHueDekhkarMuskuraoTo.Type)-1 {
        b3=b3+"'"+val3+"'"
      }else{
        b3=b3+"'"+val3+"',"
      }
    }
    sPyar:="],startDate:'"
    bPyar:=hameshaHueDekhkarMuskuraoTo.StartDate
    sPyar1:="',mobileNumber:'"
    bPyar1:=hameshaHueDekhkarMuskuraoTo.Mobile
    s4:="',pinCode:["
    b4:=""
    for i4,val4:=range hameshaHueDekhkarMuskuraoTo.PinCode {
      if i4==len(hameshaHueDekhkarMuskuraoTo.PinCode)-1 {
        b4=b4+"'"+val4+"'"
      }else{
        b4=b4+"'"+val4+"',"
      }
    }
    s5:="],jobAddressOne:'"
    b5:=hameshaHueDekhkarMuskuraoTo.AddressOne
    s6:="',jobAddressTwo:'"
    b6:=hameshaHueDekhkarMuskuraoTo.AddressTwo
    s7:="',jobOpenTime:'"
    b7:=hameshaHueDekhkarMuskuraoTo.OpenTime
    s8:="',jobCloseTime:'"
    b8:=hameshaHueDekhkarMuskuraoTo.CloseTime
    s9:="',jobWebsite:'"
    b9:=hameshaHueDekhkarMuskuraoTo.Website
    s10:="',jobRegNo:'"
    b11:=hameshaHueDekhkarMuskuraoTo.RegistrationNumber
    s12:="',jobPayScale:'"
    b12:=hameshaHueDekhkarMuskuraoTo.Payscale
    s13:="',extraBenefit:'"
    b13:=hameshaHueDekhkarMuskuraoTo.ExtraBenefit
    s14:="',transportFacility:'"
    b14:=hameshaHueDekhkarMuskuraoTo.Transport
    s15:="',healthInsurance:'"
    b15:=hameshaHueDekhkarMuskuraoTo.Insurance
    s16:="',regularPerk:'"
    b16:=hameshaHueDekhkarMuskuraoTo.RegularPerk
    s17:="',officeTrip:'"
    b17:=hameshaHueDekhkarMuskuraoTo.OfficeTrip
    s18:="',onsite:'"
    b18:=hameshaHueDekhkarMuskuraoTo.Onsite
    s19:="',onsiteCountry:'"
    b19:=hameshaHueDekhkarMuskuraoTo.OnsiteCountry
    s20:="',jobIcon:'"
    b20:=hameshaHueDekhkarMuskuraoTo.JobIcon
    s21:="'}) return collect([ee.ItemName,ee.ItemType,ee.ItemCategory,ee.ItemQuantity,ee.ItemPrice,ee.MfgDate,ee.ExpDate,ee.MaxAvailability,ee.IconArray,ee.Me,id(ee)]) as all"
    s22 := fmt.Sprint(s1,b1,s2,b2,s3,b3,sPyar,bPyar,sPyar1,bPyar1,s4,b4,s5,b5,s6,b6,s7,b7,s8,b8,s9,b9,s10,b11,s12,b12,s13,b13,s14,b14,s15,b15,s16,b16,s17,b17,s18,b18,s19,b19,s20,b20,s21)
    fmt.Println(s22)
    driver := bolt.NewDriver()
      conn, err := driver.OpenNeo("bolt://rita:b.PuhuqVThYfCn.fvurl1e25g7fzyCI@hobby-panhpmpgjildgbkepcdcklol.dbs.graphenedb.com:24786?tls=true")
      defer conn.Close()
      dur,_:=time.ParseDuration("600s")
      conn.SetTimeout(dur)
      if err != nil {
        panic(err)
      }
      dataOf, _, _, _ := conn.QueryNeoAll(s17, nil)
      m:=IndiItemGetResponseDbPratibha{
        All:dataOf[0][0].([]interface{})}
        resp, _ := json.Marshal(m)
        w.Header().Set("Content-Type", "application/json")
        w.Header().Set("charset", "utf-8")
        w.Header().Set("Access-Control-Allow-Origin", "*")
        w.Write([]byte(resp))


  }else{
    m:=HttpStatusForbiddenUnAuthorize{
    HttpStatus:403,HttpText:"Authentication Failed"}
    resp, _ := json.Marshal(m)
    w.WriteHeader(http.StatusForbidden)
    w.Header().Set("Content-Type", "application/json")
    w.Header().Set("charset", "utf-8")
    w.Header().Set("Access-Control-Allow-Origin", "*")
    w.Write([]byte(resp))
  }
}


type PyarTeraHaiJindagiMeriSamjhi struct {
  JobType []string
  JobLocation []string
  PayScale string
  Me string
  Pkey string
  Domain string
}


type ReturnableAllSeekJobHandler struct {
  All []interface{}
}


func letMeSeekMyJobPratibhaPleaseHandler(w http.ResponseWriter, r *http.Request)  {
  r.ParseForm()
  data:=r.Form.Get("data")
  fmt.Println(data)
  var hameshaHueDekhkarMuskuraoTo PyarTeraHaiJindagiMeriSamjhi
  json.Unmarshal([]byte(data),&hameshaHueDekhkarMuskuraoTo)
  fmt.Println(hameshaHueDekhkarMuskuraoTo.PayScale)

  insaan:=hameshaHueDekhkarMuskuraoTo.Me
  mainInsaan:=hameshaHueDekhkarMuskuraoTo.Pkey
  fmt.Println("aaj se tera")
  fmt.Println(mainInsaan)
  mainInsaanSaasur, _ := strconv.ParseInt(mainInsaan, 10, 64)
  rest:=authorizeThisFace(&insaan,&mainInsaanSaasur)
  if rest!="ERROR" {
    s1:="MERGE (ee:Rita) where ee.email ='"
    b:=rest
    s2:="' with (ee)-[:SEEKED_A_JOB]->(ff:seekData{jobType:["
    kiya:=""
    for i4,val4:=range hameshaHueDekhkarMuskuraoTo.JobType {
      if i4==len(hameshaHueDekhkarMuskuraoTo.JobType)-1 {
        kiya=kiya+"'"+val4+"'"
      }else{
        kiya=kiya+"'"+val4+"',"
      }
    }
    kiyaA:="],jobLocation:["
    college:=""
    for i5,val5:=range hameshaHueDekhkarMuskuraoTo.JobLocation {
      if i5==len(hameshaHueDekhkarMuskuraoTo.JobLocation)-1 {
        college=college+"'"+val5+"'"
      }else{
        college=college+"'"+val5+"',"
      }
    }
    collegeA:="],payScale:'"
    pScl:=hameshaHueDekhkarMuskuraoTo.PayScale
    pSclA:="',domain:'"
    firstTime:=hameshaHueDekhkarMuskuraoTo.Domain
    firUs:="'}) return collect(distinct [ee.email,ff.jobType,ff.jobLocation,ff.payScale,ff.domain]) as all"
    s17 := fmt.Sprint(s1,b,s2,kiya,kiyaA,college,collegeA,pScl,pSclA,firstTime,firUs)
    fmt.Println(s17)
    driver := bolt.NewDriver()
      conn, err := driver.OpenNeo("bolt://rita:b.PuhuqVThYfCn.fvurl1e25g7fzyCI@hobby-panhpmpgjildgbkepcdcklol.dbs.graphenedb.com:24786?tls=true")
      defer conn.Close()
      dur,_:=time.ParseDuration("600s")
      conn.SetTimeout(dur)
      if err != nil {
        panic(err)
      }
      dataOf, _, _, _ := conn.QueryNeoAll(s17, nil)
      //dataOf, _, _, _ := conn.QueryNeoAll(s17, nil)
      m:=ReturnableAllSeekJobHandler{
        All:dataOf[0][0].([]interface{})}
        resp, _ := json.Marshal(m)
        w.Header().Set("Content-Type", "application/json")
        w.Header().Set("charset", "utf-8")
        w.Header().Set("Access-Control-Allow-Origin", "*")
        w.Write([]byte(resp))
      // aadat<-dataOf[0][0].([]interface{})
      /*
        var sealApp []interface{}
        if len(dataOf)==0 || len(dataOf)==1 {
          sealApp = dataOf[0][0].([]interface{})
        }else{
          sealApp = append(dataOf[0][0].([]interface{}), dataOf[1][0].([]interface{})...)
        }
        *channelT<-sealApp
      */
  }else{
    m:=HttpStatusForbiddenUnAuthorize{
    HttpStatus:403,HttpText:"Authentication Failed"}
    resp, _ := json.Marshal(m)
    w.WriteHeader(http.StatusForbidden)
    w.Header().Set("Content-Type", "application/json")
    w.Header().Set("charset", "utf-8")
    w.Header().Set("Access-Control-Allow-Origin", "*")
    w.Write([]byte(resp))
  }
}


func sendOtpToTheMyUserAjay(mobNo string,msgVal string)  {
  accountSid := "AC9913262cf44e637b9438ab1759fc03a5"
  authToken := "c856440257ffb3cdc1d0f3b983eb6f6a"
  letterTemp:="https://api.twilio.com/2010-04-01/Accounts/"
  letterTemp1:="/Messages.json"
  urlStr := fmt.Sprintf("%s%s%s",letterTemp,accountSid,letterTemp1)

  // Create possible message bodies
  quotes:= make([]string, 0)
  //sl:="hello test officeTrip"
  quotes=append(quotes,msgVal)
  //quotes := [1]string{"We must be careful about what we pretend to be."}


  msgData := url.Values{}
  msgData.Set("To",mobNo)
  msgData.Set("From","+16413235765")
  msgData.Set("Body",quotes[0])
  fmt.Println(msgData)
  msgDataReader := *strings.NewReader(msgData.Encode())

  client := &http.Client{}
  req, _ := http.NewRequest("POST", urlStr, &msgDataReader)
  req.SetBasicAuth(accountSid, authToken)
  req.Header.Add("Accept", "application/json")
  req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
  //client := &http.Client{}


  resp, _ := client.Do(req)
  if (resp.StatusCode >= 200 && resp.StatusCode < 300) {
    var data map[string]interface{}
    decoder := json.NewDecoder(resp.Body)
    err := decoder.Decode(&data)
    if (err == nil) {
      fmt.Println(data["sid"])
    }
  } else {
    fmt.Println(resp.Status);
  }



}


type BilloniTeraLalGhaghra struct {
  Me string
  Pkey string
  Val []byte
  Ftype string
}

func changingBytesOfCloudinaryDpHandler(w http.ResponseWriter, r *http.Request)  {
    //var t BilloniTeraLalGhaghra
    //var billoniLalGhaghra uuid.UUID
    /*var Hohohohehe BilloniTeraLalGhaghra
    b, _ := ioutil.ReadAll(r.Body)
    json.Unmarshal(b,&Hohohohehe)
    defer r.Body.Close()
    fmt.Println(string(b))
    //ctx := context.Background()
	  //ctx = cloudinary.NewContext(ctx, "cloudinary://864654217542164:fdQqxrCeKl_OJdwR84Bw9LhuUhM@hnruvsvqz")
	  //data, _ := ioutil.ReadFile("<imageFile>")
    fftYYY:=Hohohohehe.Ftype
    nm:=uuid.New()
    tatarara:=strings.Split(fftYYY,"/")[1]
    tamarataHH:="."
    ffTT:=fmt.Sprintf("%s%s%s",nm,tamarataHH,tatarara)
    fmt.Println(ffTT)*/
    //ser, _:=cloudinary.Dial()
    //fmt.Println(ser)


	  //cloudinary.UploadStaticImage(ctx, ffTT, bytes.NewBuffer(Hohohohehe.Val))

    /*msgData := url.Values{}
    msgData.Set("file",mobNo)
    msgData.Set("From","+16413235765")
    msgData.Set("Body",quotes[0])
    fmt.Println(msgData)
    msgDataReader := *strings.NewReader(msgData.Encode())

    client := &http.Client{}
    req, _ := http.NewRequest("POST", urlStr, &msgDataReader)
    req.SetBasicAuth(accountSid, authToken)
    req.Header.Add("Accept", "application/json")
    req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
    //client := &http.Client{}


    resp, _ := client.Do(req)
    if (resp.StatusCode >= 200 && resp.StatusCode < 300) {
      var data map[string]interface{}
      decoder := json.NewDecoder(resp.Body)
      err := decoder.Decode(&data)
      if (err == nil) {
        fmt.Println(data["sid"])
      }
    } else {
      fmt.Println(resp.Status);
    }*/
}


type MainNahiManaKarRahaBabuSamajh struct {
  All []interface{}
}

type RecieveDpAuthenticationDataPratibha struct {
  Me string
  Pkey string
  ImageData string
}


func myDpWillBeEdittedPratibhaHandler(w http.ResponseWriter, r *http.Request)  {
  r.ParseForm()
  data:=r.Form.Get("data")
  fmt.Println(data)
  var hameshaHueDekhkarMuskuraoTo RecieveDpAuthenticationDataPratibha
  json.Unmarshal([]byte(data),&hameshaHueDekhkarMuskuraoTo)
  fmt.Println(hameshaHueDekhkarMuskuraoTo.ImageData)

  insaan:=hameshaHueDekhkarMuskuraoTo.Me
  mainInsaan:=hameshaHueDekhkarMuskuraoTo.Pkey
  fmt.Println("aaj se tera")
  fmt.Println(mainInsaan)
  mainInsaanSaasur, _ := strconv.ParseInt(mainInsaan, 10, 64)
  rest:=authorizeThisFace(&insaan,&mainInsaanSaasur)

  if rest!="ERROR" {
    s1:="Match (ee:Rita)-[:HAS_MANDATORY_DP]->(gg:ProfilePic) where ee.email ='"
    b:=rest
    s2:="' set gg.title='"
    bDillagi:=hameshaHueDekhkarMuskuraoTo.ImageData
    agarHamKaheTo:="' return collect(distinct [ee.email,gg.title]) as all"
    s17 := fmt.Sprint(s1,b,s2,bDillagi,agarHamKaheTo)
    fmt.Println(s17)
    driver := bolt.NewDriver()
      conn, err := driver.OpenNeo("bolt://rita:b.PuhuqVThYfCn.fvurl1e25g7fzyCI@hobby-panhpmpgjildgbkepcdcklol.dbs.graphenedb.com:24786?tls=true")
      defer conn.Close()
      dur,_:=time.ParseDuration("600s")
      conn.SetTimeout(dur)
      if err != nil {
        panic(err)
      }
      dataOf, _, _, _ := conn.QueryNeoAll(s17, nil)
      //*aadat<-dataOf[0][0].([]interface{})
      m:=MainNahiManaKarRahaBabuSamajh{
      All:dataOf[0][0].([]interface{})}
      resp, _ := json.Marshal(m)
      w.Header().Set("Content-Type", "application/json")
      w.Header().Set("charset", "utf-8")
      w.Header().Set("Access-Control-Allow-Origin", "*")
      w.Write([]byte(resp))

  }else{
    m:=HttpStatusForbiddenUnAuthorize{
    HttpStatus:403,HttpText:"Authentication Failed"}
    resp, _ := json.Marshal(m)
    w.WriteHeader(http.StatusForbidden)
    w.Header().Set("Content-Type", "application/json")
    w.Header().Set("charset", "utf-8")
    w.Header().Set("Access-Control-Allow-Origin", "*")
    w.Write([]byte(resp))
  }

}


type IwillSubmitMyWillsToYouPratibhaPlease struct {
  Me string
  Pkey string
}



type StateLoginAuthPratibhaNilOrNot struct {
  Status string
}



func deleteMyAuthHashPratibhaPleaseHelpHandler(w http.ResponseWriter, r *http.Request)  {
  r.ParseForm()
  data:=r.Form.Get("data")
  fmt.Println(data)
  var hameshaHueDekhkarMuskuraoTo IwillSubmitMyWillsToYouPratibhaPlease
  json.Unmarshal([]byte(data),&hameshaHueDekhkarMuskuraoTo)
  fmt.Println(hameshaHueDekhkarMuskuraoTo.Me)

  insaan:=hameshaHueDekhkarMuskuraoTo.Me
  mainInsaan:=hameshaHueDekhkarMuskuraoTo.Pkey
  fmt.Println("aaj se tera")
  fmt.Println(mainInsaan)
  mainInsaanSaasur, _ := strconv.ParseInt(mainInsaan, 10, 64)
  rest:=authorizeThisFace(&insaan,&mainInsaanSaasur)

  if rest!="ERROR" {
    s1:="Match (ee:Auth) where ee.hash ='"
    b:=hameshaHueDekhkarMuskuraoTo.Me
    agarHamKaheTo:="' delete ee"
    s17 := fmt.Sprint(s1,b,agarHamKaheTo)
    fmt.Println(s17)
    driver := bolt.NewDriver()
      conn, err := driver.OpenNeo("bolt://rita:b.PuhuqVThYfCn.fvurl1e25g7fzyCI@hobby-panhpmpgjildgbkepcdcklol.dbs.graphenedb.com:24786?tls=true")
      defer conn.Close()
      dur,_:=time.ParseDuration("600s")
      conn.SetTimeout(dur)
      if err != nil {
        panic(err)
      }
      var stateMaintain string
      _, _, _, err1 := conn.QueryNeoAll(s17, nil)
      //*aadat<-dataOf[0][0].([]interface{})
      if err1!=nil {
        stateMaintain="not nil"
      }else{
        stateMaintain="nil"
      }
      m:=StateLoginAuthPratibhaNilOrNot{
      Status:stateMaintain}
      resp, _ := json.Marshal(m)
      w.Header().Set("Content-Type", "application/json")
      w.Header().Set("charset", "utf-8")
      w.Header().Set("Access-Control-Allow-Origin", "*")
      w.Write([]byte(resp))
  }else{
    m:=HttpStatusForbiddenUnAuthorize{
    HttpStatus:403,HttpText:"Authentication Failed"}
    resp, _ := json.Marshal(m)
    w.WriteHeader(http.StatusForbidden)
    w.Header().Set("Content-Type", "application/json")
    w.Header().Set("charset", "utf-8")
    w.Header().Set("Access-Control-Allow-Origin", "*")
    w.Write([]byte(resp))
  }

}





func PayPalProcessAmountPratibhaPlease(cardNo *string, cardType *string, expMnth *string, expYear *string, cvvTemp *string, fName *string, lName *string, money *float64) (int,error) {
  client, err := paypalsdk.NewClient("AR6ERqJofXz37bYK4pbuUmSsH5J8B6is1LMCubv8KKaTn_wpFk7TPVkIU7R-01wcMFJXTYUhtzNXosqo", "ENiKdPNAW-DxK5cAot4mxTnFM4X0WXa_cg3gXhXC1d1EPkMBLKWypqOLohFxg5Hoo1YbusnaTOY4EHdQ", paypalsdk.APIBaseSandBox)
client.SetLog(os.Stdout) // Set log to terminal stdout

accessToken, err := client.GetAccessToken()
fmt.Println(accessToken)
if err!=nil {
  return 0,err
}

/*
*/
//fmt.Println(accessToken)
  /*var tempExpMnth string
  if sIn, errx := strconv.Atoi(*expMnth); errx == nil {
		if sIn<10 {

		}
	}else{
    return 0,errx
  }*/

  //strMoney:=fmt.Sprint(*money)
  fmt.Println(*money)
  strMoney:=strconv.FormatFloat(*money, 'f', -1, 64)
  fmt.Println("hehehhehehehhehehhehggagjhagdgjagjdgahgdhagwjdgajdwgawhjghwadjgdhagwjdgawgjd")
  fmt.Println(strMoney)
  p := paypalsdk.Payment{
  Intent: "sale",
  Payer: &paypalsdk.Payer{
      PaymentMethod: "credit_card",
      FundingInstruments: []paypalsdk.FundingInstrument{paypalsdk.FundingInstrument{
          CreditCard: &paypalsdk.CreditCard{
            Number:      *cardNo,
            Type:        *cardType,
            ExpireMonth: *expMnth,
            ExpireYear:  *expYear,
            CVV2:        *cvvTemp,
            FirstName:   *fName,
            LastName:    *lName,
          },
      }},
  },
  Transactions: []paypalsdk.Transaction{paypalsdk.Transaction{
      Amount: &paypalsdk.Amount{
          Currency: "USD",
          Total:    strMoney,
      },
      Description: "thelinkedface payment",
  }},
  RedirectURLs: &paypalsdk.RedirectURLs{
      ReturnURL: "http://localhost:5000/paypal_success",
      CancelURL: "http://localhost:5000/paypal_cancel",
  },
}
  //paymentResponse, err := client.CreatePayment(p)
  _, err1 := client.CreatePayment(p)
  if err!=nil {
    return 0,err1
  }
  //fmt.Println(paymentResponse)
  return 200,nil
}


type AttemptPaypalPaymentPratibhaPlease struct {
  Me string
  Pkey string
  CardNo string
  CVV string
  ValidMonth string
  ValidYear string
  Money string
  FirstName string
  LastName string
  CardType string
  ItemId int
  ItemQuantity string
}

func UpdateSilentlyToTheDbPls(itId *int,itQuant *string) int  {
  s1:="Match (ee:ItemContainer) where id(ee) ="
  b:=*itId
  agarHamKaheTo:=" set ee.ItemQuantity='"
  jutiAb:=*itQuant
  jutiAbc:="' return collect([ee.ItemQuantity])"
  s17 := fmt.Sprint(s1,b,agarHamKaheTo,jutiAb,jutiAbc)
  fmt.Println(s17)
  driver := bolt.NewDriver()
    conn, err := driver.OpenNeo("bolt://rita:b.PuhuqVThYfCn.fvurl1e25g7fzyCI@hobby-panhpmpgjildgbkepcdcklol.dbs.graphenedb.com:24786?tls=true")
    defer conn.Close()
    dur,_:=time.ParseDuration("600s")
    conn.SetTimeout(dur)
    if err != nil {
      //panic(err)
    }
    //var stateMaintain string
    _, _, _, errm := conn.QueryNeoAll(s17, nil)
    if errm!=nil {
      return 0
    }else{
      return 200
    }
}

func startEditTheItemInTheDbPlease(itId *int,itQuantity *string) int  {
  s1:="Match (ee:ItemContainer) where id(ee) ="
  b:=*itId
  agarHamKaheTo:=" return collect([ee.ItemQuantity])"
  s17 := fmt.Sprint(s1,b,agarHamKaheTo)
  fmt.Println(s17)
  driver := bolt.NewDriver()
    conn, err := driver.OpenNeo("bolt://rita:b.PuhuqVThYfCn.fvurl1e25g7fzyCI@hobby-panhpmpgjildgbkepcdcklol.dbs.graphenedb.com:24786?tls=true")
    defer conn.Close()
    dur,_:=time.ParseDuration("600s")
    conn.SetTimeout(dur)
    if err != nil {
      //panic(err)
    }
    //var stateMaintain string
    dataOf, _, _, _ := conn.QueryNeoAll(s17, nil)
    ss:=dataOf[0][0].([]interface{})
    secArr:=make([]string, 0)
    for _,ght:=range ss{
      for _,dnt:=range ght.([]interface{}){
        secArr=append(secArr,dnt.(string))
      }
    }
    itQOr:=secArr[0]
    //splOr:=strings.Split(itQOr," ")
    splTemp:=strings.Split(*itQuantity," ")
    splOr:=strings.Replace(itQOr, splTemp[1], "", -1)

    itorint,_:=strconv.Atoi(splOr)
    fmt.Println("fgffgfgfghghfghfhgfhgfxcvxcvx    jhjkhjkhkhkkhjkhkhfggdfgdfgd")
    fmt.Println(itorint)
    itTmpint,_:=strconv.Atoi(splTemp[0])
    var QuantOr int
    var abStrTamp string
    //sps:=" "
    if itTmpint<=itorint {
      QuantOr=itorint-itTmpint
      stTotQuant:=strconv.Itoa(QuantOr)
      tampSec:=splTemp[1]
      abStrTamp=fmt.Sprintf("%s%s",stTotQuant,tampSec)
      ittii:=*itId
      fmt.Println("fgffgfgfghghfghfhgfhgfxcvxcvx    jhjkhjkhkhkkhjkhkhfggdfgdfgd")
      fmt.Println(abStrTamp)
      retable:=UpdateSilentlyToTheDbPls(&ittii,&abStrTamp)

      return retable
    }else{
      return 0
    }
}

func addBuyedRelationshipOrderWithOrderId(email *string,itemId *int,itemQtt *string) int {
  now:=time.Now()
  nanos:=now.UnixNano()
  milis:=nanos / 1000000
  s1:="Match (ee:Rita) where ee.email='"
  hhHee:=*email
  heyHey:="' WITH ee optional match (ff:ItemContainer) where id(ff) ="
  b:=*itemId
  agarHamKaheTo:=" MERGE (ee)-[:BOUGHT_THIS_ITEM{itemQuantity:'"
  hehehoho:=*itemQtt
  hohoBuddhha:="',when:'"
  wBuss:=milis
  wBBuuddhh:="'}]-(ff)  return collect([ee.ItemQuantity])"
  s17 := fmt.Sprint(s1,hhHee,heyHey,b,agarHamKaheTo,hehehoho,hohoBuddhha,wBuss,wBBuuddhh)
  fmt.Println(s17)
  driver := bolt.NewDriver()
    conn, err := driver.OpenNeo("bolt://rita:b.PuhuqVThYfCn.fvurl1e25g7fzyCI@hobby-panhpmpgjildgbkepcdcklol.dbs.graphenedb.com:24786?tls=true")
    defer conn.Close()
    dur,_:=time.ParseDuration("600s")
    conn.SetTimeout(dur)
    if err != nil {
      return 0
    }else{
      _, _, _, emmm := conn.QueryNeoAll(s17, nil)
      if emmm!=nil {
        return 0
      }else{
        return 200
      }
    }
}


func recieveMyPaymentPratibhaPleaseYouAreOnlyHopeOfMineHandler(w http.ResponseWriter, r *http.Request)  {
  r.ParseForm()
  data:=r.Form.Get("data")
  fmt.Println(data)
  var hameshaHueDekhkarMuskuraoTo AttemptPaypalPaymentPratibhaPlease
  json.Unmarshal([]byte(data),&hameshaHueDekhkarMuskuraoTo)
  fmt.Println("gsfdghdfgfsdgkyuWIDYAIUDYUIYAIUDYIYIDYGHFHGFghfhgfhfghdfahgf")
  fmt.Println(hameshaHueDekhkarMuskuraoTo.Money)

  insaan:=hameshaHueDekhkarMuskuraoTo.Me
  mainInsaan:=hameshaHueDekhkarMuskuraoTo.Pkey
  fmt.Println("aaj se tera")
  fmt.Println(mainInsaan)
  mainInsaanSaasur, _ := strconv.ParseInt(mainInsaan, 10, 64)
  rest:=authorizeThisFace(&insaan,&mainInsaanSaasur)

  if rest!="ERROR" {
    cN:=hameshaHueDekhkarMuskuraoTo.CardNo
    cHidd:=hameshaHueDekhkarMuskuraoTo.CVV
    vfm:=hameshaHueDekhkarMuskuraoTo.ValidMonth
    vfy:=hameshaHueDekhkarMuskuraoTo.ValidYear
    mny:=hameshaHueDekhkarMuskuraoTo.Money
    fnm:=hameshaHueDekhkarMuskuraoTo.FirstName
    lnm:=hameshaHueDekhkarMuskuraoTo.LastName
    cty:=hameshaHueDekhkarMuskuraoTo.CardType
    itiid:=hameshaHueDekhkarMuskuraoTo.ItemId
    itQtt:=hameshaHueDekhkarMuskuraoTo.ItemQuantity
    mnyLol, _:=strconv.ParseFloat(mny, 64)

    _, errHappen:=PayPalProcessAmountPratibhaPlease(&cN,&cty,&vfm,&vfy,&cHidd,&fnm,&lnm,&mnyLol)

    if errHappen!=nil {
      m:=HttpStatusForbiddenUnAuthorize{
      HttpStatus:403,HttpText:"payment failed"}
      resp, _ := json.Marshal(m)
      w.WriteHeader(http.StatusForbidden)
      w.Header().Set("Content-Type", "application/json")
      w.Header().Set("charset", "utf-8")
      w.Header().Set("Access-Control-Allow-Origin", "*")
      w.Write([]byte(resp))
    }else{

      retMe:=startEditTheItemInTheDbPlease(&itiid,&itQtt)

      if retMe!=0 {
        emMiner:=rest
        itiidTamp:=hameshaHueDekhkarMuskuraoTo.ItemId
        itQttTamp:=hameshaHueDekhkarMuskuraoTo.ItemQuantity
        retaOfferDb:=addBuyedRelationshipOrderWithOrderId(&emMiner,&itiidTamp,&itQttTamp)

        if retaOfferDb!=0 {
          m:=HttpStatusForbiddenUnAuthorize{
          HttpStatus:200,HttpText:"booking done"}
          resp, _ := json.Marshal(m)
          w.Header().Set("Content-Type", "application/json")
          w.Header().Set("charset", "utf-8")
          w.Header().Set("Access-Control-Allow-Origin", "*")
          w.Write([]byte(resp))
        }else{
          m:=HttpStatusForbiddenUnAuthorize{
          HttpStatus:417,HttpText:"booking failed"}
          resp, _ := json.Marshal(m)
          w.WriteHeader(http.StatusExpectationFailed)
          w.Header().Set("Content-Type", "application/json")
          w.Header().Set("charset", "utf-8")
          w.Header().Set("Access-Control-Allow-Origin", "*")
          w.Write([]byte(resp))
        }





      }else{
        m:=HttpStatusForbiddenUnAuthorize{
        HttpStatus:417,HttpText:"booking failed"}
        resp, _ := json.Marshal(m)
        w.WriteHeader(http.StatusExpectationFailed)
        w.Header().Set("Content-Type", "application/json")
        w.Header().Set("charset", "utf-8")
        w.Header().Set("Access-Control-Allow-Origin", "*")
        w.Write([]byte(resp))
      }


    }

  }else{
    m:=HttpStatusForbiddenUnAuthorize{
    HttpStatus:403,HttpText:"Authentication Failed"}
    resp, _ := json.Marshal(m)
    w.WriteHeader(http.StatusForbidden)
    w.Header().Set("Content-Type", "application/json")
    w.Header().Set("charset", "utf-8")
    w.Header().Set("Access-Control-Allow-Origin", "*")
    w.Write([]byte(resp))
  }


}

type RecieveLoyalityFormCheck struct {
  Me string
  Pkey string
  Target [][]string
  Mentor [][]string
  EmotionalMessage string
}


func updateDbWithLoyalCheckFormPratibhaPlease(fromId *string, toId *string, emotionalMessage *string, targetTemp *string) int {
  now:=time.Now()
  nanos:=now.UnixNano()
  milis:=nanos / 1000000
  s1:="Match (ee:Rita) where ee.email='"
  hhHee:=*fromId
  heyHey:="' WITH ee optional match (ff:Rita) where ff.email ='"
  b:=*toId
  agarHamKaheTo:="' with ee,ff MERGE (ee)-[:WANT_LOYALITY_CHECK{target:'"
  hehehoho:=*targetTemp
  hohoBuddhha:="',when:'"
  wBuss:=milis
  wBBuuddhh:="',emotionalMessage:'"
  tellMeMyEm:=*emotionalMessage
  sdSwag:="'}]-(ff)  return collect([ff.email]) as all"
  s17 := fmt.Sprint(s1,hhHee,heyHey,b,agarHamKaheTo,hehehoho,hohoBuddhha,wBuss,wBBuuddhh,tellMeMyEm,sdSwag)
  fmt.Println(s17)
  driver := bolt.NewDriver()
    conn, err := driver.OpenNeo("bolt://rita:b.PuhuqVThYfCn.fvurl1e25g7fzyCI@hobby-panhpmpgjildgbkepcdcklol.dbs.graphenedb.com:24786?tls=true")
    defer conn.Close()
    dur,_:=time.ParseDuration("600s")
    conn.SetTimeout(dur)
    if err != nil {
      return 0
    }else{
      _, _, _, emmm := conn.QueryNeoAll(s17, nil)
      if emmm!=nil {
        return 0
      }else{
        return 200
      }
    }
}


type RespondableMentor struct {
  All []string
}


func submitMyLoyalityFormPratibhaPlsInDatingZoneHandler(w http.ResponseWriter, r *http.Request)  {
  r.ParseForm()
  data:=r.Form.Get("data")
  fmt.Println(data)
  var hameshaHueDekhkarMuskuraoTo RecieveLoyalityFormCheck
  json.Unmarshal([]byte(data),&hameshaHueDekhkarMuskuraoTo)
  fmt.Println("gsfdghdfgfsdgkyuWIDYAIUDYUIYAIUDYIYIDYGHFHGFghfhgfhfghdfahgf")
  fmt.Println(hameshaHueDekhkarMuskuraoTo.Me)

  insaan:=hameshaHueDekhkarMuskuraoTo.Me
  mainInsaan:=hameshaHueDekhkarMuskuraoTo.Pkey
  fmt.Println("aaj se tera")
  fmt.Println(mainInsaan)
  mainInsaanSaasur, _ := strconv.ParseInt(mainInsaan, 10, 64)
  rest:=authorizeThisFace(&insaan,&mainInsaanSaasur)
  succArrSl:=make([]string, 0)
  if rest!="ERROR" {
    for _, rrm := range hameshaHueDekhkarMuskuraoTo.Mentor {
        sBabar:=rrm[1];
        sAkbar:=hameshaHueDekhkarMuskuraoTo.Target[0][0]
        sEmotionalMessage:=hameshaHueDekhkarMuskuraoTo.EmotionalMessage
        mmRefer:=rest
        fl:=updateDbWithLoyalCheckFormPratibhaPlease(&mmRefer,&sBabar,&sEmotionalMessage,&sAkbar)
        if fl!=0 {
          succArrSl=append(succArrSl,sBabar)
        }
    }
    m:=RespondableMentor{
    All:succArrSl}
    resp, _ := json.Marshal(m)
    //w.WriteHeader(http.StatusForbidden)
    w.Header().Set("Content-Type", "application/json")
    w.Header().Set("charset", "utf-8")
    w.Header().Set("Access-Control-Allow-Origin", "*")
    w.Write([]byte(resp))
  }else {
    m:=HttpStatusForbiddenUnAuthorize{
    HttpStatus:403,HttpText:"Authentication Failed"}
    resp, _ := json.Marshal(m)
    w.WriteHeader(http.StatusForbidden)
    w.Header().Set("Content-Type", "application/json")
    w.Header().Set("charset", "utf-8")
    w.Header().Set("Access-Control-Allow-Origin", "*")
    w.Write([]byte(resp))
  }
}


type RecieveDonationRequestSeek struct {
  Me string
  Pkey string
  SeekAmount string
  EmotionalMessage string
}


type ReturnableSeekDonation struct {
  All []interface{}
}


func seekDonationForMeAjayPlsTlfDonateHandler(w http.ResponseWriter, r *http.Request)  {
  r.ParseForm()
  data:=r.Form.Get("data")
  fmt.Println(data)
  var hameshaHueDekhkarMuskuraoTo RecieveDonationRequestSeek
  json.Unmarshal([]byte(data),&hameshaHueDekhkarMuskuraoTo)
  fmt.Println("gsfdghdfgfsdgkyuWIDYAIUDYUIYAIUDYIYIDYGHFHGFghfhgfhfghdfahgf")
  fmt.Println(hameshaHueDekhkarMuskuraoTo.Me)

  insaan:=hameshaHueDekhkarMuskuraoTo.Me
  mainInsaan:=hameshaHueDekhkarMuskuraoTo.Pkey
  fmt.Println("aaj se tera")
  fmt.Println(mainInsaan)
  mainInsaanSaasur, _ := strconv.ParseInt(mainInsaan, 10, 64)
  rest:=authorizeThisFace(&insaan,&mainInsaanSaasur)

  if rest!="ERROR" {
    now:=time.Now()
    nanos:=now.UnixNano()
    milis:=nanos / 1000000
    s1:="Match (ee:Rita) where ee.email='"
    hhHee:=rest
    heyHey:="' WITH ee MERGE (ee)-[:SEEKED_DONATION]-(ff:DonationData{amount:'"
    teraPyarMera:=hameshaHueDekhkarMuskuraoTo.SeekAmount
    zzZinda:="',emotionalMessage:'"
    honiJud:=hameshaHueDekhkarMuskuraoTo.EmotionalMessage
    saathMera:="',when:'"
    dholna:=milis
    naiLagta:="'}) return collect([ee.email]) as all"

    s17 := fmt.Sprint(s1,hhHee,heyHey,teraPyarMera,zzZinda,honiJud,saathMera,dholna,naiLagta)
    fmt.Println(s17)
    driver := bolt.NewDriver()
      conn, err := driver.OpenNeo("bolt://rita:b.PuhuqVThYfCn.fvurl1e25g7fzyCI@hobby-panhpmpgjildgbkepcdcklol.dbs.graphenedb.com:24786?tls=true")
      defer conn.Close()
      dur,_:=time.ParseDuration("600s")
      conn.SetTimeout(dur)
      if err != nil {
        m:=HttpStatusForbiddenUnAuthorize{
        HttpStatus:403,HttpText:"Authentication Failed"}
        resp, _ := json.Marshal(m)
        w.WriteHeader(http.StatusForbidden)
        w.Header().Set("Content-Type", "application/json")
        w.Header().Set("charset", "utf-8")
        w.Header().Set("Access-Control-Allow-Origin", "*")
        w.Write([]byte(resp))
      }else{
        dataOf, _, _, emmm := conn.QueryNeoAll(s17, nil)
        if emmm!=nil {
          m:=HttpStatusForbiddenUnAuthorize{
          HttpStatus:403,HttpText:"Authentication Failed"}
          resp, _ := json.Marshal(m)
          w.WriteHeader(http.StatusForbidden)
          w.Header().Set("Content-Type", "application/json")
          w.Header().Set("charset", "utf-8")
          w.Header().Set("Access-Control-Allow-Origin", "*")
          w.Write([]byte(resp))
        }else{
          m:=ReturnableSeekDonation{
          All:dataOf[0][0].([]interface{})}
          resp, _ := json.Marshal(m)
          w.Header().Set("Content-Type", "application/json")
          w.Header().Set("charset", "utf-8")
          w.Header().Set("Access-Control-Allow-Origin", "*")
          w.Write([]byte(resp))
        }
      }
  }else {
    m:=HttpStatusForbiddenUnAuthorize{
    HttpStatus:403,HttpText:"Authentication Failed"}
    resp, _ := json.Marshal(m)
    w.WriteHeader(http.StatusForbidden)
    w.Header().Set("Content-Type", "application/json")
    w.Header().Set("charset", "utf-8")
    w.Header().Set("Access-Control-Allow-Origin", "*")
    w.Write([]byte(resp))
  }
}


type DirectConfessionPratibhaSinZone struct {
  Me string
  Pkey string
  ConfessTo string
  EmotionalMessage string
}


type ReturnableDirectConfessionData struct {
  All []interface{}
}


func createMyDirectConfessionPratibhaPleaseRecordMySinHandler(w http.ResponseWriter, r *http.Request)  {
  r.ParseForm()
  data:=r.Form.Get("data")
  fmt.Println(data)
  var hameshaHueDekhkarMuskuraoTo DirectConfessionPratibhaSinZone
  json.Unmarshal([]byte(data),&hameshaHueDekhkarMuskuraoTo)
  fmt.Println("gsfdghdfgfsdgkyuWIDYAIUDYUIYAIUDYIYIDYGHFHGFghfhgfhfghdfahgf")
  fmt.Println(hameshaHueDekhkarMuskuraoTo.Me)

  insaan:=hameshaHueDekhkarMuskuraoTo.Me
  mainInsaan:=hameshaHueDekhkarMuskuraoTo.Pkey
  fmt.Println("aaj se tera")
  fmt.Println(mainInsaan)
  mainInsaanSaasur, _ := strconv.ParseInt(mainInsaan, 10, 64)
  rest:=authorizeThisFace(&insaan,&mainInsaanSaasur)

  if rest!="ERROR" {
    now:=time.Now()
    nanos:=now.UnixNano()
    milis:=nanos / 1000000
    s1:="Match (ee:Rita) where ee.email='"
    hhHee:=rest
    heyHey:="' WITH ee optional Match (ff:Rita) where ff.email='"
    hhoHo:=hameshaHueDekhkarMuskuraoTo.ConfessTo
    sDuniya:="' with ee,ff MERGE (ee)-[:DIRECT_CONFESSION{emotionalMessage:'"
    sDarkar:=hameshaHueDekhkarMuskuraoTo.EmotionalMessage
    bDarkar:="',when:'"
    lDarkar:=milis
    mDarkar:="'}]-(ff)"
    naiLagta:="' return collect([ee.email]) as all"

    s17 := fmt.Sprint(s1,hhHee,heyHey,hhoHo,sDuniya,sDarkar,bDarkar,lDarkar,mDarkar,naiLagta)
    fmt.Println(s17)
    driver := bolt.NewDriver()
      conn, err := driver.OpenNeo("bolt://rita:b.PuhuqVThYfCn.fvurl1e25g7fzyCI@hobby-panhpmpgjildgbkepcdcklol.dbs.graphenedb.com:24786?tls=true")
      defer conn.Close()
      dur,_:=time.ParseDuration("600s")
      conn.SetTimeout(dur)
      if err != nil {
        m:=HttpStatusForbiddenUnAuthorize{
        HttpStatus:403,HttpText:"Authentication Failed"}
        resp, _ := json.Marshal(m)
        w.WriteHeader(http.StatusForbidden)
        w.Header().Set("Content-Type", "application/json")
        w.Header().Set("charset", "utf-8")
        w.Header().Set("Access-Control-Allow-Origin", "*")
        w.Write([]byte(resp))
      }else{
        dataOf, _, _, emmm := conn.QueryNeoAll(s17, nil)
        if emmm!=nil {
          m:=HttpStatusForbiddenUnAuthorize{
          HttpStatus:403,HttpText:"Authentication Failed"}
          resp, _ := json.Marshal(m)
          w.WriteHeader(http.StatusForbidden)
          w.Header().Set("Content-Type", "application/json")
          w.Header().Set("charset", "utf-8")
          w.Header().Set("Access-Control-Allow-Origin", "*")
          w.Write([]byte(resp))
        }else{
          m:=ReturnableDirectConfessionData{
          All:dataOf[0][0].([]interface{})}
          resp, _ := json.Marshal(m)
          w.Header().Set("Content-Type", "application/json")
          w.Header().Set("charset", "utf-8")
          w.Header().Set("Access-Control-Allow-Origin", "*")
          w.Write([]byte(resp))
        }
      }
  }else {
    m:=HttpStatusForbiddenUnAuthorize{
    HttpStatus:403,HttpText:"Authentication Failed"}
    resp, _ := json.Marshal(m)
    w.WriteHeader(http.StatusForbidden)
    w.Header().Set("Content-Type", "application/json")
    w.Header().Set("charset", "utf-8")
    w.Header().Set("Access-Control-Allow-Origin", "*")
    w.Write([]byte(resp))
  }
}

type ReturnWorkflowData struct {
  TaskNumber string
  TaskType string
  Requestor string
  TaskDescription  string
  DateOfSubmission string
  TaskStatus string
}


type WorkflowContainer struct {
  Workflow []ReturnWorkflowData
}


func serveBotTemplate(w http.ResponseWriter, r *http.Request)  {
  vars := mux.Vars(r)
  mkiloj:=vars["interactionId"]
  fmt.Println(mkiloj)
  arr:=make([]ReturnWorkflowData, 0)
  /*for i := 1; i < 10; i++ {
      mc,_:=strconv.Atoi(mkiloj)
      l:=ReturnWorkflowData{
      WorkFlowId:i*10 * mc,WorkFlowType:"anonymus",WorkFlowText:fmt.Sprintf("%s%d","template_",i),IsApproved:false}
      arr=append(arr,l)
  }*/
  l:=ReturnWorkflowData{
    TaskNumber:"WS12345678",TaskType:"Leave",Requestor:"rupsa mukherjee",TaskDescription:"Applied Leave for 06-16-18 to 06-22-18",DateOfSubmission:"06-16-18",TaskStatus:"Pending"}
  arr=append(arr,l)
  ll:=ReturnWorkflowData{
    TaskNumber:"WS12345633",TaskType:"Expense",Requestor:"Ajay Jha",TaskDescription:"Claim Request of 500.15$",DateOfSubmission:"06-16-18",TaskStatus:"Pending"}
  arr=append(arr,ll)
  m:=WorkflowContainer{
  Workflow:arr}
  resp, _ := json.Marshal(m)
  w.Header().Set("Content-Type", "application/json")
  w.Header().Set("charset", "utf-8")
  w.Header().Set("Access-Control-Allow-Origin", "*")
  w.Write([]byte(resp))
}

func serveWorkflowTemplate(w http.ResponseWriter, r *http.Request){
  vars := mux.Vars(r)
  mkiloj:=vars["workflowId"]
  //mkiloj_2:=vars["interactionId"]
  fmt.Println(mkiloj)

  //workfloId,_:=strconv.Atoi(mkiloj)
  //interactionId,_:=strconv.Atoi(mkiloj_2)


  // change 


  arr:=make([]ReturnWorkflowData, 0)
    /*for i := 1; i < 10; i++ {
      if i*10*interactionId==workfloId {
          //mc,_:=strconv.Atoi(mkiloj)
          l:=ReturnWorkflowData{
          WorkFlowId:i*10 * interactionId,WorkFlowType:"anonymus",WorkFlowText:fmt.Sprintf("%s%d","template_",i),IsApproved:true}
          arr=append(arr,l)
        }else{
          //mc,_:=strconv.Atoi(mkiloj)
          l:=ReturnWorkflowData{
          WorkFlowId:i*10 * interactionId,WorkFlowType:"anonymus",WorkFlowText:fmt.Sprintf("%s%d","template_",i),IsApproved:false}
          arr=append(arr,l)
        }
    }*/
  l:=ReturnWorkflowData{
    TaskNumber:"WS12345678",TaskType:"Leave",Requestor:"rupsa mukherjee",TaskDescription:"Applied Leave for 06-16-18 to 06-22-18",DateOfSubmission:"06-16-18",TaskStatus:"Pending"}
  arr=append(arr,l)
  ll:=ReturnWorkflowData{
    TaskNumber:"WS12345633",TaskType:"Expense",Requestor:"Ajay Jha",TaskDescription:"Claim Request of 500.15$",DateOfSubmission:"06-16-18",TaskStatus:"Pending"}
  arr=append(arr,ll)

  for i := 0; i < len(arr); i++ {
    if arr[i].TaskNumber==mkiloj {
      xx:=arr[i]
      arr=append(arr[:i], arr[i+1:]...)
      kkmm:=ReturnWorkflowData{
      TaskNumber:xx.TaskNumber,TaskType:xx.TaskType,Requestor:xx.Requestor,TaskDescription:xx.TaskDescription,DateOfSubmission:xx.DateOfSubmission,TaskStatus:"Approved"}
      arr=append(arr,kkmm)
    }
  }

  m:=WorkflowContainer{
  Workflow:arr}
  resp, _ := json.Marshal(m)
  w.Header().Set("Content-Type", "application/json")
  w.Header().Set("charset", "utf-8")
  w.Header().Set("Access-Control-Allow-Origin", "*")
  w.Write([]byte(resp))

}


type KitchatHash struct {
  Hash string `json:hash`
}


type SignHashWith struct{
  ID bson.ObjectId `bson:"_id,omitempty"`
  Name []string
  Mobile string
  Password string
  Latitude float64
  Longitude float64
  Gender string
  Me []string
  Partner []string
  Dp string
}

type SignSuccessRet struct {
  Status bool
}

type HelloLogin struct {
  Mobile string
  Password string
}


type DatingRequestFromApp struct {
  Loc bool
  Timepass bool
  Dating bool
  To string
  From string
  Message string
  Approved bool
}


type ReturnStatusRequestApi struct {
  Status bool
}


func saveDatingRequestTemplate(w http.ResponseWriter, r *http.Request) {
  decoder := json.NewDecoder(r.Body)

  var t DatingRequestFromApp
  var signal bool
  err := decoder.Decode(&t)

  if err != nil {
    //panic(err)
    signal=false
  }else{
    signal=true
  }


dialInfo := &mgo.DialInfo{
    Addrs:    []string{"rita-shard-00-00-qk9t0.mongodb.net:27017","rita-shard-00-01-qk9t0.mongodb.net:27017","rita-shard-00-02-qk9t0.mongodb.net:27017"},
    Database: "test",
    Username: "badcodercpp",
    Password: "Ajayjha93",
    ReplicaSetName: "Rita-shard-0",
    Source: "admin",
    DialServer: func(addr *mgo.ServerAddr) (net.Conn, error) {
        return tls.Dial("tcp", addr.String(), &tls.Config{})
    },
    Timeout: time.Second * 10000,
  }
  session, err := mgo.DialWithInfo(dialInfo)
  defer session.Close()
  //signal:=true
  if(err != nil){
    panic(err)
  }
  c := session.DB("test").C("Request")
    // Query One  4
    //result := SignHashWith{}
    //ml:= make(map[string]string)
    //ml["hash"]=*hash
    _,err = c.Upsert(bson.M{"to":t.To,"from":t.From},bson.M{"$set":bson.M{"loc":t.Loc,"timepass":t.Timepass,"dating":t.Dating,"to":t.To,"from":t.From,"message":t.Message,"approved":t.Approved}})
    if err != nil {
      //panic(err)
      //signal=false
      signal=false
    }else{
      signal=true
    }
    fmt.Println(signal)
    m:=ReturnStatusRequestApi{
        Status:signal}
        resp, _ := json.Marshal(m)

    w.Header().Set("Content-Type", "application/json")
  w.Header().Set("charset", "utf-8")
  w.Header().Set("Access-Control-Allow-Origin", "*")
  w.Write([]byte(resp))

}

type IncomingNotTypeDateReq struct {
  Loc bool
  Timepass bool
  Dating bool
}
type Information struct {
  To string
  From string
}
type MyAllAppNotification struct {
  Loc bool
  Timepass bool
  Dating bool
  To string
  From string
  Message string
  Approved bool
}







func alterDatingRequestTemplate(w http.ResponseWriter, r *http.Request) {
  decoder := json.NewDecoder(r.Body)

  var t DatingRequestFromApp
  var signal bool
  err := decoder.Decode(&t)

  if err != nil {
    //panic(err)
    signal=false
  }else{
    signal=true
  }


dialInfo := &mgo.DialInfo{
    Addrs:    []string{"rita-shard-00-00-qk9t0.mongodb.net:27017","rita-shard-00-01-qk9t0.mongodb.net:27017","rita-shard-00-02-qk9t0.mongodb.net:27017"},
    Database: "test",
    Username: "badcodercpp",
    Password: "Ajayjha93",
    ReplicaSetName: "Rita-shard-0",
    Source: "admin",
    DialServer: func(addr *mgo.ServerAddr) (net.Conn, error) {
        return tls.Dial("tcp", addr.String(), &tls.Config{})
    },
    Timeout: time.Second * 10000,
  }
  session, err := mgo.DialWithInfo(dialInfo)
  defer session.Close()
  //signal:=true
  if(err != nil){
    panic(err)
  }
  c := session.DB("test").C("Request")
    // Query One  4
    //result := SignHashWith{}
    //ml:= make(map[string]string)
    //ml["hash"]=*hash
    _,err = c.Upsert(bson.M{"to":t.To,"from":t.From},bson.M{"$set":bson.M{"loc":t.Loc,"timepass":t.Timepass,"dating":t.Dating,"to":t.To,"from":t.From,"message":t.Message,"approved":t.Approved}})
    if err != nil {
      //panic(err)
      //signal=false
      signal=false
    }else{
      signal=true
    }
    fmt.Println(signal)
    m:=ReturnStatusRequestApi{
        Status:signal}
        resp, _ := json.Marshal(m)

    w.Header().Set("Content-Type", "application/json")
  w.Header().Set("charset", "utf-8")
  w.Header().Set("Access-Control-Allow-Origin", "*")
  w.Write([]byte(resp))

}









func getMyAllNotification(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
  mkiloj:=vars["target"]
  dialInfo := &mgo.DialInfo{
    Addrs:    []string{"rita-shard-00-00-qk9t0.mongodb.net:27017","rita-shard-00-01-qk9t0.mongodb.net:27017","rita-shard-00-02-qk9t0.mongodb.net:27017"},
    Database: "test",
    Username: "badcodercpp",
    Password: "Ajayjha93",
    ReplicaSetName: "Rita-shard-0",
    Source: "admin",
    DialServer: func(addr *mgo.ServerAddr) (net.Conn, error) {
        return tls.Dial("tcp", addr.String(), &tls.Config{})
    },
    Timeout: time.Second * 10000,
  }
  session, err := mgo.DialWithInfo(dialInfo)
  defer session.Close()
  //signal:=true
  if(err != nil){
    panic(err)
  }
  c := session.DB("test").C("Request")
    // Query One  4
    result := make([]MyAllAppNotification,0)
    //ml:= make(map[string]string)
    //ml["Hash"]=mkiloj

    err = c.Find(bson.M{"to":bson.M{"$eq":mkiloj}}).All(&result)
    if err != nil {
      panic(err)
      //signal=false
      }

      //result.Type="Date Request"

      //hs:=result.Hash
  
  resp, _ := json.Marshal(result)
  w.Header().Set("Content-Type", "application/json")
  w.Header().Set("charset", "utf-8")
  w.Header().Set("Access-Control-Allow-Origin", "*")
  w.Write([]byte(resp))

}






func letTheUserLogin(w http.ResponseWriter, r *http.Request) {
  decoder := json.NewDecoder(r.Body)

  var t HelloLogin
  var signal bool
  err := decoder.Decode(&t)

  if err != nil {
    //panic(err)
    signal=false
  }else{
    signal=true
  }


dialInfo := &mgo.DialInfo{
    Addrs:    []string{"rita-shard-00-00-qk9t0.mongodb.net:27017","rita-shard-00-01-qk9t0.mongodb.net:27017","rita-shard-00-02-qk9t0.mongodb.net:27017"},
    Database: "test",
    Username: "badcodercpp",
    Password: "Ajayjha93",
    ReplicaSetName: "Rita-shard-0",
    Source: "admin",
    DialServer: func(addr *mgo.ServerAddr) (net.Conn, error) {
        return tls.Dial("tcp", addr.String(), &tls.Config{})
    },
    Timeout: time.Second * 10000,
  }
  session, err := mgo.DialWithInfo(dialInfo)
  defer session.Close()
  //signal:=true
  if(err != nil){
    panic(err)
  }
  c := session.DB("test").C("userdata")
    // Query One  4
    result := SignHashWith{}
    //ml:= make(map[string]string)
    //ml["hash"]=*hash
    err = c.Find(bson.M{"mobile":t.Mobile}).One(&result)
    if err != nil {
      //panic(err)
      //signal=false
      signal=false
    }else{
      signal=true
      if t.Password==result.Password {
        signal=true
      }else{
        signal=false
      }
    }
    fmt.Println(signal)
  resp, _ := json.Marshal(result)

    w.Header().Set("Content-Type", "application/json")
  w.Header().Set("charset", "utf-8")
  w.Header().Set("Access-Control-Allow-Origin", "*")
  w.Write([]byte(resp))

}

type RandomBool struct {
  Data bool
}

func signUpTheUser(w http.ResponseWriter, r *http.Request){
  decoder := json.NewDecoder(r.Body)

  var t SignHashWith
  var Status bool
  err := decoder.Decode(&t)

  if err != nil {
    //panic(err)
    Status=false
  }else{
    Status=true
  }

  chn:=make(chan bool)
  chn1:=make(chan string)
  co:=t.Mobile

  go saveMongoWithUserSignData(&chn,&t)
  go saveUserHashIfNotAvailable(&chn1,&co)

  var retHash bool
   retHash_1:=""

  for i := 0; i < 2; i++ {
      select {
        case msg := <-chn:
              retHash=msg
        case msg_1 := <-chn1:
              retHash_1=msg_1
      }
  }

  //fmt.Println(t.Dp)
  Status=retHash
  if retHash_1!="" {
    Status=true
  }else{
    Status=false
  }
  m:=SignSuccessRet{
    Status:Status}
  resp, _ := json.Marshal(m)



  w.Header().Set("Content-Type", "application/json")
  w.Header().Set("charset", "utf-8")
  w.Header().Set("Access-Control-Allow-Origin", "*")
  w.Write([]byte(resp))
}


func saveMongoWithUserSignData(mongoChannel *chan bool,hash *SignHashWith) {
  fmt.Println("executing mango")

  dialInfo := &mgo.DialInfo{
    Addrs:    []string{"rita-shard-00-00-qk9t0.mongodb.net:27017","rita-shard-00-01-qk9t0.mongodb.net:27017","rita-shard-00-02-qk9t0.mongodb.net:27017"},
    Database: "test",
    Username: "badcodercpp",
    Password: "Ajayjha93",
    ReplicaSetName: "Rita-shard-0",
    Source: "admin",
    DialServer: func(addr *mgo.ServerAddr) (net.Conn, error) {
        return tls.Dial("tcp", addr.String(), &tls.Config{})
    },
    Timeout: time.Second * 10000,
  }
  session, err := mgo.DialWithInfo(dialInfo)
  defer session.Close()


  fmt.Println(session)
  if(err != nil){
    panic(err)
  }
  c := session.DB("test").C("userdata")
  
  err = c.Insert(hash)
  fmt.Println(c)
    if err != nil {
      panic(err)
    }

    
    alllllss:=true
    *mongoChannel<-alllllss
}


func saveUserHashIfNotAvailable(mongoChannel *chan string,hash *string){
  dialInfo := &mgo.DialInfo{
    Addrs:    []string{"rita-shard-00-00-qk9t0.mongodb.net:27017","rita-shard-00-01-qk9t0.mongodb.net:27017","rita-shard-00-02-qk9t0.mongodb.net:27017"},
    Database: "test",
    Username: "badcodercpp",
    Password: "Ajayjha93",
    ReplicaSetName: "Rita-shard-0",
    Source: "admin",
    DialServer: func(addr *mgo.ServerAddr) (net.Conn, error) {
        return tls.Dial("tcp", addr.String(), &tls.Config{})
    },
    Timeout: time.Second * 10000,
  }
  session, err := mgo.DialWithInfo(dialInfo)
  defer session.Close()
  signal:=true
  if(err != nil){
    panic(err)
  }
  c := session.DB("test").C("kitchatHash")
    // Query One  4
    result := HashInserDM{}
    err = c.Find(bson.M{"Hash": *hash}).One(&result)
    if err != nil {
      //panic(err)
      //signal=false
      res:=HashInserDM{}
        res.Hash=*hash
      err1:=c.Insert(res)
      //tt,_:=time.ParseDuration("10000000h")
      if err1!=nil {
        signal=false
        
        //result=res
      }else{
        result=res
        signal=true
      }
    }else{
      signal=true
    }
    if signal {
      alllllss:=result.Hash
      *mongoChannel<-alllllss
    }else{
      alllllss:=""
      *mongoChannel<-alllllss
    }
}

type GetBoolSetterFCM struct {
  All []interface{}
}

func setFCMCandidate(w http.ResponseWriter, r *http.Request) {
  
vars := mux.Vars(r)
  mkiloj:=vars["candidate"]
  mkiloj_1:=vars["status"];
  stat:=false
  if mkiloj_1=="ok" {
    stat=true
  }else{
    stat=false
  }
  
    s1:=" merge (ee:FCMCandidate{mobile:'"
    ll:=mkiloj
    s3:="'}) with ee set ee.online = '"
    bgt:=stat
    ssb:="' return collect(ee.online) as all"
    s12 := fmt.Sprint(s1,ll,s3,bgt,ssb)
    fmt.Println(s12)
    driver := bolt.NewDriver()
      conn, err := driver.OpenNeo("bolt://rita:b.PuhuqVThYfCn.fvurl1e25g7fzyCI@hobby-panhpmpgjildgbkepcdcklol.dbs.graphenedb.com:24786?tls=true")
      defer conn.Close()
      if err != nil {
        panic(err)
      }
      data, _, _, _ := conn.QueryNeoAll(s12, nil)
      //fmt.Println("hooo")
      //fmt.Printf("COLUMNS: %#v\n", rowsMetadata["fields"].([]interface{}))  // COLUMNS: n.foo,n.bar
      //fmt.Printf("FIELDS: %s\n", data[0][0].([]interface{})) // FIELDS: 1 2.2
      m:=GetBoolSetterFCM{
        All:data[0][0].([]interface{})}
        resp, _ := json.Marshal(m)
        fmt.Println(string(resp))
      w.Header().Set("Content-Type", "application/json")
      w.Header().Set("charset", "utf-8")
      w.Header().Set("Access-Control-Allow-Origin", "*")
      w.Write([]byte(resp))
}

func getFCMCandidate(w http.ResponseWriter, r *http.Request) {
  
vars := mux.Vars(r)
  mkiloj:=vars["candidate"]
  
    s1:="Match (ee:FCMCandidate) where ee.mobile ='"
    b:=mkiloj
    s2:="' return collect(ee.online) as all"
    s12 := fmt.Sprint(s1,b,s2)
    fmt.Println(s12)
    driver := bolt.NewDriver()
      conn, err := driver.OpenNeo("bolt://rita:b.PuhuqVThYfCn.fvurl1e25g7fzyCI@hobby-panhpmpgjildgbkepcdcklol.dbs.graphenedb.com:24786?tls=true")
      defer conn.Close()
      if err != nil {
        panic(err)
      }
      data, _, _, _ := conn.QueryNeoAll(s12, nil)
      //fmt.Println("hooo")
      //fmt.Printf("COLUMNS: %#v\n", rowsMetadata["fields"].([]interface{}))  // COLUMNS: n.foo,n.bar
      //fmt.Printf("FIELDS: %s\n", data[0][0].([]interface{})) // FIELDS: 1 2.2
      m:=GetBoolSetterFCM{
        All:data[0][0].([]interface{})}
        resp, _ := json.Marshal(m)
        fmt.Println(string(resp))
      w.Header().Set("Content-Type", "application/json")
      w.Header().Set("charset", "utf-8")
      w.Header().Set("Access-Control-Allow-Origin", "*")
      w.Write([]byte(resp))
}

func stringInSlice(a string, list []string) bool {
    for _, b := range list {
        if b == a {
            return true
        }
    }
    return false
}

func getAllPeopleByCharecterName(w http.ResponseWriter, r *http.Request) {
  vars := mux.Vars(r)
  mkiloj:=vars["target"]
  dialInfo := &mgo.DialInfo{
    Addrs:    []string{"rita-shard-00-00-qk9t0.mongodb.net:27017","rita-shard-00-01-qk9t0.mongodb.net:27017","rita-shard-00-02-qk9t0.mongodb.net:27017"},
    Database: "test",
    Username: "badcodercpp",
    Password: "Ajayjha93",
    ReplicaSetName: "Rita-shard-0",
    Source: "admin",
    DialServer: func(addr *mgo.ServerAddr) (net.Conn, error) {
        return tls.Dial("tcp", addr.String(), &tls.Config{})
    },
    Timeout: time.Second * 10000,
  }
  session, err := mgo.DialWithInfo(dialInfo)
  defer session.Close()
  //signal:=true
  if(err != nil){
    panic(err)
  }
  c := session.DB("test").C("userdata")
  cc := session.DB("test").C("Request")
    // Query One  4
    result := make([]SignHashWith,0)
    res:=make([]MyAllAppNotification,0)
    //ml:= make(map[string]string)
    //ml["Hash"]=mkiloj
    xx:=make([]string,0)
    err1 := cc.Find(bson.M{"$or":[]bson.M{bson.M{"to":bson.M{"$eq":mkiloj}},bson.M{"from":bson.M{"$eq":mkiloj}}}}).All(&res)
    if err1 != nil {
      panic(err)
      //signal=false
      }

      for _, v := range res {
        if !stringInSlice(v.To,xx) {
          xx = append(xx, v.To)
        }
        if !stringInSlice(v.From,xx) {
          xx = append(xx, v.From)
        }
      }

    err = c.Find(bson.M{"$and":[]bson.M{bson.M{"mobile":bson.M{"$ne":mkiloj}},bson.M{"mobile":bson.M{"$nin":xx}}}}).All(&result)
    if err != nil {
      panic(err)
      //signal=false
      }

      //hs:=result.Hash
  
  resp, _ := json.Marshal(result)
  w.Header().Set("Content-Type", "application/json")
  w.Header().Set("charset", "utf-8")
  w.Header().Set("Access-Control-Allow-Origin", "*")
  w.Write([]byte(resp))
}

func getAllVisibleForMe(w http.ResponseWriter, r *http.Request) {
   
  vars := mux.Vars(r)
  mkiloj:=vars["target"]
  dialInfo := &mgo.DialInfo{
    Addrs:    []string{"rita-shard-00-00-qk9t0.mongodb.net:27017","rita-shard-00-01-qk9t0.mongodb.net:27017","rita-shard-00-02-qk9t0.mongodb.net:27017"},
    Database: "test",
    Username: "badcodercpp",
    Password: "Ajayjha93",
    ReplicaSetName: "Rita-shard-0",
    Source: "admin",
    DialServer: func(addr *mgo.ServerAddr) (net.Conn, error) {
        return tls.Dial("tcp", addr.String(), &tls.Config{})
    },
    Timeout: time.Second * 10000,
  }
  session, err := mgo.DialWithInfo(dialInfo)
  defer session.Close()
  //signal:=true
  if(err != nil){
    panic(err)
  }
  c := session.DB("test").C("userdata")
  cc := session.DB("test").C("Request")
    // Query One  4
    result := make([]SignHashWith,0)
    res:=make([]MyAllAppNotification,0)
    //ml:= make(map[string]string)
    //ml["Hash"]=mkiloj
    xx:=make([]string,0)
    err1 := cc.Find(bson.M{"$or":[]bson.M{bson.M{"to":bson.M{"$eq":mkiloj}},bson.M{"from":bson.M{"$eq":mkiloj}}}}).All(&res)
    if err1 != nil {
      panic(err)
      //signal=false
      }

      for _, v := range res {
        if !stringInSlice(v.To,xx) {
          xx = append(xx, v.To)
        }
        if !stringInSlice(v.From,xx) {
          xx = append(xx, v.From)
        }
      }

    err = c.Find(bson.M{"$and":[]bson.M{bson.M{"mobile":bson.M{"$ne":mkiloj}},bson.M{"mobile":bson.M{"$nin":xx}}}}).All(&result)
    if err != nil {
      panic(err)
      //signal=false
      }

      //hs:=result.Hash
  
  resp, _ := json.Marshal(result)
  w.Header().Set("Content-Type", "application/json")
  w.Header().Set("charset", "utf-8")
  w.Header().Set("Access-Control-Allow-Origin", "*")
  w.Write([]byte(resp))
}

func onlyVerifyUserHash(mongoChannel *chan string,hash *string){
  dialInfo := &mgo.DialInfo{
    Addrs:    []string{"rita-shard-00-00-qk9t0.mongodb.net:27017","rita-shard-00-01-qk9t0.mongodb.net:27017","rita-shard-00-02-qk9t0.mongodb.net:27017"},
    Database: "test",
    Username: "badcodercpp",
    Password: "Ajayjha93",
    ReplicaSetName: "Rita-shard-0",
    Source: "admin",
    DialServer: func(addr *mgo.ServerAddr) (net.Conn, error) {
        return tls.Dial("tcp", addr.String(), &tls.Config{})
    },
    Timeout: time.Second * 10000,
  }
  session, err := mgo.DialWithInfo(dialInfo)
  defer session.Close()
  signal:=true
  if(err != nil){
    panic(err)
  }
  c := session.DB("test").C("kitchatHash")
    // Query One  4
    result := HashInserDM{}
    //ml:= make(map[string]string)
    //ml["hash"]=*hash
    err = c.Find(bson.M{"hash":*hash}).One(&result)
    if err != nil {
      //panic(err)
      //signal=false
      signal=false
    }else{
      signal=true
    }
    if signal {
      alllllss:=result.Hash
      *mongoChannel<-alllllss
    }else{
      alllllss:=""
      *mongoChannel<-alllllss
    }
}



type OnlyHashDM struct {
  Hash string
}

type AllDM struct {
  All []OnlyHashDM
}

type HashInserDM struct{
  ID bson.ObjectId `bson:"_id,omitempty"`
  Hash string
}

func getAllData(w http.ResponseWriter, r *http.Request){
  dialInfo := &mgo.DialInfo{
    Addrs:    []string{"rita-shard-00-00-qk9t0.mongodb.net:27017","rita-shard-00-01-qk9t0.mongodb.net:27017","rita-shard-00-02-qk9t0.mongodb.net:27017"},
    Database: "test",
    Username: "badcodercpp",
    Password: "Ajayjha93",
    ReplicaSetName: "Rita-shard-0",
    Source: "admin",
    DialServer: func(addr *mgo.ServerAddr) (net.Conn, error) {
        return tls.Dial("tcp", addr.String(), &tls.Config{})
    },
    Timeout: time.Second * 10000,
  }
  session, err := mgo.DialWithInfo(dialInfo)
  //var result AllDM
  result:=make([]OnlyHashDM,0)
  defer session.Close()
  if(err != nil){
    panic(err)
  }
  c := session.DB("test").C("kitchatHash")
    // Query One  4
    //result := AllDM{}
    err = c.Find(nil).All(&result)
   if(err != nil){
    panic(err)
  }

  resp, _ := json.Marshal(result)
  w.Header().Set("Content-Type", "application/json")
  w.Header().Set("charset", "utf-8")
  w.Header().Set("Access-Control-Allow-Origin", "*")
  w.Write([]byte(resp))
}

type GotextTest struct {
  Res string
}

func removeHashedData(w http.ResponseWriter, r *http.Request){
  dialInfo := &mgo.DialInfo{
    Addrs:    []string{"rita-shard-00-00-qk9t0.mongodb.net:27017","rita-shard-00-01-qk9t0.mongodb.net:27017","rita-shard-00-02-qk9t0.mongodb.net:27017"},
    Database: "test",
    Username: "badcodercpp",
    Password: "Ajayjha93",
    ReplicaSetName: "Rita-shard-0",
    Source: "admin",
    DialServer: func(addr *mgo.ServerAddr) (net.Conn, error) {
        return tls.Dial("tcp", addr.String(), &tls.Config{})
    },
    Timeout: time.Second * 10000,
  }
  session, err := mgo.DialWithInfo(dialInfo)
  //var result AllDM
  //result:=make([]OnlyHashDM,0)
  defer session.Close()
  if(err != nil){
    panic(err)
  }
  c := session.DB("test").C("kitchatHash")
  m:=session.DB("test").C("userdata")
    // Query One  4
    //result := AllDM{}
    _,err = c.RemoveAll(nil)
   if(err != nil){
    panic(err)
  }
  _,err = m.RemoveAll(nil)
   if(err != nil){
    panic(err)
  }
  result:=GotextTest{Res:"ok"}
  resp, _ := json.Marshal(result)
  w.Header().Set("Content-Type", "application/json")
  w.Header().Set("charset", "utf-8")
  w.Header().Set("Access-Control-Allow-Origin", "*")
  w.Write([]byte(resp))
}

func getidDataUrl(w http.ResponseWriter, r *http.Request){
  vars := mux.Vars(r)
  mkiloj:=vars["hash"]
  //mkiloj_2:=vars["interactionId"]
  fmt.Println(mkiloj)

  dialInfo := &mgo.DialInfo{
    Addrs:    []string{"rita-shard-00-00-qk9t0.mongodb.net:27017","rita-shard-00-01-qk9t0.mongodb.net:27017","rita-shard-00-02-qk9t0.mongodb.net:27017"},
    Database: "test",
    Username: "badcodercpp",
    Password: "Ajayjha93",
    ReplicaSetName: "Rita-shard-0",
    Source: "admin",
    DialServer: func(addr *mgo.ServerAddr) (net.Conn, error) {
        return tls.Dial("tcp", addr.String(), &tls.Config{})
    },
    Timeout: time.Second * 10000,
  }
  session, err := mgo.DialWithInfo(dialInfo)
  defer session.Close()
  //signal:=true
  if(err != nil){
    panic(err)
  }
  c := session.DB("test").C("kitchatHash")
    // Query One  4
    result := HashInserDM{}
    //ml:= make(map[string]string)
    //ml["Hash"]=mkiloj

    err = c.Find(bson.M{"hash":"9836648105"}).One(&result)
    if err != nil {
      panic(err)
      //signal=false
      }

      //hs:=result.Hash
  
  resp, _ := json.Marshal(result)
  w.Header().Set("Content-Type", "application/json")
  w.Header().Set("charset", "utf-8")
  w.Header().Set("Access-Control-Allow-Origin", "*")
  w.Write([]byte(resp))

}

func serveHashTemplate(w http.ResponseWriter, r *http.Request){
  vars := mux.Vars(r)
  mkiloj:=vars["hash"]
  //mkiloj_2:=vars["interactionId"]
  fmt.Println(mkiloj)

  chn:=make(chan string)

  go onlyVerifyUserHash(&chn,&mkiloj)

  var retHash string

  for i := 0; i < 1; i++ {
      select {
        case msg := <-chn:
              retHash=msg
      }
  }

  m:=KitchatHash{
  Hash:retHash}
  resp, _ := json.Marshal(m)
  w.Header().Set("Content-Type", "application/json")
  w.Header().Set("charset", "utf-8")
  w.Header().Set("Access-Control-Allow-Origin", "*")
  w.Write([]byte(resp))

}

type  UserInner struct {
  Id string
  Name string
  Gender string 
  Interest string
  Topic []string
}

type G_USER struct {
  Data UserInner
  Zip string
}

func serveZipTemplate(w http.ResponseWriter, r *http.Request){

  decoder := json.NewDecoder(r.Body)

  var t G_USER
  var signal bool
  err := decoder.Decode(&t)

  if err != nil {
    panic(err)
    signal=false
  }else{
    signal=true
  }
  fmt.Println(signal)
  fmt.Println(t.Data)
  //vars := mux.Vars(r)
  id:=t.Data.Id
  name:=t.Data.Name
  gender:=t.Data.Gender
  interest:=t.Data.Interest
  topic:=t.Data.Topic
  zip:=t.Zip
  //mkiloj_2:=vars["interactionId"]
  fmt.Println(id)

  dialInfo := &mgo.DialInfo{
    Addrs:    []string{"rita-shard-00-00-qk9t0.mongodb.net:27017","rita-shard-00-01-qk9t0.mongodb.net:27017","rita-shard-00-02-qk9t0.mongodb.net:27017"},
    Database: "test",
    Username: "badcodercpp",
    Password: "Ajayjha93",
    ReplicaSetName: "Rita-shard-0",
    Source: "admin",
    DialServer: func(addr *mgo.ServerAddr) (net.Conn, error) {
        return tls.Dial("tcp", addr.String(), &tls.Config{})
    },
    Timeout: time.Second * 10000,
  }
  session, err := mgo.DialWithInfo(dialInfo)
  defer session.Close()
  //signal:=true
  if(err != nil){
    panic(err)
  }
  c := session.DB("test").C(zip)
    // Query One  4
  lks:=UserInner{Id:id,Name:name,Gender:gender,Interest:interest,Topic:topic}
    rs := G_USER{Data:lks,Zip:t.Zip}
    //ml:= make(map[string]string)
    //ml["Hash"]=mkiloj

    //err = c.Find(bson.M{"data":bson.M{"id":id}}).One(&rs)
    //if err != nil {
      //panic(err)
      //signal=false
      c.Upsert(bson.M{"data":bson.M{"id":id}},rs)
      //}

      //hs:=result.Hash
  
  resp, _ := json.Marshal(rs)
  w.Header().Set("Content-Type", "application/json")
  w.Header().Set("charset", "utf-8")
  w.Header().Set("Access-Control-Allow-Origin", "*")
  w.Write([]byte(resp))

}


type AllG_USER struct {
  All []string
}


func serveAllZipTemplate(w http.ResponseWriter, r *http.Request){
  vars := mux.Vars(r)
  zip:=vars["zip"]
  //mkiloj_2:=vars["interactionId"]
  fmt.Println(zip)

  dialInfo := &mgo.DialInfo{
    Addrs:    []string{"rita-shard-00-00-qk9t0.mongodb.net:27017","rita-shard-00-01-qk9t0.mongodb.net:27017","rita-shard-00-02-qk9t0.mongodb.net:27017"},
    Database: "test",
    Username: "badcodercpp",
    Password: "Ajayjha93",
    ReplicaSetName: "Rita-shard-0",
    Source: "admin",
    DialServer: func(addr *mgo.ServerAddr) (net.Conn, error) {
        return tls.Dial("tcp", addr.String(), &tls.Config{})
    },
    Timeout: time.Second * 10000,
  }
  session, err := mgo.DialWithInfo(dialInfo)
  defer session.Close()
  //signal:=true
  if(err != nil){
    panic(err)
  }
  c := session.DB("test").C(zip)
    // Query One  4
    rs := make([]G_USER,0)
    //ml:= make(map[string]string)
    //ml["Hash"]=mkiloj

    err = c.Find(nil).All(&rs)
    if err != nil {
      //rs=AllG_USER{All:[]}
      //rs.All=["error occured"]
      //panic(err)
      //signal=false
      //c.Insert(&rs)
      lll:=UserInner{}
      z:=G_USER{Data:lll,Zip:"no data found"}
      rs=append(rs,z)
      
      }
      //hs:=result.Hash
  
  resp, _ := json.Marshal(rs)
  w.Header().Set("Content-Type", "application/json")
  w.Header().Set("charset", "utf-8")
  w.Header().Set("Access-Control-Allow-Origin", "*")
  w.Write([]byte(resp))

}



type Index struct {
    Key        []string // Index key fields; prefix name with dash (-) for descending order
    Unique     bool     // Prevent two documents from having the same index key
    DropDups   bool     // Drop documents with the same index key as a previously indexed one
    Background bool     // Build index in background and return immediately
    Sparse     bool     // Only index documents containing the Key fields

    ExpireAfter time.Duration // Periodically delete docs with indexed time.Time older than that.

    Name string // Index name, computed by EnsureIndex

    Bits, Min, Max int // Properties for spatial indexes
}

//hours, _ := time.ParseDuration("10h")



func main() {

    testBooks()
    //timesOfIndia();
    /*yt := make(chan interface{})
    go youTubeVideo(&yt,"crime patrol");

    msg:=<-yt
    fmt.Println(msg)*/



    //connectMong("hello")


    // random no testing with encryption and decryption statement



    /*randTemp:="badcodercpp@gmail.com"
    mnk:=generateHashAndReplicateToDbPratibhaPlease(&randTemp)
    jjh:=mnk[0]
    jjm:=mnk[1]
    kli, _ := strconv.ParseInt(jjm, 10, 64)
    //btB:=[]byte(jjh)
    fmt.Println(jjh)
    fmt.Println(jjm)
    fmt.Println(kli)
    rest:=authorizeThisFace(&jjh,&kli)
    fmt.Println("seperate_auth")
    fmt.Println(rest)*/

    //end random no testing with encryption and decryption



    //paypal beg








    //paypal end


    //twillio start

      //sendOtpToTheMyUserAjay("+919470717982","Dhyan se dekhiya yahi hai ye ladki")

    // twillio end

    //  getMyAllNotification

    // alterDatingRequestTemplate


    port := os.Getenv("PORT")
    if port == "" {
      log.Fatal("$PORT must be set")
    }
    r := mux.NewRouter()
    r.HandleFunc("/", serveMainTemplate)
    //r.HandleFunc("/signalRTC/{userId}/{rtcId}",signalRTCHandler).Methods("GET")
    r.HandleFunc("/login", serveTemplate)
    r.HandleFunc("/saveMyDatingRequestPlease", saveDatingRequestTemplate).Methods("POST")
    r.HandleFunc("/alterDatingRequestPlease", alterDatingRequestTemplate).Methods("POST")
    r.HandleFunc("/loginMeToApp", letTheUserLogin).Methods("POST")
    r.HandleFunc("/getMyAllNotification/{target}", getMyAllNotification).Methods("GET")
    r.HandleFunc("/getAllSuggestion/{target}", getAllVisibleForMe).Methods("GET")
    r.HandleFunc("/modifyFcmCandidate/{candidate}/{status}", setFCMCandidate).Methods("GET")
    r.HandleFunc("/checkFcmCandidate/{candidate}", getFCMCandidate).Methods("GET")
    r.HandleFunc("/signup", serveSignupTemplate)
    r.HandleFunc("/getAllD", getAllData)
    r.HandleFunc("/removeAllD", removeHashedData)
    r.HandleFunc("/signupApp", signUpTheUser).Methods("POST")
    r.HandleFunc("/any/{hash}/{Pkey}/{query}", serveAnyTemplate)
    r.HandleFunc("/faces/{query}", serveFaceTemplate)

    r.HandleFunc("/testdata/{hash}", getidDataUrl)
    r.HandleFunc("/verifyHash/{hash}", serveHashTemplate)

    r.HandleFunc("/zipvsid_anddata", serveZipTemplate).Methods("POST")
    r.HandleFunc("/alluserinazipcode/{zip}", serveAllZipTemplate)

    r.HandleFunc("/bot/{interactionId}", serveBotTemplate)
    r.HandleFunc("/workflow/{interactionId}/{workflowId}", serveWorkflowTemplate)
    r.HandleFunc("/getMyCommonActionFacePratibha/{reference_id}/{Pkey}/{query}", getMyCommonActionFacePratibhaHandler)
    r.HandleFunc("/myLinkedFaces", serveMyLinedFaceTemplate)
    r.HandleFunc("/templateData", serveDataTemplate)
    r.HandleFunc("/getMyOwnShopDetailsPratibhaPleaseLoveYou", getMyOwnShopDetailsPratibhaPleaseLoveYouHandler)
    r.HandleFunc("/wowPratibhaYouLooksLikeAnAngel", serveWowPratibhaYouLooksLikeAnAngel)
    r.HandleFunc("/wowPratibhaYouLooksLikeAnAngelPratibha", serveWowPratibhaYouLooksLikeAnAngelPratibha)
    r.HandleFunc("/anyMore",serveMoreVideosYoutube)
    r.HandleFunc("/anyBooks",serveBooks)
    r.HandleFunc("/createMyJobPratibhaPleaseForMeSorry",createMyJobPratibhaPleaseForMeSorryHandler).Methods("GET")
    r.HandleFunc("/myTokenAuthTest",myTokenAuthTest).Methods("GET")
    r.HandleFunc("/createMyDirectConfessionPratibhaPleaseRecordMySin",createMyDirectConfessionPratibhaPleaseRecordMySinHandler).Methods("GET")
    r.HandleFunc("/seekDonationForMeAjayPlsTlfDonate",seekDonationForMeAjayPlsTlfDonateHandler).Methods("GET")
    r.HandleFunc("/submitMyLoyalityFormPratibhaPlsInDatingZone",submitMyLoyalityFormPratibhaPlsInDatingZoneHandler).Methods("GET")
    r.HandleFunc("/recieveMyPaymentPratibhaPleaseYouAreOnlyHopeOfMine",recieveMyPaymentPratibhaPleaseYouAreOnlyHopeOfMineHandler).Methods("GET")
    r.HandleFunc("/deleteMyAuthHashPratibhaPleaseHelp",deleteMyAuthHashPratibhaPleaseHelpHandler).Methods("GET")
    r.HandleFunc("/myDpWillBeEdittedPratibha",myDpWillBeEdittedPratibhaHandler).Methods("GET")
    r.HandleFunc("/letMeSeekMyJobPratibhaPlease",letMeSeekMyJobPratibhaPleaseHandler).Methods("GET")
    r.HandleFunc("/tumheDillaggibhulJaniParegiYaar",tumheDillaggibhulJaniParegiYaarHandler)
    r.HandleFunc("/retriveMyEventSourceForMessagePratibha",retriveMyEventSourceForMessagePratibhaHandler)
    r.HandleFunc("/saveMyMessagePratibhaPleaseHelp",saveMyMessagePratibhaPleaseHelpHandler)
    r.HandleFunc("/replicateMyMessageToDbMongo",replicateMyMessageToDbMongoHandler)
    r.HandleFunc("/getMyPostForProfilePratibhaPlease",getMyPostForProfilePratibhaPleaseHandler)
    r.HandleFunc("/sendMyDateRequestPratibhaPlease",sendMyDateRequestPratibhaPleaseHandler)
    r.HandleFunc("/letMeSeeMyPlaylistSongPratibhaPlease",letMeSeeMyPlaylistSongPratibhaPleaseHandler)
    r.HandleFunc("/satayeMenuKyonOmyYaraIloveYou",satayeMenuKyonOmyYaraIloveYouHandler)
    r.HandleFunc("/LinkMyFaceWithThemPratibhaPlease",serveLinkMyFaceWithThemPratibhaPlease)
    r.HandleFunc("/getMyRelatedItemsToBuyPratibhaPlease",getMyRelatedItemsToBuyPratibhaPleaseHandler)
    r.HandleFunc("/addItemToMyShopPratibhaPleaseLU",addItemToMyShopPratibhaPleaseLUHandler)
    r.HandleFunc("/iLoveYouPratibhaSharmaAndIWillGetYouShopCreate",iLoveYouPratibhaSharmaAndIWillGetYouShopCreateHandler)
    r.HandleFunc("/amazonQuery",serveAmazonQuery)
    r.HandleFunc("/syncAmazonMongo",syncAmazonMongoWebService)
    r.HandleFunc("/redirect",redirectHandler)
    r.HandleFunc("/newPlaylistCreation",newPlaylistCreationHandler)
    r.HandleFunc("/getMyAllRelatedPostBilla",getMyAllRelatedPostBillaHandler)
    r.HandleFunc("/getMyAllRelatedFacesBilla",getMyAllRelatedFacesBillaHandler)
    r.HandleFunc("/saveMyPostWithAttachment",saveMyPostWithAttachmentHandler).Methods("POST")
    r.HandleFunc("/mostPopularVideo/{newsType}",timesOfIndia)
    r.HandleFunc("/linkAuth", serveAuth).Methods("POST")
    r.HandleFunc("/signUpMePlease", serveAuthAndSignUp).Methods("POST")
    r.HandleFunc("/fileUploadItemIcon", uploadAndProcessMyNewDp).Methods("POST")
    r.HandleFunc("/changingBytesOfCloudinaryDp", changingBytesOfCloudinaryDpHandler).Methods("POST")
    r.HandleFunc("/linkAuth",notFound)
    r.NotFoundHandler = http.HandlerFunc(notFound)
    r.PathPrefix("/").Handler(http.FileServer(http.Dir("./public/")))
    http.Handle("/", r)
    log.Println("Listening...to all")
    http.ListenAndServe(":"+port, r)












      //Parse result
      /*if err == nil {
        aws := new(ItemLookupResponse)
        xml.Unmarshal([]byte(result), aws)
        //TODO: Use "aws" freely :-)
      }*/








}


func printIDs(sectionName string, matches map[string]string) {
        fmt.Println("%v:\n", sectionName)
        for id, _ := range matches {
                fmt.Println("[%v]\n", id)
        }
        fmt.Println("\n\n")
}
