package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"google.golang.org/api/drive/v3"
)

func createFile(service *drive.Service, name string, mimeType string, content io.Reader, parentId string) (*drive.File, error) {
	f := &drive.File{
		MimeType: mimeType,
		Name:     name,
		Parents:  []string{parentId},
	}
	file, err := service.Files.Create(f).Media(content).Do()

	if err != nil {
		log.Println("No se pudo crear el archivo: " + err.Error())
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
		log.Printf("No se pudo crear carpeta (%s): %v", name, err.Error())
	}
	return folder, nil
}

//Print Link File into the console
func PrintFile(service *drive.Service, fileID string) error {
	_, err := service.Files.Get(fileID).Do()
	if err != nil {
		fmt.Printf("Ocurrió un error: %v\n", err)
		return err
	}
	str := []string{"https://docs.google.com/presentation/d/", fileID}
	fmt.Println("\n", strings.Join(str, ""))
	return nil

}

func GetFilesID(service *drive.Service, name string) (ids []string, err error) {
	list, err := service.Files.List().Q("name='" + name + "'").Do()
	if err != nil {
		log.Printf("Ocurrió un error GetFilesID(%s): %v", name, err.Error())
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
		log.Printf("Error de Descarga : %v", err.Error())
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
		log.Printf("Ocurrió un error: %v", err.Error())
		return nil, err
	}
	return r, nil
}

func main() {
	// * Ejercicio 1 -> Hacer un código en go que cree un archivo llamado creado.docx y lo almacene en la carpeta “módulo”

	//Escribimos el archivo si no existe lo crea
	err := ioutil.WriteFile("creado.docx", []byte("Holiss"), 0755)
	if err != nil {
		log.Fatal(err)
	}

	f, err := os.Open("creado.docx")
	if err != nil {
		panic(fmt.Sprintf("No se pudo abrir el archivo: %v", err))
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
		panic(fmt.Sprintf("No se pudo crear el directorio: %v\n", err))
	}

	//give your drive folder id here in which you want to upload or create a new directory
	folderId := "1jzeyuwpA8R-WG7mcdbB8eCIm1JF8i4rJ"

	// Step 4: create the file and upload
	file, err := createFile(srv, f.Name(), "application/octet-stream", f, folderId)

	if err != nil {
		panic(fmt.Sprintf("No se pudo crear el archivo: %v\n", err))
	}

	fmt.Printf("Archivo '%s' subido correctamente", file.Name)
	fmt.Printf("\nFile Id: '%s' \n", file.Id)

	// * Ejercicio 2 -> Hacer un código en go que traiga el link hacia diapo.pptx

	err = PrintFile(srv, "1PdQELlDjren_Yt7mJyx4cXkBytt5VlbbCnqHn_tgiI0")

	if err != nil {
		panic(fmt.Sprintf("No se pudo obtener el link: %v", err.Error()))
	}

	// * Ejercicio 3 -> Hacer un código en go que edite el nombre de word.docx a texto.docx

	renamedFile, err := RenameFile(srv, "1pXhRLBmHmpVGbWqd4f8CNjz4kA9Bgq0OmtokW622gSk", "texto")
	fmt.Println("Rename word.docx file to:", renamedFile.Name)
	if err != nil {
		panic(fmt.Sprintf("No se pudo renombrar el archivo: %v", err.Error()))

	}
	//<--- This is not working --->
	//Download File
	// err = DownloadFile(srv, "1PdQELlDjren_Yt7mJyx4cXkBytt5VlbbCnqHn_tgiI0", "application/vnd.oasis.opendocument.text")
	// if err != nil {
	// 	fmt.Printf("Couldn't download file: %v", err.Error())
	// 	return
	// }
}
