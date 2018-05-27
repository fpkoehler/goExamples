package main

import (
	"errors"
	"io"
	"log"
	"net"
	"net/rpc"
	"os"
	"os/signal"
	"path/filepath"
	"time"

	"github.com/fpkoehler/goExamples/rpc/shared"
)

/*************************************************
 * RPC server
 *************************************************/

type BackupRpc struct {
	client *rpc.Client
}

var backup BackupRpc
var setId shared.SetId

func (b *BackupRpc) PushFile(path string, info os.FileInfo) error {
	fileInfo := shared.FileInfo{
		Name:    path,
		Size:    info.Size(),
		Mode:    info.Mode(),
		ModTime: info.ModTime(),
	}
	fileInfoArg := shared.FileInfoArg{
		SetId:    setId,
		FileInfo: fileInfo,
	}
	var reply shared.BoolReply
	err := b.client.Call("Backup.FileToServer", fileInfoArg, &reply)
	if err != nil {
		log.Fatal("FileToServer rpc error:", err)
	}
	if reply.Status == false {
		log.Fatal("FileToServer returned false.  file:", info.Name())
	}

	// Server go it, now update our local copy of what the server has
	//serverFiles[path] = fileInfo
	serverFiles[path] = fileInfo

	return nil
}

func (b *BackupRpc) PushFile2(path string, info os.FileInfo) error {
	var boolReply shared.BoolReply

	file, err := os.Open(path)
	if err != nil {
		log.Println("PushFile2: could not open", path, err.Error())
		return err
	}
	defer file.Close()

	var fileInitReply shared.FileInitReply
	err = b.client.Call("Backup.FileInit", shared.FileInitArg{Name: path}, &fileInitReply)
	if err != nil {
		log.Fatal("FileInit rpc error:", err)
	}
	if fileInitReply.Id == 0 {
		log.Println("FileInit returned Id=0")
		return errors.New("Server unable to accept file")
	}
	defer func() {
		log.Println("Backup.FileClose", fileInitReply.Id, "file", path)
		err = b.client.Call("Backup.FileClose", shared.FileCloseArg{Id: fileInitReply.Id}, &boolReply)
		if err != nil {
			log.Fatal("FileClose rpc error:", err)
		}
	}()

	blocks := int(info.Size() / shared.BLOCK_SIZE)
	if info.Size()%shared.BLOCK_SIZE != 0 {
		blocks += 1
	}

	log.Printf("Upload %s in %d blocks\n", path, blocks)

	fileBlockRequest := shared.FileBlockArg{Id: fileInitReply.Id, Data: make([]byte, shared.BLOCK_SIZE)}

	for blockId := 0; blockId < blocks; blockId++ {
		offset := int64(blockId) * shared.BLOCK_SIZE
		fileBlockRequest.Size, err = file.ReadAt(fileBlockRequest.Data, offset)
		if err != nil && err != io.EOF {
			return err
		}

		err = b.client.Call("Backup.FileBlock", fileBlockRequest, &boolReply)
		if err != nil {
			log.Fatal("FileBlock rpc error:", err)
		}
		if boolReply.Status == false {
			log.Fatal("FileBlock returned false.  file:", path)
		}

		//if i%((blocks-blockId)/100+1) == 0 {
		log.Println("Uploading %s [%d/%d] blocks", path, blockId+1, blocks-blockId)
		//}
	}
	log.Println("Upload %s completed", path)
	return nil
}

func (b *BackupRpc) getServerFileList() (map[string]shared.FileInfo, error) {
	var reply shared.FileInfoListReply
	err := b.client.Call("Backup.FileListFromServer", setId, &reply)
	if err != nil {
		log.Fatal("FileListFromServer rpc error:", err)
	}
	return reply.Files, err
}

/*****************************************************************************/

var serverFiles map[string]shared.FileInfo

/*****************************************************************************/

func visitFile(path string, info os.FileInfo, err error) error {
	if err != nil {
		log.Printf("prevent panic by handling failure accessing a path %q: %v\n", path, err)
		return err
	}
	if info.IsDir() {
		return nil
	}

	//	fullPath := filepath.Join(path, info.Name())
	fullPath := path
	serverFile, found := serverFiles[fullPath]
	if !found {
		log.Println("Not on server:", path)
		backup.PushFile(path, info)
		backup.PushFile2(path, info)
	} else if serverFile.ModTime.Before(info.ModTime()) {
		log.Println("Time newer from server", path)
		backup.PushFile2(path, info)
	}

	return nil
}

func main() {
	var err error

	// Tries to connect to localhost:1234 (The port the rpc server is listening to)
	conn, err := net.Dial("tcp", "localhost:1234")
	if err != nil {
		log.Fatal("Connecting:", err)
	}

	// Create a struct, that mimics all methods provided by interface.
	// It is not compulasory, we are doing it here, just to simulate a traditional method call.
	backup = BackupRpc{client: rpc.NewClient(conn)}

	setId.Host, err = os.Hostname()
	if err != nil {
		log.Fatal("Unable to get computer name:", err)
	}

	dir, err := os.Getwd()
	if err != nil {
		log.Fatal("Unable to get directory/folder path:", err)
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	serverFiles, err = backup.getServerFileList()
	if err != nil {
		log.Println("backup.getServerFileList() error:", err)
	}
	log.Println(serverFiles)

	loop := true
	for loop {
		log.Println("-----------")
		err = filepath.Walk(dir, visitFile)
		if err != nil {
			log.Printf("error walking the path %q: %v\n", dir, err)
		}

		select {
		case <-c:
			log.Println("Interrupt")
			loop = false
		case <-time.After(10 * time.Second):
		}
	}
}
