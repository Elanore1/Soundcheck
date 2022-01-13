package main

import (
	"fmt"
	"github.com/bogem/id3v2/v2"
	"log"
	_ "log"
	"os"
	"strings"
)

///TEST SI FICHIER .AAC OU .MP3
func VerifMedia(Nomfichier string) bool {

	fmt.Println("Nom du fichier :", Nomfichier)
	//utiliszation de la fct HasSuffix() pour verifier si la fin de chaine est un mp3 ou AAC
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

func VerifId(tag *id3v2.Tag) bool {
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

func main() {
	fichier := os.Args[1] //Nom du fichier à ouvrir go run
	if VerifMedia(fichier) == false {
		fmt.Println("Mauvais fichier")
	} else {
		ExistMedia(fichier)
		LectureId(fichier)
	}
}
