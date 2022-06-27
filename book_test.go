package main

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetBook(t *testing.T) {
	testcases := []struct {
		desc   string
		input  *http.Request
		output []Book
	}{
		{"get books", httptest.NewRequest(http.MethodGet, "localhost:8000/book", nil), []Book{{1, "Godin", &Author{}, "Penguin", "11/02/1988"}}},
	}

	for i, tc := range testcases {
		w := httptest.NewRecorder()
		getBook(w, tc.input)
		//res := w.Result()

		res, _ := io.ReadAll(w.Result().Body)
		resBooks := []Book{}
		err := json.Unmarshal(res, &resBooks)
		if err != nil {
			return
		}

		for p := 1; p < len(resBooks); p++ {
			if resBooks[i] != tc.output[i] {
				t.Errorf("Error test case failed")
			}
		}

	}
	//defer res.Body.Close()
}

func TestGetBookByID(t *testing.T) {
	testcasesId := []struct {
		desc    string
		input2  string
		output2 Book
	}{
		{"get books", "localhost:8000/book/1", Book{1, "Godin", &Author{id: 1}, "Penguin", "12/06/1899"}},
	}

	for _, tc := range testcasesId {
		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, tc.input2, nil)

		getBookById(w, req)
		//res := w.Result()

		res, _ := io.ReadAll(w.Result().Body)
		resBooks := Book{}
		err := json.Unmarshal(res, &resBooks)

		if err != nil {
			return
		}

		if resBooks != tc.output2 {
			t.Errorf("test caase failed in get by id case")
		}
		//fmt.Println(b)
	}

}

func TestPostByBook(t *testing.T) {
	testcases := []struct {
		desc       string
		postinput  Book
		postoutput Book
		status     int
	}{
		{"Details", Book{1, "linchpin", &Author{}, "Penguin", "18/07/2000"}, Book{3, "Linchpin", &Author{}, "Penguin", "18/07/2002"}, 200},
		{desc: "Invalid publication", postinput: Book{id: 2, Title: "RD", Author: nil, Publication: "NA", PublishedDate: "12/03/2002"}, postoutput: Book{}, status: http.StatusBadRequest},
	}

	for i, tc := range testcases {
		rw := httptest.NewRecorder()
		body, _ := json.Marshal(tc)
		req := httptest.NewRequest(http.MethodPost, "localhost:8000/Book/", bytes.NewReader(body))

		PostByBook(rw, req)

		//defer rw.Result().Body.Close()

		if rw.Result().StatusCode != tc.status {
			t.Errorf("%v test case failed at %v, with status : %v", i, tc.desc, tc.status)
		}

		res, _ := io.ReadAll(rw.Result().Body)
		resBook := Book{}
		json.Unmarshal(res, &resBook)

		if resBook != tc.postoutput {
			t.Errorf("test case failed at %v : %v", i, tc.desc)
		}

	}

}

//
//func TestPostByAuthor(t *testing.T) {
//	testcases := []struct {
//		desc       string
//		postinput  Author
//		postoutput Author
//		status     int
//	}{
//		{"Valid details", Author{FirstName: "RD", LastName: "Sharma", Dob: "2/11/1989", PenName: "Sharma"}, Author{"RD", "Sharma", "2/11/1989", "Sharma"}, http.StatusOK},
//		{"InValid details", Author{FirstName: "", LastName: "Sharma", Dob: "2/11/1989", PenName: "Sharma"}, Author{}, http.StatusBadRequest},
//		{"Author already exists", Author{FirstName: "RD", LastName: "Sharma", Dob: "2/11/1989", PenName: "Sharma"}, Author{}, http.StatusBadRequest},
//	}
//	for i, tc := range testcases {
//		rw := httptest.NewRecorder()
//		body, _ := json.Marshal(tc.postinput)
//		req := httptest.NewRequest(http.MethodPost, "localhost:8000/Book/", bytes.NewReader(body))
//
//		PostByBook(rw, req)
//
//		//defer rw.Result().Body.Close()
//
//		if rw.Result().StatusCode != tc.status {
//			t.Errorf("%v test case failed at %v, with status : %v", i, tc.desc, tc.status)
//		}
//
//		res, _ := io.ReadAll(rw.Result().Body)
//		resBook := Author{}
//		json.Unmarshal(res, &resBook)
//
//		if resBook != tc.postoutput {
//			t.Errorf("test case failed at %v : %v", i, tc.desc)
//		}
//
//	}
//
//}

//func TestPutBookById(t *testing.T) {
//	testcases := []struct {
//		desc      string
//		id        int
//		InputById Book
//		Output    Book
//		status    int
//	}{{"Details", 1, Book{Title: "Linchpin", Author: nil, Publication: "Pengiun", PublishedDate: "11/03/2002"}, Book{}, 200},
//		{"Invalid Publication", 1, Book{Title: "Rakshit", Author: nil, Publication: "Mehul", PublishedDate: "11/03/2002"}, Book{}, http.StatusBadRequest},
//		{"Published date should be between 1880 and 2022", 1, Book{Title: "Check", Author: nil, Publication: "", PublishedDate: "1/1/1870"}, Book{}, http.StatusBadRequest},
//		{"Published date should be between 1880 and 2022", 1, Book{Title: "James", Author: nil, Publication: "", PublishedDate: "1/1/2222"}, Book{}, http.StatusBadRequest},
//		{"Author should exist", 1, Book{}, Book{}, http.StatusBadRequest},
//		{"Title can't be empty", 1, Book{Title: "", Author: nil, Publication: "", PublishedDate: ""}, Book{}, http.StatusBadRequest},
//	}
//
//	for i, tc := range testcases {
//		rw := httptest.NewRecorder()
//		body, _ := json.Marshal(tc.InputById)
//		req :=
//	}
//
//}
