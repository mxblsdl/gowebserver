package models

import "time"

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Folder struct {
	Id         int64
	UserId     int
	FolderName string
	ParentId   int64
	CreatedAt  time.Time
}

// todo make this match how files are stored in the database
type File struct {
	Id        int64
	FileName  string
	Size      int64
	CreatedAt time.Time
	Content   []byte
}

type UploadFile struct {
	UserId    int
	FileName  string
	FolderId  int64
	Content   []byte
	Size      int64
	CreatedAt time.Time
}

type Item interface {
	GetName() string
	GetSize() int64
	GetCreatedAt() time.Time
	GetID() int64
	IsFolder() bool
}

// Add methods to your existing structs
func (f Folder) GetName() string         { return f.FolderName }
func (f Folder) GetSize() int64          { return 0 } // Folders don't have size
func (f Folder) GetCreatedAt() time.Time { return f.CreatedAt }
func (f Folder) GetID() int64            { return f.Id }
func (f Folder) IsFolder() bool          { return true }

func (f File) GetName() string         { return f.FileName }
func (f File) GetSize() int64          { return f.Size }
func (f File) GetCreatedAt() time.Time { return f.CreatedAt }
func (f File) GetID() int64            { return f.Id }
func (f File) IsFolder() bool          { return false }
