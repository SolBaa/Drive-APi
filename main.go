package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
)

func getDriveService() (*drive.Service, error) {
	ctx := context.Background()
	b, err := ioutil.ReadFile("credentials.json")
	if err != nil {
		fmt.Printf("Unable to read credentials.json file. Err: %v\n", err)
		return nil, err
	}

	// If you want to modify this scope, delete your previously saved token.json.
	config, err := google.ConfigFromJSON(b, drive.DriveScope)

	if err != nil {
		return nil, err
	}

	client := getClient(config)

	service, err := drive.NewService(ctx, option.WithHTTPClient(client))

	if err != nil {
		fmt.Printf("Cannot create the Google Drive service: %v\n", err)
		return nil, err
	}

	return service, err
}

// Retrieve a token, saves the token, then returns the generated client.
func getClient(config *oauth2.Config) *http.Client {
	// The file token.json stores the user's access and refresh tokens, and is
	// created automatically when the authorization flow completes for the first
	// time.
	tokFile := "token.json"
	tok, err := tokenFromFile(tokFile)
	if err != nil {
		tok = getTokenFromWeb(config)
		saveToken(tokFile, tok)
	}
	return config.Client(context.Background(), tok)
}

// Retrieves a token from a local file.
func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

// Request a token from the web, then returns the retrieved token.
func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)

	fmt.Println("Paste Authorization code here :")
	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		log.Fatalf("Unable to read authorization code %v", err)
	}

	tok, err := config.Exchange(context.TODO(), authCode)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web %v", err)
	}
	return tok
}

// Saves a token to a file path.
func saveToken(path string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}

func createFile(service *drive.Service, name string, mimeType string, content io.Reader, parentId string) (*drive.File, error) {
	f := &drive.File{
		MimeType: mimeType,
		Name:     name,
		Parents:  []string{parentId},
	}
	file, err := service.Files.Create(f).Media(content).Do()

	if err != nil {
		log.Println("Could not create file: " + err.Error())
		return nil, err
	}

	return file, nil
}

func createFolder(service *drive.Service, name string, perm string) (*drive.File, error) {
	f := &drive.File{
		Name:     name,
		MimeType: "application/vnd.google-apps.folder",
	}
	folder, err := service.Files.Create(f).Do()
	if err != nil {
		log.Printf("Could not create folder (%s): %v", name, err.Error())
	}
	return folder, nil
}

//Print Link File in the console
func PrintFile(service *drive.Service, fileID string) error {
	_, err := service.Files.Get(fileID).Do()
	if err != nil {
		fmt.Printf("Sorry an error occurred: %v\n", err)
		return err
	}
	str := []string{"https://docs.google.com/presentation/d/", fileID}
	fmt.Println("\n", strings.Join(str, ""))
	return nil

}

func GetFilesID(service *drive.Service, name string) (ids []string, err error) {
	list, err := service.Files.List().Q("name='" + name + "'").Do()
	if err != nil {
		log.Printf("An error ocurred GetFilesID(%s): %v", name, err.Error())
		return nil, err
	}
	for _, f := range list.Files {
		ids = append(ids, f.Id)
	}
	return ids, nil
}

func DownloadFile(service *drive.Service, fileId string, mimeType string) error {
	_, err := service.Files.Export(fileId, mimeType).Download()
	if err != nil {
		log.Printf("Download error : %v", err.Error())
		return err
	}
	return nil
}

func RenameFile(service *drive.Service, fileId string, newName string) (*drive.File, error) {
	f := &drive.File{
		Name: newName,
	}
	r, err := service.Files.Update(fileId, f).Do()
	if err != nil {
		log.Printf("An error ocurred: %v", err.Error())
		return nil, err
	}
	return r, nil
}

func main() {

	// * Ejercicio 1 -> Hacer un código en go que cree un archivo llamado creado.docx y lo almacene en la carpeta “módulo”

	// Step 1: Open  file
	err := ioutil.WriteFile("creado.docx", []byte("Holiss"), 0755)
	if err != nil {
		log.Fatal(err)
	}

	f, err := os.Open("creado.docx")
	if err != nil {
		panic(fmt.Sprintf("cannot open file: %v", err))
	}

	defer f.Close()
	// Step 2: Get the Google Drive service
	srv, err := getDriveService()

	if err != nil {
		log.Fatal(err)
	}

	// Step 3: Create directory
	// dir, err := createFolder(srv, "New Folder", "root")

	if err != nil {
		panic(fmt.Sprintf("Could not create dir: %v\n", err))
	}

	//give your drive folder id here in which you want to upload or create a new directory
	folderId := "1jzeyuwpA8R-WG7mcdbB8eCIm1JF8i4rJ"

	// Step 4: create the file and upload
	file, err := createFile(srv, f.Name(), "application/octet-stream", f, folderId)

	if err != nil {
		panic(fmt.Sprintf("Could not create file: %v\n", err))
	}

	fmt.Printf("File '%s' uploaded successfully", file.Name)
	fmt.Printf("\nFile Id: '%s' \n", file.Id)
	// * Ejercicio 2 -> Hacer un código en go que traiga el link hacia diapo.pptx

	err = PrintFile(srv, "1PdQELlDjren_Yt7mJyx4cXkBytt5VlbbCnqHn_tgiI0")

	if err != nil {
		panic(fmt.Sprintf("Could not get link: %v", err.Error()))
	}
	// * Ejercicio 3 -> Hacer un código en go que edite el nombre de word.docx a texto.docx

	renamedFile, err := RenameFile(srv, "1pXhRLBmHmpVGbWqd4f8CNjz4kA9Bgq0OmtokW622gSk", "texto")
	fmt.Println("Rename word.docx file to:", renamedFile)
	if err != nil {
		panic(fmt.Sprintf("Could not rename file: %v", err.Error()))

	}

	//Download File
	err = DownloadFile(srv, "1PdQELlDjren_Yt7mJyx4cXkBytt5VlbbCnqHn_tgiI0", "application/vnd.oasis.opendocument.text")
	if err != nil {
		fmt.Printf("Couldn't download file: %v", err.Error())
		return
	}
}
