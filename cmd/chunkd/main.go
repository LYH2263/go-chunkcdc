package main
import (
        "encoding/json"; "flag"; "io"; "log"; "net/http"; "os"; "path/filepath"
        "github.com/LYH2263/go-chunkcdc"
)
func main() {
        addr := flag.String("addr", ":8124", "listen")
        flag.Parse()
        eng, err := chunkcdc.New(chunkcdc.Options{})
        if err != nil { log.Fatal(err) }
        defer eng.Close()
        web := "web"
        if _, err := os.Stat(web); err != nil { web = filepath.Join("..", "..", "web") }
        mux := http.NewServeMux()
        mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir(web))))
        mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
                if r.URL.Path != "/" { http.NotFound(w, r); return }
                http.ServeFile(w, r, filepath.Join(web, "index.html"))
        })
        mux.HandleFunc("/api/status", func(w http.ResponseWriter, r *http.Request) {
                _ = json.NewEncoder(w).Encode(eng.Stats())
        })
        mux.HandleFunc("/api/trial", func(w http.ResponseWriter, r *http.Request) {
                b, _ := io.ReadAll(r.Body)
                if err := eng.Put(r.Context(), "trial", b, nil); err != nil { http.Error(w, err.Error(), 500); return }
                _ = json.NewEncoder(w).Encode(map[string]any{"ok": true})
        })
        log.Printf("%s on %s", "内容定义分块", *addr)
        log.Fatal(http.ListenAndServe(*addr, mux))
}
