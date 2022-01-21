package main

import (
	"bufio"
	"fmt"
	"github.com/bogem/id3v2/v2"
	"io"
	"log"
	_ "log"
	"net/http"
	"os"
	"strings"
)

///TEST SI FICHIER .AAC OU .MP3
func VerifMedia(Nomfichier string) bool {

	//utilisation de la fct HasSuffix() pour verifier si la fin de chaine est un mp3 ou AAC
	verif1 := strings.HasSuffix(Nomfichier, ".aac")
	if verif1 == true {
		fmt.Println("Fichier de type AAC")
		return verif1
	}
	verif2 := strings.HasSuffix(Nomfichier, ".mp3")
	if verif2 == true {
		fmt.Println("Fichier de type MP3")
	}
	return verif2
}

//SSPGRM POUR VERIF ERREUR
func VerifErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func VerifId(tag *id3v2.Tag) bool { //verif des ID non nul
	if len(tag.Title()) == 0 {
		return false
	} else if len(tag.Album()) == 0 {
		return false
	} else if len(tag.Artist()) == 0 {
		return false
	} else if len(tag.Year()) == 0 {
		return false
	} else if len(tag.Genre()) == 0 {
		return false
	} else {
		return true
	}
}

//LECTURE DES ID3TAGS
func LectureId(Nomfichier string) {
	fichier, err := id3v2.Open(Nomfichier, id3v2.Options{Parse: true})
	VerifErr(err)
	defer fichier.Close()

	if VerifId(fichier) == false {
		fmt.Println("ERREUR : Les id3tags ne sont pas present")
	} else {

		//Affichage des id3tags : https://pkg.go.dev/github.com/bogem/id3v2/v2 voir les diff fct avec id3v2
		fmt.Println("TITRE :", fichier.Title())
		fmt.Println("ALBUM :", fichier.Album())
		fmt.Println("ARTISTE :", fichier.Artist())
		fmt.Println("ANNEE DE PARUTION :", fichier.Year())
		fmt.Println("GENRE :", fichier.Genre())
	}

}

func ExistMedia(Nomfichier string) {
	if _, err := os.Stat(Nomfichier); os.IsNotExist(err) {
		// Le fichier n'existe pas
		fmt.Println("Le fichier n'est pas lisible")
	}
}

func LectureFichier(fichier string) []string {
	file, err := os.Open(fichier) //ouverture fichier
	VerifErr(err)
	var tab []string //creation tableau dynamique pour stocker les url des fichier a teleharger
	fileScanner := bufio.NewScanner(file)

	//lecture ligne par ligne
	for fileScanner.Scan() {
		tab = append(tab, fileScanner.Text())
	}
	if err := fileScanner.Err(); err != nil {
		log.Fatalf("Error while reading file: %s", err)
	}
	file.Close()
	return tab //on retourne le tableau
}

func TelechargerAudio(nv string, url string) error {
	resp, err := http.Get(url)
	VerifErr(err)
	defer resp.Body.Close()

	//Creation du nv fichier
	out, err := os.Create(nv)
	VerifErr(err)
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}

func main() {
	fichier := os.Args[1] //Nom du fichier à ouvrir go run
	tab := LectureFichier(fichier)
	for i := 0; i < len(tab); i++ {
		err := TelechargerAudio("audio.mp3", tab[i])
		if err != nil {
			panic(err)
		}
		if VerifMedia("audio.mp3") == false {
			fmt.Println("Mauvais fichier")
		} else {
			fmt.Print("Fichier Audio N°")
			fmt.Println(i + 1)
			ExistMedia("audio.mp3")
			LectureId("audio.mp3")
		}
	}
}
