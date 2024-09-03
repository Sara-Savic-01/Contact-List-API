package middleware
import(
	"net/http"
	"net/http/httptest"
	"testing"
	
)
func TestAuthMiddleware(t *testing.T){
	t.Run("Valid token", func(t *testing.T){
		handler:=http.HandlerFunc(func(w http.ResponseWriter,r *http.Request){
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("OK"))
		})
		token:="Axf2FVAusahoXmKMLZih7LrhBwmYLVmyLDiMoYizPGReJTKEaseAb12oGYvbLleS"
		middleware:=AuthMiddleware(token, handler)
		req, err:=http.NewRequest("GET","/lists", nil)
		if err!=nil{
			t.Fatalf("Could not create request:%v",err)
		}
		req.Header.Set("Authorization", "Bearer "+token)
		rr:=httptest.NewRecorder()
		middleware.ServeHTTP(rr, req)
		if status:=rr.Code; status!=http.StatusOK{
			t.Errorf("Expected status code %d, got %d",http.StatusOK, status)
		}
		if body:=rr.Body.String(); body!="OK"{
			t.Errorf("Expected body 'OK', got '%s'", body)
		}
	})
	t.Run("Invalid token", func(t *testing.T){
		handler:=http.HandlerFunc(func(w http.ResponseWriter,r *http.Request){
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("OK"))
		})
		token:="Axf2FVAusahoXmKMLZih7LrhBwmYLVmyLDiMoYizPGReJTKEaseAb12oGYvbLleS"
		middleware:=AuthMiddleware(token, handler)
		req, err:=http.NewRequest("GET","/lists", nil)
		if err!=nil{
			t.Fatalf("Could not create request:%v",err)
		}
		req.Header.Set("Authorization", "Bearer invalid token")
		rr:=httptest.NewRecorder()
		middleware.ServeHTTP(rr, req)
		if status:=rr.Code; status!=http.StatusForbidden{
			t.Errorf("Expected status code %d, got %d",http.StatusForbidden, status)
		}
		if body:=rr.Body.String(); body!="Forbidden\n"{
			t.Errorf("Expected body 'OK', got '%s'", body)
		}
	})
	t.Run("Missing Authorization header", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		    w.WriteHeader(http.StatusOK)
		    w.Write([]byte("OK"))
		})
		token := "Axf2FVAusahoXmKMLZih7LrhBwmYLVmyLDiMoYizPGReJTKEaseAb12oGYvbLleS"
		middleware := AuthMiddleware(token, handler)
		req, err := http.NewRequest("GET", "/lists", nil)
		if err != nil {
		    t.Fatalf("Could not create request: %v", err)
		}
		
		rr := httptest.NewRecorder()
		middleware.ServeHTTP(rr, req)
		if status := rr.Code; status != http.StatusUnauthorized {
		    t.Errorf("Expected status code %d, got %d", http.StatusUnauthorized, status)
		}
		if body := rr.Body.String(); body != "Unauthorized\n" {
		    t.Errorf("Expected body 'Unauthorized', got '%s'", body)
		}
    	})
}
	