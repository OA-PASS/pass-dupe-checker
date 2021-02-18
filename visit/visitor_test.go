package visit

import (
	"dupe-checker/model"
	"dupe-checker/retriever"
	"github.com/stretchr/testify/assert"
	"log"
	"net/http"
	"testing"
	"time"
)

func Test_VisitSimple(t *testing.T) {
	client := &http.Client{
		Timeout: 120 * time.Second,
	}

	maxSimultaneousReqs := 2

	underTest := Visitor{
		retriever:  retriever.New(client, "fedoraAdmin", "moo", "Test_VisitSimple"),
		semaphore:  make(chan int, maxSimultaneousReqs),
		uris:       make(chan string),
		containers: make(chan model.LdpContainer),
		errors:     make(chan error),
	}

	go underTest.visit()

	//wg := sync.WaitGroup{}
	//wg.Add(1)

	go func() {
		underTest.uris <- "http://fcrepo:8080/fcrepo/rest/funders"
		underTest.uris <- "http://fcrepo:8080/fcrepo/rest/repositoryCopies"
		underTest.uris <- "http://fcrepo:8080/fcrepo/rest/publishers"
		close(underTest.uris)
		//wg.Done()
	}()

	//wg.Wait()

	for result := range underTest.containers {
		assert.NotNil(t, result)
		assert.True(t, len(result.Contains()) > 0)
		assert.NotNil(t, result.Uri())

		ok, passResource := result.IsPassResource()

		assert.False(t, ok)
		assert.Equal(t, "", passResource)
	}

}

func TestVisitor_Walk(t *testing.T) {
	client := &http.Client{
		Timeout: 120 * time.Second,
	}

	maxSimultaneousReqs := 2

	underTest := Visitor{
		retriever:  retriever.New(client, "fedoraAdmin", "moo", "TestVisitor_Walk"),
		semaphore:  make(chan int, maxSimultaneousReqs),
		uris:       make(chan string),
		containers: make(chan model.LdpContainer),
		errors:     make(chan error),
	}

	underTest.Walk("http://fcrepo:8080/fcrepo/rest/repositoryCopies", func(container model.LdpContainer) bool {
		if ok, _ := container.IsPassResource(); !ok {
			log.Printf("visit: recursing non-PASS resource %s", container.Uri())
			return true
		}
		log.Printf("visit: terminal PASS resource %s", container.Uri())
		return false
	})
}
