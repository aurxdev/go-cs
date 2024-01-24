package main

import (
	"context"
	"flag"
	"log"
	"time"

	"github.com/gotk3/gotk3/gtk"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	pb "google.golang.org/grpc/examples/helloworld/helloworld"
)

var addr = flag.String("addr", "localhost:50051", "the address to connect to")

// sendButton est appelée lors du clic sur le bouton
func sendButton(builder *gtk.Builder, c pb.GreeterClient) {
	// Récupération de l'objet GtkEntry
	entryObj, err := builder.GetObject("nomEntry")
	if err != nil {
		log.Fatal("Erreur lors de la récupération de GtkEntry :", err)
	}
	entryText, err := entryObj.(*gtk.Entry).GetText()
	if err != nil {
		log.Fatal("Erreur lors de la récupération du texte de GtkEntry :", err)
	}
	// Récupération de l'objet GtkComboBoxText
	comboBoxObj, err := builder.GetObject("listeEntry")
	if err != nil {
		log.Fatal("Erreur lors de la récupération de GtkComboBoxText :", err)
	}

	// Vérification que l'objet est bien un GtkComboBoxText
	comboBox, ok := comboBoxObj.(*gtk.ComboBoxText)
	if !ok {
		log.Fatal("Erreur lors de la conversion de l'objet en *gtk.ComboBoxText")
	}

	// Récupération du texte actif de GtkComboBoxText
	comboBoxText := comboBox.GetActiveText()

	// on contacte le serveur
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// Création d'une requête avec les données à envoyer
	req := &pb.HelloRequest{
		Name: entryText,
		List: comboBoxText,
	}

	// Appel de la méthode RPC SayHello
	r, err := c.SayHello(ctx, req)
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}

	log.Printf("%s", r.GetMessage())
} // sendButton

// fonction principale
func main() {
	// on initialise GTK
	gtk.Init(nil)
	builder, err := gtk.BuilderNew()
	if err != nil {
		log.Fatal("Erreur lors de la création du builder :", err)
	}

	err = builder.AddFromFile("fenetre.glade")
	if err != nil {
		log.Fatal("Erreur lors du chargement du fichier Glade :", err)
	}

	obj, err := builder.GetObject("main_window")
	if err != nil {
		log.Fatal("Erreur lors de la récupération de la fenêtre principale :", err)
	}

	win, ok := obj.(*gtk.Window)
	if !ok {
		log.Fatal("Erreur lors de la conversion de l'objet en *gtk.Window")
	}

	// initialisation serveur

	flag.Parse()
	// Set up a connection to the server.
	conn, err := grpc.Dial(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewGreeterClient(conn)

	// Récupération du bouton
	buttonObj, err := builder.GetObject("sendButton")
	if err != nil {
		log.Fatal("Erreur lors de la récupération du bouton :", err)
	}

	buttonObj.(*gtk.Button).Connect("clicked", func() {
		sendButton(builder, c)
	})

	win.Connect("destroy", func() {
		gtk.MainQuit()
	})

	win.ShowAll() // on affiche la fenêtre principale

	gtk.Main() // on lance la boucle principale
} // main
