package main

import (
	"log"

	"net/http"

	"github.com/cloudnautique/go-vol/volumes"
	"github.com/emicklei/go-restful"
	"github.com/emicklei/go-restful/swagger"
	"github.com/rancherio/sherdock/config"
	"github.com/rancherio/sherdock/containers"
	"github.com/rancherio/sherdock/images"
	"github.com/samalba/dockerclient"
)

type Response struct {
}

type DockerResource struct {
	url string
}

func (u DockerResource) Register(container *restful.Container) {
	ws := new(restful.WebService)
	ws.
		Path("/api/images").
		Doc("Show Images").
		Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON)

	ws.Route(ws.GET("/").To(u.getImages).
		Operation("findUser").
		Writes(Response{}))

	ws.Route(ws.DELETE("/{id}").To(u.deleteImage).
		Operation("findUser").
		Param(ws.PathParameter("id", "identifier of the image").DataType("string")).
		Writes(Response{}))

	container.Add(ws)

	ws = new(restful.WebService)
	ws.
		Path("/api/containers").
		Doc("Show Containers").
		Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON)

	ws.Route(ws.GET("/").To(u.getContainers).
		Operation("findUser").
		Writes(Response{}))

	container.Add(ws)

	ws = new(restful.WebService)
	ws.
		Path("/api/volumes").
		Doc("Show Volumes").
		Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON)

	ws.Route(ws.GET("/").To(u.getVolumes).
		Operation("findUser").
		Writes(Response{}))

	container.Add(ws)
}

func (u DockerResource) getImages(request *restful.Request, response *restful.Response) {

	// Init the client
	docker, err := dockerclient.NewDockerClient(u.url, nil)

	if err != nil {
		log.Fatal("Couldn't connect to docker client")
	}

	images, err := images.ListImagesDetailed(docker)
	if err != nil {
		log.Println(err)
		log.Fatal("Unable to fetch running containers")
	}
	response.WriteEntity(images)
}

func (u DockerResource) deleteImage(request *restful.Request, response *restful.Response) {

	id := request.PathParameter("id")

	docker, err := dockerclient.NewDockerClient(u.url, nil)

	if err != nil {
		log.Fatal("Couldn't connect to docker client")
	}

	_, err = docker.RemoveImage(id)

	if err != nil {
		response.WriteErrorString(http.StatusInternalServerError, err.Error())
	}
}

func (u DockerResource) getContainers(request *restful.Request, response *restful.Response) {

	// Init the client
	docker, err := dockerclient.NewDockerClient(u.url, nil)

	if err != nil {
		log.Fatal("Couldn't connect to docker client")
	}

	if request.QueryParameter("detailed") == "false" {
		containers, err := docker.ListContainers(true, false, "")
		if err != nil {
			log.Println(err)
			log.Fatal("Unable to fetch running containers")
		}
		response.WriteEntity(containers)
	} else {
		containers, err := containers.ListContainersDetailed(docker)
		if err != nil {
			log.Println(err)
			log.Fatal("Unable to fetch running containers")
		}
		response.WriteEntity(containers)
	}
}

func (u DockerResource) getVolumes(request *restful.Request, response *restful.Response) {

	vols := &volumes.Volumes{}
	err := vols.GetVolumes("/var/lib/docker/volumes")
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	response.WriteEntity(vols)

}

func main() {
	err := config.LoadGlobalConfig()
	if err != nil {
		log.Fatal("Failed to load config", err)
	}

	// to see what happens in the package, uncomment the following
	//restful.TraceLogger(log.New(os.Stdout, "[restful] ", log.LstdFlags|log.Lshortfile))

	wsContainer := restful.NewContainer()
	u := DockerResource{url: "unix:///var/run/docker.sock"}
	u.Register(wsContainer)

	cors := restful.CrossOriginResourceSharing{
		ExposeHeaders:  []string{"X-My-Header"},
		AllowedHeaders: []string{"Content-Type"},
		CookiesAllowed: false,
		Container:      wsContainer}
	wsContainer.Filter(cors.Filter)

	wsContainer.Filter(wsContainer.OPTIONSFilter)

	// Optionally, you can install the Swagger Service which provides a nice Web UI on your REST API
	// You need to download the Swagger HTML5 assets and change the FilePath location in the config below.
	// Open http://localhost:8080/apidocs and enter http://localhost:8080/apidocs.json in the api input field.
	config := swagger.Config{
		WebServices:    wsContainer.RegisteredWebServices(), // you control what services are visible
		WebServicesUrl: "http://localhost:8080",
		ApiPath:        "/apidocs.json",

		// Optionally, specifiy where the UI is located
		SwaggerPath:     "/apidocs/",
		SwaggerFilePath: "/Users/emicklei/xProjects/swagger-ui/dist"}
	swagger.RegisterSwaggerService(config, wsContainer)

	log.Printf("start listening on localhost:8080")
	server := &http.Server{Addr: ":8080", Handler: wsContainer}
	log.Fatal(server.ListenAndServe())
}
