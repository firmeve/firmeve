package server

//func TestAdminOnly(t *testing.T) {
//	r := mux.NewRouter()
//	r.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
//		writer.Write([]byte("abc"))
//	}).Name("abc").Methods("Get")
//
//	http.ListenAndServe(":11221",r)
//}