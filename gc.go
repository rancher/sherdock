package main

import (
	"bytes"
	"crypto/tls"
	"log"
	"os"

	"github.com/samalba/dockerclient"
)

// Callback used to listen to Docker's events
func eventCallback(event *dockerclient.Event, ec chan error, args ...interface{}) {
	log.Printf("Received event: %#v\n", *event)
}

func main1() {
	var tlsc tls.Config
	var certPath = os.Getenv("DOCKER_CERT_PATH")
	var certPem bytes.Buffer
	certPem.WriteString(certPath)
	certPem.WriteString("/cert.pem")

	var keyPem bytes.Buffer
	keyPem.WriteString(certPath)
	keyPem.WriteString("/key.pem")

	cert, err := tls.LoadX509KeyPair(certPem.String(), keyPem.String())
	tlsc.Certificates = append(tlsc.Certificates, cert)
	tlsc.InsecureSkipVerify = true
	docker, err := dockerclient.NewDockerClient(os.Getenv("DOCKER_HOST"), &tlsc)

	// list all the containers
	containers, err := docker.ListContainers(true, false, "")
	if err != nil {
		log.Fatal(err)
	}

	type usedFor string

	const used = usedFor("used")
	const unused = usedFor("unused")

	allImages := make(map[string]usedFor)
	parentImages := make(map[string]bool)

	images, err := docker.ListImages()
	if err != nil {
		log.Fatal(err)
	}

	for _, i := range images {
		allImages[i.Id] = unused
		parentImages[i.ParentId] = true
	}

	for _, c := range containers {
		info, _ := docker.InspectContainer(c.Id)
		allImages[info.Image] = used
	}

	imagesToDelete := make(map[string]bool)

	for image, i := range allImages {
		//FIXME - add only when image is a leaf (not a parent to any other image)
		if i == unused {
			imagesToDelete[image] = true
		}
	}

	for image, _ := range imagesToDelete {
		var removeMsg bytes.Buffer
		removeMsg.WriteString("Removing unused image: ")
		removeMsg.WriteString(image)
		log.Println(removeMsg.String())
		_, err := docker.RemoveImage(image)
		if err != nil {
			log.Fatal(err)
		}
	}
}
