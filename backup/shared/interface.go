package shared

type Backup interface {
	FileListFromServer(args *SetId, reply *FileInfoListReply)
	FileToServer(args *FileInfo, replay *BoolReply)

	FileInit(args *FileInitArg, reply *FileInitReply)
	FileBlock(args *FileBlockArg, reply *BoolReply)
	FileClose(args *FileCloseArg, reply *BoolReply)
}
