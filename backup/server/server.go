package main

import (
	"log"
	"net"
	"net/rpc"
	"os"
	"path/filepath"

	"github.com/fpkoehler/goExamples/rpc/shared"
)

type Set struct {
	id    shared.SetId
	files map[string]shared.FileInfo
}

var sets map[shared.SetId]Set

type Session struct {
	/* for uploading a file */
	file *os.File
	name string
}

var sessions map[int]Session

/*************************************************
 * RPC server
 *************************************************/

type Backup int

func (b *Backup) FileListFromServer(args *shared.SetId, reply *shared.FileInfoListReply) error {
	_, found := sets[*args]
	if found == false {
		sets[*args] = Set{
			id:    *args,
			files: make(map[string]shared.FileInfo),
		}
		log.Println("Created set", *args)
	}

	log.Println("FileListFromServer #files", len(sets[*args].files))
	reply.Files = sets[*args].files

	return nil
}

func (b *Backup) FileToServer(args *shared.FileInfoArg, reply *shared.BoolReply) error {
	_, found := sets[args.SetId]
	if found == false {
		sets[args.SetId] = Set{
			id:    args.SetId,
			files: make(map[string]shared.FileInfo),
		}
		log.Println("Created set", args.SetId)
	}

	log.Println("set", args.SetId, "updated:", args.FileInfo)
	sets[args.SetId].files[args.FileInfo.Name] = args.FileInfo

	reply.Status = true
	return nil
}

func (b *Backup) FileInit(args *shared.FileInitArg, reply *shared.FileInitReply) error {
	/* TODO: need a lock */
	sessionId := 1
	_, found := sessions[sessionId]
	if found == true {
		log.Println("FileInit: no more sessions available")
		reply.Id = 0
		return nil
	}

	localFileName := "/tmp/" + filepath.Base(args.Name)
	file, err := os.OpenFile(localFileName, os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Println("FileInit: could not open", localFileName, err.Error())
		reply.Id = 0
		return nil
	}
	session := Session{name: localFileName, file: file}

	sessions[sessionId] = session

	reply.Id = sessionId
	return nil
}

func (b *Backup) FileBlock(args *shared.FileBlockArg, reply *shared.BoolReply) error {
	session, found := sessions[args.Id]
	if found == false {
		reply.Status = false
		log.Println("FileBlock no session id", args.Id)
		return nil
	}

	session.file.Write(args.Data[:args.Size])
	log.Println("FileBlock session", args.Id, "wrote", len(args.Data[:args.Size]), "bytes to", session.name)
	reply.Status = true
	return nil
}

func (b *Backup) FileClose(args *shared.FileCloseArg, reply *shared.BoolReply) error {
	session, found := sessions[args.Id]
	if found == false {
		log.Println("File.Close could not find session id", args.Id)
		reply.Status = false
		return nil
	}
	session.file.Close()
	log.Println("FileClose removed session id", args.Id, "name", session.name)
	delete(sessions, args.Id)
	reply.Status = true
	return nil
}

/*************************************************
 * main
 *************************************************/

func main() {

	sets = make(map[shared.SetId]Set)
	sessions = make(map[int]Session)

	// Create an instance of struct which implements Backup interface
	backup := new(Backup)

	server := rpc.NewServer()
	server.RegisterName("Backup", backup)

	l, e := net.Listen("tcp", ":1234")
	if e != nil {
		log.Fatal("listen error:", e)
	}

	server.Accept(l)
}
