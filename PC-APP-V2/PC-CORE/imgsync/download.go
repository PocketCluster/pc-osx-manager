package imgsync

import (
    "archive/tar"
    "compress/gzip"
    "io"
    "net/http"
    "os"
    "path/filepath"
    "strings"
    "sync"

    log "github.com/Sirupsen/logrus"
    "github.com/pkg/errors"
    "xi2.org/x/xz"
)

func localGzipUncompress() {
    // write files to directory
    out, err := os.Create("protobuf-src-2.5.0.tar")
    if err != nil {
        log.Fatal(errors.WithStack(err))
    }
    defer out.Close()

    // write files to directory
    in, err := os.Open("protobuf-src-2.5.0.tar.gz")
    if err != nil {
        log.Fatal(errors.WithStack(err))
    }
    defer in.Close()

    gr, err := gzip.NewReader(in)
    if err != nil {
        log.Fatal(errors.WithStack(err))
    }
    defer gr.Close()

    n, err := io.Copy(out, gr)
    log.Printf("written  %v %v", n, err)
}

func ioPipedDownload() {
    var wg sync.WaitGroup

    // write files to directory
    out, err := os.Create("protobuf-src-2.5.0.tar.gz")
    if err != nil {
        log.Fatal(errors.WithStack(err))
    }
    defer out.Close()

    // http reads with pipe
    pr, pw := io.Pipe()
    defer pw.Close()
    defer pr.Close()

    go func() {
        wg.Add(1)
        defer wg.Done()

        n, err := io.Copy(out, pr)
        log.Printf("written  %v %v", n, err)
    }()

    resp, err := http.Get("http://localhost:8080/protobuf-src-2.5.0.tar.gz")
    if err != nil {
        log.Error(errors.WithStack(err))
    }

    log.Print("Write request to pipe")
    n, err := io.Copy(pw, resp.Body)
    log.Printf("read  %v %v", n, err)
    resp.Body.Close()
    pw.Close()

    log.Print("Waiting to finish BG job...")
    wg.Wait()
}

func ioPipedClientDownload() {
    var wg sync.WaitGroup

    // write files to directory
    out, err := os.Create("protobuf-src-2.5.0.tar.gz")
    if err != nil {
        log.Fatal(errors.WithStack(err))
    }
    defer out.Close()

    // http reads with pipe
    pr, pw := io.Pipe()
    defer pr.Close()

    go func() {
        log.Print("Going in Background...")
        wg.Add(1)
        defer wg.Done()

        n, err := io.Copy(out, pr)
        log.Printf("written  %v %v", n, err)
    }()

    log.Print("Prep HTTP client")
    client := new(http.Client)
    request, err := http.NewRequest("GET", "http://localhost:8080/protobuf-src-2.5.0.tar.gz", nil)
    if err != nil {
        log.Error(errors.WithStack(err))
    }
    response, err := client.Do(request)
    if err != nil {
        log.Error(errors.WithStack(err))
    }

    io.Copy(pw, response.Body)
    response.Body.Close()
    pw.Close()

    log.Print("Waiting to finish BG job...")
    wg.Wait()
}

func pipedHttpGzipUncompress(client *http.Client, url string) error {
    var (
        outFilename string
        out *os.File
        pr *io.PipeReader
        pw *io.PipeWriter
        response *http.Response
        wg sync.WaitGroup
        written int64
        err error
    )

    // write files to directory
    urlcomp := strings.Split(url, "/")
    outFilename = strings.Replace(urlcomp[len(urlcomp) - 1], "gz", "", -1)
    out, err = os.Create(outFilename)
    if err != nil {
        return err
    }
    defer out.Close()

    // http reads with pipe
    pr, pw = io.Pipe()
    defer pr.Close()

    go func() {
        wg.Add(1)
        defer wg.Done()

        gr, err := gzip.NewReader(pr)
        if err != nil {
            log.Fatal(errors.WithStack(err))
        }
        defer gr.Close()

        written, err := io.Copy(out, gr)
        log.Printf("written  %v %v", written, err)
    }()

    response, err = client.Get(url)
    if err != nil {
        return err
    }

    written, err = io.Copy(pw, response.Body)
    log.Printf("read  %v %v", written, err)
    response.Body.Close()
    pw.Close()

    wg.Wait()
    return err
}

func streamHttpGzipUncompress(client *http.Client, url string) error {
    var (
        outFilename string
        out *os.File
        reader io.ReadCloser
        request *http.Request
        response *http.Response
        written int64
        err error
    )

    request, err = http.NewRequest("GET", url, nil)
    if err != nil {
        return err
    }
    request.Header.Add("Accept-Encoding", "gzip")
    response, err = client.Do(request)
    if err != nil {
        return err
    }
    defer response.Body.Close()

    // Check that the server actually sent compressed data
    reader, err = gzip.NewReader(response.Body)
    if err != nil {
        return err
    }
    defer reader.Close()

    urlcomp := strings.Split(url, "/")
    outFilename = strings.Replace(urlcomp[len(urlcomp) - 1], "gz", "", -1)
    out, err = os.Create(outFilename)
    if err != nil {
        return err
    }
    defer out.Close()

    // io.Copy : Copy copies from src to dst until either EOF is reached on src or an error occurs.
    // It returns the number of bytes copied and the first error encountered while copying, if any.
    written, err = io.Copy(out, reader)
    log.Printf("written  %v %v", written, err)
    return err
}

func streamHttpXzUncompress(client *http.Client, url, uncompPath string) error {
    var (
        reader *xz.Reader
        request *http.Request
        response *http.Response
        unarchive *tar.Reader
        err error
    )

    request, err = http.NewRequest("GET", url, nil)
    if err != nil {
        return err
    }
    request.Header.Add("Accept-Encoding", "xz")
    response, err = client.Do(request)
    if err != nil {
        return err
    }
    defer response.Body.Close()

    // Check that the server actually sent compressed data
    reader, err = xz.NewReader(response.Body, 0)
    if err != nil {
        return err
    }
    //defer reader.Close()

    unarchive = tar.NewReader(reader)
    for {
        header, err := unarchive.Next()
        if err == io.EOF {
            break
        } else if err != nil {
            return err
        }

        path := filepath.Join(uncompPath, header.Name)
        info := header.FileInfo()
        if info.IsDir() {
            if err = os.MkdirAll(path, info.Mode()); err != nil {
                return err
            }
            continue
        }
        file, err := os.OpenFile(path, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, info.Mode())
        if err != nil {
            return err
        }
        written, err := io.Copy(file, unarchive)
        if err != nil {
            file.Close()
            return err
        } else {
            log.Printf("written  %v", written)
        }
        err = file.Close()
        if err != nil {
            return err
        }
    }
    return nil
}

func DownSyncImageToPath(imageURL, targetPath string) {
    // Disable auto-decompression (http://stackoverflow.com/questions/13130341/reading-gzipped-http-response-in-go)
    transport := &http.Transport{
        //TLSClientConfig: tlsConfig,
        DisableCompression: true,
    }
    client := &http.Client{Transport: transport}


    err := streamHttpXzUncompress(client, imageURL, targetPath)
    if err != nil {
        log.Fatal(errors.WithStack(err))
    }
}