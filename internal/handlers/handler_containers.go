package handlers

import(
    "fmt"
    "log"
    "bytes"
    "regexp"
    "net/http"
    "github.com/neophenix/lxdepot/internal/lxd"
    "github.com/neophenix/lxdepot/internal/config"
)

func ContainerListHandler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "text/html")

    containerInfo, err := lxd.GetContainers("", "", true)
    if err != nil {
        log.Printf("Could not get container list %s\n", err.Error())
    }

    tmpl := readTemplate("container_list.tmpl")

    var out bytes.Buffer
    tmpl.ExecuteTemplate(&out, "base", map[string]interface{}{
        "Page": "containers",
        "Containers": containerInfo,
    })

    fmt.Fprintf(w, string(out.Bytes()))
}

func ContainerHandler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "text/html")

    reg := regexp.MustCompile("/container/(?P<Host>[^:]+):(?P<Name>.+)")
    match := reg.FindStringSubmatch(r.URL.Path)

    if len(match) != 3 {
        FourOhFourHandler(w,r)
        return
    }

    containerInfo, err := lxd.GetContainers(match[1], match[2], true)
    if err != nil {
        log.Printf("Could not get container list %s\n", err.Error())
    }
    if len(containerInfo) == 0 {
        FourOhFourHandler(w,r)
        return
    }

    var host *config.LXDhost
    for _, lxdh := range Conf.LXDhosts {
        if lxdh.Host == containerInfo[0].Host {
            host = lxdh
        }
    }

    tmpl := readTemplate("container.tmpl")

    var out bytes.Buffer
    tmpl.ExecuteTemplate(&out, "base", map[string]interface{}{
        "Page": "containers",
        "Host": host,
        "Container": containerInfo[0],
    })

    fmt.Fprintf(w, string(out.Bytes()))
}

func NewContainerHandler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "text/html")

    images, err := lxd.GetImages("")
    if err != nil {
        log.Printf("Could not get image list %s\n", err.Error())
    }

    tmpl := readTemplate("container_new.tmpl")

    var out bytes.Buffer
    tmpl.ExecuteTemplate(&out, "base", map[string]interface{}{
        "Page": "containers",
        "Hosts": Conf.LXDhosts,
        "Images": images,
    })

    fmt.Fprintf(w, string(out.Bytes()))
}
