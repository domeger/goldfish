package notes_test

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"

	"github.com/lucas-clemente/notes/notes"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type mockRepo struct {
	files map[string]string
}

func (r *mockRepo) ReadFile(path string) (io.ReadCloser, error) {
	if c, ok := r.files[path]; ok {
		return ioutil.NopCloser(bytes.NewBufferString(c)), nil
	}
	return nil, notes.NotFoundError{}
}
func (r *mockRepo) StoreFile(path string, reader io.Reader) error {
	c, err := ioutil.ReadAll(reader)
	if err != nil {
		return err
	}
	r.files[path] = string(c)
	return nil
}

var _ = Describe("Handler", func() {
	var (
		repo *mockRepo
		resp *httptest.ResponseRecorder
	)

	BeforeEach(func() {
		files := map[string]string{
			"/foo/bar.md": "foobar",
		}
		repo = &mockRepo{files: files}
		resp = httptest.NewRecorder()
	})

	It("GETs pages", func() {
		handler := notes.NewHandler(repo, "/v1")
		req, err := http.NewRequest("GET", "/v1/foo/bar.md", nil)
		Expect(err).To(BeNil())
		handler.ServeHTTP(resp, req)
		Expect(resp.Code).To(Equal(http.StatusOK))
		Expect(resp.Body.String()).To(Equal("foobar"))
	})

	It("404s", func() {
		handler := notes.NewHandler(repo, "/v1")
		req, err := http.NewRequest("GET", "/v1/noooooooo", nil)
		Expect(err).To(BeNil())
		handler.ServeHTTP(resp, req)
		Expect(resp.Code).To(Equal(http.StatusNotFound))
	})

	It("POST updates pages", func() {
		handler := notes.NewHandler(repo, "/v1")
		req, err := http.NewRequest("POST", "/v1/foo/bar.md", bytes.NewBufferString("foobaz"))
		Expect(err).To(BeNil())
		handler.ServeHTTP(resp, req)
		Expect(resp.Code).To(Equal(http.StatusNoContent))
		Expect(repo.files["/foo/bar.md"]).To(Equal("foobaz"))
	})

	It("POSTs new pages", func() {
		handler := notes.NewHandler(repo, "/v1")
		req, err := http.NewRequest("POST", "/v1/new", bytes.NewBufferString("foobaz"))
		Expect(err).To(BeNil())
		handler.ServeHTTP(resp, req)
		Expect(resp.Code).To(Equal(http.StatusNoContent))
		Expect(repo.files["/new"]).To(Equal("foobaz"))
	})
})
