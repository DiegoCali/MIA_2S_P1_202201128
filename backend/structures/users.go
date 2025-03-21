package structures

import (
	"backend/utils"
	"fmt"
	"strings"
)

func (spBlock *SuperBlock) ValidateUser(name string, pass string, path string) (bool, error) {
	// Look for /users.txt
	inode, err := spBlock.getUsersInode(path)
	if err != nil {
		return false, err
	}
	// Read Inode FBlocks
	fullString, err := spBlock.getInodeContent(inode, path)
	if err != nil {
		return false, err
	}
	// Get lines
	lines := strings.Split(fullString, "\n")
	// Iterate each line
	for i := 0; i < len(lines); i++ {
		// Get fields
		fields := strings.Split(lines[i], ",")
		// Check if fields are empty
		if len(fields) < 3 {
			break
		}
		if fields[0] == "0" {
			continue
		}
		// Check if user exists and password is correct
		if fields[1] == "U" && fields[2] == name && fields[4] == pass {
			return true, nil
		}
	}
	return false, nil
}

func (spBlock *SuperBlock) ValidateGroup(name string, path string) (bool, error) {
	inode, err := spBlock.getUsersInode(path)
	if err != nil {
		return false, err
	}
	fullString, err := spBlock.getInodeContent(inode, path)
	if err != nil {
		return false, err
	}
	lines := strings.Split(fullString, "\n")
	for i := 0; i < len(lines); i++ {
		fields := strings.Split(lines[i], ",")
		if len(fields) < 3 {
			break
		}
		if fields[0] == "0" {
			continue
		}
		if fields[1] == "G" && fields[2] == name {
			return true, nil
		}
	}
	return false, nil
}

func (spBlock *SuperBlock) ChangeGroup(name string, group string, path string) error {
	inode, err := spBlock.getUsersInode(path)
	if err != nil {
		return err
	}
	fullString, err := spBlock.getInodeContent(inode, path)
	if err != nil {
		return err
	}
	lines := strings.Split(fullString, "\n")
	newString := ""
	for i := 0; i < len(lines); i++ {
		fields := strings.Split(lines[i], ",")
		if len(fields) < 3 {
			break
		}
		if fields[0] == "0" {
			newString += strings.Join(fields, ",") + "\n"
			continue
		}
		if fields[1] == "U" && fields[2] == name {
			fmt.Println("Changing group to group: ", group)
			fields[3] = group
		}
		newString += strings.Join(fields, ",") + "\n"
	}
	fmt.Println(newString)
	err = spBlock.writeInode(inode, path, newString, int(spBlock.InodeStart+100)) // first inode is root, +100 is users.txt
	if err != nil {
		return err
	}
	return nil
}

func (spBlock *SuperBlock) CreateGroup(name string, path string) error {
	// Validate length of name
	if len(name) > 10 {
		return fmt.Errorf("name must be less than 10 characters")
	}
	groupText := "1,G," + name + "\n"
	inode, err := spBlock.getUsersInode(path)
	if err != nil {
		return err
	}
	fullString, err := spBlock.getInodeContent(inode, path)
	if err != nil {
		return err
	}
	fullString += groupText
	err = spBlock.writeInode(inode, path, fullString, int(spBlock.InodeStart+100)) // first inode is root, +100 is users.txt
	if err != nil {
		return err
	}
	return nil
}

func (spBlock *SuperBlock) CreateUser(name string, pass string, group string, path string) error {
	// Validate length of name and pass
	if len(name) > 10 || len(pass) > 10 {
		return fmt.Errorf("name and pass must be less than 10 characters")
	}
	userText := "1,U," + name + "," + group + "," + pass + "\n"
	inode, err := spBlock.getUsersInode(path)
	if err != nil {
		return err
	}
	fullString, err := spBlock.getInodeContent(inode, path)
	if err != nil {
		return err
	}
	fullString += userText
	err = spBlock.writeInode(inode, path, fullString, int(spBlock.InodeStart+100)) // first inode is root, +100 is users.txt
	if err != nil {
		return err
	}
	return nil
}

func (spBlock *SuperBlock) RemoveGroup(name string, path string) error {
	inode, err := spBlock.getUsersInode(path)
	if err != nil {
		return err
	}
	fullString, err := spBlock.getInodeContent(inode, path)
	if err != nil {
		return err
	}
	lines := strings.Split(fullString, "\n")
	newString := ""
	// To remove a group, we need to only chance the first number of the line to 0, fields[0] = 0
	for i := 0; i < len(lines); i++ {
		fields := strings.Split(lines[i], ",")
		if len(fields) < 3 {
			break
		}
		if fields[0] == "0" {
			newString += strings.Join(fields, ",") + "\n"
			continue
		}
		if fields[1] == "G" && fields[2] == name {
			fields[0] = "0"
		}
		newString += strings.Join(fields, ",") + "\n"
	}
	err = spBlock.writeInode(inode, path, newString, int(spBlock.InodeStart+100)) // first inode is root, +100 is users.txt
	if err != nil {
		return err
	}
	return nil
}

func (spBlock *SuperBlock) RemoveUser(name string, path string) error {
	inode, err := spBlock.getUsersInode(path)
	if err != nil {
		return err
	}
	fullString, err := spBlock.getInodeContent(inode, path)
	if err != nil {
		return err
	}
	lines := strings.Split(fullString, "\n")
	newString := ""
	// To remove a user, we need to only chance the first number of the line to 0, fields[0] = 0
	for i := 0; i < len(lines); i++ {
		fields := strings.Split(lines[i], ",")
		if len(fields) < 3 {
			break
		}
		if fields[0] == "0" {
			newString += strings.Join(fields, ",") + "\n"
			continue
		}
		if fields[1] == "U" && fields[2] == name {
			fields[0] = "0"
		}
		newString += strings.Join(fields, ",") + "\n"
	}
	err = spBlock.writeInode(inode, path, newString, int(spBlock.InodeStart+100)) // first inode is root, +100 is users.txt
	if err != nil {
		return err
	}
	return nil
}

func (spBlock *SuperBlock) getUsersInode(path string) (*Inode, error) {
	inode := &Inode{}
	err := utils.Deserialize(inode, path, int(spBlock.InodeStart+100)) // first inode is root, +100 is users.txt
	if err != nil {
		return nil, err
	}
	return inode, nil
}
